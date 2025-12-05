package org.liquibase.lpm.model;

import com.fasterxml.jackson.annotation.JsonProperty;
import org.semver4j.Semver;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.Optional;

/**
 * Represents a package in the package registry.
 * <p>
 * A package has a name, category (extension, driver, or pro), and a list of
 * available versions.
 *
 * @param name     the package name (e.g., "liquibase-postgresql")
 * @param category the package category ("extension", "driver", or "pro")
 * @param versions the list of available versions
 */
public record Package(
        @JsonProperty("name") String name,
        @JsonProperty("category") String category,
        @JsonProperty("versions") List<PackageVersion> versions
) {

    /**
     * Category constant for extension packages.
     */
    public static final String CATEGORY_EXTENSION = "extension";

    /**
     * Category constant for driver packages.
     */
    public static final String CATEGORY_DRIVER = "driver";

    /**
     * Category constant for pro packages.
     */
    public static final String CATEGORY_PRO = "pro";

    /**
     * Returns a defensive copy of versions.
     */
    @Override
    public List<PackageVersion> versions() {
        return versions == null ? Collections.emptyList() : Collections.unmodifiableList(versions);
    }

    /**
     * Gets the latest version compatible with the given Liquibase version.
     * <p>
     * For drivers, all versions are considered compatible (Liquibase version is ignored).
     * For extensions and pro packages, only versions where the package's liquibaseCore
     * requirement is less than or equal to the installed Liquibase version are considered.
     *
     * @param liquibaseVersion the installed Liquibase version
     * @return the latest compatible version, or an empty Optional if none found
     */
    public Optional<PackageVersion> getLatestVersion(Semver liquibaseVersion) {
        if (versions == null || versions.isEmpty()) {
            return Optional.empty();
        }

        Semver latestVersionSemver = new Semver("0.0.0");
        PackageVersion latestVersion = null;

        for (PackageVersion v : versions) {
            // For non-driver packages, check Liquibase compatibility
            if (!CATEGORY_DRIVER.equals(category)) {
                String requiredCore = v.liquibaseCore();
                if (requiredCore != null && !requiredCore.isBlank()) {
                    try {
                        Semver required = new Semver(requiredCore);
                        if (liquibaseVersion.isLowerThan(required)) {
                            continue; // Skip versions that require newer Liquibase
                        }
                    } catch (Exception e) {
                        // Invalid semver, skip this version
                        continue;
                    }
                }
            }

            // Compare version tags
            try {
                Semver versionSemver = new Semver(v.tag());
                if (versionSemver.isGreaterThan(latestVersionSemver)) {
                    latestVersionSemver = versionSemver;
                    latestVersion = v;
                }
            } catch (Exception e) {
                // Invalid semver tag, skip this version
            }
        }

        return Optional.ofNullable(latestVersion);
    }

    /**
     * Gets a specific version by tag.
     *
     * @param versionTag the version tag to find
     * @return the matching version, or an empty Optional if not found
     */
    public Optional<PackageVersion> getVersion(String versionTag) {
        if (versions == null || versionTag == null) {
            return Optional.empty();
        }
        return versions.stream()
                .filter(v -> versionTag.equals(v.tag()))
                .findFirst();
    }

    /**
     * Gets the installed version from the classpath files.
     *
     * @param classpathFilenames list of filenames in the classpath
     * @return the installed version, or an empty Optional if not installed
     */
    public Optional<PackageVersion> getInstalledVersion(List<String> classpathFilenames) {
        if (versions == null || classpathFilenames == null) {
            return Optional.empty();
        }
        return versions.stream()
                .filter(v -> v.isInClasspathByFilename(classpathFilenames))
                .findFirst();
    }

    /**
     * Checks if any version of this package is installed in the classpath.
     *
     * @param classpathFilenames list of filenames in the classpath
     * @return true if any version is installed
     */
    public boolean isInClasspath(List<String> classpathFilenames) {
        return getInstalledVersion(classpathFilenames).isPresent();
    }

    /**
     * Returns a new list of versions with the specified version removed.
     *
     * @param versionToDelete the version to remove
     * @return a new list without the specified version
     */
    public List<PackageVersion> deleteVersion(PackageVersion versionToDelete) {
        if (versions == null || versionToDelete == null) {
            return versions();
        }
        List<PackageVersion> result = new ArrayList<>(versions);
        result.removeIf(v -> versionToDelete.tag().equals(v.tag()));
        return Collections.unmodifiableList(result);
    }

    /**
     * Checks if this is an empty/null package (used for "not found" scenarios).
     *
     * @return true if the name is null or empty
     */
    public boolean isEmpty() {
        return name == null || name.isBlank();
    }

    /**
     * Creates an empty package (used as a sentinel value).
     *
     * @return an empty Package
     */
    public static Package empty() {
        return new Package(null, null, null);
    }
}
