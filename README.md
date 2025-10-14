# lpm - Liquibase Package Manager

![GitHub Release Date](https://img.shields.io/github/release-date/liquibase/liquibase-package-manager?style=flat-square)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/liquibase/liquibase-package-manager?style=flat-square)
![GitHub all releases](https://img.shields.io/github/downloads/liquibase/liquibase-package-manager/total?style=flat-square)
![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/liquibase/liquibase-package-manager/nightly-e2e-tests.yml?label=E2E%20Tests&style=flat-square)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/liquibase/liquibase-package-manager?style=flat-square)

A Community-supported capability to easily manage external dependencies for Liquibase Community 5.0+, enabling searching, listing, installing, and uninstalling Liquibase drivers, extensions, and utilities.

## lpm is experimental and community-supported

lpm is an experimental project. Issues can be reported [here](https://github.com/liquibase/liquibase-package-manager/issues), but there is no guarantee of support.

## Integrated into Liquibase Community 5.0+

lpm is integrated and available for use directly from within the Community CLI experience via the `liquibase lpm` command. Learn more with `liquibase lpm --help`


## Alternately, install the Binary for pre-5.0 versions

lpm is distributed as a single binary. Install lpm by downloading, unzipping, and moving it to a directory included in your system's PATH. Releases are available [here](https://github.com/liquibase/liquibase-package-manager/releases).

## Setup

lpm will make a best effort to locate the location of the `liquibase_lib` directory to store lpm-managed dependencies. It is recommended to set the LIQUIBASE_HOME environment variable.

Examples:

```shell
export LIQUIBASE_HOME=/usr/local/opt/liquibase/libexec
echo  'export LIQUIBASE_HOME=/usr/local/opt/liquibase/libexec' >> ~/.bashrc 
```

## Usage from *within* liquibase 5.0+

```shell
liquibase lpm <command>
```

### Available Commands

Getting Started for Liquibase Community 5.0+ only
* `liquibase lpm` or `liquibase lpm --download=true`
* [Optional] `liquibase lpm --lpmHome=/the/location/to/store/lpm/command <command>`

Note: If run as an independent binary outside liquibase, or with pre-5.0 versions, drop the leading `liquibase` for the following commands
* `liquibase lpm add`
* `liquibase lpm completion`
* `liquibase lpm dedupe`
* `liquibase lpm help`
* `liquibase lpm install`
* `liquibase lpm list`
* `liquibase lpm remove`
* `liquibase lpm search`
* `liquibase lpm update`
* `liquibase lpm upgrade`

### Important Clarifications
- `liquibase lpm add` = add packages to the `liquibase.json` file and to this Liquibase installation

- `liquibase lpm update` = updates all packages listed in the `liquibase.json` file
- `liquibase lpm update <package>` = updates a specific packages listed in the `liquibase.json` file

- `liquibase lpm install` = install to this Liquibase project all packages listed in the `liquibase.json` file
- `liquibase lpm install <package>` = install to this Liquibase project a specific packages listed in the `liquibase.json` file

- `liquibase lpm upgrade` = upgrade lpm binary
- `liquibase lpm upgrade --all` = upgrade lpm + all packages
- `liquibase lpm upgrade <package>` = upgrade specific package


## Usage *not within* Liquibase Community

```shell
lpm <command>
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
* upgrade




## Autocompletion scripts for pre-Liquibase 5.0+ versions

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

## License

This project is licensed under the Apache License 2.0. See the LICENSE file for details.

This project also includes software governed by the Mozilla Public License, v. 2.0.

* `github.com/hashicorp/go-version v1.6.0` (<http://mozilla.org/MPL/2.0/>)

## Releasing LPM

LPM uses an **automated release process** with minimal manual intervention.

### Quick Release Guide for Maintainers

**3-Step Release Process:**

1. **Bump Version** - Run the "Bump Version" GitHub Action
   - Select `patch`, `minor`, or `major`
   - Automatically updates VERSION file and creates tag

2. **Build Artifacts** - Run the "Attach Artifact to Release" GitHub Action
   - Builds for all platforms (darwin, linux, windows, s390x)
   - Generates checksums for all platforms
   - Uploads artifacts to draft release

3. **Publish** - Review and publish the draft release on GitHub
   - VERSION file automatically syncs after publishing

**That's it!** ✨ The entire process is now mostly automated.

### For Contributors: PR Labeling Guidelines

To help with automated changelog generation, please label your PRs:

- `feature` or `enhancement` - New features (🚀 Features section)
- `bug` or `fix` - Bug fixes (🐛 Bug Fixes section)
- `documentation` - Documentation changes (📚 Documentation section)
- `dependencies` - Dependency updates (📦 Dependencies section)
- `security` - Security fixes (🔒 Security section)

Most labels are automatically applied based on file paths and PR titles!

### Detailed Documentation

For complete release process documentation, troubleshooting, and best practices, see [RELEASING.md](RELEASING.md).

## Adding an Extension

https://github.com/liquibase/liquibase-package-manager/blob/master/CONTRIBUTING.md
