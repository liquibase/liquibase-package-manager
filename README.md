# lpm - Liquibase Package Manager

![GitHub Release Date](https://img.shields.io/github/release-date/liquibase/liquibase-package-manager?style=flat-square)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/liquibase/liquibase-package-manager?style=flat-square)
![GitHub all releases](https://img.shields.io/github/downloads/liquibase/liquibase-package-manager/total?style=flat-square)
![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/liquibase/liquibase-package-manager/nightly-e2e-tests.yml?label=E2E%20Tests&style=flat-square)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/liquibase/liquibase-package-manager?style=flat-square)

Easily manage external dependencies for Database Development. Search for, install, and uninstall liquibase drivers, extensions, and utilities.

## lpm is experimental and not officially supported

lpm is an experimental project. Issues can be reported [here](https://github.com/liquibase/liquibase-package-manager/issues), but there is no guarantee of support.

## Installation

lpm is distributed as a single binary. Install lpm by downloading, unzipping, and moving it to a directory included in your system's PATH. Releases are available [here](https://github.com/liquibase/liquibase-package-manager/releases).

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

### Self-Update Examples

```shell
# Update lpm itself to the latest version
lpm upgrade

# Check for lpm updates without installing
lpm upgrade --check

# Update lpm to a specific version
lpm upgrade --version=1.2.3

# Update both lpm and all installed packages
lpm upgrade --all

# Update specific packages only
lpm upgrade package1 package2
```

### Available Commands

* add
* completion
* dedupe
* help
* install
* list
* remove
* search
* update
* upgrade (self-update lpm or packages)

## Autocompletion

lpm can generate shell completions for multiple shells. The following shells are available:

### bash

Generate the autocompletion script for lpm for the bash shell.
To load completions in your current shell session:
`source <(lpm completion bash)`

To load completions for every new session, execute once:

* Linux:

```shell
lpm completion bash > /etc/bash_completion.d/lpm
```

* MacOS:

```shell
lpm completion bash > /usr/local/etc/bash_completion.d/lpm
```

### zsh

To load completions in your current shell session:
`source <(lpm completion zsh)`

To load completions for every new session, execute once:
`lpm completion zsh > "${fpath[1]}/_lpm"`

### fish

To load completions in your current shell session:
`lpm completion fish | source`

To load completions for every new session, execute once:
`lpm completion fish > ~/.config/fish/completions/lpm.fish`

You will need to start a new shell for this setup to take effect.

## Security Features

lpm includes several security measures to ensure safe operation:

* **TLS 1.2+ Only**: All downloads use secure HTTPS connections with minimum TLS 1.2
* **Certificate Verification**: SSL certificates are always verified (no insecure connections)
* **Checksum Verification**: Downloaded files are verified using SHA-256 checksums when available
* **Atomic Updates**: Self-updates are performed atomically with automatic rollback on failure
* **Backup & Recovery**: Automatic backup creation before updates with rollback capability
* **Secure Temp Files**: Temporary files are created with restricted permissions

## License

This project is licensed under the Apache License 2.0. See the LICENSE file for details.

This project also includes software governed by the Mozilla Public License, v. 2.0.

* `github.com/hashicorp/go-version v1.6.0` (<http://mozilla.org/MPL/2.0/>)

## Releasing LPM

For instructions on releasing LPM, see [RELEASING.md](RELEASING.md).

## Adding an Extension

https://github.com/liquibase/liquibase-package-manager/blob/master/CONTRIBUTING.md

