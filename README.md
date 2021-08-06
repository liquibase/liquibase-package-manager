# lpm - Liquibase Package Manager
![GitHub Release Date](https://img.shields.io/github/release-date/liquibase/liquibase-package-manager?style=flat-square)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/liquibase/liquibase-package-manager?style=flat-square)
![GitHub all releases](https://img.shields.io/github/downloads/liquibase/liquibase-package-manager/total?style=flat-square)
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
### Available Commands
* add
* completion
* help
* install
* list
* remove
* search
* update

## Autocompletion
lpm can generate shell completions for multiple shells. The following shells are available:

### bash
Generate the autocompletion script for lpm for the bash shell.
To load completions in your current shell session:
`source <(lpm completion bash)`

To load completions for every new session, execute once:
- Linux:
```shell
lpm completion bash > /etc/bash_completion.d/lpm
```
- MacOS:
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
