package org.liquibase.lpm.commands;

import org.liquibase.lpm.commands.mixins.GlobalMixin;
import org.liquibase.lpm.exception.LpmException;
import org.liquibase.lpm.model.Dependency;
import org.liquibase.lpm.model.DependencyManifest;
import org.liquibase.lpm.model.Package;
import org.liquibase.lpm.model.PackageRegistry;
import org.liquibase.lpm.model.PackageVersion;
import org.liquibase.lpm.service.PackageService;
import org.liquibase.lpm.util.ConsoleOutput;
import picocli.CommandLine.Command;
import picocli.CommandLine.Mixin;
import picocli.CommandLine.Parameters;
import picocli.CommandLine.ParentCommand;

import java.util.List;
import java.util.Optional;
import java.util.concurrent.Callable;

/**
 * Command to upgrade installed packages to their latest versions.
 */
@Command(
        name = "upgrade",
        aliases = {"up"},
        description = "Upgrades Installed Packages to the Latest Versions"
)
public class UpgradeCommand implements Callable<Integer> {

    @ParentCommand
    private RootCommand parent;

    @Mixin
    private GlobalMixin globalMixin;

    @Parameters(arity = "0..*",
            description = "Specific packages to upgrade (upgrades all if not specified)")
    private List<String> packages;

    private final PackageService packageService;

    /**
     * Creates a new upgrade command.
     */
    public UpgradeCommand() {
        this.packageService = new PackageService();
    }

    /**
     * Creates an upgrade command with a custom package service (for testing).
     *
     * @param packageService the package service to use
     */
    UpgradeCommand(PackageService packageService) {
        this.packageService = packageService;
    }

    @Override
    public Integer call() {
        try {
            boolean global = globalMixin != null && globalMixin.isGlobal();
            boolean dryRun = globalMixin != null && globalMixin.isDryRun();

            // Initialize service
            packageService.initialize(global);

            // Print classpath location
            ConsoleOutput.printClasspath(packageService.getClasspathService().getClasspathDir().toString());

            // Get outdated packages
            PackageRegistry outdated = packageService.getOutdatedPackages();

            if (outdated.isEmpty()) {
                ConsoleOutput.printOutdatedCount(0);
                return 0;
            }

            // Display outdated packages
            List<String> output = ConsoleOutput.formatOutdatedPackages(
                    outdated.getPackages(),
                    packageService.getClasspathFilenames(),
                    pkg -> pkg.getLatestVersion(packageService.getLiquibaseInfo().version())
            );
            ConsoleOutput.printLines(output);
            System.out.println();
            ConsoleOutput.printOutdatedCount(outdated.size());

            // Stop here if dry-run
            if (dryRun) {
                return 0;
            }

            // Load dependency manifest for local mode
            DependencyManifest manifest = null;
            if (!global) {
                manifest = new DependencyManifest();
                if (manifest.fileExists()) {
                    manifest.read();
                }
            }

            // Upgrade packages
            System.out.println();
            for (Package pkg : outdated.getPackages()) {
                try {
                    Optional<PackageVersion> installedOpt = pkg.getInstalledVersion(packageService.getClasspathFilenames());
                    Optional<PackageVersion> latestOpt = pkg.getLatestVersion(packageService.getLiquibaseInfo().version());

                    if (installedOpt.isEmpty() || latestOpt.isEmpty()) {
                        continue;
                    }

                    PackageVersion installed = installedOpt.get();
                    PackageVersion latest = latestOpt.get();

                    // Remove old version
                    packageService.getClasspathService().removeVersion(installed);
                    ConsoleOutput.printRemoveSuccess(installed.getFilename());

                    // Update manifest
                    if (!global && manifest != null) {
                        manifest.remove(pkg.name());
                    }

                    // Install new version
                    packageService.getClasspathService().installVersion(latest);
                    ConsoleOutput.printInstallSuccess(latest.getFilename());

                    // Update manifest
                    if (!global && manifest != null) {
                        manifest.add(new Dependency(pkg.name(), latest.tag()));
                    }
                } catch (LpmException e) {
                    System.err.println(e.getMessage());
                    // Continue with other packages
                }
            }

            // Write updated manifest for local mode
            if (!global && manifest != null && manifest.fileExists()) {
                manifest.write();
            }

            // Refresh classpath cache
            packageService.getClasspathService().refreshClasspathFilenames();

            return 0;
        } catch (Exception e) {
            System.err.println(e.getMessage());
            return 1;
        }
    }
}
