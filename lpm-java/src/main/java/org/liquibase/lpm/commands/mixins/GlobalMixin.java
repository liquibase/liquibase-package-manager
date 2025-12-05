package org.liquibase.lpm.commands.mixins;

import picocli.CommandLine.Option;

/**
 * Mixin providing common flags shared across multiple commands.
 */
public class GlobalMixin {

    @Option(names = {"-g", "--global"},
            description = "Use global classpath (LIQUIBASE_HOME/lib) instead of local (./liquibase_libs)")
    private boolean global;

    @Option(names = "--dry-run",
            description = "Show what would be done without making changes")
    private boolean dryRun;

    /**
     * Returns whether to use the global classpath.
     *
     * @return true if global mode is enabled
     */
    public boolean isGlobal() {
        return global;
    }

    /**
     * Returns whether dry-run mode is enabled.
     *
     * @return true if dry-run mode is enabled
     */
    public boolean isDryRun() {
        return dryRun;
    }
}
