package org.liquibase.lpm.service;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.io.TempDir;
import org.liquibase.lpm.exception.ManifestException;
import org.liquibase.lpm.model.PackageVersion;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.List;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;
import static org.mockito.Mockito.*;

class ClasspathServiceTest {

    @TempDir
    Path tempDir;

    private DownloadService downloadService;
    private ChecksumService checksumService;
    private ClasspathService classpathService;

    @BeforeEach
    void setUp() {
        downloadService = mock(DownloadService.class);
        checksumService = mock(ChecksumService.class);
        classpathService = new ClasspathService(downloadService, checksumService);
    }

    @Test
    void setClasspathDir_setsDirectory() {
        classpathService.setClasspathDir(tempDir);

        assertThat(classpathService.getClasspathDir()).isEqualTo(tempDir);
    }

    @Test
    void useGlobalClasspath_setsLibSubdirectory() {
        classpathService.useGlobalClasspath(tempDir);

        assertThat(classpathService.getClasspathDir()).isEqualTo(tempDir.resolve("lib"));
    }

    @Test
    void classpathExists_existingDirectory_returnsTrue() {
        classpathService.setClasspathDir(tempDir);

        assertThat(classpathService.classpathExists()).isTrue();
    }

    @Test
    void classpathExists_nonExistingDirectory_returnsFalse() {
        classpathService.setClasspathDir(tempDir.resolve("nonexistent"));

        assertThat(classpathService.classpathExists()).isFalse();
    }

    @Test
    void classpathExists_nullDirectory_returnsFalse() {
        assertThat(classpathService.classpathExists()).isFalse();
    }

    @Test
    void ensureClasspathExists_createsDirectory() throws IOException {
        Path newDir = tempDir.resolve("new-classpath");
        classpathService.setClasspathDir(newDir);

        classpathService.ensureClasspathExists();

        assertThat(Files.isDirectory(newDir)).isTrue();
    }

    @Test
    void ensureClasspathExists_existingDirectory_doesNotFail() {
        classpathService.setClasspathDir(tempDir);

        // Should not throw
        classpathService.ensureClasspathExists();

        assertThat(Files.isDirectory(tempDir)).isTrue();
    }

    @Test
    void ensureClasspathExists_nullDirectory_throwsException() {
        assertThatThrownBy(() -> classpathService.ensureClasspathExists())
                .isInstanceOf(IllegalStateException.class)
                .hasMessageContaining("not set");
    }

    @Test
    void getClasspathFilenames_returnsListOfFiles() throws IOException {
        Files.createFile(tempDir.resolve("file1.jar"));
        Files.createFile(tempDir.resolve("file2.jar"));
        Files.createDirectory(tempDir.resolve("subdirectory"));

        classpathService.setClasspathDir(tempDir);

        List<String> filenames = classpathService.getClasspathFilenames();

        assertThat(filenames).containsExactlyInAnyOrder("file1.jar", "file2.jar");
    }

    @Test
    void getClasspathFilenames_emptyDirectory_returnsEmptyList() {
        classpathService.setClasspathDir(tempDir);

        List<String> filenames = classpathService.getClasspathFilenames();

        assertThat(filenames).isEmpty();
    }

    @Test
    void getClasspathFilenames_nonExistentDirectory_returnsEmptyList() {
        classpathService.setClasspathDir(tempDir.resolve("nonexistent"));

        List<String> filenames = classpathService.getClasspathFilenames();

        assertThat(filenames).isEmpty();
    }

    @Test
    void getClasspathFilenames_isCached() throws IOException {
        Files.createFile(tempDir.resolve("file1.jar"));
        classpathService.setClasspathDir(tempDir);

        List<String> first = classpathService.getClasspathFilenames();

        // Add another file
        Files.createFile(tempDir.resolve("file2.jar"));

        // Should return cached result
        List<String> second = classpathService.getClasspathFilenames();

        assertThat(first).isSameAs(second);
        assertThat(first).containsExactly("file1.jar");
    }

    @Test
    void refreshClasspathFilenames_clearsCache() throws IOException {
        Files.createFile(tempDir.resolve("file1.jar"));
        classpathService.setClasspathDir(tempDir);

        List<String> first = classpathService.getClasspathFilenames();
        assertThat(first).containsExactly("file1.jar");

        Files.createFile(tempDir.resolve("file2.jar"));

        List<String> refreshed = classpathService.refreshClasspathFilenames();

        assertThat(refreshed).containsExactlyInAnyOrder("file1.jar", "file2.jar");
    }

    @Test
    void installVersion_httpPath_downloadsAndWritesFile() throws IOException {
        PackageVersion version = new PackageVersion(
                "1.0.0",
                "https://example.com/test-1.0.0.jar",
                "SHA1",
                "checksum123",
                null
        );
        byte[] content = "jar content".getBytes();
        when(downloadService.downloadAndVerify(
                "https://example.com/test-1.0.0.jar",
                "checksum123",
                "SHA1",
                "test-1.0.0.jar"
        )).thenReturn(content);

        classpathService.setClasspathDir(tempDir);
        classpathService.installVersion(version);

        Path installedFile = tempDir.resolve("test-1.0.0.jar");
        assertThat(Files.exists(installedFile)).isTrue();
        assertThat(Files.readAllBytes(installedFile)).isEqualTo(content);
    }

