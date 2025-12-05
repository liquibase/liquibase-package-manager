package org.liquibase.lpm.model;

import org.junit.jupiter.api.Test;
import org.semver4j.Semver;

import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.Optional;

import static org.assertj.core.api.Assertions.assertThat;

class PackageTest {

    private static final PackageVersion V1_0_0 = new PackageVersion(
            "1.0.0", "https://example.com/pkg-1.0.0.jar", "SHA1", "abc", "4.0.0");
    private static final PackageVersion V1_5_0 = new PackageVersion(
            "1.5.0", "https://example.com/pkg-1.5.0.jar", "SHA1", "def", "4.6.0");
    private static final PackageVersion V2_0_0 = new PackageVersion(
            "2.0.0", "https://example.com/pkg-2.0.0.jar", "SHA1", "ghi", "4.16.0");

    @Test
    void getLatestVersion_returnsLatestForExtension() {
        Package pkg = new Package("test-extension", "extension", Arrays.asList(V1_0_0, V1_5_0, V2_0_0));

        Optional<PackageVersion> latest = pkg.getLatestVersion(new Semver("4.20.0"));

        assertThat(latest).isPresent();
        assertThat(latest.get().tag()).isEqualTo("2.0.0");
    }

    @Test
    void getLatestVersion_filtersIncompatibleVersions() {
        Package pkg = new Package("test-extension", "extension", Arrays.asList(V1_0_0, V1_5_0, V2_0_0));

        // Liquibase 4.8.0 is lower than V2_0_0's requirement (4.16.0)
        Optional<PackageVersion> latest = pkg.getLatestVersion(new Semver("4.8.0"));

        assertThat(latest).isPresent();
        assertThat(latest.get().tag()).isEqualTo("1.5.0");
    }

    @Test
    void getLatestVersion_returnsOnlyCompatibleVersion() {
        Package pkg = new Package("test-extension", "extension", Arrays.asList(V1_0_0, V1_5_0, V2_0_0));

        // Liquibase 4.2.0 is only compatible with V1_0_0
        Optional<PackageVersion> latest = pkg.getLatestVersion(new Semver("4.2.0"));

        assertThat(latest).isPresent();
        assertThat(latest.get().tag()).isEqualTo("1.0.0");
    }

    @Test
    void getLatestVersion_returnsEmptyWhenNoCompatibleVersion() {
        Package pkg = new Package("test-extension", "extension", Arrays.asList(V1_0_0, V1_5_0, V2_0_0));

        // Liquibase 3.0.0 is too old for any version
        Optional<PackageVersion> latest = pkg.getLatestVersion(new Semver("3.0.0"));

        assertThat(latest).isEmpty();
    }

    @Test
    void getLatestVersion_ignoresLiquibaseVersionForDrivers() {
        Package pkg = new Package("test-driver", "driver", Arrays.asList(V1_0_0, V1_5_0, V2_0_0));

        // Even with old Liquibase, drivers return latest version
        Optional<PackageVersion> latest = pkg.getLatestVersion(new Semver("3.0.0"));

        assertThat(latest).isPresent();
        assertThat(latest.get().tag()).isEqualTo("2.0.0");
    }

    @Test
    void getLatestVersion_returnsEmptyForEmptyVersionsList() {
        Package pkg = new Package("test", "extension", Collections.emptyList());

        Optional<PackageVersion> latest = pkg.getLatestVersion(new Semver("4.0.0"));

        assertThat(latest).isEmpty();
    }

    @Test
    void getLatestVersion_returnsEmptyForNullVersionsList() {
        Package pkg = new Package("test", "extension", null);

        Optional<PackageVersion> latest = pkg.getLatestVersion(new Semver("4.0.0"));

        assertThat(latest).isEmpty();
    }

    @Test
    void getVersion_returnsMatchingVersion() {
        Package pkg = new Package("test", "extension", Arrays.asList(V1_0_0, V1_5_0, V2_0_0));

        Optional<PackageVersion> version = pkg.getVersion("1.5.0");

        assertThat(version).isPresent();
        assertThat(version.get().tag()).isEqualTo("1.5.0");
    }

