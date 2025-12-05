package org.liquibase.lpm.service;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.io.TempDir;
import org.liquibase.lpm.exception.LiquibaseDetectionException;
import org.semver4j.Semver;

import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.jar.JarEntry;
import java.util.jar.JarOutputStream;
import java.util.zip.ZipEntry;
import java.util.zip.ZipOutputStream;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;

class LiquibaseDetectorTest {

    @TempDir
    Path tempDir;

    private LiquibaseDetector detector;

    @BeforeEach
    void setUp() {
        detector = new LiquibaseDetector();
    }

    @Test
    void detectFromPath_validDirectoryWithJar_returnsInfo() throws IOException {
        // Create a fake liquibase.jar with build properties
        Path jarPath = tempDir.resolve("liquibase.jar");
        createLiquibaseJar(jarPath, "4.25.1");

        LiquibaseDetector.LiquibaseInfo info = detector.detectFromPath(tempDir);

        assertThat(info.homePath()).isEqualTo(tempDir.normalize());
        assertThat(info.version()).isEqualTo(new Semver("4.25.1"));
        assertThat(info.buildProperties()).containsEntry("build.version", "4.25.1");
    }

    @Test
    void detectFromPath_directoryWithoutJar_returnsDefaultVersion() throws IOException {
        // Create empty directory
        Path lbDir = tempDir.resolve("liquibase");
        Files.createDirectories(lbDir);

        LiquibaseDetector.LiquibaseInfo info = detector.detectFromPath(lbDir);

        assertThat(info.homePath()).isEqualTo(lbDir.normalize());
        assertThat(info.version()).isEqualTo(new Semver("0.0.0"));
        assertThat(info.buildProperties()).isEmpty();
    }

    @Test
    void detectFromPath_nonExistentDirectory_throwsException() {
        Path nonExistent = tempDir.resolve("nonexistent");

        assertThatThrownBy(() -> detector.detectFromPath(nonExistent))
                .isInstanceOf(LiquibaseDetectionException.class)
                .hasMessageContaining("does not exist");
    }

    @Test
    void detectFromPath_internalLibStructure_findsCoreJar() throws IOException {
        // Create internal/lib directory structure
        Path internalLib = tempDir.resolve("internal").resolve("lib");
        Files.createDirectories(internalLib);

        Path coreJar = internalLib.resolve("liquibase-core.jar");
        createLiquibaseJar(coreJar, "4.26.0");

        LiquibaseDetector.LiquibaseInfo info = detector.detectFromPath(tempDir);

        assertThat(info.version()).isEqualTo(new Semver("4.26.0"));
    }

    @Test
    void detectFromPath_internalLibWithCommercialJar_findsCommercialJar() throws IOException {
        // Create internal/lib directory structure
        Path internalLib = tempDir.resolve("internal").resolve("lib");
        Files.createDirectories(internalLib);

        Path commercialJar = internalLib.resolve("liquibase-commercial.jar");
        createLiquibaseJar(commercialJar, "4.27.0-commercial");

        LiquibaseDetector.LiquibaseInfo info = detector.detectFromPath(tempDir);

        assertThat(info.version()).isEqualTo(new Semver("4.27.0-commercial"));
    }

    @Test
    void detectFromPath_prefersDirect_liquibaseJar() throws IOException {
        // Create both liquibase.jar and internal/lib structure
        Path directJar = tempDir.resolve("liquibase.jar");
        createLiquibaseJar(directJar, "4.25.0");

        Path internalLib = tempDir.resolve("internal").resolve("lib");
        Files.createDirectories(internalLib);
        Path coreJar = internalLib.resolve("liquibase-core.jar");
        createLiquibaseJar(coreJar, "4.26.0");

        LiquibaseDetector.LiquibaseInfo info = detector.detectFromPath(tempDir);

        // Should prefer direct liquibase.jar
        assertThat(info.version()).isEqualTo(new Semver("4.25.0"));
    }

    @Test
    void detectFromPath_invalidVersionInJar_returnsFallbackVersion() throws IOException {
        // Create JAR with invalid version
        Path jarPath = tempDir.resolve("liquibase.jar");
        createLiquibaseJarWithProperties(jarPath, "build.version=not-a-valid-version\n");

        LiquibaseDetector.LiquibaseInfo info = detector.detectFromPath(tempDir);

        assertThat(info.version()).isEqualTo(new Semver("0.0.0"));
    }

