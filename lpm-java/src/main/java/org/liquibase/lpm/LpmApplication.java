package org.liquibase.lpm;

import org.liquibase.lpm.commands.RootCommand;
import picocli.CommandLine;

/**
 * Main entry point for the Liquibase Package Manager CLI.
 * <p>
 * This application provides package management capabilities for Liquibase,
 * allowing users to search, install, update, and remove extensions, drivers,
 * and utilities.
 */
public class LpmApplication {

    /**
     * Application entry point.
     *
     * @param args command line arguments
     */
    public static void main(String[] args) {
        int exitCode = new CommandLine(new RootCommand())
                .setCaseInsensitiveEnumValuesAllowed(true)
                .execute(args);
        System.exit(exitCode);
    }
}
