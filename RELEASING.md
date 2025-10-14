# How to Release LPM

This document describes the **automated release process** for LPM. The process has been significantly streamlined to reduce manual steps and prevent version inconsistencies.

## Quick Start (TL;DR)

**For a standard patch release:**

1. Run the "Bump Version" workflow (select `patch`)
2. Run the "Attach Artifact to Release" workflow (leave version empty)
3. Publish the draft release on GitHub

That's it! ✨

---

## Release Process Overview

```
┌─────────────────────────────────────────────────────────────────┐
│  1. Bump Version Workflow (Manual)                              │
│     • Calculates new version (patch/minor/major)                │
│     • Updates VERSION file                                      │
│     • Creates git tag                                           │
│     • Pushes to master                                          │
└────────────────┬────────────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────────┐
│  2. Create Release Workflow (Automatic on push to master)       │
│     • Creates/updates draft release                             │
│     • Generates changelog from merged PRs                       │
│     • Updates VERSION file with draft version                   │
└────────────────┬────────────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────────┐
│  3. Attach Artifact Workflow (Manual)                           │
│     • Reads VERSION from file (or accepts custom version)       │
│     • Builds artifacts for all platforms                        │
│     • Generates SHA256 checksums for ALL platforms              │
│     • Uploads artifacts + checksums.txt to release              │
└────────────────┬────────────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────────┐
│  4. Review and Publish (Manual)                                 │
│     • Review the draft release on GitHub                        │
│     • Verify all artifacts are attached                         │
│     • Click "Publish release"                                   │
└────────────────┬────────────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────────┐
│  5. Sync VERSION Workflow (Automatic on release publish)        │
│     • Syncs VERSION file with published release tag             │
│     • Updates internal/app/VERSION                              │
│     • Commits changes to master                                 │
└────────────────┬────────────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────────┐
│  6. Update Docker Repo Workflow (Automatic on release publish)  │
│     • Extracts checksums for linux-amd64 and linux-arm64        │
│     • Updates Dockerfile, Dockerfile.alpine, DockerfileSecure   │
│     • Creates PR in liquibase/docker repository                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## Detailed Step-by-Step Guide

### Step 1: Bump the Version

**Option A: Automated Bump (Recommended)**

1. Go to **Actions** → **Bump Version**
2. Click **Run workflow**
3. Select version bump type:
   - `patch` - Bug fixes (0.2.9 → 0.2.10)
   - `minor` - New features (0.2.9 → 0.3.0)
   - `major` - Breaking changes (0.2.9 → 1.0.0)
4. Or enter a custom version (optional)
5. Click **Run workflow**

This will:
- ✅ Update the `VERSION` file
- ✅ Create a git commit
- ✅ Create and push a git tag (e.g., `v0.2.10`)
- ✅ Trigger the release drafter

**Option B: Manual Bump (Not Recommended)**

If you must bump the version manually:

1. Update the `VERSION` file
2. Commit: `git commit -m "Bump version to X.Y.Z"`
3. Tag: `git tag -a vX.Y.Z -m "Release vX.Y.Z"`
4. Push: `git push origin master --tags`

⚠️ **Warning:** Manual version bumps will be validated by the `Validate VERSION` workflow on PR.

---

### Step 2: Draft Release Created Automatically

**What happens automatically:**

When you push to `master` (either through the Bump Version workflow or a merged PR):

1. The **Create Release** workflow runs automatically
2. A draft release is created/updated with:
   - Automatically categorized changelog (🚀 Features, 🐛 Bug Fixes, etc.)
   - List of contributors
   - Version calculated from PR labels
3. The `VERSION` file is synced with the draft release version

**No action required** - just verify the draft release looks correct on GitHub.

---

### Step 3: Build and Attach Artifacts

1. Go to **Actions** → **Attach Artifact to Release**
2. Click **Run workflow**
3. Leave version **empty** (it will read from VERSION file) OR enter a custom version
4. Leave release-id **empty** (it will use the latest draft)
5. Click **Run workflow**

This workflow will:
- ✅ Build artifacts for all platforms:
  - `darwin` (macOS Intel)
  - `darwin-arm64` (macOS Apple Silicon)
  - `linux` (Linux AMD64)
  - `linux-arm64` (Linux ARM64)
  - `s390x` (IBM System z)
  - `windows` (Windows AMD64)
- ✅ Calculate SHA256 checksums for **ALL** platforms
- ✅ Upload all artifacts to the draft release
- ✅ Generate and upload a `checksums.txt` file
- ✅ Display a comprehensive summary with all checksums

**Build time:** ~5-10 minutes (runs in parallel)

---

### Step 4: Review and Publish the Release

1. Go to the [Releases page](https://github.com/liquibase/liquibase-package-manager/releases)
2. Find the draft release (should be at the top)
3. **Review:**
   - ✅ Version number is correct
   - ✅ Changelog is accurate and well-formatted
   - ✅ All 6 platform artifacts are attached
   - ✅ `checksums.txt` file is attached
4. Edit the release notes if needed
5. Click **Publish release**

---

### Step 5: VERSION File Synced Automatically

**What happens automatically:**

When you publish the release:

1. The **Sync VERSION on Release Publish** workflow runs automatically
2. It extracts the version from the release tag
3. It updates the `VERSION` file if different
4. It updates `internal/app/VERSION`
5. It commits and pushes to master

**No action required** - the VERSION file will always stay in sync with published releases.

---

### Step 6: Docker Repository Updated Automatically

**What happens automatically:**

When you publish the release:

1. The **Update Docker Repository** workflow runs automatically
2. It downloads the `checksums.txt` from the release
3. It extracts SHA256 checksums for:
   - `linux-amd64` (x86_64)
   - `linux-arm64` (aarch64)
4. It checks out the `liquibase/docker` repository
5. It updates three Dockerfiles:
   - `Dockerfile`
   - `Dockerfile.alpine`
   - `DockerfileSecure`
6. It creates a PR in the docker repository with:
   - Updated LPM version
   - Updated SHA256 checksums for both architectures
   - Professional PR description with all changes
   - Labels: `lpm`, `dependencies`, `automated`

**Action required:**

1. Review the PR in the [liquibase/docker repository](https://github.com/liquibase/docker/pulls)
2. Verify the checksums match the release
3. Merge the PR to trigger Docker image builds

**Note:** This requires a `BOT_TOKEN` secret to be configured with write access to the `liquibase/docker` repository.

---

## Workflow Files Reference

| Workflow | Trigger | Purpose |
|----------|---------|--------|
| `bump-version.yml` | Manual | Bump VERSION file and create tag |
| `create-release.yml` | Push to master | Create/update draft release |
| `attach-artifact-release.yml` | Manual | Build and upload release artifacts |
| `validate-version.yml` | PR with VERSION changes | Validate VERSION file changes |
| `publish-release.yml` | Release published | Sync VERSION file after publish |
| `update-docker-repo.yml` | Release published/Manual | Update LPM in docker repository |
| `test.yml` | PR/Push to master | Run tests and quality checks |
| `nightly-e2e-tests.yml` | Schedule/Manual | Run end-to-end tests |
| `nightly-update-packages.yml` | Schedule/Manual/PR | Update packages.json |

---

## PR Labeling for Better Changelogs

The release drafter automatically categorizes PRs based on labels. Use these labels on your PRs:

| Label | Category | Version Bump |
|-------|----------|-------------|
| `feature`, `enhancement` | 🚀 Features | `minor` |
| `bug`, `bugfix`, `fix` | 🐛 Bug Fixes | `patch` |
| `documentation`, `docs` | 📚 Documentation | - |
| `dependencies` | 📦 Dependencies | `patch` |
| `security` | 🔒 Security | `patch` |
| `breaking`, `major` | - | `major` |
| `skip-changelog` | (excluded) | - |

**Auto-labeling** is configured for:
- PRs modifying `go.mod`/`go.sum` → `dependencies`
- PRs with "fix" or "bug" in title → `bug`
- PRs with "feat" or "feature" in title → `feature`
- PRs modifying `*.md` files → `documentation`

---

## Troubleshooting

### VERSION file out of sync

**Problem:** VERSION file doesn't match the latest release.

**Solution:**
```bash
# Get the latest release tag
LATEST_TAG=$(git describe --tags --abbrev=0)
VERSION=${LATEST_TAG#v}

# Update VERSION file
echo "$VERSION" > VERSION

# Commit and push
git add VERSION
git commit -m "Sync VERSION to $VERSION"
git push origin master
```

### Build artifacts failed

**Problem:** "Attach Artifact to Release" workflow failed.

**Solution:**
1. Check the workflow logs for specific errors
2. Common issues:
   - Go module download failures → Retry the workflow
   - Build errors → Fix the code and re-run
3. Re-run the workflow after fixing issues

### Draft release not created

**Problem:** No draft release after pushing to master.

**Solution:**
1. Check the "Create Release" workflow logs
2. Ensure you have proper permissions
3. Verify `release-drafter.yml` configuration is valid
4. Manually trigger the workflow if needed

### Version validation failed on PR

**Problem:** PR with VERSION changes fails validation.

**Solution:**
1. Read the validation error comment on the PR
2. Common issues:
   - Version format incorrect (must be X.Y.Z)
   - Version decreased (not allowed)
   - Tag already exists
3. Fix the VERSION file or use the "Bump Version" workflow instead

---

## Manual Override (Emergency)

If automation fails and you need to release manually:

1. **Update VERSION file:**
   ```bash
   echo "X.Y.Z" > VERSION
   git add VERSION
   git commit -m "Bump version to X.Y.Z"
   ```

2. **Create and push tag:**
   ```bash
   git tag -a vX.Y.Z -m "Release vX.Y.Z"
   git push origin master --tags
   ```

3. **Build artifacts locally:**
   ```bash
   make release
   ```

4. **Create release manually on GitHub:**
   - Go to Releases → Draft a new release
   - Set tag to `vX.Y.Z`
   - Upload artifacts from `bin/` directories
   - Generate checksums: `sha256sum bin/*/*.zip > checksums.txt`
   - Upload `checksums.txt`
   - Publish the release

---

## Best Practices

✅ **DO:**
- Use the "Bump Version" workflow for version bumps
- Label PRs appropriately for better changelogs
- Review draft releases before publishing
- Verify all artifacts are attached before publishing
- Test releases in non-production environments first

❌ **DON'T:**
- Manually edit VERSION file (use the workflow)
- Skip the artifact attachment step
- Publish releases without reviewing
- Use random version numbers (follow semver)
- Delete tags without good reason

---

## Support

If you encounter issues with the release process:

1. Check workflow logs in the Actions tab
2. Review this documentation
3. Open an issue with:
   - Workflow name and run ID
   - Error messages from logs
   - Steps to reproduce

---

## Changelog

**2025-10-14:** Complete automation overhaul
- Added automated VERSION bumping workflow
- Enhanced release drafter with categorization
- Replaced deprecated GitHub Actions
- Added SHA256 checksums for all platforms
- Added VERSION file validation
- Added automatic VERSION sync on release publish
- Reduced manual steps from 10+ to 3
