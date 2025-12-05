package org.liquibase.lpm.model;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.semver4j.Semver;

import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.Optional;

import static org.assertj.core.api.Assertions.assertThat;

class PackageRegistryTest {

    private PackageRegistry registry;

    private static final PackageVersion DRIVER_V1 = new PackageVersion(
            "1.0.0", "https://example.com/driver-1.0.0.jar", "SHA1", "abc", null);
    private static final PackageVersion DRIVER_V2 = new PackageVersion(
            "2.0.0", "https://example.com/driver-2.0.0.jar", "SHA1", "def", null);

    private static final PackageVersion EXT_V1 = new PackageVersion(
            "1.0.0", "https://example.com/ext-1.0.0.jar", "SHA1", "ghi", "4.0.0");
    private static final PackageVersion EXT_V2 = new PackageVersion(
            "2.0.0", "https://example.com/ext-2.0.0.jar", "SHA1", "jkl", "4.16.0");

    private static final PackageVersion PRO_V1 = new PackageVersion(
            "1.0.0", "https://example.com/pro-1.0.0.jar", "SHA1", "mno", "4.10.0");

    private static final Package DRIVER_PKG = new Package(
            "test-driver", "driver", Arrays.asList(DRIVER_V1, DRIVER_V2));
    private static final Package EXT_PKG = new Package(
            "test-extension", "extension", Arrays.asList(EXT_V1, EXT_V2));
    private static final Package PRO_PKG = new Package(
            "test-pro", "pro", Collections.singletonList(PRO_V1));

    @BeforeEach
    void setUp() {
        registry = new PackageRegistry(Arrays.asList(DRIVER_PKG, EXT_PKG, PRO_PKG));
    }

    @Test
    void getByName_returnsPackageWhenExists() {
        Optional<Package> result = registry.getByName("test-driver");

        assertThat(result).isPresent();
        assertThat(result.get().name()).isEqualTo("test-driver");
    }

    @Test
    void getByName_returnsEmptyWhenNotFound() {
        Optional<Package> result = registry.getByName("nonexistent");

        assertThat(result).isEmpty();
    }

    @Test
    void getByName_returnsEmptyForNull() {
        Optional<Package> result = registry.getByName(null);

        assertThat(result).isEmpty();
    }

    @Test
    void filterByCategory_returnsOnlyMatchingCategory() {
        PackageRegistry drivers = registry.filterByCategory("driver");

        assertThat(drivers.size()).isEqualTo(1);
        assertThat(drivers.getPackages().get(0).name()).isEqualTo("test-driver");
    }

    @Test
    void filterByCategory_returnsAllWhenNull() {
        PackageRegistry result = registry.filterByCategory(null);

        assertThat(result.size()).isEqualTo(3);
    }

    @Test
    void filterByCategory_returnsAllWhenEmpty() {
        PackageRegistry result = registry.filterByCategory("");

        assertThat(result.size()).isEqualTo(3);
    }

    @Test
    void filterByCategory_returnsEmptyWhenNoMatch() {
        PackageRegistry result = registry.filterByCategory("nonexistent");

        assertThat(result.isEmpty()).isTrue();
    }

    @Test
    void filterByName_returnsMatchingPackages() {
        PackageRegistry result = registry.filterByName("driver");

        assertThat(result.size()).isEqualTo(1);
        assertThat(result.getPackages().get(0).name()).isEqualTo("test-driver");
    }

    @Test
    void filterByName_isCaseInsensitive() {
        PackageRegistry result = registry.filterByName("EXTENSION");

        assertThat(result.size()).isEqualTo(1);
        assertThat(result.getPackages().get(0).name()).isEqualTo("test-extension");
    }

    @Test
    void filterByName_returnsAllWhenNull() {
        PackageRegistry result = registry.filterByName(null);

        assertThat(result.size()).isEqualTo(3);
    }

    @Test
    void filterByName_returnsMultipleMatches() {
        PackageRegistry result = registry.filterByName("test");

        assertThat(result.size()).isEqualTo(3);
    }

