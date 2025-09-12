# lpm - Liquibase Package Manager

![GitHub Release Date](https://img.shields.io/github/release-date/liquibase/liquibase-package-manager?style=flat-square)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/liquibase/liquibase-package-manager?style=flat-square)
![GitHub all releases](https://img.shields.io/github/downloads/liquibase/liquibase-package-manager/total?style=flat-square)
![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/liquibase/liquibase-package-manager/nightly-e2e-tests.yml?label=E2E%20Tests&style=flat-square)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/liquibase/liquibase-package-manager?style=flat-square)

**lpm** makes it easy to add Liquibase extensions, database drivers, and utilities to your projects. Search, install, and manage package versions just like you would with npm or pip.

## ⚠️ Experimental Project

lpm is an experimental project. Issues can be reported [here](https://github.com/liquibase/liquibase-package-manager/issues), but there is no guarantee of support.

> **📝 Documentation Note:** This documentation was generated with assistance from Claude Code and may contain inaccuracies. Please verify commands and workflows in your environment and report any issues.

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Package Types](#package-types)
- [Essential Commands](#essential-commands)
- [Project Integration](#project-integration)
- [Version Management](#version-management)
- [Documentation](#documentation)

## Installation

lpm is distributed as a single binary. Follow these steps to install:

1. **Download the latest release** from [GitHub releases](https://github.com/liquibase/liquibase-package-manager/releases)
2. **Extract the binary** from the zip file
3. **Move to your PATH** (e.g., `/usr/local/bin/` on macOS/Linux)
4. **Verify installation**: `lpm --version`

### Setting up LIQUIBASE_HOME

lpm will automatically detect your Liquibase installation, but it's recommended to set the `LIQUIBASE_HOME` environment variable:

```bash
# Find your Liquibase installation
which liquibase

# Set LIQUIBASE_HOME (example paths)
export LIQUIBASE_HOME=/usr/local/opt/liquibase/libexec  # Homebrew on macOS
export LIQUIBASE_HOME=/opt/liquibase                    # Linux
export LIQUIBASE_HOME="C:\liquibase"                    # Windows

# Make it permanent
echo 'export LIQUIBASE_HOME=/usr/local/opt/liquibase/libexec' >> ~/.bashrc
```

## Quick Start

Here's how to add your first Liquibase extension to a project:

```bash
# 1. Search for available extensions
lpm search postgres

# 2. Add a PostgreSQL extension to your project
lpm add liquibase-postgresql

# 3. See what's installed
lpm list

# 4. Check your project files
ls -la ./liquibase_libs/    # Your packages
cat liquibase.json          # Your dependency file
```

**What just happened?**
- lpm downloaded the extension to `./liquibase_libs/`  
- lpm created a `liquibase.json` file to track your dependencies
- Your Liquibase installation can now use the PostgreSQL extension

## Package Types

lpm manages three types of packages:

- **🔧 Extensions** - Add functionality to Liquibase (database-specific features, custom change types)
- **🔌 Drivers** - JDBC drivers for database connectivity (PostgreSQL, MySQL, Oracle, etc.)
- **💼 Pro** - Commercial Liquibase Pro features (requires license)

## Essential Commands

### Search for packages
```bash
# Search all packages
lpm search

# Search for specific packages
lpm search postgres
lpm search mysql

# Search by category
lpm search --category extension
lpm search --category driver
```

### Add packages to your project
```bash
# Add latest compatible version
lpm add liquibase-postgresql

# Add specific version  
lpm add liquibase-postgresql@4.2.0

# Add multiple packages
lpm add liquibase-postgresql mysql-connector-java

# Add globally (to LIQUIBASE_HOME)
lpm add liquibase-postgresql --global
```

### List installed packages
```bash
# List packages in current project
lpm list

# List global packages
lpm list --global
```

### Remove packages
```bash
# Remove from current project
lpm remove liquibase-postgresql

# Remove multiple packages
lpm remove liquibase-postgresql mysql-connector-java

# Remove globally
lpm remove liquibase-postgresql --global
```

### Install from dependency file
```bash
# Install all packages listed in liquibase.json
lpm install
```

### Other useful commands
```bash
# Show outdated packages and upgrade them
lpm upgrade

# Remove duplicate package versions
lpm dedupe

# Update package manifest
lpm update

# Get help
lpm help
lpm help search
```

## Project Integration

### Dependency Management

lpm creates a `liquibase.json` file to track your project's package dependencies:

```json
{
  "dependencies": [
    {"liquibase-postgresql": "4.2.0"},
    {"mysql-connector-java": "8.0.33"}
  ]
}
```

**Best Practices:**
- ✅ **Commit `liquibase.json`** to version control (like `package.json`)
- ✅ **Add `liquibase_libs/` to `.gitignore`** (like `node_modules/`)
- ✅ **Use `lpm install`** in CI/CD to restore dependencies
- ✅ **Pin versions** for reproducible builds

### Team Workflow

```bash
# Developer 1: Add packages and commit dependency file
lpm add liquibase-postgresql@4.2.0
git add liquibase.json
git commit -m "Add PostgreSQL extension"

# Developer 2: Install dependencies from file
git pull
lpm install  # Installs packages listed in liquibase.json
```

### .gitignore Example

```
# Liquibase packages (similar to node_modules)
liquibase_libs/

# Keep the dependency file
!liquibase.json
```

## Version Management

### Specifying Versions

```bash
# Latest compatible version (recommended)
lpm add liquibase-postgresql

# Specific version
lpm add liquibase-postgresql@4.2.0

# Latest version (might not be compatible)
lpm add liquibase-postgresql@latest
```

### Compatibility Checking

lpm automatically checks compatibility between:
- **Package versions** and your **Liquibase version**
- **Dependencies** to prevent conflicts

```bash
# lpm will warn you if versions are incompatible
lpm add some-extension@5.0.0
# Error: some-extension@5.0.0 requires Liquibase 4.20.0+, but you have 4.18.0
```

### Upgrading Packages

```bash
# See what can be upgraded
lpm upgrade --dry-run

# Upgrade all packages to latest compatible versions
lpm upgrade

# Upgrade specific packages only
lpm upgrade liquibase-postgresql
```

## Documentation

- **[CLI Reference](CLI_REFERENCE.md)** - Complete command reference with all flags and options
- **[Advanced Usage](ADVANCED_USAGE.md)** - Global packages, manifest updates, and complex scenarios  
- **[Troubleshooting](TROUBLESHOOTING.md)** - Common errors and solutions

## Shell Completion

Enable tab completion for your shell:

```bash
# Bash
source <(lpm completion bash)

# Zsh  
source <(lpm completion zsh)

# Fish
lpm completion fish | source
```

See [CLI Reference](CLI_REFERENCE.md#shell-completion) for permanent setup instructions.

## Contributing

- **[Adding Extensions](CONTRIBUTING.md)** - How to add new packages to the registry
- **[Releasing LPM](RELEASING.md)** - Release process for maintainers

## License

This project is licensed under the Apache License 2.0. See the LICENSE file for details.

This project also includes software governed by the Mozilla Public License, v. 2.0:
* `github.com/hashicorp/go-version v1.6.0` ([MPL 2.0](http://mozilla.org/MPL/2.0/))

