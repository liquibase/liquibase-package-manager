package org.liquibase.lpm.commands;

import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.junit.jupiter.api.io.TempDir;
import org.liquibase.lpm.exception.PackageNotFoundException;
import org.liquibase.lpm.exception.PackageNotInstalledException;
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
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class RemoveCommandTest {

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
    void remove_singlePackage_removesSuccessfully() {
        RemoveCommand removeCommand = new RemoveCommand(packageService);

        PackageVersion removed = new PackageVersion("4.5.0", "http://example.com/pkg-4.5.0.jar", "SHA1", "abc", null);
        when(packageService.getClasspathService()).thenReturn(classpathService);
        when(classpathService.getClasspathDir()).thenReturn(tempDir);
        when(packageService.removePackage("liquibase-postgresql")).thenReturn(removed);

        int exitCode = removeCommand.remove(new String[]{"liquibase-postgresql"});

        assertThat(exitCode).isEqualTo(0);
        verify(packageService).removePackage("liquibase-postgresql");
        assertThat(outContent.toString()).contains("Removed");
    }

    @Test
    void remove_multiplePackages_removesAll() {
        RemoveCommand removeCommand = new RemoveCommand(packageService);

        PackageVersion pkg1 = new PackageVersion("1.0.0", "http://example.com/pkg1-1.0.0.jar", "SHA1", "abc", null);
        PackageVersion pkg2 = new PackageVersion("2.0.0", "http://example.com/pkg2-2.0.0.jar", "SHA1", "def", null);

        when(packageService.getClasspathService()).thenReturn(classpathService);
        when(classpathService.getClasspathDir()).thenReturn(tempDir);
        when(packageService.removePackage("pkg1")).thenReturn(pkg1);
        when(packageService.removePackage("pkg2")).thenReturn(pkg2);

        int exitCode = removeCommand.remove(new String[]{"pkg1", "pkg2"});

        assertThat(exitCode).isEqualTo(0);
        verify(packageService).removePackage("pkg1");
        verify(packageService).removePackage("pkg2");
    }

    @Test
    void remove_packageNotFound_continuesWithOthers() {
        RemoveCommand removeCommand = new RemoveCommand(packageService);

        when(packageService.getClasspathService()).thenReturn(classpathService);
        when(classpathService.getClasspathDir()).thenReturn(tempDir);
        when(packageService.removePackage("nonexistent"))
                .thenThrow(new PackageNotFoundException("nonexistent"));
        PackageVersion pkg2 = new PackageVersion("2.0.0", "http://example.com/pkg2-2.0.0.jar", "SHA1", "def", null);
        when(packageService.removePackage("pkg2")).thenReturn(pkg2);

        int exitCode = removeCommand.remove(new String[]{"nonexistent", "pkg2"});

        assertThat(exitCode).isEqualTo(0);
        assertThat(errContent.toString()).contains("not found");
        verify(packageService).removePackage("pkg2");
    }

    @Test
    void remove_notInstalled_printsWarning() {
        RemoveCommand removeCommand = new RemoveCommand(packageService);

        when(packageService.getClasspathService()).thenReturn(classpathService);
        when(classpathService.getClasspathDir()).thenReturn(tempDir);
        when(packageService.removePackage("not-installed"))
                .thenThrow(new PackageNotInstalledException("not-installed"));

        int exitCode = removeCommand.remove(new String[]{"not-installed"});

        assertThat(exitCode).isEqualTo(0);
        assertThat(errContent.toString()).contains("not installed");
    }

    @Test
    void remove_noPackages_printsHelp() {
        RemoveCommand removeCommand = new RemoveCommand(packageService);

        int exitCode = removeCommand.remove(null);

        assertThat(exitCode).isEqualTo(0);
        assertThat(outContent.toString()).contains("specify at least one package");
    }

    @Test
    void remove_viaCli_parsesPackages() {
        RemoveCommand removeCommand = new RemoveCommand(packageService);

        PackageVersion removed = new PackageVersion("1.0.0", "http://example.com/pkg-1.0.0.jar", "SHA1", "abc", null);
        when(packageService.getClasspathService()).thenReturn(classpathService);
        when(classpathService.getClasspathDir()).thenReturn(tempDir);
        when(packageService.removePackage(anyString())).thenReturn(removed);

        CommandLine cmd = new CommandLine(new RootCommand());
        cmd.addSubcommand("remove", removeCommand);

        int exitCode = cmd.execute("remove", "pkg1", "pkg2");

        assertThat(exitCode).isEqualTo(0);
        verify(packageService).removePackage("pkg1");
        verify(packageService).removePackage("pkg2");
    }

    @Test
    void remove_rmAlias_works() {
        RemoveCommand removeCommand = new RemoveCommand(packageService);

        PackageVersion removed = new PackageVersion("1.0.0", "http://example.com/pkg-1.0.0.jar", "SHA1", "abc", null);
        when(packageService.getClasspathService()).thenReturn(classpathService);
        when(classpathService.getClasspathDir()).thenReturn(tempDir);
        when(packageService.removePackage("pkg")).thenReturn(removed);

        CommandLine cmd = new CommandLine(new RootCommand());
        cmd.addSubcommand("remove", removeCommand);

        // Test with 'rm' alias
        int exitCode = cmd.execute("rm", "pkg");

        assertThat(exitCode).isEqualTo(0);
        verify(packageService).removePackage("pkg");
    }
}
