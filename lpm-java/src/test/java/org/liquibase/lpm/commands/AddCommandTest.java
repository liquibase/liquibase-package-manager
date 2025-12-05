package org.liquibase.lpm.commands;

import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.junit.jupiter.api.io.TempDir;
import org.liquibase.lpm.exception.PackageAlreadyInstalledException;
import org.liquibase.lpm.exception.PackageNotFoundException;
import org.liquibase.lpm.model.PackageVersion;
import org.liquibase.lpm.service.ClasspathService;
import org.liquibase.lpm.service.PackageService;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import picocli.CommandLine;

import java.io.ByteArrayOutputStream;
import java.io.PrintStream;
import java.nio.file.Path;

import static org.assertj.core.api.Assertions.assertThat;
import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class AddCommandTest {

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
    void add_singlePackage_installsSuccessfully() {
        AddCommand addCommand = new AddCommand(packageService);

        PackageVersion installed = new PackageVersion("4.5.0", "http://example.com/pkg-4.5.0.jar", "SHA1", "abc", null);
        when(packageService.getClasspathService()).thenReturn(classpathService);
        when(classpathService.getClasspathDir()).thenReturn(tempDir);
        when(packageService.resolveAndInstall("liquibase-postgresql")).thenReturn(installed);

        int exitCode = addCommand.add(new String[]{"liquibase-postgresql"});

        assertThat(exitCode).isEqualTo(0);
        verify(packageService).resolveAndInstall("liquibase-postgresql");
        assertThat(outContent.toString()).contains("pkg-4.5.0.jar");
    }

    @Test
    void add_packageWithVersion_installsSpecificVersion() {
        AddCommand addCommand = new AddCommand(packageService);

        PackageVersion installed = new PackageVersion("4.5.0", "http://example.com/pkg-4.5.0.jar", "SHA1", "abc", null);
        when(packageService.getClasspathService()).thenReturn(classpathService);
        when(classpathService.getClasspathDir()).thenReturn(tempDir);
        when(packageService.resolveAndInstall("liquibase-postgresql@4.5.0")).thenReturn(installed);

        int exitCode = addCommand.add(new String[]{"liquibase-postgresql@4.5.0"});

        assertThat(exitCode).isEqualTo(0);
        verify(packageService).resolveAndInstall("liquibase-postgresql@4.5.0");
    }

    @Test
    void add_multiplePackages_installsAll() {
        AddCommand addCommand = new AddCommand(packageService);

        PackageVersion pkg1 = new PackageVersion("1.0.0", "http://example.com/pkg1-1.0.0.jar", "SHA1", "abc", null);
        PackageVersion pkg2 = new PackageVersion("2.0.0", "http://example.com/pkg2-2.0.0.jar", "SHA1", "def", null);

        when(packageService.getClasspathService()).thenReturn(classpathService);
        when(classpathService.getClasspathDir()).thenReturn(tempDir);
        when(packageService.resolveAndInstall("pkg1")).thenReturn(pkg1);
        when(packageService.resolveAndInstall("pkg2")).thenReturn(pkg2);

        int exitCode = addCommand.add(new String[]{"pkg1", "pkg2"});

        assertThat(exitCode).isEqualTo(0);
        verify(packageService).resolveAndInstall("pkg1");
        verify(packageService).resolveAndInstall("pkg2");
    }

    @Test
    void add_packageNotFound_continuesWithOthers() {
        AddCommand addCommand = new AddCommand(packageService);

        when(packageService.getClasspathService()).thenReturn(classpathService);
        when(classpathService.getClasspathDir()).thenReturn(tempDir);
        when(packageService.resolveAndInstall("nonexistent"))
                .thenThrow(new PackageNotFoundException("nonexistent"));
        PackageVersion pkg2 = new PackageVersion("2.0.0", "http://example.com/pkg2-2.0.0.jar", "SHA1", "def", null);
        when(packageService.resolveAndInstall("pkg2")).thenReturn(pkg2);

        int exitCode = addCommand.add(new String[]{"nonexistent", "pkg2"});

        assertThat(exitCode).isEqualTo(0);
        assertThat(errContent.toString()).contains("not found");
        verify(packageService).resolveAndInstall("pkg2");
    }

    @Test
    void add_alreadyInstalled_printsWarning() {
        AddCommand addCommand = new AddCommand(packageService);

        when(packageService.getClasspathService()).thenReturn(classpathService);
        when(classpathService.getClasspathDir()).thenReturn(tempDir);
        when(packageService.resolveAndInstall("already-there"))
                .thenThrow(new PackageAlreadyInstalledException("already-there", "1.0.0"));

        int exitCode = addCommand.add(new String[]{"already-there"});

        assertThat(exitCode).isEqualTo(0);
        assertThat(errContent.toString()).contains("already installed");
    }

    @Test
    void add_noPackages_printsHelp() {
        AddCommand addCommand = new AddCommand(packageService);

        int exitCode = addCommand.add(null);

        assertThat(exitCode).isEqualTo(0);
        assertThat(outContent.toString()).contains("specify at least one package");
    }

    @Test
    void add_emptyPackages_printsHelp() {
        AddCommand addCommand = new AddCommand(packageService);

        int exitCode = addCommand.add(new String[]{});

        assertThat(exitCode).isEqualTo(0);
        assertThat(outContent.toString()).contains("specify at least one package");
    }

    @Test
    void add_viaCli_parsesPackages() {
        AddCommand addCommand = new AddCommand(packageService);

        PackageVersion installed = new PackageVersion("1.0.0", "http://example.com/pkg-1.0.0.jar", "SHA1", "abc", null);
        when(packageService.getClasspathService()).thenReturn(classpathService);
        when(classpathService.getClasspathDir()).thenReturn(tempDir);
        when(packageService.resolveAndInstall(anyString())).thenReturn(installed);

        CommandLine cmd = new CommandLine(new RootCommand());
        cmd.addSubcommand("add", addCommand);

        int exitCode = cmd.execute("add", "pkg1", "pkg2@1.0.0");

        assertThat(exitCode).isEqualTo(0);
        verify(packageService).resolveAndInstall("pkg1");
        verify(packageService).resolveAndInstall("pkg2@1.0.0");
    }
}
