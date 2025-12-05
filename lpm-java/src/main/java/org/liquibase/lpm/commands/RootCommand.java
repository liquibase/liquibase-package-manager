package org.liquibase.lpm.commands;

import org.liquibase.lpm.LpmVersionProvider;
import picocli.CommandLine.Command;
import picocli.CommandLine.Option;

import java.util.concurrent.Callable;

/**
 * Root command for the Liquibase Package Manager CLI.
 * <p>
 * This is the main entry point for all CLI operations. Subcommands are
 * registered here and global options (like --category) are defined.
 */
@Command(
        name = "lpm",
        description = "Liquibase Package Manager - Easily manage external dependencies for Database Development.%n" +
                      "Search for, install, and uninstall liquibase drivers, extensions, and utilities.",
        mixinStandardHelpOptions = true,
        versionProvider = LpmVersionProvider.class,
        subcommands = {
                AddCommand.class,
                InstallCommand.class,
                SearchCommand.class,
                ListCommand.class,
                RemoveCommand.class,
                UpdateCommand.class,
                UpgradeCommand.class,
                DedupeCommand.class,
        }
)
public class RootCommand implements Callable<Integer> {

    @Option(names = "--category",
            description = "Filter packages by category: extension, driver, or pro")
    private String category;

    /**
     * Gets the category filter.
     *
     * @return the category filter, or null if not set
     */
    public String getCategory() {
        return category;
    }

    @Override
    public Integer call() {
        // When called without subcommand, show help
        System.out.println("Liquibase Package Manager");
        System.out.println("Use --help to see available commands.");
        return 0;
    }
}
