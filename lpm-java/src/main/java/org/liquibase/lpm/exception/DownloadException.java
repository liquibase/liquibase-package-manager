package org.liquibase.lpm.exception;

/**
 * Thrown when a package download fails.
 */
public final class DownloadException extends LpmException {

    /**
     * Creates a new exception for a download failure.
     *
     * @param url   the URL that failed to download
     * @param cause the underlying cause
     */
    public DownloadException(String url, Throwable cause) {
        super("Failed to download from: " + url, 1, cause);
    }

    /**
     * Creates a new exception for a download failure with HTTP status.
     *
     * @param url        the URL that failed to download
     * @param statusCode the HTTP status code received
     */
    public DownloadException(String url, int statusCode) {
        super(String.format("Failed to download from '%s'. HTTP status: %d", url, statusCode), 1);
    }
}
