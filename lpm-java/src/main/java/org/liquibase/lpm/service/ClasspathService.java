package org.liquibase.lpm.service;

import org.liquibase.lpm.exception.LpmException;
import org.liquibase.lpm.exception.ManifestException;
import org.liquibase.lpm.model.PackageVersion;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.StandardOpenOption;
import java.util.Collections;
import java.util.List;
import java.util.stream.Collectors;
import java.util.stream.Stream;

/**
 * Service for managing the classpath where packages are installed.
 * <p>
 * Supports both global classpath (LIQUIBASE_HOME/lib/) and local classpath
 * (./liquibase_libs/).
 */
public class ClasspathService {

    /**
     * Default local classpath directory name.
     */
    public static final String LOCAL_CLASSPATH_DIR = "liquibase_libs";

    private final DownloadService downloadService;
    private final ChecksumService checksumService;

    private Path classpathDir;
    private List<String> classpathFilenames;

    /**
     * Creates a new classpath service.
     */
    public ClasspathService() {
        this(new DownloadService(), new ChecksumService());
    }

    /**
     * Creates a new classpath service with custom dependencies.
     *
     * @param downloadService the download service to use
     * @param checksumService the checksum service to use
     */
    public ClasspathService(DownloadService downloadService, ChecksumService checksumService) {
        this.downloadService = downloadService;
        this.checksumService = checksumService;
    }

    /**
     * Sets the classpath directory.
     *
     * @param classpathDir the classpath directory path
     */
    public void setClasspathDir(Path classpathDir) {
        this.classpathDir = classpathDir;
        this.classpathFilenames = null; // Reset cached filenames
    }

    /**
     * Gets the current classpath directory.
     *
     * @return the classpath directory path
     */
    public Path getClasspathDir() {
        return classpathDir;
    }

    /**
     * Sets up the classpath for global mode.
     *
     * @param liquibaseHome the Liquibase home directory
     */
    public void useGlobalClasspath(Path liquibaseHome) {
        setClasspathDir(liquibaseHome.resolve("lib"));
    }

    /**
     * Sets up the classpath for local mode.
     */
    public void useLocalClasspath() {
        Path pwd = Path.of(System.getProperty("user.dir"));
        setClasspathDir(pwd.resolve(LOCAL_CLASSPATH_DIR));
    }

    /**
     * Checks if the classpath directory exists.
     *
     * @return true if the directory exists
     */
    public boolean classpathExists() {
        return classpathDir != null && Files.isDirectory(classpathDir);
    }

    /**
     * Creates the classpath directory if it doesn't exist.
     */
    public void ensureClasspathExists() {
        if (classpathDir == null) {
            throw new IllegalStateException("Classpath directory not set");
        }

        if (!Files.exists(classpathDir)) {
            try {
                Files.createDirectories(classpathDir);
            } catch (IOException e) {
                throw new ManifestException("Failed to create classpath directory: " + classpathDir, e);
            }
        }
    }

    /**
     * Gets the list of filenames in the classpath directory.
     * <p>
     * Results are cached until the classpath directory is changed.
     *
     * @return list of filenames in the classpath
     */
    public List<String> getClasspathFilenames() {
        if (classpathFilenames != null) {
            return classpathFilenames;
        }

        if (!classpathExists()) {
            classpathFilenames = Collections.emptyList();
            return classpathFilenames;
        }

        try (Stream<Path> files = Files.list(classpathDir)) {
            classpathFilenames = files
                    .filter(Files::isRegularFile)
                    .map(p -> p.getFileName().toString())
                    .sorted()
                    .collect(Collectors.toList());
            return classpathFilenames;
        } catch (IOException e) {
            return Collections.emptyList();
        }
    }

    /**
     * Refreshes the cached list of classpath filenames.
     *
     * @return the refreshed list
     */
    public List<String> refreshClasspathFilenames() {
        classpathFilenames = null;
        return getClasspathFilenames();
    }

    /**
     * Installs a package version to the classpath.
     * <p>
     * If the version path is HTTP, downloads and verifies the file.
     * Otherwise, copies from the local path.
     *
     * @param version the package version to install
     */
    public void installVersion(PackageVersion version) {
        ensureClasspathExists();

        String filename = version.getFilename();
        Path targetPath = classpathDir.resolve(filename);

        byte[] content;
        if (version.isHttpPath()) {
            content = downloadService.downloadAndVerify(
                    version.path(),
                    version.checksum(),
                    version.algorithm(),
                    filename
            );
            System.out.println("Checksum verified. Installing " + filename + " to " + classpathDir);
        } else {
            // Local file - read and optionally verify
            try {
                content = Files.readAllBytes(Path.of(version.path()));
                if (version.checksum() != null && !version.checksum().isBlank()) {
                    checksumService.verifyOrThrow(content, version.checksum(), version.algorithm(), filename);
                }
            } catch (IOException e) {
                throw new ManifestException("Failed to read local file: " + version.path(), e);
            }
        }

        writeToClasspath(targetPath, content, filename);
        classpathFilenames = null; // Invalidate cache
    }

    /**
     * Removes a file from the classpath.
     *
     * @param filename the filename to remove
     */
    public void removeFile(String filename) {
        if (classpathDir == null) {
            throw new IllegalStateException("Classpath directory not set");
        }

        Path targetPath = classpathDir.resolve(filename);
        try {
            Files.deleteIfExists(targetPath);
            classpathFilenames = null; // Invalidate cache
        } catch (IOException e) {
            throw new ManifestException("Failed to remove file: " + filename, e);
        }
    }

    /**
     * Removes a package version from the classpath.
     *
     * @param version the version to remove
     */
    public void removeVersion(PackageVersion version) {
        removeFile(version.getFilename());
    }

    /**
     * Writes content to the classpath.
     *
     * @param targetPath the target file path
     * @param content    the content to write
     * @param filename   the filename (for error messages)
     */
    private void writeToClasspath(Path targetPath, byte[] content, String filename) {
        try {
            Files.write(targetPath, content,
                    StandardOpenOption.CREATE,
                    StandardOpenOption.TRUNCATE_EXISTING,
                    StandardOpenOption.WRITE);
        } catch (IOException e) {
            throw new ManifestException("Failed to write " + filename + " to classpath", e);
        }
    }

    /**
     * Writes the packages.json manifest to the classpath.
     *
     * @param content the manifest content
     */
    public void writePackagesJson(byte[] content) {
        ensureClasspathExists();
        writeToClasspath(classpathDir.resolve("packages.json"), content, "packages.json");
    }

    /**
     * Checks if packages.json exists in the classpath.
     *
     * @return true if packages.json exists
     */
    public boolean packagesJsonExists() {
        if (classpathDir == null) {
            return false;
        }
        return Files.exists(classpathDir.resolve("packages.json"));
    }

    /**
     * Reads packages.json from the classpath.
     *
     * @return the content of packages.json
     */
    public byte[] readPackagesJson() {
        if (classpathDir == null) {
            throw new IllegalStateException("Classpath directory not set");
        }

        Path packagesJson = classpathDir.resolve("packages.json");
        try {
            return Files.readAllBytes(packagesJson);
        } catch (IOException e) {
            throw new ManifestException("Failed to read packages.json from classpath", e);
        }
    }
}
