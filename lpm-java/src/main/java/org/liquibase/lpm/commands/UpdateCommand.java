package org.liquibase.lpm.commands;

import org.liquibase.lpm.service.PackageService;
import picocli.CommandLine.Command;
import picocli.CommandLine.Option;
import picocli.CommandLine.ParentCommand;

import java.util.concurrent.Callable;

/**
 * Command to update the package manifest (packages.json).
 */
@Command(
        name = "update",
        description = "Updates the Package Manifest"
)
public class UpdateCommand implements Callable<Integer> {

    private static final String DEFAULT_MANIFEST_URL =
            "https://raw.githubusercontent.com/liquibase/liquibase-package-manager/master/internal/app/packages.json";

    @ParentCommand
    private RootCommand parent;

    @Option(names = {"-p", "--path"},
            description = "Path or URL to the new packages.json",
            defaultValue = DEFAULT_MANIFEST_URL)
    private String path;

    private final PackageService packageService;

    /**
     * Creates a new update command.
     */
    public UpdateCommand() {
        this.packageService = new PackageService();
    }

    /**
     * Creates an update command with a custom package service (for testing).
     *
     * @param packageService the package service to use
     */
    UpdateCommand(PackageService packageService) {
        this.packageService = packageService;
    }

    @Override
    public Integer call() {
        try {
            // Initialize with global classpath (manifest lives in global classpath)
            packageService.initialize(true);

            // Update the manifest
            packageService.updateManifest(path);

            System.out.println("Package manifest updated from: " + path);

            return 0;
        } catch (Exception e) {
            System.err.println(e.getMessage());
            return 1;
        }
    }
}
