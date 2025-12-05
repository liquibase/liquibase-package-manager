package org.liquibase.lpm.service;

import org.liquibase.lpm.exception.*;
import org.liquibase.lpm.model.*;
import org.liquibase.lpm.service.LiquibaseDetector.LiquibaseInfo;
import org.semver4j.Semver;

import java.util.List;
import java.util.Optional;

/**
 * Main service for package management operations.
 * <p>
 * Coordinates between the package registry, classpath service, and Liquibase detector
 * to provide high-level package operations.
 */
public class PackageService {

    private final ClasspathService classpathService;
    private final LiquibaseDetector liquibaseDetector;

    private PackageRegistry packageRegistry;
    private LiquibaseInfo liquibaseInfo;

    /**
     * Creates a new package service with default dependencies.
     */
    public PackageService() {
        this(new ClasspathService(), new LiquibaseDetector());
    }

    /**
     * Creates a new package service with custom dependencies.
     *
     * @param classpathService  the classpath service
     * @param liquibaseDetector the Liquibase detector
     */
    public PackageService(ClasspathService classpathService, LiquibaseDetector liquibaseDetector) {
        this.classpathService = classpathService;
        this.liquibaseDetector = liquibaseDetector;
    }

    /**
     * Initializes the service by detecting Liquibase and loading packages.
     *
     * @param global whether to use global classpath
     */
    public void initialize(boolean global) {
        // Detect Liquibase installation
        liquibaseInfo = liquibaseDetector.detect();

        // Set up classpath
        if (global) {
            classpathService.useGlobalClasspath(liquibaseInfo.homePath());
        } else {
            classpathService.useLocalClasspath();
        }

        // Ensure packages.json is in classpath
        if (!classpathService.packagesJsonExists()) {
            // Copy embedded packages.json to classpath
            packageRegistry = PackageRegistry.loadFromResource();
            classpathService.writePackagesJson(packageRegistry.toJson());
        } else {
            // Load from classpath
            packageRegistry = PackageRegistry.loadFromBytes(classpathService.readPackagesJson());
        }
    }

    /**
     * Gets the package registry.
     *
     * @return the package registry
     */
    public PackageRegistry getPackageRegistry() {
        return packageRegistry;
    }

    /**
     * Gets the Liquibase info.
     *
     * @return the Liquibase info
     */
    public LiquibaseInfo getLiquibaseInfo() {
        return liquibaseInfo;
    }

    /**
     * Gets the classpath service.
     *
     * @return the classpath service
     */
    public ClasspathService getClasspathService() {
        return classpathService;
    }

    /**
     * Gets the list of filenames in the classpath.
     *
     * @return the classpath filenames
     */
    public List<String> getClasspathFilenames() {
        return classpathService.getClasspathFilenames();
    }

    /**
     * Resolves a package specification and installs it.
     *
     * @param packageSpec the package specification (name or name@version)
     * @return the installed version
     */
    public PackageVersion resolveAndInstall(String packageSpec) {
        Dependency dep = Dependency.fromSpec(packageSpec);
        return resolveAndInstall(dep);
    }

    /**
     * Resolves a dependency and installs it.
     *
     * @param dependency the dependency to install
     * @return the installed version
     */
    public PackageVersion resolveAndInstall(Dependency dependency) {
        // Find the package
        Package pkg = packageRegistry.getByName(dependency.name())
                .orElseThrow(() -> new PackageNotFoundException(dependency.name()));

        // Get the version to install
        PackageVersion version;
        if (dependency.hasVersion()) {
            // Specific version requested
            version = pkg.getVersion(dependency.version())
                    .orElseThrow(() -> new VersionNotFoundException(dependency.name(), dependency.version()));
        } else {
            // Get latest compatible version
            version = pkg.getLatestVersion(liquibaseInfo.version())
                    .orElseThrow(() -> new VersionIncompatibleException(
                            dependency.name(), "any",
                            "unknown", liquibaseInfo.version().toString()));
        }

        // Check Liquibase compatibility for non-driver packages
        if (!Package.CATEGORY_DRIVER.equals(pkg.category())) {
            checkCompatibility(pkg.name(), version);
        }

        // Check if already installed
        Optional<PackageVersion> installed = pkg.getInstalledVersion(getClasspathFilenames());
        if (installed.isPresent()) {
            throw new PackageAlreadyInstalledException(pkg.name(), installed.get().tag());
        }

        // Install
        classpathService.installVersion(version);

        return version;
    }