    @Test
    void installVersion_localPath_copiesFile() throws IOException {
        // Create a local file to copy from
        Path localFile = tempDir.resolve("source.jar");
        byte[] content = "local jar content".getBytes();
        Files.write(localFile, content);

        PackageVersion version = new PackageVersion(
                "1.0.0",
                localFile.toString(),
                "SHA1",
                "checksum123",
                null
        );

        // Set a different directory as classpath
        Path classpathDir = tempDir.resolve("lib");
        Files.createDirectory(classpathDir);
        classpathService.setClasspathDir(classpathDir);

        classpathService.installVersion(version);

        Path installedFile = classpathDir.resolve("source.jar");
        assertThat(Files.exists(installedFile)).isTrue();
        assertThat(Files.readAllBytes(installedFile)).isEqualTo(content);
        verify(checksumService).verifyOrThrow(content, "checksum123", "SHA1", "source.jar");
    }

    @Test
    void installVersion_localPathWithoutChecksum_copiesWithoutVerification() throws IOException {
        Path localFile = tempDir.resolve("source.jar");
        byte[] content = "local jar content".getBytes();
        Files.write(localFile, content);

        PackageVersion version = new PackageVersion(
                "1.0.0",
                localFile.toString(),
                null,
                "",
                null
        );

        Path classpathDir = tempDir.resolve("lib");
        Files.createDirectory(classpathDir);
        classpathService.setClasspathDir(classpathDir);

        classpathService.installVersion(version);

        verifyNoInteractions(checksumService);
    }

    @Test
    void removeFile_existingFile_deletesFile() throws IOException {
        Path fileToRemove = tempDir.resolve("file.jar");
        Files.createFile(fileToRemove);

        classpathService.setClasspathDir(tempDir);
        classpathService.removeFile("file.jar");

        assertThat(Files.exists(fileToRemove)).isFalse();
    }

    @Test
    void removeFile_nonExistingFile_doesNotThrow() {
        classpathService.setClasspathDir(tempDir);

        // Should not throw
        classpathService.removeFile("nonexistent.jar");
    }

    @Test
    void removeFile_nullClasspath_throwsException() {
        assertThatThrownBy(() -> classpathService.removeFile("file.jar"))
                .isInstanceOf(IllegalStateException.class);
    }

    @Test
    void removeVersion_removesFileByFilename() throws IOException {
        Path fileToRemove = tempDir.resolve("test-1.0.0.jar");
        Files.createFile(fileToRemove);

        PackageVersion version = new PackageVersion(
                "1.0.0",
                "https://example.com/test-1.0.0.jar",
                "SHA1",
                "checksum",
                null
        );

        classpathService.setClasspathDir(tempDir);
        classpathService.removeVersion(version);

        assertThat(Files.exists(fileToRemove)).isFalse();
    }

    @Test
    void writePackagesJson_writesContent() throws IOException {
        byte[] content = "{\"packages\":[]}".getBytes();
        classpathService.setClasspathDir(tempDir);

        classpathService.writePackagesJson(content);

        Path packagesJson = tempDir.resolve("packages.json");
        assertThat(Files.exists(packagesJson)).isTrue();
        assertThat(Files.readAllBytes(packagesJson)).isEqualTo(content);
    }

    @Test
    void packagesJsonExists_exists_returnsTrue() throws IOException {
        Files.createFile(tempDir.resolve("packages.json"));
        classpathService.setClasspathDir(tempDir);

        assertThat(classpathService.packagesJsonExists()).isTrue();
    }

    @Test
    void packagesJsonExists_notExists_returnsFalse() {
        classpathService.setClasspathDir(tempDir);

        assertThat(classpathService.packagesJsonExists()).isFalse();
    }

    @Test
    void packagesJsonExists_nullClasspath_returnsFalse() {
        assertThat(classpathService.packagesJsonExists()).isFalse();
    }

    @Test
    void readPackagesJson_readsContent() throws IOException {
        byte[] content = "{\"packages\":[]}".getBytes();
        Files.write(tempDir.resolve("packages.json"), content);
        classpathService.setClasspathDir(tempDir);

        byte[] result = classpathService.readPackagesJson();

        assertThat(result).isEqualTo(content);
    }

    @Test
    void readPackagesJson_notExists_throwsException() {
        classpathService.setClasspathDir(tempDir);

        assertThatThrownBy(() -> classpathService.readPackagesJson())
                .isInstanceOf(ManifestException.class);
    }

    @Test
    void readPackagesJson_nullClasspath_throwsException() {
        assertThatThrownBy(() -> classpathService.readPackagesJson())
                .isInstanceOf(IllegalStateException.class);
    }

    @Test
    void installVersion_invalidatesFilenameCache() throws IOException {
        Files.createFile(tempDir.resolve("existing.jar"));
        classpathService.setClasspathDir(tempDir);

        // Cache the filenames
        List<String> before = classpathService.getClasspathFilenames();
        assertThat(before).containsExactly("existing.jar");

        // Install a new version
        PackageVersion version = new PackageVersion(
                "1.0.0",
                "https://example.com/new-1.0.0.jar",
                "SHA1",
                "checksum",
                null
        );
        when(downloadService.downloadAndVerify(any(), any(), any(), any()))
                .thenReturn("content".getBytes());

        classpathService.installVersion(version);

        // Get filenames again - should reflect the new file
        List<String> after = classpathService.getClasspathFilenames();
        assertThat(after).containsExactlyInAnyOrder("existing.jar", "new-1.0.0.jar");
    }

    @Test
    void removeFile_invalidatesFilenameCache() throws IOException {
        Files.createFile(tempDir.resolve("file1.jar"));
        Files.createFile(tempDir.resolve("file2.jar"));
        classpathService.setClasspathDir(tempDir);

        // Cache the filenames
        List<String> before = classpathService.getClasspathFilenames();
        assertThat(before).containsExactlyInAnyOrder("file1.jar", "file2.jar");

        // Remove a file
        classpathService.removeFile("file1.jar");

        // Get filenames again - should reflect removal
        List<String> after = classpathService.getClasspathFilenames();
        assertThat(after).containsExactly("file2.jar");
    }
}
