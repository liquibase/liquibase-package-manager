package org.liquibase.lpm.commands;

import org.liquibase.lpm.commands.mixins.GlobalMixin;
import org.liquibase.lpm.exception.LpmException;
import org.liquibase.lpm.model.Dependency;
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
 * Command to add packages to the liquibase.json file and install them.
 */
@Command(
        name = "add",
        description = "Add packages to the liquibase.json file and to this Liquibase installation"
)
public class AddCommand implements Callable<Integer> {

    @ParentCommand
    private RootCommand parent;

    @Mixin
    private GlobalMixin globalMixin;

    @Parameters(arity = "1..*",
            description = "Packages to add (format: name or name@version)")
    private List<String> packages;

    private final PackageService packageService;

    /**
     * Creates a new add command.
     */
    public AddCommand() {
        this.packageService = new PackageService();
    }

    /**
     * Creates an add command with a custom package service (for testing).
     *
     * @param packageService the package service to use
     */
    AddCommand(PackageService packageService) {
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
            for (String packageSpec : packages) {
                try {
                    Dependency dep = Dependency.fromSpec(packageSpec);

                    // Install the package
                    PackageVersion installed = packageService.resolveAndInstall(dep);

                    // Print success
                    ConsoleOutput.printInstallSuccess(installed.getFilename());

                    // Track dependency in manifest for local mode
                    if (!global && manifest != null) {
                        manifest.add(new Dependency(dep.name(), installed.tag()));
                    }
                } catch (LpmException e) {
                    System.err.println(e.getMessage());
                    // Continue with other packages
                }
            }

            // Write manifest for local mode
            if (!global && manifest != null) {
                manifest.createIfNotExists();
                manifest.write();

                // Show JAVA_OPTS instructions if needed
                if (packageService.requiresJavaOptsInstructions()) {
                    ConsoleOutput.printJavaOptsInstructions(
                            packageService.getClasspathService().getClasspathDir().toString());
                }
            }

            return 0;
        } catch (Exception e) {
            System.err.println(e.getMessage());
            return 1;
        }
    }
}
