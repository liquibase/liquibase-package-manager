# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is **lpm** (Liquibase Package Manager) - a Go-based CLI tool for managing Liquibase extensions, drivers, and utilities. The project allows users to search, install, update, and remove packages for Liquibase database development.

## Key Architecture

### Core Components
- **CLI Layer**: Built with Cobra framework in `cmd/lpm/` with platform-specific entry points
- **Application Layer**: Core business logic in `internal/app/` with embedded package metadata
- **Package Management**: Handles package discovery, installation, and lifecycle in `internal/app/packages/`
- **Liquibase Integration**: Auto-detection and version management in `internal/app/utils/Liquibase.go`
- **Populator**: Automated package discovery from Maven/GitHub repositories in `cmd/populator/`

### Package System
- **Package Registry**: Centralized metadata in `internal/app/packages.json` (embedded at build time)
- **Module Definitions**: Package sources defined in `cmd/populator/modules.go`
- **Categories**: Packages organized as `extension`, `driver`, or `pro`
- **Artifactory Support**: Maven Central and GitHub release integration

### Installation Strategy
- **Global Mode**: Installs to `$LIQUIBASE_HOME/lib/` (default behavior)
- **Local Mode**: Installs to `./liquibase_libs/` in current directory
- **Auto-detection**: Automatically locates Liquibase installations via environment variables or common paths

## Common Development Commands

### Build and Test
```bash
# Build for current platform
make build

# Run unit tests with coverage
make test

# View test coverage report
make cover

# Run end-to-end tests
make e2e

# Generate cross-platform releases
make release
```

### Development Workflow
```bash
# Update package metadata from sources
make generateExtensionPackages

# Build single binary for local testing
make build

# Test specific command end-to-end
make test-search  # or test-add, test-install, etc.
```

### Adding New Packages

1. **Add module definition** in `cmd/populator/modules.go`:
   ```go
   {
       name:        "package-name",
       category:    Extension, // or Driver, Pro
       url:         "https://repo1.maven.org/maven2/...",
       artifactory: Maven{}, // or Github{}
   }
   ```

2. **Add package stub** in `internal/app/packages.json`:
   ```json
   {
       "name": "package-name",
       "category": "extension",
       "versions": []
   }
   ```

3. **Generate updated metadata**:
   ```bash
   make generateExtensionPackages
   ```

## Testing

### Unit Tests
- Located in `internal/app/` with `_test.go` suffix
- Use `go test` with coverage reporting
- Static analysis with `staticcheck`

### End-to-End Tests
- YAML-based test definitions in `tests/endtoend/`
- Execute with `vexrun` test runner
- Cover all major CLI commands and workflows

### Mock Data
- Test fixtures in `tests/mocks/`
- Includes sample Liquibase installations and package files

## Build System

### Multi-Platform Support
- Cross-compilation for Darwin (Intel/ARM), Linux (x86_64/ARM64), Windows, and s390x
- Automated artifact packaging with version-specific naming
- Release automation via GitHub Actions

### Dependencies
- **Go 1.23.8**: Primary language runtime
- **Cobra**: CLI framework and command structure
- **go-version**: Semantic version handling
- **oauth2**: GitHub API authentication for package discovery

## Package Categories

- **Extensions**: Community and commercial Liquibase extensions for specific databases
- **Drivers**: JDBC drivers for database connectivity
- **Pro**: Commercial Liquibase Pro features and integrations

## Environment Variables

- `LIQUIBASE_HOME`: Primary installation directory (recommended to set)
- Used for auto-detection of Liquibase installations and lib directories