package org.liquibase.lpm.util;

import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.liquibase.lpm.model.Package;
import org.liquibase.lpm.model.PackageVersion;

import java.io.ByteArrayOutputStream;
import java.io.PrintStream;
import java.util.Collections;
import java.util.List;
import java.util.Optional;

import static org.assertj.core.api.Assertions.assertThat;

class ConsoleOutputTest {

    private ByteArrayOutputStream outContent;
    private PrintStream originalOut;

    @BeforeEach
    void setUp() {
        outContent = new ByteArrayOutputStream();
        originalOut = System.out;
        System.setOut(new PrintStream(outContent));
    }

    @AfterEach
    void tearDown() {
        System.setOut(originalOut);
    }

    @Test
    void formatPackages_emptyList_returnsHeaderOnly() {
        List<String> result = ConsoleOutput.formatPackages(Collections.emptyList(), Collections.emptyList());

        assertThat(result).hasSize(1);
        assertThat(result.get(0)).contains("Package").contains("Category");
    }

    @Test
    void formatPackages_singlePackage_usesLastBranch() {
        Package pkg = new Package("test-pkg", "extension", Collections.emptyList());

        List<String> result = ConsoleOutput.formatPackages(List.of(pkg), Collections.emptyList());

        assertThat(result).hasSize(2);
        assertThat(result.get(1)).contains(ConsoleOutput.LAST_BRANCH);
        assertThat(result.get(1)).contains("test-pkg");
        assertThat(result.get(1)).contains("extension");
    }

    @Test
    void formatPackages_multiplePackages_usesCorrectBranches() {
        Package pkg1 = new Package("pkg1", "extension", Collections.emptyList());
        Package pkg2 = new Package("pkg2", "driver", Collections.emptyList());
        Package pkg3 = new Package("pkg3", "pro", Collections.emptyList());

        List<String> result = ConsoleOutput.formatPackages(List.of(pkg1, pkg2, pkg3), Collections.emptyList());

        assertThat(result).hasSize(4);
        assertThat(result.get(1)).contains(ConsoleOutput.BRANCH);
        assertThat(result.get(2)).contains(ConsoleOutput.BRANCH);
        assertThat(result.get(3)).contains(ConsoleOutput.LAST_BRANCH);
    }

    @Test
    void formatPackages_installedPackage_showsVersion() {
        PackageVersion version = new PackageVersion(
                "1.0.0", "http://example.com/pkg-1.0.0.jar", "SHA1", "abc", null
        );
        Package pkg = new Package("test-pkg", "extension", List.of(version));

        List<String> result = ConsoleOutput.formatPackages(
                List.of(pkg),
                List.of("pkg-1.0.0.jar")
        );

        assertThat(result.get(1)).contains("test-pkg@1.0.0");
    }

    @Test
    void formatOutdatedPackages_showsInstalledAndAvailable() {
        PackageVersion v1 = new PackageVersion("1.0.0", "http://example.com/pkg-1.0.0.jar", "SHA1", "abc", null);
        PackageVersion v2 = new PackageVersion("2.0.0", "http://example.com/pkg-2.0.0.jar", "SHA1", "def", null);
        Package pkg = new Package("test-pkg", "extension", List.of(v1, v2));

        List<String> result = ConsoleOutput.formatOutdatedPackages(
                List.of(pkg),
                List.of("pkg-1.0.0.jar"),
                p -> Optional.of(v2)
        );

        assertThat(result).hasSize(2);
        assertThat(result.get(0)).contains("Installed").contains("Available");
        assertThat(result.get(1)).contains("1.0.0").contains("2.0.0");
    }

    @Test
    void formatVersionsForDedupe_formatsCorrectly() {
        PackageVersion v1 = new PackageVersion("2.0.0", "http://example.com/pkg-2.0.0.jar", "SHA1", "abc", null);
        PackageVersion v2 = new PackageVersion("1.0.0", "http://example.com/pkg-1.0.0.jar", "SHA1", "def", null);

        List<String> result = ConsoleOutput.formatVersionsForDedupe("test-pkg", List.of(v1, v2));

        assertThat(result).hasSize(3);
        assertThat(result.get(0)).isEqualTo("Package");
        assertThat(result.get(1)).contains(ConsoleOutput.BRANCH).contains("test-pkg@2.0.0");
        assertThat(result.get(2)).contains(ConsoleOutput.LAST_BRANCH).contains("test-pkg@1.0.0");
    }

    @Test
    void printLines_printsAllLines() {
        List<String> lines = List.of("line1", "line2", "line3");

        ConsoleOutput.printLines(lines);

        String output = outContent.toString();
        assertThat(output).contains("line1").contains("line2").contains("line3");
    }

    @Test
    void printInstallSuccess_formatsCorrectly() {
        ConsoleOutput.printInstallSuccess("pkg-1.0.0.jar");

        assertThat(outContent.toString()).contains("pkg-1.0.0.jar")
                .contains("successfully installed");
    }

    @Test
    void printRemoveSuccess_formatsCorrectly() {
        ConsoleOutput.printRemoveSuccess("pkg-1.0.0.jar");

        assertThat(outContent.toString()).contains("pkg-1.0.0.jar")
                .contains("successfully removed");
    }

    @Test
    void printClasspath_printsPath() {
        ConsoleOutput.printClasspath("/home/user/liquibase/lib");

        assertThat(outContent.toString()).contains("Classpath:")
                .contains("/home/user/liquibase/lib");
    }

    @Test
    void printOutdatedCount_zeroPackages_printsUpToDate() {
        ConsoleOutput.printOutdatedCount(0);

        assertThat(outContent.toString()).contains("All packages are up to date");
    }

    @Test
    void printOutdatedCount_onePackage_usesSingular() {
        ConsoleOutput.printOutdatedCount(1);

        assertThat(outContent.toString()).contains("1 package can be upgraded");
    }

    @Test
    void printOutdatedCount_multiplePackages_usesPlural() {
        ConsoleOutput.printOutdatedCount(5);

        assertThat(outContent.toString()).contains("5 packages can be upgraded");
    }

    @Test
    void printJavaOptsInstructions_includesClasspath() {
        ConsoleOutput.printJavaOptsInstructions("./liquibase_libs");

        String output = outContent.toString();
        assertThat(output).contains("JAVA_OPTS");
        assertThat(output).contains("./liquibase_libs");
        assertThat(output).contains("4.6.2");
    }

    @Test
    void constants_areCorrect() {
        assertThat(ConsoleOutput.BRANCH).isEqualTo("├──");
        assertThat(ConsoleOutput.LAST_BRANCH).isEqualTo("└──");
    }
}
