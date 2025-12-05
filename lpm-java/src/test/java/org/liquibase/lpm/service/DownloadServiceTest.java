package org.liquibase.lpm.service;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.liquibase.lpm.exception.ChecksumValidationException;
import org.liquibase.lpm.exception.DownloadException;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.io.IOException;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.nio.charset.StandardCharsets;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class DownloadServiceTest {

    @Mock
    private HttpClient httpClient;

    @Mock
    private HttpResponse<byte[]> httpResponse;

    @Mock
    private ChecksumService checksumService;

    private DownloadService downloadService;

    @BeforeEach
    void setUp() {
        downloadService = new DownloadService(httpClient, checksumService);
    }

    @Test
    void download_successfulRequest_returnsContent() throws Exception {
        byte[] expectedContent = "test content".getBytes(StandardCharsets.UTF_8);
        when(httpResponse.statusCode()).thenReturn(200);
        when(httpResponse.body()).thenReturn(expectedContent);
        when(httpClient.send(any(HttpRequest.class), any(HttpResponse.BodyHandler.class)))
                .thenReturn(httpResponse);

        byte[] result = downloadService.download("https://example.com/file.jar");

        assertThat(result).isEqualTo(expectedContent);
        verify(httpClient).send(any(HttpRequest.class), any(HttpResponse.BodyHandler.class));
    }

    @Test
    void download_notFoundResponse_throwsException() throws Exception {
        when(httpResponse.statusCode()).thenReturn(404);
        when(httpClient.send(any(HttpRequest.class), any(HttpResponse.BodyHandler.class)))
                .thenReturn(httpResponse);

        assertThatThrownBy(() -> downloadService.download("https://example.com/notfound.jar"))
                .isInstanceOf(DownloadException.class)
                .hasMessageContaining("404");
    }

    @Test
    void download_serverError_throwsException() throws Exception {
        when(httpResponse.statusCode()).thenReturn(500);
        when(httpClient.send(any(HttpRequest.class), any(HttpResponse.BodyHandler.class)))
                .thenReturn(httpResponse);

        assertThatThrownBy(() -> downloadService.download("https://example.com/error.jar"))
                .isInstanceOf(DownloadException.class)
                .hasMessageContaining("500");
    }

    @Test
    void download_ioException_throwsException() throws Exception {
        when(httpClient.send(any(HttpRequest.class), any(HttpResponse.BodyHandler.class)))
                .thenThrow(new IOException("Connection refused"));

        assertThatThrownBy(() -> downloadService.download("https://example.com/file.jar"))
                .isInstanceOf(DownloadException.class)
                .hasMessageContaining("example.com");
    }

    @Test
    void download_interruptedException_throwsExceptionAndSetsInterruptFlag() throws Exception {
        when(httpClient.send(any(HttpRequest.class), any(HttpResponse.BodyHandler.class)))
                .thenThrow(new InterruptedException("Thread interrupted"));

        Thread.currentThread().interrupted(); // Clear any existing interrupt flag

        assertThatThrownBy(() -> downloadService.download("https://example.com/file.jar"))
                .isInstanceOf(DownloadException.class);

        // Verify the interrupt flag was set
        assertThat(Thread.currentThread().isInterrupted()).isTrue();
        Thread.interrupted(); // Clear the flag for other tests
    }

    @Test
    void downloadString_successfulRequest_returnsTrimmedString() throws Exception {
        byte[] content = "  checksum value  \n".getBytes(StandardCharsets.UTF_8);
        when(httpResponse.statusCode()).thenReturn(200);
        when(httpResponse.body()).thenReturn(content);
        when(httpClient.send(any(HttpRequest.class), any(HttpResponse.BodyHandler.class)))
                .thenReturn(httpResponse);

        String result = downloadService.downloadString("https://example.com/checksum.txt");

        assertThat(result).isEqualTo("checksum value");
    }

    @Test
    void downloadAndVerify_validChecksum_returnsContent() throws Exception {
        byte[] content = "file content".getBytes(StandardCharsets.UTF_8);
        when(httpResponse.statusCode()).thenReturn(200);
        when(httpResponse.body()).thenReturn(content);
        when(httpClient.send(any(HttpRequest.class), any(HttpResponse.BodyHandler.class)))
                .thenReturn(httpResponse);
        // checksumService.verifyOrThrow doesn't throw when valid

        byte[] result = downloadService.downloadAndVerify(
                "https://example.com/file.jar",
                "abc123",
                "SHA1",
                "file.jar"
        );

        assertThat(result).isEqualTo(content);
        verify(checksumService).verifyOrThrow(content, "abc123", "SHA1", "file.jar");
    }

    @Test
    void downloadAndVerify_invalidChecksum_throwsException() throws Exception {
        byte[] content = "file content".getBytes(StandardCharsets.UTF_8);
        when(httpResponse.statusCode()).thenReturn(200);
        when(httpResponse.body()).thenReturn(content);
        when(httpClient.send(any(HttpRequest.class), any(HttpResponse.BodyHandler.class)))
                .thenReturn(httpResponse);
        doThrow(new ChecksumValidationException("file.jar", "expected", "actual"))
                .when(checksumService).verifyOrThrow(any(), any(), any(), any());

        assertThatThrownBy(() -> downloadService.downloadAndVerify(
                "https://example.com/file.jar",
                "wrongchecksum",
                "SHA1",
                "file.jar"
        )).isInstanceOf(ChecksumValidationException.class);
    }

    @Test
    void download_correctUrlIsParsed() throws Exception {
        byte[] content = "test".getBytes(StandardCharsets.UTF_8);
        when(httpResponse.statusCode()).thenReturn(200);
        when(httpResponse.body()).thenReturn(content);
        when(httpClient.send(any(HttpRequest.class), any(HttpResponse.BodyHandler.class)))
                .thenAnswer(invocation -> {
                    HttpRequest request = invocation.getArgument(0);
                    assertThat(request.uri().toString()).isEqualTo("https://repo.maven.org/file.jar");
                    return httpResponse;
                });

        downloadService.download("https://repo.maven.org/file.jar");
    }
}
