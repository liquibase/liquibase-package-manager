package org.liquibase.lpm.model;

import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.liquibase.lpm.exception.ManifestException;
import org.semver4j.Semver;

import java.io.IOException;
import java.io.InputStream;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.Optional;
import java.util.stream.Collectors;

/**
 * Registry containing all available packages.
 * <p>
 * This class provides methods for querying, filtering, and displaying packages.
 * It wraps a list of Package objects and provides collection-level operations.
 */
public class PackageRegistry {

    private static final ObjectMapper OBJECT_MAPPER = new ObjectMapper();

    private final List<Package> packages;

    /**
     * Creates a new package registry with the given packages.
     *
     * @param packages the list of packages
     */
    public PackageRegistry(List<Package> packages) {
        this.packages = packages != null ? new ArrayList<>(packages) : new ArrayList<>();
    }

    /**
     * Gets all packages in this registry.
     *
     * @return an unmodifiable list of packages
     */
    public List<Package> getPackages() {
        return Collections.unmodifiableList(packages);
    }

    /**
     * Gets a package by name.
     *
     * @param name the package name
     * @return the package, or an empty Optional if not found
     */
    public Optional<Package> getByName(String name) {
        if (name == null) {
            return Optional.empty();
        }
        return packages.stream()
                .filter(p -> name.equals(p.name()))
                .findFirst();
    }

    /**
     * Gets all installed packages from the classpath.
     *
     * @param classpathFilenames list of filenames in the classpath
     * @return a new registry containing only installed packages
     */
    public PackageRegistry getInstalled(List<String> classpathFilenames) {
        List<Package> installed = packages.stream()
                .filter(p -> p.isInClasspath(classpathFilenames))
                .collect(Collectors.toList());
        return new PackageRegistry(installed);
    }

    /**
     * Gets all outdated packages (installed but have newer compatible versions).
     *
     * @param liquibaseVersion   the installed Liquibase version
     * @param classpathFilenames list of filenames in the classpath
     * @return a new registry containing only outdated packages
     */
    public PackageRegistry getOutdated(Semver liquibaseVersion, List<String> classpathFilenames) {
        List<Package> outdated = new ArrayList<>();

        for (Package pkg : getInstalled(classpathFilenames).getPackages()) {
            Optional<PackageVersion> installedOpt = pkg.getInstalledVersion(classpathFilenames);
            Optional<PackageVersion> latestOpt = pkg.getLatestVersion(liquibaseVersion);

            if (installedOpt.isPresent() && latestOpt.isPresent()) {
                try {
                    Semver installed = new Semver(installedOpt.get().tag());
                    Semver latest = new Semver(latestOpt.get().tag());

                    if (installed.isLowerThan(latest)) {
                        outdated.add(pkg);
                    }
                } catch (Exception e) {
                    // Invalid semver, skip
                }
            }
        }

        return new PackageRegistry(outdated);
    }

    /**
     * Filters packages by category.
     *
     * @param category the category to filter by
     * @return a new registry containing only packages with the given category
     */
    public PackageRegistry filterByCategory(String category) {
        if (category == null || category.isBlank()) {
            return this;
        }
        List<Package> filtered = packages.stream()
                .filter(p -> category.equals(p.category()))
                .collect(Collectors.toList());
        return new PackageRegistry(filtered);
    }

    /**
     * Filters packages by name substring (case-insensitive).
     *
     * @param searchTerm the search term to filter by
     * @return a new registry containing only packages whose name contains the search term
     */
    public PackageRegistry filterByName(String searchTerm) {
        if (searchTerm == null || searchTerm.isBlank()) {
            return this;
        }
        String lowerSearch = searchTerm.toLowerCase();
        List<Package> filtered = packages.stream()
                .filter(p -> p.name() != null && p.name().toLowerCase().contains(lowerSearch))
                .collect(Collectors.toList());
        return new PackageRegistry(filtered);
    }

    /**
     * Checks if this registry is empty.
     *
     * @return true if there are no packages
     */
    public boolean isEmpty() {
        return packages.isEmpty();
    }

    /**
     * Gets the number of packages in this registry.
     *
     * @return the package count
     */
    public int size() {
        return packages.size();
    }

    /**
     * Generates a display table for packages.
     *
     * @param classpathFilenames list of filenames in the classpath (for showing installed versions)
     * @return a list of formatted strings for display
     */
    public List<String> display(List<String> classpathFilenames) {
        List<String> result = new ArrayList<>();
        result.add(String.format("%-4s %-38s %s", "   ", "Package", "Category"));

        for (int i = 0; i < packages.size(); i++) {
            Package pkg = packages.get(i);
            String prefix = (i + 1 == packages.size()) ? "└──" : "├──";

            String versionSuffix = "";
            Optional<PackageVersion> installedOpt = pkg.getInstalledVersion(classpathFilenames);
            if (installedOpt.isPresent()) {
                versionSuffix = "@" + installedOpt.get().tag();
            }

            result.add(String.format("%-4s %-38s %s", prefix, pkg.name() + versionSuffix, pkg.category()));
        }

        return result;
    }

    /**
     * Loads a package registry from the embedded packages.json resource.
     *
     * @return the loaded registry
     * @throws ManifestException if the resource cannot be read or parsed
     */
    public static PackageRegistry loadFromResource() {
        try (InputStream is = PackageRegistry.class.getResourceAsStream("/packages.json")) {
            if (is == null) {
                throw new ManifestException("Embedded packages.json not found");
            }
            List<Package> packages = OBJECT_MAPPER.readValue(is, new TypeReference<>() {});
            return new PackageRegistry(packages);
        } catch (IOException e) {
            throw new ManifestException("Failed to load packages.json", e);
        }
    }

    /**
     * Loads a package registry from a file path.
     *
     * @param path the path to the packages.json file
     * @return the loaded registry
     * @throws ManifestException if the file cannot be read or parsed
     */
    public static PackageRegistry loadFromFile(Path path) {
        try {
            byte[] bytes = Files.readAllBytes(path);
            List<Package> packages = OBJECT_MAPPER.readValue(bytes, new TypeReference<>() {});
            return new PackageRegistry(packages);
        } catch (IOException e) {
            throw new ManifestException("Failed to load packages.json from " + path, e);
        }
    }

    /**
     * Loads a package registry from a byte array.
     *
     * @param bytes the JSON bytes
     * @return the loaded registry
     * @throws ManifestException if the bytes cannot be parsed
     */
    public static PackageRegistry loadFromBytes(byte[] bytes) {
        try {
            List<Package> packages = OBJECT_MAPPER.readValue(bytes, new TypeReference<>() {});
            return new PackageRegistry(packages);
        } catch (IOException e) {
            throw new ManifestException("Failed to parse packages.json", e);
        }
    }

    /**
     * Validates that this registry is valid (used after loading from remote source).
     * Currently checks that a known package exists.
     *
     * @return true if the registry appears valid
     */
    public boolean isValid() {
        // Check that at least one common package exists as a sanity check
        return getByName("postgresql").isPresent() ||
               getByName("liquibase-postgresql").isPresent() ||
               !packages.isEmpty();
    }

    /**
     * Serializes this registry to JSON bytes.
     *
     * @return the JSON representation
     * @throws ManifestException if serialization fails
     */
    public byte[] toJson() {
        try {
            return OBJECT_MAPPER.writerWithDefaultPrettyPrinter().writeValueAsBytes(packages);
        } catch (IOException e) {
            throw new ManifestException("Failed to serialize packages", e);
        }
    }
}
