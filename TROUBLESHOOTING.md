# LPM Troubleshooting Guide

Solutions for common issues and errors when using Liquibase Package Manager.

> **📝 Documentation Note:** This documentation was generated with assistance from Claude Code and may contain inaccuracies. Please verify commands and workflows in your environment and report any issues.

## Table of Contents

- [Installation Issues](#installation-issues)
- [Package Management Errors](#package-management-errors)
- [Version Compatibility Issues](#version-compatibility-issues)
- [Permission Problems](#permission-problems)
- [Network and Download Issues](#network-and-download-issues)
- [Environment Configuration](#environment-configuration)
- [Diagnostic Commands](#diagnostic-commands)

## Installation Issues

### lpm command not found

**Error:**
```bash
$ lpm --version
bash: lpm: command not found
```

**Solutions:**

1. **Verify installation:**
   ```bash
   # Check if lpm is in your PATH
   which lpm
   
   # If not found, check where you placed it
   ls -la /usr/local/bin/lpm
   ```

2. **Add to PATH:**
   ```bash
   # Temporary (current session)
   export PATH=$PATH:/path/to/lpm/directory
   
   # Permanent (add to ~/.bashrc or ~/.zshrc)
   echo 'export PATH=$PATH:/usr/local/bin' >> ~/.bashrc
   source ~/.bashrc
   ```

3. **Reinstall lpm:**
   ```bash
   # Download latest release
   curl -L https://github.com/liquibase/liquibase-package-manager/releases/latest/download/lpm-linux.zip -o lpm.zip
   unzip lpm.zip
   sudo mv lpm /usr/local/bin/
   chmod +x /usr/local/bin/lpm
   ```

### Permission denied when installing

**Error:**
```bash
$ sudo mv lpm /usr/local/bin/
mv: cannot move 'lpm' to '/usr/local/bin/lpm': Permission denied
```

**Solutions:**

1. **Use proper sudo permissions:**
   ```bash
   sudo mv lpm /usr/local/bin/
   sudo chmod +x /usr/local/bin/lpm
   ```

2. **Install to user directory:**
   ```bash
   mkdir -p ~/bin
   mv lpm ~/bin/
   chmod +x ~/bin/lpm
   echo 'export PATH=$PATH:$HOME/bin' >> ~/.bashrc
   ```

## Package Management Errors

### Package not found

**Error:**
```bash
$ lpm add unknown-package
Package 'unknown-package' not found.
```

**Solutions:**

1. **Search for correct package name:**
   ```bash
   lpm search postgres  # Find PostgreSQL-related packages
   lpm search mysql     # Find MySQL-related packages
   ```

2. **Update package registry:**
   ```bash
   lpm update
   lpm search unknown-package
   ```

3. **Check category filters:**
   ```bash
   lpm search --category extension
   lpm search --category driver
   lpm search --category pro
   ```

### Version not available

**Error:**
```bash
$ lpm add liquibase-postgresql@9.9.9
Version '9.9.9' not available.
```

**Solutions:**

1. **Check available versions:**
   ```bash
   lpm search liquibase-postgresql
   # Shows all available versions
   ```

2. **Use latest compatible version:**
   ```bash
   lpm add liquibase-postgresql  # No version specified
   ```

3. **Check version format:**
   ```bash
   # Correct format
   lpm add package@1.2.3
   
   # Incorrect formats
   lpm add package:1.2.3  # ❌
   lpm add package-1.2.3  # ❌
   ```

### Package already installed

**Error:**
```bash
$ lpm add liquibase-postgresql
liquibase-postgresql@4.2.0 is already installed.
liquibase-postgresql can not be installed.
Consider running `lpm upgrade`.
```

**Solutions:**

1. **Upgrade to newer version:**
   ```bash
   lpm upgrade liquibase-postgresql
   ```

2. **Force reinstall by removing first:**
   ```bash
   lpm remove liquibase-postgresql
   lpm add liquibase-postgresql@4.2.0
   ```

3. **Install different version:**
   ```bash
   lpm add liquibase-postgresql@4.1.0
   ```

### Can't install from liquibase.json globally

**Error:**
```bash
$ lpm install --global
Can not install packages from liquibase.json globally
```

**Solution:**
```bash
# Use install without --global flag
lpm install

# Or add packages globally one by one
lpm add package-name --global
```

## Version Compatibility Issues

### Package incompatible with Liquibase version

**Error:**
```bash
$ lpm add some-extension@5.0.0
some-extension@5.0.0 is not compatible with liquibase v4.18.0. Please consider updating liquibase.
```

**Solutions:**

1. **Use latest compatible version:**
   ```bash
   lpm add some-extension  # lpm will find compatible version
   ```

2. **Update Liquibase:**
   ```bash
   # Check current Liquibase version
   liquibase --version
   
   # Update Liquibase (method depends on installation)
   # Homebrew: brew upgrade liquibase
   # Manual: download latest from liquibase.org
   ```

3. **Check compatibility matrix:**
   ```bash
   lpm search some-extension
   # Shows compatible Liquibase versions for each package version
   ```

### Multiple package versions causing conflicts

**Error:**
```bash
$ lpm list
Multiple versions of 'liquibase-postgresql' found:
├── liquibase-postgresql@4.2.0
└── liquibase-postgresql@4.1.0
```

**Solutions:**

1. **Remove duplicates:**
   ```bash
   lpm dedupe
   ```

2. **Manual cleanup:**
   ```bash
   lpm remove liquibase-postgresql
   lpm add liquibase-postgresql@4.2.0
   ```

3. **Preview before fixing:**
   ```bash
   lpm dedupe --dry-run  # See what would be removed
   lpm dedupe           # Actually remove duplicates
   ```

## Permission Problems

### Unable to write to directory

**Error:**
```bash
$ lpm add liquibase-postgresql --global
Unable to remove old-package.jar from classpath.
```

**Solutions:**

1. **Check LIQUIBASE_HOME permissions:**
   ```bash
   ls -la $LIQUIBASE_HOME/lib/
   # Should be writable by your user
   ```

2. **Fix permissions:**
   ```bash
   sudo chown -R $USER:$USER $LIQUIBASE_HOME/lib/
   chmod -R 755 $LIQUIBASE_HOME/lib/
   ```

3. **Use local installation instead:**
   ```bash
   # Install locally instead of globally
   lpm add liquibase-postgresql  # No --global flag
   ```

### Permission denied creating liquibase_libs

**Error:**
```bash
$ lpm add package-name
mkdir: cannot create directory './liquibase_libs': Permission denied
```

**Solutions:**

1. **Check directory permissions:**
   ```bash
   ls -la .
   # Current directory should be writable
   ```

2. **Change to writable directory:**
   ```bash
   cd ~/my-project  # Move to project directory
   lpm add package-name
   ```

3. **Fix directory permissions:**
   ```bash
   sudo chown -R $USER:$USER .
   ```

## Network and Download Issues

### Failed to download package

**Error:**
```bash
$ lpm add some-package
Error downloading package from https://example.com/package.jar
```

**Solutions:**

1. **Check network connectivity:**
   ```bash
   curl -I https://repo1.maven.org/maven2/
   ping github.com
   ```

2. **Check proxy settings:**
   ```bash
   echo $HTTP_PROXY
   echo $HTTPS_PROXY
   
   # Set proxy if needed
   export HTTP_PROXY=http://proxy.company.com:8080
   export HTTPS_PROXY=https://proxy.company.com:8080
   ```

3. **Retry installation:**
   ```bash
   lpm add some-package
   ```

### Package registry update failed

**Error:**
```bash
$ lpm update
Unable to download package manifest from remote URL
```

**Solutions:**

1. **Check network connectivity:**
   ```bash
   curl -I https://raw.githubusercontent.com/liquibase/liquibase-package-manager/master/internal/app/packages.json
   ```

2. **Use cached registry:**
   ```bash
   # lpm will use existing registry if update fails
   lpm search  # Should still work with cached data
   ```

3. **Manual registry update:**
   ```bash
   # Download manually and update from file
   curl -o packages.json https://raw.githubusercontent.com/liquibase/liquibase-package-manager/master/internal/app/packages.json
   lpm update --path ./packages.json
   ```

## Environment Configuration

### LIQUIBASE_HOME not set or incorrect

**Error:**
```bash
$ lpm add package-name --global
Error: Unable to locate Liquibase installation
```

**Solutions:**

1. **Find Liquibase installation:**
   ```bash
   which liquibase
   # Usually points to: /usr/local/bin/liquibase (symlink)
   
   # Find actual installation
   readlink $(which liquibase)
   # Example: /usr/local/opt/liquibase/bin/liquibase
   ```

2. **Set LIQUIBASE_HOME:**
   ```bash
   # For the example above
   export LIQUIBASE_HOME=/usr/local/opt/liquibase/libexec
   
   # Make permanent
   echo 'export LIQUIBASE_HOME=/usr/local/opt/liquibase/libexec' >> ~/.bashrc
   ```

3. **Common installation paths:**
   ```bash
   # Homebrew (macOS)
   export LIQUIBASE_HOME=/usr/local/opt/liquibase/libexec
   
   # Manual installation
   export LIQUIBASE_HOME=/opt/liquibase
   
   # Windows
   export LIQUIBASE_HOME="C:\\liquibase"
   ```

### liquibase.json format errors

**Error:**
```bash
$ lpm install
Error parsing liquibase.json: invalid JSON format
```

**Solutions:**

1. **Validate JSON format:**
   ```bash
   cat liquibase.json | python -m json.tool
   # Or use online JSON validator
   ```

2. **Fix common JSON errors:**
   ```json
   {
     "dependencies": [
       {"package-name": "1.0.0"},     ✅ Correct
       {"package-name": "1.0.0",}     ❌ Trailing comma
       {package-name: "1.0.0"}        ❌ Missing quotes
       {"package-name": 1.0.0}        ❌ Version should be string
     ]
   }
   ```

3. **Recreate liquibase.json:**
   ```bash
   rm liquibase.json
   lpm add package1 package2  # Recreates file with correct format
   ```

## Diagnostic Commands

### Environment Information

```bash
# Check lpm version and installation
lpm --version
which lpm

# Check Liquibase configuration
liquibase --version
echo $LIQUIBASE_HOME

# Check Java environment
java -version
echo $JAVA_HOME
```

### Package Diagnostics

```bash
# List all packages (installed and available)
lpm search

# Show installed packages
lpm list
lpm list --global

# Check for outdated packages
lpm upgrade --dry-run

# Verify package registry
ls -la $LIQUIBASE_HOME/lib/packages.json
```

### File System Diagnostics

```bash
# Check local package directory
ls -la ./liquibase_libs/
du -sh ./liquibase_libs/

# Check global package directory
ls -la $LIQUIBASE_HOME/lib/
du -sh $LIQUIBASE_HOME/lib/

# Check dependency file
cat liquibase.json
```

### Network Diagnostics

```bash
# Test package registry connectivity
curl -I https://raw.githubusercontent.com/liquibase/liquibase-package-manager/master/internal/app/packages.json

# Test Maven Central connectivity
curl -I https://repo1.maven.org/maven2/

# Check proxy settings
echo $HTTP_PROXY $HTTPS_PROXY
```

## Getting Help

If you're still experiencing issues:

1. **Search existing issues:** [GitHub Issues](https://github.com/liquibase/liquibase-package-manager/issues)

2. **Create a new issue with:**
   - lpm version (`lpm --version`)
   - Operating system and version
   - Complete error message
   - Steps to reproduce
   - Environment details (LIQUIBASE_HOME, Java version, etc.)

3. **Use verbose output if available:**
   ```bash
   # Some commands may support verbose flags
   lpm command --verbose  # Check if available
   ```

Remember: lpm is experimental software. While these troubleshooting steps address common issues, some problems may require updates to lpm itself.