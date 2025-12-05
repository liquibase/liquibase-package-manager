package org.liquibase.lpm.exception;

/**
 * Base exception for all LPM-related errors.
 * <p>
 * This is a sealed class hierarchy that provides type-safe error handling
 * with specific exit codes for each error type.
 */
public sealed class LpmException extends RuntimeException
        permits PackageNotFoundException, VersionNotFoundException,
                VersionIncompatibleException, ChecksumValidationException,
                DownloadException, PackageAlreadyInstalledException,
                PackageNotInstalledException, ManifestException,
                LiquibaseDetectionException {

    private final int exitCode;

    /**
     * Creates a new LPM exception.
     *
     * @param message  the error message
     * @param exitCode the exit code to use when this exception causes program termination
     */
    protected LpmException(String message, int exitCode) {
        super(message);
        this.exitCode = exitCode;
    }

    /**
     * Creates a new LPM exception with a cause.
     *
     * @param message  the error message
     * @param exitCode the exit code to use when this exception causes program termination
     * @param cause    the underlying cause
     */
    protected LpmException(String message, int exitCode, Throwable cause) {
        super(message, cause);
        this.exitCode = exitCode;
    }

    /**
     * Returns the exit code associated with this exception.
     *
     * @return the exit code
     */
    public int getExitCode() {
        return exitCode;
    }
}
