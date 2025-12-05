package org.liquibase.lpm.exception;

/**
 * Thrown when there are errors reading or writing manifest files (packages.json, liquibase.json).
 */
public final class ManifestException extends LpmException {

    /**
     * Creates a new exception for manifest read/write errors.
     *
     * @param message the error message
     * @param cause   the underlying cause
     */
    public ManifestException(String message, Throwable cause) {
        super(message, 1, cause);
    }

    /**
     * Creates a new exception for manifest validation errors.
     *
     * @param message the error message
     */
    public ManifestException(String message) {
        super(message, 1);
    }
}
