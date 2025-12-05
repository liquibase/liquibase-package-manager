package org.liquibase.lpm.commands;

import org.liquibase.lpm.model.PackageRegistry;
import org.liquibase.lpm.service.PackageService;
import org.liquibase.lpm.util.ConsoleOutput;
import picocli.CommandLine.Command;
import picocli.CommandLine.Parameters;
import picocli.CommandLine.ParentCommand;

import java.util.List;
import java.util.concurrent.Callable;

/**
 * Command to search for packages in the registry.
 */
@Command(
        name = "search",
        description = "Search for Packages"
)
public class SearchCommand implements Callable<Integer> {

    @ParentCommand
    private RootCommand parent;

    @Parameters(index = "0", arity = "0..1",
            description = "Package name to search for (minimum 3 characters)")
    private String searchTerm;

    private final PackageService packageService;

    /**
     * Creates a new search command.
     */
    public SearchCommand() {
        this.packageService = new PackageService();
    }

    /**
     * Creates a search command with a custom package service (for testing).
     *
     * @param packageService the package service to use
     */
    SearchCommand(PackageService packageService) {
        this.packageService = packageService;
    }

    @Override
    public Integer call() {
        try {
            // Initialize with global classpath (search doesn't modify anything)
            packageService.initialize(true);

            // Validate search term if provided
            if (searchTerm != null && !searchTerm.isBlank() && searchTerm.length() < 3) {
                System.err.println("Search term must be at least 3 characters.");
                return 1;
            }

            // Get packages, optionally filtered
            PackageRegistry registry = packageService.getPackageRegistry();

            // Apply category filter if specified
            String category = parent != null ? parent.getCategory() : null;
            if (category != null && !category.isBlank()) {
                registry = registry.filterByCategory(category);
            }

            // Apply search filter if specified
            if (searchTerm != null && !searchTerm.isBlank()) {
                registry = registry.filterByName(searchTerm);
            }

            // Check results
            if (registry.isEmpty()) {
                System.out.println("No results found.");
                return 0;
            }

            // Display results
            List<String> output = registry.display(packageService.getClasspathFilenames());
            ConsoleOutput.printLines(output);

            return 0;
        } catch (Exception e) {
            System.err.println(e.getMessage());
            return 1;
        }
    }
}
