package org.liquibase.lpm.exception;

/**
 * Thrown when a requested version does not exist for a package.
 */
public final class VersionNotFoundException extends LpmException {

    /**
     * Creates a new exception for a missing version.
     *
     * @param packageName the name of the package
     * @param version     the version that was not found
     */
    public VersionNotFoundException(String packageName, String version) {
        super("Version '" + version + "' not found for package '" + packageName + "'.", 1);
    }
}
