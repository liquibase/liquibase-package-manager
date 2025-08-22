# LPM Advanced Usage Guide

Advanced workflows and scenarios for power users and enterprise environments.

> **📝 Documentation Note:** This documentation was generated with assistance from Claude Code and may contain inaccuracies. Please verify commands and workflows in your environment and report any issues.

## Table of Contents

- [Global vs Local Package Management](#global-vs-local-package-management)
- [Custom Package Registries](#custom-package-registries)
- [Batch Operations](#batch-operations)
- [CI/CD Integration](#cicd-integration)
- [Enterprise Workflows](#enterprise-workflows)
- [Package Development](#package-development)
- [Performance Optimization](#performance-optimization)

## Global vs Local Package Management

### When to Use Global vs Local

**Local Packages (Default - Recommended)**
- ✅ Project-specific dependencies
- ✅ Version isolation between projects
- ✅ Team collaboration via `liquibase.json`
- ✅ CI/CD reproducible builds

**Global Packages**
- ✅ System-wide utilities used by all projects
- ✅ Common drivers (PostgreSQL, MySQL, etc.)
- ✅ Shared development environments
- ⚠️ Can cause version conflicts between projects

### Global Package Management

```bash
# Install packages globally
lpm add liquibase-postgresql --global
lpm add mysql-connector-java --global

# List global packages
lpm list --global

# Remove global packages
lpm remove liquibase-postgresql --global

# Upgrade global packages
lpm upgrade --global
```

**Global Installation Directory:**
```
$LIQUIBASE_HOME/lib/
├── liquibase-postgresql-4.2.0.jar
├── mysql-connector-java-8.0.33.jar
└── packages.json
```

### Mixed Environments

You can use both global and local packages simultaneously:

```bash
# Global drivers for database connectivity
lpm add mysql-connector-java --global
lpm add postgresql --global

# Project-specific extensions locally
lpm add liquibase-mongodb
lpm add liquibase-cosmosdb
```

**Resolution Order:**
1. Local packages (`./liquibase_libs/`) - highest priority
2. Global packages (`$LIQUIBASE_HOME/lib/`) - fallback

## Custom Package Registries

### Using Custom Package Sources

For organizations with private packages or custom repositories:

```bash
# Update to custom registry
lpm update --path https://internal-registry.company.com/packages.json

# Update from local file
lpm update --path /path/to/custom-packages.json

# Verify current registry
lpm search  # Shows packages from current registry
```

### Creating Custom Package Registry

Create a custom `packages.json` manifest:

```json
{
  "packages": [
    {
      "name": "company-liquibase-extension",
      "category": "extension", 
      "versions": [
        {
          "tag": "1.0.0",
          "liquibaseCore": "4.20.0",
          "url": "https://internal-registry.company.com/releases/extension-1.0.0.jar"
        }
      ]
    }
  ]
}
```

**Registry Requirements:**
- Must be valid JSON format
- Must be accessible via HTTP/HTTPS or local file path
- Should follow the same schema as official registry

### Registry Management

```bash
# Backup current registry
cp $LIQUIBASE_HOME/lib/packages.json ./backup-packages.json

# Switch to custom registry
lpm update --path https://internal.company.com/packages.json

# Restore official registry  
lpm update  # Uses default official source
```

## Batch Operations

### Managing Multiple Packages

```bash
# Add multiple packages with specific versions
lpm add \
  liquibase-postgresql@4.2.0 \
  mysql-connector-java@8.0.33 \
  liquibase-oracle@4.1.0

# Remove multiple packages
lpm remove \
  liquibase-postgresql \
  mysql-connector-java \
  liquibase-oracle

# Upgrade specific packages only
lpm upgrade liquibase-postgresql mysql-connector-java
```

### Dependency File Templates

Create project templates with common dependencies:

**backend-api-template/liquibase.json:**
```json
{
  "dependencies": [
    {"postgresql": "42.5.0"},
    {"liquibase-postgresql": "4.2.0"},
    {"liquibase-hibernate": "4.1.0"}
  ]
}
```

**Usage:**
```bash
# Start new project from template
cp backend-api-template/liquibase.json ./
lpm install
```

### Scripted Package Management

```bash
#!/bin/bash
# setup-project.sh - Automated project setup

echo "Setting up Liquibase project..."

# Install core dependencies
lpm add postgresql liquibase-postgresql

# Install development utilities
lpm add liquibase-hibernate liquibase-test-harness

# Verify installation
lpm list

echo "Project setup complete!"
```

## CI/CD Integration

### Docker Integration

**Dockerfile example:**
```dockerfile
FROM liquibase/liquibase:4.20.0

# Install lpm
RUN curl -L https://github.com/liquibase/liquibase-package-manager/releases/latest/download/lpm-linux.zip -o lpm.zip && \
    unzip lpm.zip && \
    mv lpm /usr/local/bin/ && \
    chmod +x /usr/local/bin/lpm

# Copy dependency file
COPY liquibase.json ./

# Install project dependencies
RUN lpm install

# Your application files
COPY changelog/ ./changelog/
```

### GitHub Actions

```yaml
name: Database Migration

on: [push, pull_request]

jobs:
  migrate:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Setup Liquibase
      uses: liquibase-github-actions/install@v1
      with:
        version: 4.20.0
        
    - name: Install LPM
      run: |
        curl -L https://github.com/liquibase/liquibase-package-manager/releases/latest/download/lpm-linux.zip -o lpm.zip
        unzip lpm.zip && sudo mv lpm /usr/local/bin/
        
    - name: Install Dependencies
      run: lpm install
      
    - name: Run Migration
      run: liquibase update
      env:
        LIQUIBASE_COMMAND_URL: ${{ secrets.DB_URL }}
        LIQUIBASE_COMMAND_USERNAME: ${{ secrets.DB_USERNAME }}
        LIQUIBASE_COMMAND_PASSWORD: ${{ secrets.DB_PASSWORD }}
```

### Jenkins Pipeline

```groovy
pipeline {
    agent any
    
    stages {
        stage('Setup') {
            steps {
                sh 'curl -L https://github.com/liquibase/liquibase-package-manager/releases/latest/download/lpm-linux.zip -o lpm.zip'
                sh 'unzip lpm.zip && sudo mv lpm /usr/local/bin/'
            }
        }
        
        stage('Install Dependencies') {
            steps {
                sh 'lpm install'
            }
        }
        
        stage('Database Migration') {
            steps {
                sh 'liquibase update'
            }
        }
    }
}
```

## Enterprise Workflows

### Team Collaboration Patterns

**1. Centralized Dependency Management:**
```bash
# Senior developer manages dependencies
lpm add new-extension@1.2.0
git add liquibase.json
git commit -m "Add new extension v1.2.0"

# Team members sync
git pull
lpm install
```

**2. Environment-Specific Dependencies:**
```bash
# development.liquibase.json
{
  "dependencies": [
    {"liquibase-test-harness": "1.0.0"},
    {"postgresql": "42.5.0"}
  ]
}

# production.liquibase.json  
{
  "dependencies": [
    {"postgresql": "42.5.0"}
  ]
}

# Usage
cp development.liquibase.json liquibase.json && lpm install  # dev
cp production.liquibase.json liquibase.json && lpm install   # prod
```

### Version Pinning Strategy

```bash
# Always pin to specific versions in production
lpm add liquibase-postgresql@4.2.0  # ✅ Good
lpm add liquibase-postgresql         # ⚠️ Risky - version may change

# Use latest for development, pin for production
# Development
lpm add liquibase-postgresql

# Production (pin the version that worked in dev)
lpm add liquibase-postgresql@4.2.0
```

### Security Considerations

**Package Source Verification:**
```bash
# Verify package registry source
lpm update --dry-run  # Shows what would be updated

# Use internal registry for sensitive environments
lpm update --path https://secure-internal-registry.company.com/packages.json
```

**Dependency Auditing:**
```bash
# Regular dependency review
lpm list > installed-packages.txt
lpm upgrade --dry-run > outdated-packages.txt

# Review and approve updates
# Then apply: lpm upgrade
```

## Package Development

### Testing Custom Packages

```bash
# Install package from local file for testing
lpm update --path ./test-packages.json
lpm add my-custom-extension

# Test with existing project
lpm install
liquibase validate
```

### Package Registry Development

Create and test custom package definitions:

```json
{
  "packages": [
    {
      "name": "my-extension",
      "category": "extension",
      "description": "Custom database extension",
      "versions": [
        {
          "tag": "1.0.0-beta",
          "liquibaseCore": "4.20.0", 
          "url": "file:///path/to/local/extension.jar"
        }
      ]
    }
  ]
}
```

## Performance Optimization

### Caching Strategies

**Global Package Caching:**
- Global packages are shared across all projects
- Reduces download time and disk usage
- Consider for frequently used drivers

```bash
# Install commonly used drivers globally once
lpm add postgresql mysql-connector-java oracle-jdbc --global
```

**Local Package Benefits:**
- Faster project setup (no global conflicts)
- Isolated environments
- Reproducible builds

### Network Optimization

**Batch Downloads:**
```bash
# Download all dependencies at once
lpm install  # More efficient than individual lpm add commands
```

**Offline Mode:**
```bash
# Pre-download packages for offline use
lpm add package1 package2 package3
# Later: packages are cached and available offline
```

### Storage Management

**Clean Up Old Versions:**
```bash
# Remove duplicate package versions
lpm dedupe

# Manual cleanup of old global packages
ls $LIQUIBASE_HOME/lib/*.jar
# Remove old versions manually if needed
```

**Monitor Disk Usage:**
```bash
# Check local package size
du -sh ./liquibase_libs/

# Check global package size  
du -sh $LIQUIBASE_HOME/lib/
```