    @Test
    void detectFromPath_jarWithoutBuildProperties_returnsDefaultVersion() throws IOException {
        // Create JAR without build properties
        Path jarPath = tempDir.resolve("liquibase.jar");
        createEmptyJar(jarPath);

        LiquibaseDetector.LiquibaseInfo info = detector.detectFromPath(tempDir);

        assertThat(info.version()).isEqualTo(new Semver("0.0.0"));
        assertThat(info.buildProperties()).isEmpty();
    }

    @Test
    void detectFromPath_jarWithMultipleProperties_readsAll() throws IOException {
        Path jarPath = tempDir.resolve("liquibase.jar");
        String properties = """
            build.version=4.28.0
            build.timestamp=2024-01-15T10:30:00Z
            build.branch=main
            build.number=12345
            """;
        createLiquibaseJarWithProperties(jarPath, properties);

        LiquibaseDetector.LiquibaseInfo info = detector.detectFromPath(tempDir);

        assertThat(info.version()).isEqualTo(new Semver("4.28.0"));
        assertThat(info.buildProperties())
                .containsEntry("build.version", "4.28.0")
                .containsEntry("build.timestamp", "2024-01-15T10:30:00Z")
                .containsEntry("build.branch", "main")
                .containsEntry("build.number", "12345");
    }

    @Test
    void getLibPath_returnsCorrectPath() throws IOException {
        Path jarPath = tempDir.resolve("liquibase.jar");
        createLiquibaseJar(jarPath, "4.25.0");

        LiquibaseDetector.LiquibaseInfo info = detector.detectFromPath(tempDir);

        assertThat(info.getLibPath()).isEqualTo(tempDir.normalize().resolve("lib"));
    }

    @Test
    void detectFromPath_handlesEmptyProperties() throws IOException {
        Path jarPath = tempDir.resolve("liquibase.jar");
        createLiquibaseJarWithProperties(jarPath, "");

        LiquibaseDetector.LiquibaseInfo info = detector.detectFromPath(tempDir);

        assertThat(info.version()).isEqualTo(new Semver("0.0.0"));
    }

    @Test
    void detectFromPath_handlesPropertiesWithEmptyValues() throws IOException {
        Path jarPath = tempDir.resolve("liquibase.jar");
        String properties = """
            build.version=4.29.0
            empty.value=
            """;
        createLiquibaseJarWithProperties(jarPath, properties);

        LiquibaseDetector.LiquibaseInfo info = detector.detectFromPath(tempDir);

        assertThat(info.buildProperties())
                .containsEntry("build.version", "4.29.0")
                .containsEntry("empty.value", "");
    }

    /**
     * Helper method to create a fake liquibase JAR with build properties.
     */
    private void createLiquibaseJar(Path jarPath, String version) throws IOException {
        String properties = "build.version=" + version + "\n";
        createLiquibaseJarWithProperties(jarPath, properties);
    }

    /**
     * Helper method to create a fake liquibase JAR with custom properties content.
     */
    private void createLiquibaseJarWithProperties(Path jarPath, String propertiesContent) throws IOException {
        try (ZipOutputStream zos = new ZipOutputStream(Files.newOutputStream(jarPath))) {
            // Add build properties entry
            ZipEntry entry = new ZipEntry("liquibase.build.properties");
            zos.putNextEntry(entry);
            zos.write(propertiesContent.getBytes());
            zos.closeEntry();

            // Add a dummy class file to make it look like a real JAR
            ZipEntry classEntry = new ZipEntry("liquibase/Liquibase.class");
            zos.putNextEntry(classEntry);
            zos.write(new byte[0]);
            zos.closeEntry();
        }
    }

    /**
     * Helper method to create an empty JAR file.
     */
    private void createEmptyJar(Path jarPath) throws IOException {
        try (ZipOutputStream zos = new ZipOutputStream(Files.newOutputStream(jarPath))) {
            // Add just a dummy entry
            ZipEntry entry = new ZipEntry("META-INF/MANIFEST.MF");
            zos.putNextEntry(entry);
            zos.write("Manifest-Version: 1.0\n".getBytes());
            zos.closeEntry();
        }
    }
}
