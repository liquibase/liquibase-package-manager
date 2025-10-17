# Docker Repository Auto-Update

This workflow automatically updates the LPM version in the `liquibase/docker` repository when a new LPM release is published.

## How It Works

### Automatic Trigger (Recommended)

When you **publish** a release in the `liquibase-package-manager` repository:

1. The workflow automatically runs
2. Extracts the version from the release tag
3. Downloads the `checksums.txt` file from the release
4. Extracts SHA256 checksums for `linux-amd64` and `linux-arm64`
5. Updates three Dockerfiles in the `liquibase/docker` repository:
   - `Dockerfile`
   - `Dockerfile.alpine`
   - `DockerfileSecure`
6. Creates a PR in the docker repository

### Manual Trigger

You can also manually trigger the workflow:

1. Go to **Actions** → **Update Docker Repository with New LPM Version**
2. Click **Run workflow**
3. Enter:
   - **lpm-version**: Version without 'v' prefix (e.g., `0.2.15`)
   - **sha256-amd64**: (Optional) Leave empty to auto-fetch from release
   - **sha256-arm64**: (Optional) Leave empty to auto-fetch from release
4. Click **Run workflow**

## Prerequisites

### Required Secret

The workflow requires a `BOT_TOKEN` secret with:
- Write access to the `liquibase/docker` repository
- Ability to create branches and pull requests

To set up the token:

1. Create a GitHub Personal Access Token (classic) or Fine-grained token with:
   - Repository access: `liquibase/docker`
   - Permissions: `Contents: Read and Write`, `Pull Requests: Read and Write`
2. Add it as a repository secret named `BOT_TOKEN` in the `liquibase-package-manager` repository

## What Gets Updated

### Dockerfile Variables

In each Dockerfile, three ARG variables are updated:

```dockerfile
ARG LPM_VERSION=0.2.15
ARG LPM_SHA256=<sha256 for linux-amd64>
ARG LPM_SHA256_ARM=<sha256 for linux-arm64>
```

### Files Updated

- `Dockerfile` - Standard Ubuntu-based image
- `Dockerfile.alpine` - Alpine Linux-based image
- `DockerfileSecure` - Secure/licensed image

## PR Details

The automatically created PR includes:

- **Title**: `Update LPM to v{version}`
- **Labels**: `lpm`, `dependencies`, `automated`
- **Description**: 
  - Version being updated
  - SHA256 checksums for both architectures
  - List of files changed
  - Link to the LPM release
  - Verification notes

## Troubleshooting

### "checksums.txt not found"

**Cause**: The checksums file wasn't uploaded to the release.

**Solution**: Ensure the "Attach Artifact to Release" workflow completes successfully before publishing.

### "Invalid SHA256"

**Cause**: Checksum doesn't match expected format (64 characters).

**Solution**: 
1. Check that `checksums.txt` contains valid checksums
2. Manually trigger the workflow with correct checksums

### "No changes detected"

**Cause**: Docker repository already has the same version/checksums.

**Solution**: This is expected if the version was already updated. No action needed.

### "PR creation failed"

**Cause**: Missing or invalid `BOT_TOKEN` secret.

**Solution**:
1. Verify `BOT_TOKEN` is configured in repository secrets
2. Check token has correct permissions
3. Ensure token hasn't expired

## Manual Update (Fallback)

If automation fails, you can manually update the docker repository:

```bash
# Clone docker repository
git clone https://github.com/liquibase/docker.git
cd docker

# Update version and checksums
VERSION="0.2.15"
SHA256_AMD64="..."  # Get from release
SHA256_ARM64="..."  # Get from release

# Update all Dockerfiles
for file in Dockerfile Dockerfile.alpine DockerfileSecure; do
  sed -i "s/^ARG LPM_VERSION=.*/ARG LPM_VERSION=${VERSION}/" $file
  sed -i "s/^ARG LPM_SHA256=.*/ARG LPM_SHA256=${SHA256_AMD64}/" $file
  sed -i "s/^ARG LPM_SHA256_ARM=.*/ARG LPM_SHA256_ARM=${SHA256_ARM64}/" $file
done

# Create branch and PR
git checkout -b update-lpm-${VERSION}
git add Dockerfile Dockerfile.alpine DockerfileSecure
git commit -m "Update LPM to version ${VERSION}"
git push origin update-lpm-${VERSION}
# Then create PR on GitHub
```

## Testing

To test the workflow without creating a real PR:

1. Fork the `liquibase/docker` repository to your account
2. Update the workflow to use your fork:
   ```yaml
   repository: YOUR_USERNAME/docker
   ```
3. Run the workflow manually
4. Verify the PR is created in your fork
5. Review changes and delete the test PR

## Related Workflows

- `attach-artifact-release.yml` - Generates the checksums file
- `publish-release.yml` - Syncs VERSION file after release
- Both run before this workflow in the release process