    @Test
    void getVersion_returnsEmptyWhenNotFound() {
        Package pkg = new Package("test", "extension", Arrays.asList(V1_0_0, V1_5_0, V2_0_0));

        Optional<PackageVersion> version = pkg.getVersion("3.0.0");

        assertThat(version).isEmpty();
    }

    @Test
    void getVersion_returnsEmptyForNullTag() {
        Package pkg = new Package("test", "extension", Arrays.asList(V1_0_0, V1_5_0, V2_0_0));

        Optional<PackageVersion> version = pkg.getVersion(null);

        assertThat(version).isEmpty();
    }

    @Test
    void getInstalledVersion_returnsInstalledVersion() {
        Package pkg = new Package("test", "extension", Arrays.asList(V1_0_0, V1_5_0, V2_0_0));
        List<String> filenames = Arrays.asList("other.jar", "pkg-1.5.0.jar");

        Optional<PackageVersion> installed = pkg.getInstalledVersion(filenames);

        assertThat(installed).isPresent();
        assertThat(installed.get().tag()).isEqualTo("1.5.0");
    }

    @Test
    void getInstalledVersion_returnsEmptyWhenNotInstalled() {
        Package pkg = new Package("test", "extension", Arrays.asList(V1_0_0, V1_5_0, V2_0_0));
        List<String> filenames = Arrays.asList("other.jar", "different.jar");

        Optional<PackageVersion> installed = pkg.getInstalledVersion(filenames);

        assertThat(installed).isEmpty();
    }

    @Test
    void isInClasspath_returnsTrueWhenInstalled() {
        Package pkg = new Package("test", "extension", Arrays.asList(V1_0_0, V1_5_0));
        List<String> filenames = Arrays.asList("pkg-1.0.0.jar");

        assertThat(pkg.isInClasspath(filenames)).isTrue();
    }

    @Test
    void isInClasspath_returnsFalseWhenNotInstalled() {
        Package pkg = new Package("test", "extension", Arrays.asList(V1_0_0, V1_5_0));
        List<String> filenames = Arrays.asList("other.jar");

        assertThat(pkg.isInClasspath(filenames)).isFalse();
    }

    @Test
    void deleteVersion_removesSpecifiedVersion() {
        Package pkg = new Package("test", "extension", Arrays.asList(V1_0_0, V1_5_0, V2_0_0));

        List<PackageVersion> remaining = pkg.deleteVersion(V1_5_0);

        assertThat(remaining).hasSize(2);
        assertThat(remaining).extracting(PackageVersion::tag)
                .containsExactly("1.0.0", "2.0.0");
    }

    @Test
    void deleteVersion_returnsOriginalWhenVersionNotFound() {
        Package pkg = new Package("test", "extension", Arrays.asList(V1_0_0, V1_5_0));
        PackageVersion nonExistent = new PackageVersion("9.9.9", "path", "SHA1", "abc", null);

        List<PackageVersion> remaining = pkg.deleteVersion(nonExistent);

        assertThat(remaining).hasSize(2);
    }

    @Test
    void isEmpty_returnsTrueForNullName() {
        Package pkg = new Package(null, "extension", Collections.emptyList());

        assertThat(pkg.isEmpty()).isTrue();
    }

    @Test
    void isEmpty_returnsFalseForValidName() {
        Package pkg = new Package("test", "extension", Collections.emptyList());

        assertThat(pkg.isEmpty()).isFalse();
    }

    @Test
    void versions_returnsUnmodifiableList() {
        Package pkg = new Package("test", "extension", Arrays.asList(V1_0_0));

        List<PackageVersion> versions = pkg.versions();

        assertThat(versions).hasSize(1);
        // Verify it's unmodifiable
        org.junit.jupiter.api.Assertions.assertThrows(
                UnsupportedOperationException.class,
                () -> versions.add(V2_0_0)
        );
    }

    @Test
    void versions_returnsEmptyListForNull() {
        Package pkg = new Package("test", "extension", null);

        assertThat(pkg.versions()).isEmpty();
    }
}
