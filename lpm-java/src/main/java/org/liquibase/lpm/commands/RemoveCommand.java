package org.liquibase.lpm.commands;

import org.liquibase.lpm.commands.mixins.GlobalMixin;
import org.liquibase.lpm.exception.LpmException;
import org.liquibase.lpm.model.DependencyManifest;
import org.liquibase.lpm.model.PackageVersion;
import org.liquibase.lpm.service.PackageService;
import org.liquibase.lpm.util.ConsoleOutput;
import picocli.CommandLine.Command;
import picocli.CommandLine.Mixin;
import picocli.CommandLine.Parameters;
import picocli.CommandLine.ParentCommand;

import java.util.List;
import java.util.concurrent.Callable;

/**
 * Command to remove installed packages.
 */
@Command(
        name = "remove",
        aliases = {"rm"},
        description = "Removes Package"
)
public class RemoveCommand implements Callable<Integer> {

    @ParentCommand
    private RootCommand parent;

    @Mixin
    private GlobalMixin globalMixin;

    @Parameters(arity = "1..*",
            description = "Packages to remove")
    private List<String> packages;

    private final PackageService packageService;

    /**
     * Creates a new remove command.
     */
    public RemoveCommand() {
        this.packageService = new PackageService();
    }

    /**
     * Creates a remove command with a custom package service (for testing).
     *
     * @param packageService the package service to use
     */
    RemoveCommand(PackageService packageService) {
        this.packageService = packageService;
    }

    @Override
    public Integer call() {
        try {
            boolean global = globalMixin != null && globalMixin.isGlobal();

            // Initialize service
            packageService.initialize(global);

            // Load dependency manifest for local mode
            DependencyManifest manifest = null;
            if (!global) {
                manifest = new DependencyManifest();
                if (manifest.fileExists()) {
                    manifest.read();
                }
            }

            // Process each package
            for (String packageName : packages) {
                try {
                    // Remove the package
                    PackageVersion removed = packageService.removePackage(packageName);

                    // Print success
                    ConsoleOutput.printRemoveSuccess(removed.getFilename());

                    // Remove from manifest for local mode
                    if (!global && manifest != null) {
                        manifest.remove(packageName);
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

            return 0;
        } catch (Exception e) {
            System.err.println(e.getMessage());
            return 1;
        }
    }
}
