package org.liquibase.lpm.exception;

/**
 * Thrown when attempting to install a package that is already installed.
 */
public final class PackageAlreadyInstalledException extends LpmException {

    /**
     * Creates a new exception for an already installed package.
     *
     * @param packageName      the name of the package
     * @param installedVersion the version already installed
     */
    public PackageAlreadyInstalledException(String packageName, String installedVersion) {
        super(String.format("Package '%s' is already installed (version %s).", packageName, installedVersion), 1);
    }
}
