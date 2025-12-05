package org.liquibase.lpm.commands;

import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import picocli.CommandLine;

import java.io.ByteArrayOutputStream;
import java.io.PrintStream;

import static org.assertj.core.api.Assertions.assertThat;

class RootCommandTest {

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
    void rootCommand_noArgs_printsHelp() {
        RootCommand rootCommand = new RootCommand();

        int exitCode = rootCommand.call();

        assertThat(exitCode).isEqualTo(0);
        assertThat(outContent.toString()).contains("Liquibase Package Manager");
        assertThat(outContent.toString()).contains("--help");
    }

    @Test
    void rootCommand_helpFlag_printsDetailedHelp() {
        CommandLine cmd = new CommandLine(new RootCommand());

        int exitCode = cmd.execute("--help");

        assertThat(exitCode).isEqualTo(0);
        String output = outContent.toString();
        assertThat(output).contains("lpm");
        assertThat(output).contains("add");
        assertThat(output).contains("install");
        assertThat(output).contains("search");
        assertThat(output).contains("list");
        assertThat(output).contains("remove");
        assertThat(output).contains("update");
        assertThat(output).contains("upgrade");
        assertThat(output).contains("dedupe");
    }

    @Test
    void rootCommand_versionFlag_printsVersion() {
        CommandLine cmd = new CommandLine(new RootCommand());

        int exitCode = cmd.execute("--version");

        assertThat(exitCode).isEqualTo(0);
        // Should contain some version info (either from file or unknown)
        assertThat(outContent.toString()).isNotEmpty();
    }

    @Test
    void rootCommand_categoryOption_isAvailable() {
        RootCommand rootCommand = new RootCommand();
        CommandLine cmd = new CommandLine(rootCommand);

        // Parse with --category option
        cmd.parseArgs("--category", "extension");

        assertThat(rootCommand.getCategory()).isEqualTo("extension");
    }

    @Test
    void rootCommand_categoryOption_differentValues() {
        RootCommand rootCommand = new RootCommand();
        CommandLine cmd = new CommandLine(rootCommand);

        // Test driver category
        cmd.parseArgs("--category", "driver");
        assertThat(rootCommand.getCategory()).isEqualTo("driver");

        // Reset and test pro category
        rootCommand = new RootCommand();
        cmd = new CommandLine(rootCommand);
        cmd.parseArgs("--category", "pro");
        assertThat(rootCommand.getCategory()).isEqualTo("pro");
    }

    @Test
    void rootCommand_hasAllSubcommands() {
        CommandLine cmd = new CommandLine(new RootCommand());

        assertThat(cmd.getSubcommands()).containsKeys(
                "add", "install", "search", "list", "remove", "update", "upgrade", "dedupe"
        );
    }

    @Test
    void rootCommand_listHasAlias() {
        CommandLine cmd = new CommandLine(new RootCommand());

        // ls should be an alias for list
        assertThat(cmd.getSubcommands()).containsKey("ls");
    }

    @Test
    void rootCommand_removeHasAlias() {
        CommandLine cmd = new CommandLine(new RootCommand());

        // rm should be an alias for remove
        assertThat(cmd.getSubcommands()).containsKey("rm");
    }

    @Test
    void rootCommand_upgradeHasAlias() {
        CommandLine cmd = new CommandLine(new RootCommand());

        // up should be an alias for upgrade
        assertThat(cmd.getSubcommands()).containsKey("up");
    }

    @Test
    void rootCommand_unknownCommand_printsError() {
        CommandLine cmd = new CommandLine(new RootCommand());

        int exitCode = cmd.execute("unknowncommand");

        assertThat(exitCode).isNotEqualTo(0);
        assertThat(errContent.toString()).contains("unknowncommand");
    }
}
