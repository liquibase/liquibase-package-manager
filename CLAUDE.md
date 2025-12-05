# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is **lpm** (Liquibase Package Manager) - a Go-based CLI tool for managing Liquibase extensions, drivers, and utilities. The project allows users to search, install, update, and remove packages for Liquibase database development.

## Key Architecture

### Core Components
- **CLI Layer**: Built with Cobra framework in `cmd/lpm/` - `darwin.go` serves as the main entry point for Unix-like systems and Linux, `windows.go` for Windows
- **Application Layer**: Core business logic in `internal/app/` with embedded package metadata
- **Commands**: Individual commands (add, install, remove, search, etc.) in `internal/app/commands/`
- **Package Management**: Package structs and operations in `internal/app/packages/`
- **Liquibase Integration**: Auto-detection and version management in `internal/app/utils/Liquibase.go`
- **Populator**: Automated package discovery from Maven/GitHub repositories in `cmd/populator/`

### Package System
- **Package Registry**: `internal/app/packages.json` is embedded at build time via Go's `//go:embed` directive
- **Module Definitions**: Package sources defined in `cmd/populator/modules.go` with `Module` structs
- **Categories**: `Extension`, `Driver`, or `Pro` (enum in populator, string in JSON)
- **Artifactory Types**: `Maven{}` for Maven Central or `Github{}` for GitHub releases

### Installation Strategy
- **Global Mode** (default): Installs to `$LIQUIBASE_HOME/lib/`
- **Local Mode**: Installs to `./liquibase_libs/` in current directory
- **Auto-detection**: Locates Liquibase via `LIQUIBASE_HOME` env var or common paths

## Common Development Commands

### Build and Test
```bash
# Build for current platform (outputs to ./bin/lpm)
make build

# Run unit tests with coverage and static analysis
make test

# View test coverage in browser
make cover

# Run all end-to-end tests
make e2e

# Run specific e2e test
make test-search   # also: test-add, test-install, test-list, test-remove, test-version, test-completion, test-help

# Run a single unit test
go test -v ./internal/app/packages -run TestPackageName

# Generate cross-platform release artifacts
make release
```

### Package Metadata Updates
```bash
# Regenerate packages.json from Maven/GitHub sources
make generateExtensionPackages
```

### Adding New Packages

1. **Add module definition** in `cmd/populator/modules.go` inside `func init()`:
   ```go
   {
       name:        "package-name",
       category:    Extension, // or Driver, Pro
       url:         "https://repo1.maven.org/maven2/...",
       artifactory: Maven{}, // or Github{} with owner/repo fields
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

3. **Regenerate metadata**: Push changes and a CI workflow auto-populates version details, or run `make generateExtensionPackages` locally.

## Testing

### Unit Tests
- Located alongside source files with `_test.go` suffix in `internal/app/`
- Run with `go test -v ./internal/app/...`
- Static analysis: `staticcheck ./internal/app/...`

### End-to-End Tests
- YAML-based test definitions in `tests/endtoend/`
- Uses `vexrun` test runner (downloaded automatically by `make test-setup`)
- Test files: `add.yml`, `completion.yml`, `help.yml`, `install.yml`, `list.yml`, `remove.yml`, `search.yml`, `version.yml`

### Mock Data
- Test fixtures in `tests/mocks/`

## Build System

### Version Management
- Version stored in `VERSION` file at project root
- Copied to `internal/app/VERSION` during build (embedded via `//go:embed`)

### Multi-Platform Targets
- `darwin_amd64`, `darwin_arm64`: macOS Intel and Apple Silicon
- `linux_amd64`, `linux_arm64`, `s390x`: Linux variants
- `windows`: Windows x64

### Key Dependencies
- **Go 1.25+**: See go.mod for exact version
- **Cobra** (`github.com/spf13/cobra`): CLI framework
- **go-version** (`github.com/hashicorp/go-version`): Semantic versioning
- **gopom** (`github.com/vifraa/gopom`): Maven POM parsing
- **go-github** + **oauth2**: GitHub API for release discovery

## Environment Variables

- `LIQUIBASE_HOME`: Primary Liquibase installation directory (strongly recommended to set)