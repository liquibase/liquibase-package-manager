# Contributing

## Git Fork Workflow

1. Fork the repository
1. Enable public access for your fork
1. Make your code changes
1. Commit your changes (`git commit -am 'Add some feature/folder'`)
1. Push your changes (`git push origin master`)
1. Create a Pull Request to [master](https://github.com/liquibase/liquibase-package-manager) branch.

## Important Guidelines!

- No acronyms in folder name in order to provide clarity and avoid collisions.
- Do not use camel case or underscores for folder names.
- Arrange your folders in alphabetical order.
- Folder name or folder description should not contain special characters.
- 
## Add an Extension

1. Edit `cmd/populator/modules.go` under `func init() `  with the extension details.
```
{
    name:        "liquibase-commercial-dynamodb",
    category:    Extension,
    url:         "https://repo1.maven.org/maven2/org/liquibase/ext/liquibase-commercial-dynamodb",
    artifactory: Maven{},
}
```
2. Edit `internal/app/packages.json` with the extension details.
```
{
    "name": "extension-name",
    "category": "extension", // or "driver" or "pro" Categories should match with the populator (which is an enum not a string)
    "versions": []
  }
```
3. Push your changes and a PR will automatically be created with edits to `internal/app/packages.json` file with values for `algorithm`, `checksum`, `liquibaseCore`, `path` 