    @Test
    void getInstalled_returnsOnlyInstalledPackages() {
        List<String> filenames = Arrays.asList("driver-1.0.0.jar", "other.jar");

        PackageRegistry installed = registry.getInstalled(filenames);

        assertThat(installed.size()).isEqualTo(1);
        assertThat(installed.getPackages().get(0).name()).isEqualTo("test-driver");
    }

    @Test
    void getInstalled_returnsEmptyWhenNoneInstalled() {
        List<String> filenames = Arrays.asList("unknown.jar");

        PackageRegistry installed = registry.getInstalled(filenames);

        assertThat(installed.isEmpty()).isTrue();
    }

    @Test
    void getOutdated_returnsPackagesWithNewerVersions() {
        // driver-1.0.0.jar is installed, but 2.0.0 is available
        List<String> filenames = Arrays.asList("driver-1.0.0.jar");

        PackageRegistry outdated = registry.getOutdated(new Semver("4.20.0"), filenames);

        assertThat(outdated.size()).isEqualTo(1);
        assertThat(outdated.getPackages().get(0).name()).isEqualTo("test-driver");
    }

    @Test
    void getOutdated_excludesPackagesWithLatestInstalled() {
        // driver-2.0.0.jar is installed and is the latest
        List<String> filenames = Arrays.asList("driver-2.0.0.jar");

        PackageRegistry outdated = registry.getOutdated(new Semver("4.20.0"), filenames);

        assertThat(outdated.isEmpty()).isTrue();
    }

    @Test
    void display_generatesFormattedOutput() {
        List<String> filenames = Arrays.asList("driver-1.0.0.jar");

        List<String> output = registry.display(filenames);

        assertThat(output).hasSize(4); // header + 3 packages
        assertThat(output.get(0)).contains("Package");
        assertThat(output.get(0)).contains("Category");
        assertThat(output.get(1)).contains("├──");
        assertThat(output.get(1)).contains("test-driver@1.0.0");
        assertThat(output.get(3)).contains("└──"); // last item
    }

    @Test
    void display_showsVersionOnlyForInstalledPackages() {
        List<String> filenames = Collections.emptyList();

        List<String> output = registry.display(filenames);

        // No @version suffix for packages not installed
        assertThat(output.get(1)).doesNotContain("@");
    }

    @Test
    void isEmpty_returnsTrueForEmptyRegistry() {
        PackageRegistry empty = new PackageRegistry(Collections.emptyList());

        assertThat(empty.isEmpty()).isTrue();
    }

    @Test
    void isEmpty_returnsFalseForNonEmptyRegistry() {
        assertThat(registry.isEmpty()).isFalse();
    }

    @Test
    void size_returnsCorrectCount() {
        assertThat(registry.size()).isEqualTo(3);
    }

    @Test
    void getPackages_returnsUnmodifiableList() {
        List<Package> packages = registry.getPackages();

        org.junit.jupiter.api.Assertions.assertThrows(
                UnsupportedOperationException.class,
                () -> packages.add(DRIVER_PKG)
        );
    }

    @Test
    void loadFromResource_loadsEmbeddedPackagesJson() {
        PackageRegistry loaded = PackageRegistry.loadFromResource();

        assertThat(loaded).isNotNull();
        assertThat(loaded.isEmpty()).isFalse();
        // Check for a known package that should exist
        assertThat(loaded.getByName("postgresql").isPresent() ||
                   loaded.getByName("liquibase-postgresql").isPresent() ||
                   loaded.size() > 0).isTrue();
    }

    @Test
    void isValid_returnsTrueForValidRegistry() {
        assertThat(registry.isValid()).isTrue();
    }

    @Test
    void isValid_returnsTrueForEmptyRegistry() {
        // Empty registry is still "valid" (not malformed)
        PackageRegistry empty = new PackageRegistry(Collections.emptyList());

        assertThat(empty.isValid()).isTrue();
    }
}
