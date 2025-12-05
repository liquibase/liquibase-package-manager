package org.liquibase.lpm.commands;

import org.liquibase.lpm.commands.mixins.GlobalMixin;
import org.liquibase.lpm.model.Package;
import org.liquibase.lpm.model.PackageRegistry;
import org.liquibase.lpm.model.PackageVersion;
import org.liquibase.lpm.service.PackageService;
import org.liquibase.lpm.util.ConsoleOutput;
import org.semver4j.Semver;
import picocli.CommandLine.Command;
import picocli.CommandLine.Mixin;
import picocli.CommandLine.ParentCommand;

import java.util.ArrayList;
import java.util.Comparator;
import java.util.List;
import java.util.concurrent.Callable;

/**
 * Command to remove duplicate versions of installed packages.
 * <p>
 * When multiple versions of the same package are installed, this command
 * removes all but the latest version.
 */
@Command(
        name = "dedupe",
        description = "Deduplicate Packages"
)
public class DedupeCommand implements Callable<Integer> {

    @ParentCommand
    private RootCommand parent;

    @Mixin
    private GlobalMixin globalMixin;

    private final PackageService packageService;

    /**
     * Creates a new dedupe command.
     */
    public DedupeCommand() {
        this.packageService = new PackageService();
    }

    /**
     * Creates a dedupe command with a custom package service (for testing).
     *
     * @param packageService the package service to use
     */
    DedupeCommand(PackageService packageService) {
        this.packageService = packageService;
    }

    @Override
    public Integer call() {
        try {
            boolean global = globalMixin != null && globalMixin.isGlobal();
            boolean dryRun = globalMixin != null && globalMixin.isDryRun();

            // Initialize service
            packageService.initialize(global);

            // Get installed packages
            PackageRegistry installed = packageService.getInstalledPackages();

            if (installed.isEmpty()) {
                System.out.println("No packages installed.");
                return 0;
            }

            List<String> classpathFilenames = packageService.getClasspathFilenames();
            boolean foundDuplicates = false;

            // Check each package for duplicates
            for (Package pkg : installed.getPackages()) {
                // Find all installed versions of this package
                List<PackageVersion> installedVersions = findAllInstalledVersions(pkg, classpathFilenames);

                if (installedVersions.size() <= 1) {
                    continue; // No duplicates
                }

                foundDuplicates = true;

                // Sort by version (descending) - keep the highest
                installedVersions.sort((v1, v2) -> {
                    try {
                        Semver s1 = new Semver(v1.tag());
                        Semver s2 = new Semver(v2.tag());
                        return s2.compareTo(s1); // Descending order
                    } catch (Exception e) {
                        return v2.tag().compareTo(v1.tag());
                    }
                });

                // Display duplicates
                List<String> output = ConsoleOutput.formatVersionsForDedupe(pkg.name(), installedVersions);
                ConsoleOutput.printLines(output);
                System.out.println();

                // Remove all except the first (highest version)
                if (!dryRun) {
                    for (int i = 1; i < installedVersions.size(); i++) {
                        PackageVersion toRemove = installedVersions.get(i);
                        packageService.getClasspathService().removeVersion(toRemove);
                        ConsoleOutput.printRemoveSuccess(toRemove.getFilename());
                    }
                }
            }

            if (!foundDuplicates) {
                System.out.println("No duplicate packages found.");
            } else {
                // Refresh classpath cache
                packageService.getClasspathService().refreshClasspathFilenames();
                System.out.println("Deduplication complete.");
            }

            return 0;
        } catch (Exception e) {
            System.err.println(e.getMessage());
            return 1;
        }
    }

    /**
     * Finds all installed versions of a package in the classpath.
     */
    private List<PackageVersion> findAllInstalledVersions(Package pkg, List<String> classpathFilenames) {
        List<PackageVersion> installed = new ArrayList<>();

        for (PackageVersion version : pkg.versions()) {
            if (version.isInClasspathByFilename(classpathFilenames)) {
                installed.add(version);
            }
        }

        return installed;
    }
}
