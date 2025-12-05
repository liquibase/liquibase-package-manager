package org.liquibase.lpm.exception;

/**
 * Thrown when attempting to remove a package that is not installed.
 */
public final class PackageNotInstalledException extends LpmException {

    /**
     * Creates a new exception for a package that is not installed.
     *
     * @param packageName the name of the package
     */
    public PackageNotInstalledException(String packageName) {
        super("Package '" + packageName + "' is not installed.", 1);
    }
}
