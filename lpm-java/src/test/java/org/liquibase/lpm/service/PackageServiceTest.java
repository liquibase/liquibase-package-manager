package org.liquibase.lpm.service;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.junit.jupiter.api.io.TempDir;
import org.liquibase.lpm.exception.*;
import org.liquibase.lpm.model.*;
import org.liquibase.lpm.service.LiquibaseDetector.LiquibaseInfo;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.semver4j.Semver;

import java.nio.file.Path;
import java.util.Collections;
import java.util.List;
import java.util.Map;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class PackageServiceTest {

    @TempDir
    Path tempDir;

    @Mock
    private ClasspathService classpathService;

    @Mock
    private LiquibaseDetector liquibaseDetector;

    private PackageService packageService;

    private LiquibaseInfo defaultLiquibaseInfo;

    @BeforeEach
    void setUp() {
        packageService = new PackageService(classpathService, liquibaseDetector);
        defaultLiquibaseInfo = new LiquibaseInfo(
                tempDir,
                new Semver("4.25.0"),
                Map.of("build.version", "4.25.0")
        );
    }

    @Test
    void initialize_global_usesGlobalClasspath() {
        when(liquibaseDetector.detect()).thenReturn(defaultLiquibaseInfo);
        when(classpathService.packagesJsonExists()).thenReturn(false);

        packageService.initialize(true);

        verify(classpathService).useGlobalClasspath(tempDir);
        verify(classpathService).writePackagesJson(any());
    }

    @Test
    void initialize_local_usesLocalClasspath() {
        when(liquibaseDetector.detect()).thenReturn(defaultLiquibaseInfo);
        when(classpathService.packagesJsonExists()).thenReturn(false);

        packageService.initialize(false);

        verify(classpathService).useLocalClasspath();
    }

    @Test
    void initialize_existingPackagesJson_loadsFromClasspath() {
        when(liquibaseDetector.detect()).thenReturn(defaultLiquibaseInfo);
        when(classpathService.packagesJsonExists()).thenReturn(true);
        String packagesJson = """
            {
                "packages": [
                    {
                        "name": "test-package",
                        "category": "extension",
                        "versions": []
                    }
                ]
            }
            """;
        when(classpathService.readPackagesJson()).thenReturn(packagesJson.getBytes());

        packageService.initialize(true);

        verify(classpathService, never()).writePackagesJson(any());
        assertThat(packageService.getPackageRegistry().getByName("test-package")).isPresent();
    }

    @Test
    void resolveAndInstall_specificVersion_installsVersion() {
        setupWithPackage("test-pkg", "extension", List.of(
                new PackageVersion("1.0.0", "http://example.com/test-1.0.0.jar", "SHA1", "checksum", null),
                new PackageVersion("2.0.0", "http://example.com/test-2.0.0.jar", "SHA1", "checksum", null)
        ));
        when(classpathService.getClasspathFilenames()).thenReturn(Collections.emptyList());

        PackageVersion installed = packageService.resolveAndInstall("test-pkg@1.0.0");

        assertThat(installed.tag()).isEqualTo("1.0.0");
        verify(classpathService).installVersion(installed);
    }

    @Test
    void resolveAndInstall_noVersion_installsLatestCompatible() {
        setupWithPackage("test-pkg", "extension", List.of(
                new PackageVersion("1.0.0", "http://example.com/test-1.0.0.jar", "SHA1", "checksum", null),
                new PackageVersion("2.0.0", "http://example.com/test-2.0.0.jar", "SHA1", "checksum", null)
        ));
        when(classpathService.getClasspathFilenames()).thenReturn(Collections.emptyList());

        PackageVersion installed = packageService.resolveAndInstall("test-pkg");

        assertThat(installed.tag()).isEqualTo("2.0.0");
        verify(classpathService).installVersion(installed);
    }

    @Test
    void resolveAndInstall_packageNotFound_throwsException() {
        setupWithPackage("test-pkg", "extension", List.of(
                new PackageVersion("1.0.0", "http://example.com/test-1.0.0.jar", "SHA1", "checksum", null)
        ));

        assertThatThrownBy(() -> packageService.resolveAndInstall("nonexistent"))
                .isInstanceOf(PackageNotFoundException.class)
                .hasMessageContaining("nonexistent");
    }

    @Test
    void resolveAndInstall_versionNotFound_throwsException() {
        setupWithPackage("test-pkg", "extension", List.of(
                new PackageVersion("1.0.0", "http://example.com/test-1.0.0.jar", "SHA1", "checksum", null)
        ));
        when(classpathService.getClasspathFilenames()).thenReturn(Collections.emptyList());

        assertThatThrownBy(() -> packageService.resolveAndInstall("test-pkg@9.9.9"))
                .isInstanceOf(VersionNotFoundException.class)
                .hasMessageContaining("9.9.9");
    }

    @Test
    void resolveAndInstall_alreadyInstalled_throwsException() {
        setupWithPackage("test-pkg", "extension", List.of(
                new PackageVersion("1.0.0", "http://example.com/test-pkg-1.0.0.jar", "SHA1", "checksum", null)
        ));
        when(classpathService.getClasspathFilenames()).thenReturn(List.of("test-pkg-1.0.0.jar"));

        assertThatThrownBy(() -> packageService.resolveAndInstall("test-pkg"))
                .isInstanceOf(PackageAlreadyInstalledException.class)
                .hasMessageContaining("test-pkg");
    }

    @Test
    void resolveAndInstall_incompatibleLiquibase_throwsException() {
        // Package requires Liquibase 4.30.0
        setupWithPackage("test-pkg", "extension", List.of(
                new PackageVersion("1.0.0", "http://example.com/test-1.0.0.jar", "SHA1", "checksum", "4.30.0")
        ));
        when(classpathService.getClasspathFilenames()).thenReturn(Collections.emptyList());

        assertThatThrownBy(() -> packageService.resolveAndInstall("test-pkg@1.0.0"))
                .isInstanceOf(VersionIncompatibleException.class)
                .hasMessageContaining("4.30.0");
    }

    @Test
    void resolveAndInstall_driverCategory_skipsCompatibilityCheck() {
        // Driver with high requirement should still install
        setupWithPackage("postgres-driver", "driver", List.of(
                new PackageVersion("1.0.0", "http://example.com/postgres-1.0.0.jar", "SHA1", "checksum", "99.0.0")
        ));
        when(classpathService.getClasspathFilenames()).thenReturn(Collections.emptyList());

        // Should not throw
        PackageVersion installed = packageService.resolveAndInstall("postgres-driver@1.0.0");

        assertThat(installed.tag()).isEqualTo("1.0.0");
    }

    @Test
    void removePackage_installedPackage_removesIt() {
        setupWithPackage("test-pkg", "extension", List.of(
                new PackageVersion("1.0.0", "http://example.com/test-pkg-1.0.0.jar", "SHA1", "checksum", null)
        ));
        when(classpathService.getClasspathFilenames()).thenReturn(List.of("test-pkg-1.0.0.jar"));

        PackageVersion removed = packageService.removePackage("test-pkg");

        assertThat(removed.tag()).isEqualTo("1.0.0");
        verify(classpathService).removeVersion(removed);
    }

    @Test
    void removePackage_packageNotFound_throwsException() {
        setupWithPackage("test-pkg", "extension", List.of(
                new PackageVersion("1.0.0", "http://example.com/test-1.0.0.jar", "SHA1", "checksum", null)
        ));

        assertThatThrownBy(() -> packageService.removePackage("nonexistent"))
                .isInstanceOf(PackageNotFoundException.class);
    }

    @Test
    void removePackage_notInstalled_throwsException() {
        setupWithPackage("test-pkg", "extension", List.of(
                new PackageVersion("1.0.0", "http://example.com/test-1.0.0.jar", "SHA1", "checksum", null)
        ));
        when(classpathService.getClasspathFilenames()).thenReturn(Collections.emptyList());

        assertThatThrownBy(() -> packageService.removePackage("test-pkg"))
                .isInstanceOf(PackageNotInstalledException.class)
                .hasMessageContaining("test-pkg");
    }

    @Test
    void getOutdatedPackages_returnsOutdated() {
        setupWithPackage("test-pkg", "extension", List.of(
                new PackageVersion("1.0.0", "http://example.com/test-pkg-1.0.0.jar", "SHA1", "checksum", null),
                new PackageVersion("2.0.0", "http://example.com/test-pkg-2.0.0.jar", "SHA1", "checksum", null)
        ));
        when(classpathService.getClasspathFilenames()).thenReturn(List.of("test-pkg-1.0.0.jar"));

        PackageRegistry outdated = packageService.getOutdatedPackages();

        assertThat(outdated.getPackages()).hasSize(1);
        assertThat(outdated.getPackages().getFirst().name()).isEqualTo("test-pkg");
    }

    @Test
    void getInstalledPackages_returnsInstalled() {
        setupWithPackage("test-pkg", "extension", List.of(
                new PackageVersion("1.0.0", "http://example.com/test-pkg-1.0.0.jar", "SHA1", "checksum", null)
        ));
        when(classpathService.getClasspathFilenames()).thenReturn(List.of("test-pkg-1.0.0.jar"));

        PackageRegistry installed = packageService.getInstalledPackages();

        assertThat(installed.getPackages()).hasSize(1);
    }

    @Test
    void requiresJavaOptsInstructions_oldLiquibase_returnsTrue() {
        LiquibaseInfo oldInfo = new LiquibaseInfo(
                tempDir,
                new Semver("4.5.0"),
                Map.of("build.version", "4.5.0")
        );
        when(liquibaseDetector.detect()).thenReturn(oldInfo);
        when(classpathService.packagesJsonExists()).thenReturn(false);

        packageService.initialize(true);

        assertThat(packageService.requiresJavaOptsInstructions()).isTrue();
    }

    @Test
    void requiresJavaOptsInstructions_newLiquibase_returnsFalse() {
        when(liquibaseDetector.detect()).thenReturn(defaultLiquibaseInfo); // 4.25.0
        when(classpathService.packagesJsonExists()).thenReturn(false);

        packageService.initialize(true);

        assertThat(packageService.requiresJavaOptsInstructions()).isFalse();
    }

    @Test
    void filterByCategory_filtersCorrectly() {
        String packagesJson = """
            {
                "packages": [
                    {"name": "ext1", "category": "extension", "versions": []},
                    {"name": "drv1", "category": "driver", "versions": []}
                ]
            }
            """;
        when(liquibaseDetector.detect()).thenReturn(defaultLiquibaseInfo);
        when(classpathService.packagesJsonExists()).thenReturn(true);
        when(classpathService.readPackagesJson()).thenReturn(packagesJson.getBytes());

        packageService.initialize(true);
        PackageRegistry filtered = packageService.filterByCategory("driver");

        assertThat(filtered.getPackages()).hasSize(1);
        assertThat(filtered.getPackages().getFirst().name()).isEqualTo("drv1");
    }

    @Test
    void searchPackages_findsMatches() {
        String packagesJson = """
            {
                "packages": [
                    {"name": "liquibase-postgresql", "category": "extension", "versions": []},
                    {"name": "liquibase-mysql", "category": "extension", "versions": []},
                    {"name": "some-other", "category": "extension", "versions": []}
                ]
            }
            """;
        when(liquibaseDetector.detect()).thenReturn(defaultLiquibaseInfo);
        when(classpathService.packagesJsonExists()).thenReturn(true);
        when(classpathService.readPackagesJson()).thenReturn(packagesJson.getBytes());

        packageService.initialize(true);
        PackageRegistry results = packageService.searchPackages("postgres");

        assertThat(results.getPackages()).hasSize(1);
        assertThat(results.getPackages().getFirst().name()).isEqualTo("liquibase-postgresql");
    }

    @Test
    void getLiquibaseInfo_returnsDetectedInfo() {
        when(liquibaseDetector.detect()).thenReturn(defaultLiquibaseInfo);
        when(classpathService.packagesJsonExists()).thenReturn(false);

        packageService.initialize(true);

        assertThat(packageService.getLiquibaseInfo()).isEqualTo(defaultLiquibaseInfo);
    }

    @Test
    void getClasspathService_returnsService() {
        assertThat(packageService.getClasspathService()).isEqualTo(classpathService);
    }

    /**
     * Helper to set up the package service with a single package.
     */
    private void setupWithPackage(String name, String category, List<PackageVersion> versions) {
        Package pkg = new Package(name, category, versions);
        String packagesJson = String.format("""
            {
                "packages": [
                    {
                        "name": "%s",
                        "category": "%s",
                        "versions": %s
                    }
                ]
            }
            """, name, category, versionsToJson(versions));

        when(liquibaseDetector.detect()).thenReturn(defaultLiquibaseInfo);
        when(classpathService.packagesJsonExists()).thenReturn(true);
        when(classpathService.readPackagesJson()).thenReturn(packagesJson.getBytes());

        packageService.initialize(true);
    }

    private String versionsToJson(List<PackageVersion> versions) {
        StringBuilder sb = new StringBuilder("[");
        for (int i = 0; i < versions.size(); i++) {
            PackageVersion v = versions.get(i);
            if (i > 0) sb.append(",");
            sb.append(String.format("""
                {
                    "tag": "%s",
                    "path": "%s",
                    "algorithm": %s,
                    "checksum": %s,
                    "liquibaseCore": %s
                }
                """,
                    v.tag(),
                    v.path(),
                    v.algorithm() != null ? "\"" + v.algorithm() + "\"" : "null",
                    v.checksum() != null ? "\"" + v.checksum() + "\"" : "null",
                    v.liquibaseCore() != null ? "\"" + v.liquibaseCore() + "\"" : "null"
            ));
        }
        sb.append("]");
        return sb.toString();
    }
}
