# LPM CLI Reference

Complete reference for all lpm commands, flags, and options.

> **📝 Documentation Note:** This documentation was generated with assistance from Claude Code and may contain inaccuracies. Please verify commands and workflows in your environment and report any issues.

## Table of Contents

- [Global Flags](#global-flags)
- [Commands](#commands)
  - [search](#search)
  - [add](#add)
  - [install](#install)
  - [list](#list)
  - [remove](#remove)
  - [upgrade](#upgrade)
  - [update](#update)
  - [dedupe](#dedupe)
  - [completion](#completion)
  - [help](#help)
- [Shell Completion](#shell-completion)
- [Exit Codes](#exit-codes)

## Global Flags

These flags can be used with any command:

| Flag | Description |
|------|-------------|
| `--category <type>` | Filter by package category: `extension`, `driver`, or `pro` |
| `--version` | Show lpm version |
| `--help` | Show help information |

## Commands

### search

Search for available packages in the registry.

**Usage:**
```bash
lpm search [PACKAGE] [flags]
```

**Arguments:**
- `PACKAGE` - Package name to search for (minimum 3 characters)

**Flags:**
- `--category <type>` - Filter by category: `extension`, `driver`, or `pro`

**Examples:**
```bash
# Search all packages
lpm search

# Search for PostgreSQL-related packages  
lpm search postgres

# Search only extensions
lpm search --category extension

# Search only drivers
lpm search postgres --category driver
```

**Output Format:**
```
Package Name                    Version    Category    Description
postgresql-extension            4.2.0      extension   PostgreSQL support for Liquibase
mysql-connector-java            8.0.33     driver      MySQL JDBC driver
```

---

### add

Add packages to your project or globally.

**Usage:**
```bash
lpm add [PACKAGE]... [flags]
```

**Arguments:**
- `PACKAGE` - One or more package names, optionally with version (`package@version`)

**Flags:**
- `--global`, `-g` - Install packages globally to `$LIQUIBASE_HOME/lib`

**Version Syntax:**
- `package` - Latest compatible version
- `package@1.2.3` - Specific version
- `package@latest` - Latest version (may not be compatible)

**Examples:**
```bash
# Add latest compatible version
lpm add liquibase-postgresql

# Add specific version
lpm add liquibase-postgresql@4.2.0

# Add multiple packages
lpm add liquibase-postgresql mysql-connector-java

# Install globally
lpm add liquibase-postgresql --global

# Mix local and specific versions
lpm add package1@1.0.0 package2
```

**Behavior:**
- Creates `liquibase.json` dependency file if it doesn't exist
- Downloads packages to `./liquibase_libs/` (local) or `$LIQUIBASE_HOME/lib` (global)
- Checks compatibility with your Liquibase version
- Prevents duplicate installations

---

### install

Install all packages listed in the `liquibase.json` dependency file.

**Usage:**
```bash
lpm install
```

**Arguments:** None

**Flags:** None

**Examples:**
```bash
# Install dependencies from liquibase.json
lpm install
```

**Behavior:**
- Reads `liquibase.json` from current directory
- Installs exact versions specified in the file
- Only works for local packages (not global)
- Fails if `liquibase.json` doesn't exist

**Error Cases:**
- File not found: `liquibase.json not found in current directory`
- Global mode: `Can not install packages from liquibase.json globally`

---

### list

List installed packages in current project or globally.

**Usage:**
```bash
lpm list [flags]
```

**Aliases:** `ls`

**Arguments:** None

**Flags:**
- `--global`, `-g` - List global packages instead of local

**Examples:**
```bash
# List packages in current project
lpm list

# List global packages  
lpm list --global

# Using alias
lpm ls
```

**Output Format:**
```
/path/to/project/liquibase_libs
├── liquibase-postgresql@4.2.0
└── mysql-connector-java@8.0.33
```

---

### remove

Remove packages from your project or globally.

**Usage:**
```bash
lpm remove [PACKAGE]... [flags]
```

**Aliases:** `rm`

**Arguments:**
- `PACKAGE` - One or more package names to remove

**Flags:**
- `--global`, `-g` - Remove packages globally from `$LIQUIBASE_HOME/lib`

**Examples:**
```bash
# Remove package from current project
lpm remove liquibase-postgresql

# Remove multiple packages
lpm remove liquibase-postgresql mysql-connector-java

# Remove globally
lpm remove liquibase-postgresql --global

# Using alias
lpm rm package-name
```

**Behavior:**
- Removes package files from filesystem
- Updates `liquibase.json` (local mode only)
- Fails if package is not installed

---

### upgrade

Show outdated packages and optionally upgrade them to latest compatible versions.

**Usage:**
```bash
lpm upgrade [PACKAGE]... [flags]
```

**Aliases:** `up`

**Arguments:**
- `PACKAGE` - Specific packages to upgrade (optional, defaults to all)

**Flags:**
- `--global`, `-g` - Upgrade global packages
- `--dry-run` - Show what would be upgraded without making changes

**Examples:**
```bash
# Show what can be upgraded
lpm upgrade --dry-run

# Upgrade all packages
lpm upgrade

# Upgrade specific packages
lpm upgrade liquibase-postgresql

# Upgrade global packages
lpm upgrade --global

# Using alias
lpm up --dry-run
```

**Output Format:**
```
You have 2 outdated package(s) installed.
/path/to/project/liquibase_libs
    Package                  Installed    Available
├── liquibase-postgresql     4.1.0        4.2.0
└── mysql-connector-java     8.0.32       8.0.33
```

---

### update

Update the package manifest (packages.json) from a remote or local source.

**Usage:**
```bash
lpm update [flags]
```

**Arguments:** None

**Flags:**
- `--path <path>`, `-p <path>` - Path to new packages.json manifest (default: official GitHub repository)

**Default Source:**
```
https://raw.githubusercontent.com/liquibase/liquibase-package-manager/master/internal/app/packages.json
```

**Examples:**
```bash
# Update from official repository
lpm update

# Update from custom URL
lpm update --path https://my-company.com/packages.json

# Update from local file
lpm update --path ./custom-packages.json
```

**Behavior:**
- Downloads and validates the new manifest
- Replaces the current package registry
- Affects package search and installation
- Use with caution in production environments

---

### dedupe

Remove duplicate versions of installed packages, keeping only the latest version of each.

**Usage:**
```bash
lpm dedupe [flags]
```

**Arguments:** None

**Flags:**
- `--dry-run` - Show what would be removed without making changes

**Examples:**
```bash
# Show duplicate packages
lpm dedupe --dry-run

# Remove duplicate packages
lpm dedupe
```

**Output Format:**
```
/path/to/project/liquibase_libs
    Package                  Installed
├── liquibase-postgresql     4.2.0
└── liquibase-postgresql     4.1.0  ← will be removed
```

**Behavior:**
- Scans for packages with multiple versions installed
- Keeps the highest version of each package
- Removes older versions
- Only affects local packages

---

### completion

Generate shell completion scripts for bash, zsh, or fish.

**Usage:**
```bash
lpm completion [shell]
```

**Arguments:**
- `shell` - Target shell: `bash`, `zsh`, or `fish`

**Examples:**
```bash
# Generate bash completion
lpm completion bash

# Generate zsh completion  
lpm completion zsh

# Generate fish completion
lpm completion fish
```

See [Shell Completion](#shell-completion) section for installation instructions.

---

### help

Show help information for lpm or specific commands.

**Usage:**
```bash
lpm help [command]
```

**Arguments:**
- `command` - Show help for specific command (optional)

**Examples:**
```bash
# Show general help
lpm help

# Show help for specific command
lpm help search
lpm help add
```

## Shell Completion

### Bash

**Temporary (current session):**
```bash
source <(lpm completion bash)
```

**Permanent:**
```bash
# Linux
lpm completion bash | sudo tee /etc/bash_completion.d/lpm

# macOS (with Homebrew)
lpm completion bash > /usr/local/etc/bash_completion.d/lpm
```

### Zsh

**Temporary (current session):**
```bash
source <(lpm completion zsh)
```

**Permanent:**
```bash
lpm completion zsh > "${fpath[1]}/_lpm"
```

### Fish

**Temporary (current session):**
```bash
lpm completion fish | source
```

**Permanent:**
```bash
lpm completion fish > ~/.config/fish/completions/lpm.fish
```

**Note:** You'll need to restart your shell or source your shell configuration after permanent installation.

## Exit Codes

lpm uses standard exit codes:

| Code | Meaning |
|------|---------|
| `0` | Success |
| `1` | General error (package not found, installation failed, etc.) |

**Common Error Scenarios:**
- Package not found: `Package 'xyz' not found.`
- Version not available: `Version '1.0.0' not available.`
- Already installed: `package@version is already installed.`
- Compatibility error: `package is not compatible with liquibase version.`
- Permission error: `Unable to write to directory.`

## Environment Variables

| Variable | Description |
|----------|-------------|
| `LIQUIBASE_HOME` | Path to Liquibase installation (recommended to set) |

**Example:**
```bash
export LIQUIBASE_HOME=/usr/local/opt/liquibase/libexec
```

## File Locations

| File/Directory | Purpose |
|----------------|---------|
| `./liquibase.json` | Project dependency file |
| `./liquibase_libs/` | Local package installation directory |
| `$LIQUIBASE_HOME/lib/` | Global package installation directory |
| `$LIQUIBASE_HOME/lib/packages.json` | Package registry manifest |