# LPM - Liquibase Package Manager
Easily manage external dependencies for Database Development. Search for, install, and uninstall liquibase drivers, extensions, and utilities.

## Installation
lpm is distributed as a single binary. Install lpm by downloading, unzipping, and moving it to a directory included in your system's PATH. Releases are available [here](https://github.com/mcred/liquibase-package-manager/releases).

## Setup
lpm will make a best effort to locate the location of the liquibase lib directory. It is recommended to set the LIQUIBASE_HOME environment variable.

Examples:
```shell
export LIQUIBASE_HOME=/usr/local/opt/liquibase/libexec
echo  'export LIQUIBASE_HOME=/usr/local/opt/liquibase/libexec' >> ~/.bashrc 
```

## Usage
```shell
lpm <command>
```
### Available Commands
* help        Help about any command
* install     Install Packages
* ls          List Installed Packages
* search      Search for Packages
* uninstall   Uninstall Package
* update      Updates the Package Manifest
