package org.liquibase.lpm.commands;

import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.junit.jupiter.api.io.TempDir;
import org.liquibase.lpm.model.Package;
import org.liquibase.lpm.model.PackageRegistry;
import org.liquibase.lpm.model.PackageVersion;
import org.liquibase.lpm.service.ClasspathService;
import org.liquibase.lpm.service.LiquibaseDetector.LiquibaseInfo;
import org.liquibase.lpm.service.PackageService;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.semver4j.Semver;

import java.io.ByteArrayOutputStream;
import java.io.PrintStream;
import java.nio.file.Path;
import java.util.Collections;
import java.util.List;
import java.util.Map;

import static org.assertj.core.api.Assertions.assertThat;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class ListCommandTest {

    @TempDir
    Path tempDir;

    @Mock
    private PackageService packageService;

    @Mock
    private ClasspathService classpathService;

    private ByteArrayOutputStream outContent;
    private ByteArrayOutputStream errContent;
    private PrintStream originalOut;
    private PrintStream originalErr;

    @BeforeEach
    void setUp() {
        outContent = new ByteArrayOutputStream();
        errContent = new ByteArrayOutputStream();
        originalOut = System.out;
        originalErr = System.err;
        System.setOut(new PrintStream(outContent));
        System.setErr(new PrintStream(errContent));
    }

    @AfterEach
    void tearDown() {
        System.setOut(originalOut);
        System.setErr(originalErr);
    }

    @Test
    void list_noInstalledPackages_printsNoPackages() {
        ListCommand listCommand = new ListCommand(packageService);

        when(packageService.getClasspathService()).thenReturn(classpathService);
        when(classpathService.getClasspathDir()).thenReturn(tempDir);
        when(packageService.getInstalledPackages()).thenReturn(new PackageRegistry(Collections.emptyList()));

        int exitCode = listCommand.call();

        assertThat(exitCode).isEqualTo(0);
        assertThat(outContent.toString()).contains("No packages installed");
    }

    @Test
    void list_withInstalledPackages_printsPackages() {
        ListCommand listCommand = new ListCommand(packageService);

        Package pkg = new Package(
                "liquibase-postgresql",
                "extension",
                List.of(new PackageVersion("4.5.0", "http://example.com/liquibase-postgresql-4.5.0.jar", "SHA1", "abc", null))
        );
        PackageRegistry installed = new PackageRegistry(List.of(pkg));

        when(packageService.getClasspathService()).thenReturn(classpathService);
        when(classpathService.getClasspathDir()).thenReturn(tempDir);
        when(packageService.getInstalledPackages()).thenReturn(installed);
        when(packageService.getClasspathFilenames()).thenReturn(List.of("liquibase-postgresql-4.5.0.jar"));
        when(packageService.getLiquibaseInfo()).thenReturn(
                new LiquibaseInfo(tempDir, new Semver("4.25.0"), Map.of())
        );

        int exitCode = listCommand.call();

        assertThat(exitCode).isEqualTo(0);
        assertThat(outContent.toString()).contains("liquibase-postgresql");
    }

    @Test
    void list_exception_returnsErrorCode() {
        ListCommand listCommand = new ListCommand(packageService);

        when(packageService.getClasspathService()).thenReturn(classpathService);
        when(classpathService.getClasspathDir()).thenReturn(tempDir);
        when(packageService.getInstalledPackages()).thenThrow(new RuntimeException("Test error"));

        int exitCode = listCommand.call();

        assertThat(exitCode).isEqualTo(1);
        assertThat(errContent.toString()).contains("Test error");
    }
}
