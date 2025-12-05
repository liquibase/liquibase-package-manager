package org.liquibase.lpm.service;

import org.liquibase.lpm.exception.DownloadException;

import java.io.IOException;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.time.Duration;

/**
 * Service for downloading files from HTTP/HTTPS URLs.
 * <p>
 * Uses Java's built-in HttpClient for HTTP operations.
 */
public class DownloadService {

    private static final Duration DEFAULT_TIMEOUT = Duration.ofSeconds(30);
    private static final Duration DEFAULT_CONNECT_TIMEOUT = Duration.ofSeconds(10);

    private final HttpClient httpClient;
    private final ChecksumService checksumService;

    /**
     * Creates a new download service with default settings.
     */
    public DownloadService() {
        this(HttpClient.newBuilder()
                .followRedirects(HttpClient.Redirect.NORMAL)
                .connectTimeout(DEFAULT_CONNECT_TIMEOUT)
                .build(),
             new ChecksumService());
    }

    /**
     * Creates a new download service with a custom HTTP client.
     *
     * @param httpClient      the HTTP client to use
     * @param checksumService the checksum service to use
     */
    public DownloadService(HttpClient httpClient, ChecksumService checksumService) {
        this.httpClient = httpClient;
        this.checksumService = checksumService;
    }

    /**
     * Downloads content from a URL and returns it as bytes.
     *
     * @param url the URL to download from
     * @return the downloaded content
     * @throws DownloadException if the download fails
     */
    public byte[] download(String url) {
        try {
            HttpRequest request = HttpRequest.newBuilder()
                    .uri(URI.create(url))
                    .timeout(DEFAULT_TIMEOUT)
                    .GET()
                    .build();

            HttpResponse<byte[]> response = httpClient.send(request, HttpResponse.BodyHandlers.ofByteArray());

            if (response.statusCode() != 200) {
                throw new DownloadException(url, response.statusCode());
            }

            return response.body();
        } catch (IOException e) {
            throw new DownloadException(url, e);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new DownloadException(url, e);
        }
    }

    /**
     * Downloads content and verifies its checksum.
     *
     * @param url              the URL to download from
     * @param expectedChecksum the expected checksum
     * @param algorithm        the checksum algorithm (SHA1 or SHA256)
     * @param filename         the filename for error messages
     * @return the downloaded content (verified)
     * @throws DownloadException                               if the download fails
     * @throws org.liquibase.lpm.exception.ChecksumValidationException if the checksum doesn't match
     */
    public byte[] downloadAndVerify(String url, String expectedChecksum, String algorithm, String filename) {
        byte[] content = download(url);
        checksumService.verifyOrThrow(content, expectedChecksum, algorithm, filename);
        return content;
    }

    /**
     * Downloads a string from a URL (e.g., for checksum files).
     *
     * @param url the URL to download from
     * @return the content as a string
     * @throws DownloadException if the download fails
     */
    public String downloadString(String url) {
        byte[] bytes = download(url);
        return new String(bytes).trim();
    }
}
