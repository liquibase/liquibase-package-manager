# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Liquibase Package Manager (lpm) is a Go-based CLI tool for managing Liquibase extensions, drivers, and utilities. It allows users to search, install, and manage external dependencies for database development with Liquibase.

## Essential Commands

### Development
```bash
# Build the project
make build

# Run tests
make test

# Run end-to-end tests  
make e2e

# Build for all platforms
make release

# Update packages.json from artifact repositories
go run cmd/populator

# View test coverage
make cover

# Run a specific test
go test -v ./internal/app/dependencies -run TestDependencyParse
```

### Self-Update Testing
```bash
# Test self-update functionality (dry run)
lpm upgrade --dry-run

# Check for available updates
lpm upgrade --check

# Test version-specific updates
lpm upgrade --version=X.Y.Z --dry-run

# Test combined package and lpm updates
lpm upgrade --all --dry-run
```

### Linting and Code Quality
```bash
# Run static analysis
staticcheck ./...

# Run go vet
go vet ./...
```

## Architecture Overview

### Core Structure
- **CLI Framework**: Cobra (github.com/spf13/cobra) for command handling
- **Commands** (`internal/app/commands/`): Each subcommand (add, install, list, remove, search, upgrade, etc.) has its own file
- **Package Registry**: `internal/app/packages.json` contains all available packages with metadata
- **Dependency Management**: `liquibase.json` tracks local project dependencies
- **Self-Update System**: Secure binary replacement with GitHub release integration

### Key Components
1. **Package Categories**:
   - `extension` - Liquibase extensions
   - `driver` - JDBC drivers  
   - `pro` - Commercial/Pro extensions

2. **Installation Paths**:
   - Global: `$LIQUIBASE_HOME/lib/`
   - Local: `./liquibase_libs/`

3. **Version Compatibility**: Packages specify required liquibaseCore version; lpm automatically selects compatible versions

### Adding New Packages
1. Add entry to `cmd/populator/modules.go`:
   ```go
   {
       name:        "package-name",
       category:    Extension,  // or Driver, Pro
       url:         "maven-url",
       artifactory: Maven{},    // or Github{}
   }
   ```

2. Add to `internal/app/packages.json`:
   ```json
   {
       "name": "package-name",
       "category": "extension",
       "versions": []
   }
   ```

3. Run populator: `go run cmd/populator`

## Important Implementation Details

### Error Handling
- Custom error types in `internal/app/errors/`
- User-friendly error messages with actionable guidance

### HTTP Operations
- Utility functions in `internal/app/utils/Http.go`
- Checksum verification for all downloads
- Secure TLS 1.2+ connections with certificate verification

### Self-Update System
- **GitHub Integration** (`internal/app/utils/GitHub.go`): Release fetching and version checking
- **Secure Downloads** (`internal/app/utils/Download.go`): Progress tracking, checksum verification, secure HTTP client
- **Atomic Updates** (`internal/app/utils/Update.go`): Cross-platform binary replacement with rollback capability
- **Platform Detection** (`internal/app/utils/Platform.go`): Architecture and OS detection for binary selection
- **Version Management** (`internal/app/utils/Version.go`): Semantic version comparison and validation

### Testing
- Unit tests alongside implementation files
- End-to-end tests in `tests/endtoend/` using YAML definitions
- Mock data in `tests/mocks/` for isolated testing

### Cross-Platform Support
- Platform-specific entry points in `cmd/lpm/`
- Builds for Darwin (AMD64/ARM64), Linux (AMD64/ARM64/s390x), Windows (AMD64)

## Key Files to Understand

### Core Application
- `internal/app/App.go` - Main application logic and initialization
- `internal/app/commands/root.go` - CLI setup and global flags
- `internal/app/commands/upgrade.go` - Self-update and package upgrade logic
- `internal/app/dependencies/Dependencies.go` - liquibase.json parsing and management
- `internal/app/packages/` - Package installation and management logic
- `internal/app/utils/Liquibase.go` - Liquibase installation detection and version checking

### Self-Update Utilities
- `internal/app/utils/GitHub.go` - GitHub API integration for release management
- `internal/app/utils/Download.go` - Secure download functionality with progress tracking
- `internal/app/utils/Update.go` - Atomic binary replacement and rollback mechanisms
- `internal/app/utils/Platform.go` - Cross-platform executable and architecture detection
- `internal/app/utils/Version.go` - Semantic version parsing and comparison

## Development Notes

- The project is marked as **experimental** and not officially supported
- Never modify git config settings
- `internal/app/packages.json` is automatically updated by CI - don't manually edit version information
- Use existing error patterns and user messaging styles when adding new features
- Follow the established command structure when adding new subcommands

### Self-Update Development Guidelines
- Always test updates with `--dry-run` flag first
- The update system creates automatic backups - test rollback scenarios
- Cross-platform compatibility is critical - test on Windows, macOS, and Linux
- GitHub releases must follow semantic versioning (vX.Y.Z format)
- Security is paramount - never disable TLS verification or checksum validation
- Update utilities should handle network failures and partial downloads gracefully