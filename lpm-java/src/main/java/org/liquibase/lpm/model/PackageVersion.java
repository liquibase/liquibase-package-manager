package org.liquibase.lpm.model;

import com.fasterxml.jackson.annotation.JsonProperty;

import java.nio.file.Path;
import java.util.List;

/**
 * Represents a specific version of a package.
 * <p>
 * This record contains all metadata needed to identify, download, and validate
 * a package version, including its download path, checksum, and Liquibase
 * compatibility information.
 *
 * @param tag           the version tag (e.g., "1.0.0")
 * @param path          the download URL or local file path
 * @param algorithm     the checksum algorithm ("SHA1" or "SHA256")
 * @param checksum      the expected checksum value
 * @param liquibaseCore the minimum Liquibase version required (may be null for drivers)
 */
public record PackageVersion(
        @JsonProperty("tag") String tag,
        @JsonProperty("path") String path,
        @JsonProperty("algorithm") String algorithm,
        @JsonProperty("checksum") String checksum,
        @JsonProperty("liquibaseCore") String liquibaseCore
) {

    /**
     * Extracts the filename from the path.
     *
     * @return the filename portion of the path
     */
    public String getFilename() {
        if (path == null || path.isBlank()) {
            return "";
        }
        // Handle both URL paths and file paths
        String normalizedPath = path.replace('\\', '/');
        int lastSlash = normalizedPath.lastIndexOf('/');
        return lastSlash >= 0 ? normalizedPath.substring(lastSlash + 1) : normalizedPath;
    }

    /**
     * Checks if the path is a remote HTTP/HTTPS URL.
     *
     * @return true if the path starts with "http", false otherwise
     */
    public boolean isHttpPath() {
        return path != null && path.toLowerCase().startsWith("http");
    }

    /**
     * Checks if this version is installed in the given classpath.
     *
     * @param classpathFiles list of file paths in the classpath
     * @return true if this version's file exists in the classpath
     */
    public boolean isInClasspath(List<Path> classpathFiles) {
        if (classpathFiles == null || classpathFiles.isEmpty()) {
            return false;
        }
        String filename = getFilename();
        return classpathFiles.stream()
                .anyMatch(p -> p.getFileName().toString().equals(filename));
    }

    /**
     * Checks if this version is installed in the given classpath (using string filenames).
     *
     * @param filenames list of filenames in the classpath
     * @return true if this version's file exists in the classpath
     */
    public boolean isInClasspathByFilename(List<String> filenames) {
        if (filenames == null || filenames.isEmpty()) {
            return false;
        }
        String filename = getFilename();
        return filenames.contains(filename);
    }

    /**
     * Checks if this is an empty/null version (used for "not found" scenarios).
     *
     * @return true if the tag is null or empty
     */
    public boolean isEmpty() {
        return tag == null || tag.isBlank();
    }

    /**
     * Creates an empty version (used as a sentinel value).
     *
     * @return an empty PackageVersion
     */
    public static PackageVersion empty() {
        return new PackageVersion(null, null, null, null, null);
    }
}
