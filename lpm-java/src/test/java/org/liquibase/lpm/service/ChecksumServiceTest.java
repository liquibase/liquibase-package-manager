package org.liquibase.lpm.service;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.liquibase.lpm.exception.ChecksumValidationException;

import java.nio.charset.StandardCharsets;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;

class ChecksumServiceTest {

    private ChecksumService checksumService;

    @BeforeEach
    void setUp() {
        checksumService = new ChecksumService();
    }

    @Test
    void calculate_sha1_returnsCorrectChecksum() {
        byte[] data = "Hello, World!".getBytes(StandardCharsets.UTF_8);

        String checksum = checksumService.calculate(data, "SHA1");

        // Expected SHA1 for "Hello, World!"
        assertThat(checksum).isEqualTo("0a0a9f2a6772942557ab5355d76af442f8f65e01");
    }

    @Test
    void calculate_sha256_returnsCorrectChecksum() {
        byte[] data = "Hello, World!".getBytes(StandardCharsets.UTF_8);

        String checksum = checksumService.calculate(data, "SHA256");

        // Expected SHA256 for "Hello, World!"
        assertThat(checksum).isEqualTo("dffd6021bb2bd5b0af676290809ec3a53191dd81c7f70a4b28688a362182986f");
    }

    @Test
    void calculate_sha1_caseInsensitive() {
        byte[] data = "test".getBytes(StandardCharsets.UTF_8);

        String checksum1 = checksumService.calculate(data, "SHA1");
        String checksum2 = checksumService.calculate(data, "sha1");

        assertThat(checksum1).isEqualTo(checksum2);
    }

    @Test
    void calculate_sha256_caseInsensitive() {
        byte[] data = "test".getBytes(StandardCharsets.UTF_8);

        String checksum1 = checksumService.calculate(data, "SHA256");
        String checksum2 = checksumService.calculate(data, "sha256");

        assertThat(checksum1).isEqualTo(checksum2);
    }

    @Test
    void calculate_unknownAlgorithm_throwsException() {
        byte[] data = "test".getBytes(StandardCharsets.UTF_8);

        assertThatThrownBy(() -> checksumService.calculate(data, "MD5"))
                .isInstanceOf(ChecksumValidationException.class)
                .hasMessageContaining("MD5");
    }

    @Test
    void calculate_nullAlgorithm_throwsException() {
        byte[] data = "test".getBytes(StandardCharsets.UTF_8);

        assertThatThrownBy(() -> checksumService.calculate(data, null))
                .isInstanceOf(ChecksumValidationException.class);
    }

    @Test
    void verify_matchingChecksum_returnsTrue() {
        byte[] data = "Hello, World!".getBytes(StandardCharsets.UTF_8);
        String expectedChecksum = "0a0a9f2a6772942557ab5355d76af442f8f65e01";

        boolean result = checksumService.verify(data, expectedChecksum, "SHA1");

        assertThat(result).isTrue();
    }

    @Test
    void verify_nonMatchingChecksum_returnsFalse() {
        byte[] data = "Hello, World!".getBytes(StandardCharsets.UTF_8);
        String wrongChecksum = "0000000000000000000000000000000000000000";

        boolean result = checksumService.verify(data, wrongChecksum, "SHA1");

        assertThat(result).isFalse();
    }

    @Test
    void verify_caseInsensitiveChecksum() {
        byte[] data = "Hello, World!".getBytes(StandardCharsets.UTF_8);
        String uppercaseChecksum = "0A0A9F2A6772942557AB5355D76AF442F8F65E01";

        boolean result = checksumService.verify(data, uppercaseChecksum, "SHA1");

        assertThat(result).isTrue();
    }

    @Test
    void verify_truncatesLongerExpectedChecksum() {
        byte[] data = "Hello, World!".getBytes(StandardCharsets.UTF_8);
        // SHA1 checksum with extra characters appended
        String extendedChecksum = "0a0a9f2a6772942557ab5355d76af442f8f65e01  filename.jar";

        boolean result = checksumService.verify(data, extendedChecksum, "SHA1");

        assertThat(result).isTrue();
    }

    @Test
    void verifyOrThrow_matchingChecksum_doesNotThrow() {
        byte[] data = "Hello, World!".getBytes(StandardCharsets.UTF_8);
        String expectedChecksum = "0a0a9f2a6772942557ab5355d76af442f8f65e01";

        // Should not throw
        checksumService.verifyOrThrow(data, expectedChecksum, "SHA1", "test.jar");
    }

    @Test
    void verifyOrThrow_nonMatchingChecksum_throwsException() {
        byte[] data = "Hello, World!".getBytes(StandardCharsets.UTF_8);
        String wrongChecksum = "0000000000000000000000000000000000000000";

        assertThatThrownBy(() ->
                checksumService.verifyOrThrow(data, wrongChecksum, "SHA1", "test.jar"))
                .isInstanceOf(ChecksumValidationException.class)
                .hasMessageContaining("test.jar")
                .hasMessageContaining("Expected")
                .hasMessageContaining("Got");
    }

    @Test
    void calculate_emptyData_returnsValidChecksum() {
        byte[] data = new byte[0];

        String checksum = checksumService.calculate(data, "SHA1");

        // SHA1 of empty data
        assertThat(checksum).isEqualTo("da39a3ee5e6b4b0d3255bfef95601890afd80709");
    }

    @Test
    void calculate_largeData_returnsValidChecksum() {
        // 1MB of data
        byte[] data = new byte[1024 * 1024];
        for (int i = 0; i < data.length; i++) {
            data[i] = (byte) (i % 256);
        }

        String checksum = checksumService.calculate(data, "SHA256");

        assertThat(checksum).hasSize(64); // SHA256 produces 64 hex characters
    }
}
