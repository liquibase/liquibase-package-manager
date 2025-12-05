package org.liquibase.lpm.commands;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.junit.jupiter.api.io.TempDir;
import org.liquibase.lpm.model.Package;
import org.liquibase.lpm.model.PackageRegistry;
import org.liquibase.lpm.model.PackageVersion;
import org.liquibase.lpm.service.LiquibaseDetector.LiquibaseInfo;
import org.liquibase.lpm.service.PackageService;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.semver4j.Semver;
import picocli.CommandLine;

import java.io.ByteArrayOutputStream;
import java.io.PrintStream;
import java.nio.file.Path;
import java.util.Collections;
import java.util.List;
import java.util.Map;

import static org.assertj.core.api.Assertions.assertThat;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class SearchCommandTest {

    @TempDir
    Path tempDir;

    @Mock
    private PackageService packageService;

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

    void tearDown() {
        System.setOut(originalOut);
        System.setErr(originalErr);
    }

    @Test
    void search_noResults_printsNoPackagesFound() {
        SearchCommand searchCommand = new SearchCommand(packageService);

        when(packageService.searchPackages("nonexistent"))
                .thenReturn(new PackageRegistry(Collections.emptyList()));

        int exitCode = searchCommand.search("nonexistent");

        assertThat(exitCode).isEqualTo(0);
        assertThat(outContent.toString()).contains("No packages found");
        tearDown();
    }

    @Test
    void search_withResults_printsPackages() {
        SearchCommand searchCommand = new SearchCommand(packageService);

        Package pkg = new Package(
                "liquibase-postgresql",
                "extension",
                List.of(new PackageVersion("4.5.0", "http://example.com/pkg.jar", "SHA1", "abc", null))
        );
        PackageRegistry results = new PackageRegistry(List.of(pkg));

        when(packageService.searchPackages("postgres")).thenReturn(results);
        when(packageService.getClasspathFilenames()).thenReturn(Collections.emptyList());
        when(packageService.getLiquibaseInfo()).thenReturn(
                new LiquibaseInfo(tempDir, new Semver("4.25.0"), Map.of())
        );

        int exitCode = searchCommand.search("postgres");

        assertThat(exitCode).isEqualTo(0);
        assertThat(outContent.toString()).contains("liquibase-postgresql");
        tearDown();
    }

    @Test
    void search_exception_returnsErrorCode() {
        SearchCommand searchCommand = new SearchCommand(packageService);

        when(packageService.searchPackages(any()))
                .thenThrow(new RuntimeException("Test error"));

        int exitCode = searchCommand.search("test");

        assertThat(exitCode).isEqualTo(1);
        assertThat(errContent.toString()).contains("Test error");
        tearDown();
    }

    @Test
    void search_viaCli_parsesSearchTerm() {
        SearchCommand searchCommand = new SearchCommand(packageService);
        when(packageService.searchPackages("postgres"))
                .thenReturn(new PackageRegistry(Collections.emptyList()));

        CommandLine cmd = new CommandLine(new RootCommand());
        // Override the search command with our mock-injected version
        cmd.addSubcommand("search", searchCommand);

        int exitCode = cmd.execute("search", "postgres");

        assertThat(exitCode).isEqualTo(0);
        verify(packageService).searchPackages("postgres");
        tearDown();
    }
}