    /**
     * Removes a package from the classpath.
     *
     * @param packageName the package name to remove
     * @return the removed version
     */
    public PackageVersion removePackage(String packageName) {
        // Find the package
        Package pkg = packageRegistry.getByName(packageName)
                .orElseThrow(() -> new PackageNotFoundException(packageName));

        // Check if installed
        PackageVersion installed = pkg.getInstalledVersion(getClasspathFilenames())
                .orElseThrow(() -> new PackageNotInstalledException(packageName));

        // Remove
        classpathService.removeVersion(installed);

        return installed;
    }

    /**
     * Gets all outdated packages.
     *
     * @return registry of outdated packages
     */
    public PackageRegistry getOutdatedPackages() {
        return packageRegistry.getOutdated(liquibaseInfo.version(), getClasspathFilenames());
    }

    /**
     * Gets all installed packages.
     *
     * @return registry of installed packages
     */
    public PackageRegistry getInstalledPackages() {
        return packageRegistry.getInstalled(getClasspathFilenames());
    }

    /**
     * Upgrades a package to the latest compatible version.
     *
     * @param packageName the package to upgrade
     * @return the new version
     */
    public PackageVersion upgradePackage(String packageName) {
        Package pkg = packageRegistry.getByName(packageName)
                .orElseThrow(() -> new PackageNotFoundException(packageName));

        // Get currently installed version
        PackageVersion installed = pkg.getInstalledVersion(getClasspathFilenames())
                .orElseThrow(() -> new PackageNotInstalledException(packageName));

        // Get latest version
        PackageVersion latest = pkg.getLatestVersion(liquibaseInfo.version())
                .orElseThrow(() -> new VersionIncompatibleException(
                        packageName, "latest",
                        "unknown", liquibaseInfo.version().toString()));

        // Remove old version
        classpathService.removeVersion(installed);

        // Install new version
        classpathService.installVersion(latest);

        return latest;
    }

    /**
     * Updates the packages.json manifest from a URL or file path.
     *
     * @param source the source URL or file path
     */
    public void updateManifest(String source) {
        byte[] content;

        if (source.toLowerCase().startsWith("http")) {
            DownloadService downloadService = new DownloadService();
            content = downloadService.download(source);
        } else {
            try {
                content = java.nio.file.Files.readAllBytes(java.nio.file.Path.of(source));
            } catch (java.io.IOException e) {
                throw new ManifestException("Failed to read manifest from: " + source, e);
            }
        }

        // Validate the content
        PackageRegistry newRegistry = PackageRegistry.loadFromBytes(content);
        if (!newRegistry.isValid()) {
            throw new ManifestException("Invalid packages.json content");
        }

        // Write to classpath
        classpathService.writePackagesJson(content);

        // Reload registry
        packageRegistry = newRegistry;
    }

    /**
     * Checks if a version is compatible with the installed Liquibase.
     */
    private void checkCompatibility(String packageName, PackageVersion version) {
        String requiredCore = version.liquibaseCore();
        if (requiredCore == null || requiredCore.isBlank()) {
            return; // No requirement specified
        }

        try {
            Semver required = new Semver(requiredCore);
            if (liquibaseInfo.version().isLowerThan(required)) {
                throw new VersionIncompatibleException(
                        packageName, version.tag(),
                        requiredCore, liquibaseInfo.version().toString());
            }
        } catch (VersionIncompatibleException e) {
            throw e;
        } catch (Exception e) {
            // Invalid semver, allow it
        }
    }

    /**
     * Checks if the Liquibase version requires JAVA_OPTS for local classpath.
     *
     * @return true if JAVA_OPTS instructions should be shown
     */
    public boolean requiresJavaOptsInstructions() {
        try {
            Semver threshold = new Semver("4.6.2");
            return liquibaseInfo.version().isLowerThan(threshold);
        } catch (Exception e) {
            return false;
        }
    }

    /**
     * Filters the package registry by category.
     *
     * @param category the category to filter by
     * @return filtered registry
     */
    public PackageRegistry filterByCategory(String category) {
        return packageRegistry.filterByCategory(category);
    }

    /**
     * Searches packages by name.
     *
     * @param searchTerm the search term
     * @return matching packages
     */
    public PackageRegistry searchPackages(String searchTerm) {
        return packageRegistry.filterByName(searchTerm);
    }
}
