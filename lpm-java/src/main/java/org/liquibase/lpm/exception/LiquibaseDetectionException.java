package org.liquibase.lpm.exception;

/**
 * Thrown when Liquibase installation cannot be detected.
 */
public final class LiquibaseDetectionException extends LpmException {

    /**
     * Creates a new exception when Liquibase cannot be found.
     */
    public LiquibaseDetectionException() {
        super("Unable to locate Liquibase installation. " +
              "Please set the LIQUIBASE_HOME environment variable.", 1);
    }

    /**
     * Creates a new exception with a specific message.
     *
     * @param message the error message
     */
    public LiquibaseDetectionException(String message) {
        super(message, 1);
    }

    /**
     * Creates a new exception with a cause.
     *
     * @param message the error message
     * @param cause   the underlying cause
     */
    public LiquibaseDetectionException(String message, Throwable cause) {
        super(message, 1, cause);
    }
}
