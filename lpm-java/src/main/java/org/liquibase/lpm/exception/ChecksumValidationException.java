package org.liquibase.lpm.exception;

/**
 * Thrown when checksum verification fails during package download.
 */
public final class ChecksumValidationException extends LpmException {

    /**
     * Creates a new exception for checksum validation failure.
     *
     * @param filename         the name of the file being verified
     * @param expectedChecksum the expected checksum value
     * @param actualChecksum   the actual computed checksum value
     */
    public ChecksumValidationException(String filename, String expectedChecksum, String actualChecksum) {
        super(String.format(
                "Checksum validation failed for '%s'. Expected: %s, Got: %s",
                filename, expectedChecksum, actualChecksum), 1);
    }

    /**
     * Creates a new exception for an unsupported checksum algorithm.
     *
     * @param algorithm the unsupported algorithm name
     */
    public ChecksumValidationException(String algorithm) {
        super("Unsupported checksum algorithm: " + algorithm, 1);
    }
}
