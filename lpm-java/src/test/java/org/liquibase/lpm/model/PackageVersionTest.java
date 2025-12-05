package org.liquibase.lpm.model;

import org.junit.jupiter.api.Test;

import java.util.Arrays;
import java.util.Collections;
import java.util.List;

import static org.assertj.core.api.Assertions.assertThat;

class PackageVersionTest {

    @Test
    void getFilename_returnsFilenameFromHttpUrl() {
        PackageVersion version = new PackageVersion(
                "1.0.0",
                "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-postgresql-1.0.0.jar",
                "SHA1",
                "abc123",
                "4.0.0"
        );

        assertThat(version.getFilename()).isEqualTo("liquibase-postgresql-1.0.0.jar");
    }

    @Test
    void getFilename_returnsFilenameFromLocalPath() {
        PackageVersion version = new PackageVersion(
                "1.0.0",
                "/path/to/file/driver-1.0.0.jar",
                "SHA1",
                "abc123",
                null
        );

        assertThat(version.getFilename()).isEqualTo("driver-1.0.0.jar");
    }

    @Test
    void getFilename_handlesWindowsPaths() {
        PackageVersion version = new PackageVersion(
                "1.0.0",
                "C:\\path\\to\\file\\driver-1.0.0.jar",
                "SHA1",
                "abc123",
                null
        );

        assertThat(version.getFilename()).isEqualTo("driver-1.0.0.jar");
    }

    @Test
    void getFilename_returnsEmptyForNullPath() {
        PackageVersion version = new PackageVersion("1.0.0", null, "SHA1", "abc123", null);

        assertThat(version.getFilename()).isEmpty();
    }

    @Test
    void isHttpPath_returnsTrueForHttpUrl() {
        PackageVersion version = new PackageVersion(
                "1.0.0",
                "http://example.com/file.jar",
                "SHA1",
                "abc123",
                null
        );

        assertThat(version.isHttpPath()).isTrue();
    }

    @Test
    void isHttpPath_returnsTrueForHttpsUrl() {
        PackageVersion version = new PackageVersion(
                "1.0.0",
                "https://example.com/file.jar",
                "SHA1",
                "abc123",
                null
        );

        assertThat(version.isHttpPath()).isTrue();
    }

    @Test
    void isHttpPath_returnsFalseForLocalPath() {
        PackageVersion version = new PackageVersion(
                "1.0.0",
                "/path/to/file.jar",
                "SHA1",
                "abc123",
                null
        );

        assertThat(version.isHttpPath()).isFalse();
    }

    @Test
    void isHttpPath_returnsFalseForNullPath() {
        PackageVersion version = new PackageVersion("1.0.0", null, "SHA1", "abc123", null);

        assertThat(version.isHttpPath()).isFalse();
    }

    @Test
    void isInClasspathByFilename_returnsTrueWhenFileExists() {
        PackageVersion version = new PackageVersion(
                "1.0.0",
                "https://example.com/driver-1.0.0.jar",
                "SHA1",
                "abc123",
                null
        );

        List<String> filenames = Arrays.asList("other.jar", "driver-1.0.0.jar", "another.jar");

        assertThat(version.isInClasspathByFilename(filenames)).isTrue();
    }

    @Test
    void isInClasspathByFilename_returnsFalseWhenFileNotExists() {
        PackageVersion version = new PackageVersion(
                "1.0.0",
                "https://example.com/driver-1.0.0.jar",
                "SHA1",
                "abc123",
                null
        );

        List<String> filenames = Arrays.asList("other.jar", "driver-2.0.0.jar", "another.jar");

        assertThat(version.isInClasspathByFilename(filenames)).isFalse();
    }

    @Test
    void isInClasspathByFilename_returnsFalseForEmptyList() {
        PackageVersion version = new PackageVersion(
                "1.0.0",
                "https://example.com/driver-1.0.0.jar",
                "SHA1",
                "abc123",
                null
        );

        assertThat(version.isInClasspathByFilename(Collections.emptyList())).isFalse();
    }

    @Test
    void isInClasspathByFilename_returnsFalseForNullList() {
        PackageVersion version = new PackageVersion(
                "1.0.0",
                "https://example.com/driver-1.0.0.jar",
                "SHA1",
                "abc123",
                null
        );

        assertThat(version.isInClasspathByFilename(null)).isFalse();
    }

    @Test
    void isEmpty_returnsTrueForNullTag() {
        PackageVersion version = new PackageVersion(null, "path", "SHA1", "abc", null);

        assertThat(version.isEmpty()).isTrue();
    }

    @Test
    void isEmpty_returnsTrueForBlankTag() {
        PackageVersion version = new PackageVersion("  ", "path", "SHA1", "abc", null);

        assertThat(version.isEmpty()).isTrue();
    }

    @Test
    void isEmpty_returnsFalseForValidTag() {
        PackageVersion version = new PackageVersion("1.0.0", "path", "SHA1", "abc", null);

        assertThat(version.isEmpty()).isFalse();
    }

    @Test
    void empty_createsEmptyVersion() {
        PackageVersion version = PackageVersion.empty();

        assertThat(version.isEmpty()).isTrue();
        assertThat(version.tag()).isNull();
        assertThat(version.path()).isNull();
    }
}
