package org.liquibase.lpm.exception;

/**
 * Thrown when a package version is incompatible with the installed Liquibase version.
 */
public final class VersionIncompatibleException extends LpmException {

    /**
     * Creates a new exception for an incompatible version.
     *
     * @param packageName      the name of the package
     * @param packageVersion   the version of the package
     * @param requiredLbCore   the minimum Liquibase version required by the package
     * @param installedLbCore  the currently installed Liquibase version
     */
    public VersionIncompatibleException(String packageName, String packageVersion,
                                        String requiredLbCore, String installedLbCore) {
        super(String.format(
                "Package '%s@%s' requires Liquibase %s or higher, but installed version is %s.",
                packageName, packageVersion, requiredLbCore, installedLbCore), 1);
    }
}
