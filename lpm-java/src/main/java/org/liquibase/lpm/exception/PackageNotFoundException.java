package org.liquibase.lpm.exception;

/**
 * Thrown when a requested package does not exist in the package registry.
 */
public final class PackageNotFoundException extends LpmException {

    /**
     * Creates a new exception for a missing package.
     *
     * @param packageName the name of the package that was not found
     */
    public PackageNotFoundException(String packageName) {
        super("Package '" + packageName + "' not found.", 1);
    }
}
