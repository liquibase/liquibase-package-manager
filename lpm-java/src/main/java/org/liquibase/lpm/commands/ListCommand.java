package org.liquibase.lpm.commands;

import org.liquibase.lpm.commands.mixins.GlobalMixin;
import org.liquibase.lpm.model.PackageRegistry;
import org.liquibase.lpm.service.PackageService;
import org.liquibase.lpm.util.ConsoleOutput;
import picocli.CommandLine.Command;
import picocli.CommandLine.Mixin;
import picocli.CommandLine.ParentCommand;

import java.util.List;
import java.util.concurrent.Callable;

/**
 * Command to list installed packages.
 */
@Command(
        name = "list",
        aliases = {"ls"},
        description = "List Installed Packages"
)
public class ListCommand implements Callable<Integer> {

    @ParentCommand
    private RootCommand parent;

    @Mixin
    private GlobalMixin globalMixin;

    private final PackageService packageService;

    /**
     * Creates a new list command.
     */
    public ListCommand() {
        this.packageService = new PackageService();
    }

    /**
     * Creates a list command with a custom package service (for testing).
     *
     * @param packageService the package service to use
     */
    ListCommand(PackageService packageService) {
        this.packageService = packageService;
    }

    @Override
    public Integer call() {
        try {
            boolean global = globalMixin != null && globalMixin.isGlobal();

            // Initialize service
            packageService.initialize(global);

            // Print classpath location
            ConsoleOutput.printClasspath(packageService.getClasspathService().getClasspathDir().toString());

            // Get installed packages
            PackageRegistry installed = packageService.getInstalledPackages();

            // Apply category filter if specified
            String category = parent != null ? parent.getCategory() : null;
            if (category != null && !category.isBlank()) {
                installed = installed.filterByCategory(category);
            }

            // Check results
            if (installed.isEmpty()) {
                System.out.println("No packages installed.");
                return 0;
            }

            // Display results
            List<String> output = installed.display(packageService.getClasspathFilenames());
            ConsoleOutput.printLines(output);

            return 0;
        } catch (Exception e) {
            System.err.println(e.getMessage());
            return 1;
        }
    }
}
