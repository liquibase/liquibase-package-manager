package org.liquibase.lpm.service;

import org.liquibase.lpm.exception.ChecksumValidationException;

import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.util.HexFormat;

/**
 * Service for calculating and verifying checksums.
 * <p>
 * Supports SHA1 and SHA256 algorithms as used by the package registry.
 */
public class ChecksumService {

    /**
     * Algorithm constant for SHA-1.
     */
    public static final String SHA1 = "SHA1";

    /**
     * Algorithm constant for SHA-256.
     */
    public static final String SHA256 = "SHA256";

    /**
     * Calculates a checksum for the given data.
     *
     * @param data      the data to checksum
     * @param algorithm the algorithm to use (SHA1 or SHA256)
     * @return the checksum as a lowercase hex string
     * @throws ChecksumValidationException if the algorithm is not supported
     */
    public String calculate(byte[] data, String algorithm) {
        String javaAlgorithm = mapAlgorithm(algorithm);

        try {
            MessageDigest digest = MessageDigest.getInstance(javaAlgorithm);
            byte[] hash = digest.digest(data);
            return HexFormat.of().formatHex(hash);
        } catch (NoSuchAlgorithmException e) {
            throw new ChecksumValidationException(algorithm);
        }
    }

    /**
     * Verifies that the data matches the expected checksum.
     *
     * @param data             the data to verify
     * @param expectedChecksum the expected checksum (hex string)
     * @param algorithm        the algorithm to use (SHA1 or SHA256)
     * @return true if the checksums match
     */
    public boolean verify(byte[] data, String expectedChecksum, String algorithm) {
        String actualChecksum = calculate(data, algorithm);
        // Compare only the relevant portion (some registries store extra chars)
        int expectedLength = actualChecksum.length();
        String normalizedExpected = normalizeChecksum(expectedChecksum, expectedLength);
        return actualChecksum.equalsIgnoreCase(normalizedExpected);
    }

    /**
     * Verifies the checksum and throws an exception if it doesn't match.
     *
     * @param data             the data to verify
     * @param expectedChecksum the expected checksum
     * @param algorithm        the algorithm to use
     * @param filename         the filename (for error messages)
     * @throws ChecksumValidationException if the checksums don't match
     */
    public void verifyOrThrow(byte[] data, String expectedChecksum, String algorithm, String filename) {
        String actualChecksum = calculate(data, algorithm);
        int expectedLength = actualChecksum.length();
        String normalizedExpected = normalizeChecksum(expectedChecksum, expectedLength);

        if (!actualChecksum.equalsIgnoreCase(normalizedExpected)) {
            throw new ChecksumValidationException(filename, normalizedExpected, actualChecksum);
        }
    }

    /**
     * Maps our algorithm names to Java's MessageDigest algorithm names.
     */
    private String mapAlgorithm(String algorithm) {
        if (algorithm == null) {
            throw new ChecksumValidationException("null");
        }

        return switch (algorithm.toUpperCase()) {
            case "SHA1" -> "SHA-1";
            case "SHA256" -> "SHA-256";
            default -> throw new ChecksumValidationException(algorithm);
        };
    }

    /**
     * Normalizes a checksum to the expected length.
     * <p>
     * Some sources provide checksums with extra characters or different lengths.
     * This method truncates to match the expected hash length.
     *
     * @param checksum       the checksum to normalize
     * @param expectedLength the expected length
     * @return the normalized checksum
     */
    private String normalizeChecksum(String checksum, int expectedLength) {
        if (checksum == null) {
            return "";
        }
        String trimmed = checksum.trim();
        if (trimmed.length() > expectedLength) {
            return trimmed.substring(0, expectedLength);
        }
        return trimmed;
    }
}
