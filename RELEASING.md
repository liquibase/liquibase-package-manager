# How to Release LPM

## Workflow Definitions

- **Release Workflow**:
  - `create-release.yml` creates/updates the draft release when a push to master is made.
- **Attach Artifact Workflow**:
  - `attach-artifact-release.yml` is a manually executed workflow with an input value to set the LPM version for the draft release.

## Attach Artifact Workflow Details

1. **Cleanup**: Retrieves the latest draft release and removes any old artifacts from previous runs.
2. **Set Release Name**: Uses the version input parameter to set the draft release name.
3. **Set Version File**: Updates the `VERSION` file in the repository with the version input parameter.
4. **Run Release Command**: Executes the `make release` command.
5. **Upload Artifacts**: Uploads each zip file to the draft release.
6. **SHA256 Summary**: Prints the sha256 checksums for `linux` and `linux-arm64` in the action run summary.

## Workflow Process

- **Automatic Draft Creation**: Every merged PR to master will create/update the draft release.
- **Manual Artifact Attachment**: When ready to release a new version, manually trigger the `Attach Artifact to Release` workflow with the new version.

### Final Steps

If everything looks good, finalize the draft release, and GitHub will attach the source code tar and zip files as usual.
