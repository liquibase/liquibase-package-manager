package org.liquibase.lpm.service;

import org.liquibase.lpm.exception.LiquibaseDetectionException;
import org.semver4j.Semver;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.HashMap;
import java.util.Map;
import java.util.Optional;
import java.util.Properties;
import java.util.concurrent.TimeUnit;
import java.util.zip.ZipEntry;
import java.util.zip.ZipFile;

/**
 * Service for detecting Liquibase installations and extracting version information.
 * <p>
 * Detection strategy:
 * 1. Check LIQUIBASE_HOME environment variable
 * 2. Use `which liquibase` (Unix) or `where liquibase` (Windows) command
 * 3. Extract version from JAR's liquibase.build.properties
 */
public class LiquibaseDetector {

    private static final boolean IS_WINDOWS = System.getProperty("os.name")
            .toLowerCase().contains("win");

    private static final String BUILD_PROPERTIES_ENTRY = "liquibase.build.properties";
    private static final String BUILD_VERSION_KEY = "build.version";
    private static final Semver DEFAULT_VERSION = new Semver("0.0.0");

    /**
     * Result of Liquibase detection containing home path and version.
     *
     * @param homePath        the path to the Liquibase installation
     * @param version         the detected Liquibase version
     * @param buildProperties all properties from liquibase.build.properties
     */
    public record LiquibaseInfo(
            Path homePath,
            Semver version,
            Map<String, String> buildProperties
    ) {
        /**
         * Gets the lib directory path.
         *
         * @return the path to the lib directory
         */
        public Path getLibPath() {
            return homePath.resolve("lib");
        }
    }

    /**
     * Detects the Liquibase installation.
     *
     * @return information about the detected installation
     * @throws LiquibaseDetectionException if Liquibase cannot be found
     */
    public LiquibaseInfo detect() {
        // 1. Check LIQUIBASE_HOME environment variable
        String envHome = System.getenv("LIQUIBASE_HOME");
        if (envHome != null && !envHome.isBlank()) {
            Path homePath = Path.of(envHome);
            if (Files.isDirectory(homePath)) {
                return loadFromPath(homePath);
            }
        }

        // 2. Try to find via which/where command
        Optional<Path> commandPath = findViaCommand();
        if (commandPath.isPresent()) {
            return loadFromPath(commandPath.get());
        }

        throw new LiquibaseDetectionException();
    }

    /**
     * Detects Liquibase from a specific path.
     *
     * @param homePath the path to check
     * @return information about the installation
     */
    public LiquibaseInfo detectFromPath(Path homePath) {
        if (!Files.isDirectory(homePath)) {
            throw new LiquibaseDetectionException("Path does not exist: " + homePath);
        }
        return loadFromPath(homePath);
    }

    /**
     * Loads Liquibase information from a home path.
     */
    private LiquibaseInfo loadFromPath(Path homePath) {
        // Normalize the path (resolve symlinks, etc.)
        Path normalizedPath = homePath.normalize();

        // Try to find and read the JAR
        Optional<Path> jarPath = findLiquibaseJar(normalizedPath);
        if (jarPath.isEmpty()) {
            // No JAR found, return with default version
            return new LiquibaseInfo(normalizedPath, DEFAULT_VERSION, new HashMap<>());
        }

        // Extract build properties from JAR
        Map<String, String> buildProperties = extractBuildProperties(jarPath.get());

        // Parse version
        String versionString = buildProperties.getOrDefault(BUILD_VERSION_KEY, "0.0.0");
        Semver version;
        try {
            version = new Semver(versionString);
        } catch (Exception e) {
            System.err.println("Warning: Error parsing version '" + versionString +
                             "'. Falling back to version '0.0.0'.");
            version = DEFAULT_VERSION;
        }

        return new LiquibaseInfo(normalizedPath, version, buildProperties);
    }

    /**
     * Finds the Liquibase JAR file in the home path.
     * <p>
     * Search order:
     * 1. liquibase.jar (direct)
     * 2. internal/lib/liquibase-core.jar
     * 3. internal/lib/liquibase-commercial.jar
     */
    private Optional<Path> findLiquibaseJar(Path homePath) {
        // Try liquibase.jar first
        Path directJar = homePath.resolve("liquibase.jar");
        if (Files.exists(directJar)) {
            return Optional.of(directJar);
        }

        // Try internal/lib/liquibase-core.jar
        Path coreJar = homePath.resolve("internal").resolve("lib").resolve("liquibase-core.jar");
        if (Files.exists(coreJar)) {
            return Optional.of(coreJar);
        }

        // Try internal/lib/liquibase-commercial.jar
        Path commercialJar = homePath.resolve("internal").resolve("lib").resolve("liquibase-commercial.jar");
        if (Files.exists(commercialJar)) {
            return Optional.of(commercialJar);
        }

        return Optional.empty();
    }

    /**
     * Extracts build properties from a Liquibase JAR.
     */
    private Map<String, String> extractBuildProperties(Path jarPath) {
        Map<String, String> properties = new HashMap<>();

        try (ZipFile zipFile = new ZipFile(jarPath.toFile())) {
            ZipEntry entry = zipFile.getEntry(BUILD_PROPERTIES_ENTRY);
            if (entry == null) {
                return properties;
            }

            try (InputStream is = zipFile.getInputStream(entry);
                 BufferedReader reader = new BufferedReader(new InputStreamReader(is))) {

                String line;
                while ((line = reader.readLine()) != null) {
                    int equalIndex = line.indexOf('=');
                    if (equalIndex >= 0) {
                        String key = line.substring(0, equalIndex).trim();
                        String value = equalIndex < line.length() - 1
                                ? line.substring(equalIndex + 1).trim()
                                : "";
                        if (!key.isEmpty()) {
                            properties.put(key, value);
                        }
                    }
                }
            }
        } catch (IOException e) {
            // Return empty properties on error
            System.err.println("Warning: Could not read JAR properties: " + e.getMessage());
        }

        return properties;
    }

    /**
     * Tries to find Liquibase installation via which/where command.
     */
    private Optional<Path> findViaCommand() {
        String command = IS_WINDOWS ? "where" : "which";

        try {
            ProcessBuilder pb = new ProcessBuilder(command, "liquibase");
            pb.redirectErrorStream(true);
            Process process = pb.start();

            String output;
            try (BufferedReader reader = new BufferedReader(
                    new InputStreamReader(process.getInputStream()))) {
                output = reader.readLine();
            }

            boolean completed = process.waitFor(5, TimeUnit.SECONDS);
            if (!completed) {
                process.destroyForcibly();
                return Optional.empty();
            }

            if (process.exitValue() == 0 && output != null && !output.isBlank()) {
                Path execPath = Path.of(output.trim());

                // Resolve symlinks
                Path realPath = resolveSymlinks(execPath);

                // Get the parent directory (liquibase home)
                Path parent = realPath.getParent();
                if (parent != null) {
                    // Check if we're in a bin directory
                    if (parent.getFileName().toString().equals("bin")) {
                        parent = parent.getParent();
                    }
                    return Optional.of(parent);
                }
            }
        } catch (IOException | InterruptedException e) {
            // Ignore and return empty
            if (e instanceof InterruptedException) {
                Thread.currentThread().interrupt();
            }
        }

        return Optional.empty();
    }

    /**
     * Resolves symbolic links to get the real path.
     */
    private Path resolveSymlinks(Path path) {
        try {
            if (Files.isSymbolicLink(path)) {
                return Files.readSymbolicLink(path).toAbsolutePath().normalize();
            }
            return path.toAbsolutePath().normalize();
        } catch (IOException e) {
            return path;
        }
    }
}
