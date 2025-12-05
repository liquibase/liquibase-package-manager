package org.liquibase.lpm.commands;

import org.liquibase.lpm.exception.LpmException;
import org.liquibase.lpm.model.Dependency;
import org.liquibase.lpm.model.DependencyManifest;
import org.liquibase.lpm.model.PackageVersion;
import org.liquibase.lpm.service.PackageService;
import org.liquibase.lpm.util.ConsoleOutput;
import picocli.CommandLine.Command;
import picocli.CommandLine.ParentCommand;

import java.util.concurrent.Callable;

/**
 * Command to install packages listed in liquibase.json file.
 */
@Command(
        name = "install",
        description = "Install packages listed in liquibase.json file"
)
public class InstallCommand implements Callable<Integer> {

    @ParentCommand
    private RootCommand parent;

    private final PackageService packageService;

    /**
     * Creates a new install command.
     */
    public InstallCommand() {
        this.packageService = new PackageService();
    }

    /**
     * Creates an install command with a custom package service (for testing).
     *
     * @param packageService the package service to use
     */
    InstallCommand(PackageService packageService) {
        this.packageService = packageService;
    }

    @Override
    public Integer call() {
        try {
            // Install command always uses local classpath
            packageService.initialize(false);

            // Load dependency manifest
            DependencyManifest manifest = new DependencyManifest();
            if (!manifest.fileExists()) {
                System.out.println("No liquibase.json found. Nothing to install.");
                return 0;
            }

            manifest.read();

            if (manifest.isEmpty()) {
                System.out.println("No dependencies in liquibase.json. Nothing to install.");
                return 0;
            }

            // Install each dependency
            for (Dependency dep : manifest.getDependencies()) {
                try {
                    PackageVersion installed = packageService.resolveAndInstall(dep);
                    ConsoleOutput.printInstallSuccess(installed.getFilename());
                } catch (LpmException e) {
                    System.err.println(e.getMessage());
                    // Continue with other packages
                }
            }

            // Show JAVA_OPTS instructions if needed
            if (packageService.requiresJavaOptsInstructions()) {
                ConsoleOutput.printJavaOptsInstructions(
                        packageService.getClasspathService().getClasspathDir().toString());
            }

            return 0;
        } catch (Exception e) {
            System.err.println(e.getMessage());
            return 1;
        }
    }
}
