# How to Release LPM

This document describes the **fully automated release process** for LPM using Release Drafter and GitHub Actions.

## Quick Start (TL;DR)

**To create a release:**

1. Merge PRs to `master` (labeled appropriately: `feature`, `bug`, `breaking`, etc.)
2. Review the auto-generated draft release on GitHub
3. Click **Publish release**
4. Wait ~5-10 minutes for artifacts to build and attach automatically

That's it! Everything else happens automatically. ✨

---

## Release Process Overview

```
┌─────────────────────────────────────────────────────────────────┐
│  1. Merge PRs to master                                         │
│     • Ensure PRs are labeled (feature, bug, breaking, etc.)     │
│     • Labels determine version bump (major/minor/patch)         │
└────────────────┬────────────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────────┐
│  2. Release Drafter Updates Draft (Automatic on push)           │
│     • Creates/updates draft release                             │
│     • Calculates next version from PR labels                    │
│     • Generates categorized changelog                           │
│     • NO artifacts yet - just changelog                         │
└────────────────┬────────────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────────┐
│  3. Review and Publish (Manual - ONLY manual step)              │
│     • Review the draft release on GitHub                        │
│     • Verify changelog is accurate                              │
│     • Click "Publish release"                                   │
└────────────────┬────────────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────────┐
│  4. Post-Publish Automation (Automatic - runs in parallel)      │
│     ┌───────────────────────────────────────────────────────┐   │
│     │ A. Build & Attach Artifacts (~5-10 min)              │   │
│     │    • Builds for all 6 platforms                       │   │
│     │    • Generates SHA256 checksums                       │   │
│     │    • Uploads artifacts to published release           │   │
│     └───────────────────────────────────────────────────────┘   │
│     ┌───────────────────────────────────────────────────────┐   │
│     │ B. Sync VERSION File (~10 sec)                        │   │
│     │    • Updates VERSION file to match release            │   │
│     │    • Updates internal/app/VERSION                     │   │
│     │    • Commits to master                                │   │
│     └───────────────────────────────────────────────────────┘   │
│     ┌───────────────────────────────────────────────────────┐   │
│     │ C. Update Docker Repository (~30 sec)                 │   │
│     │    • Extracts checksums from release                  │   │
│     │    • Updates Dockerfiles in liquibase/docker          │   │
│     │    • Creates PR with changes                          │   │
│     └───────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

---

## Detailed Step-by-Step Guide

### Step 1: Merge PRs to Master

**How versioning works:**

The version is automatically determined by PR labels:

| Label | Version Bump | Example |
|-------|-------------|---------|
| `breaking`, `major` | Major | 0.2.9 → 1.0.0 |
| `feature`, `enhancement`, `minor` | Minor | 0.2.9 → 0.3.0 |
| `bug`, `bugfix`, `fix`, `patch`, `dependencies`, `security` | Patch | 0.2.9 → 0.2.10 |
| None | Patch (default) | 0.2.9 → 0.2.10 |

**Before merging your PR:**
1. Add appropriate labels to your PR
2. Verify the PR title and description are clear
3. Merge to `master`

**Pro tip:** You can merge multiple PRs before releasing. Release Drafter will accumulate all changes in the draft release.

---

### Step 2: Draft Release Maintained Automatically

**What happens automatically:**

When you push to `master`:

1. The **Create Release** workflow runs automatically
2. Release Drafter creates/updates a draft release with:
   - Next version calculated from PR labels
   - Automatically categorized changelog (🚀 Features, 🐛 Bug Fixes, etc.)
   - List of contributors

**Important**: The draft release contains **ONLY the changelog** at this point. Artifacts are built AFTER you publish.

**No action required** - Release Drafter maintains a single draft release that gets updated with each merge to master.

---

### Step 3: Review and Publish the Release (ONLY MANUAL STEP)

This is the **only manual step** in the entire process.

1. Go to the [Releases page](https://github.com/liquibase/liquibase-package-manager/releases)
2. Find the draft release (should be at the top)
3. **Review:**
   - ✅ Version number is correct
   - ✅ Changelog is accurate and well-formatted
   - ✅ All PRs are properly categorized
4. Edit the release notes if needed
5. Click **Publish release**

**Note**: There will be no artifacts on the draft yet - they're built after you publish.

---

### Step 4: Post-Publish Automation (Runs Automatically)

**What happens automatically:**

When you click "Publish release", **three workflows** run in parallel:

#### A. Build and Attach Artifacts (~5-10 minutes)

The **Build and Attach Artifacts** workflow:
1. Builds artifacts for all platforms:
   - `darwin` (macOS Intel)
   - `darwin-arm64` (macOS Apple Silicon)
   - `linux` (Linux AMD64)
   - `linux-arm64` (Linux ARM64)
   - `s390x` (IBM System z)
   - `windows` (Windows AMD64)
2. Calculates SHA256 checksums for **ALL** platforms
3. Uploads all artifacts to the published release
4. Generates and uploads a `checksums.txt` file

**Duration**: ~5-10 minutes (builds run in parallel)

**Result**: All 6 platform binaries + checksums.txt attached to the published release

#### B. Sync VERSION File (~10 seconds)

The **Sync VERSION on Release Publish** workflow:
1. Extracts version from release tag
2. Updates `VERSION` file if different
3. Updates `internal/app/VERSION`
4. Commits and pushes to master

**Duration**: ~10 seconds

**Result**: VERSION file in master always matches the latest published release

#### C. Update Docker Repository (~30 seconds)

The **Update Docker Repository** workflow:
1. Downloads `checksums.txt` from the release
2. Extracts SHA256 checksums for linux-amd64 and linux-arm64
3. Updates three Dockerfiles in `liquibase/docker` repository:
   - `Dockerfile`
   - `Dockerfile.alpine`
   - `DockerfileSecure`
4. Creates a PR with all changes

**Duration**: ~30 seconds

**Result**: PR created in [liquibase/docker](https://github.com/liquibase/docker/pulls) ready for review

---

### Step 5: Review Docker PR

**Action required (optional but recommended):**

1. Go to the [liquibase/docker repository](https://github.com/liquibase/docker/pulls)
2. Find the PR titled "Update LPM to vX.Y.Z"
3. Verify the checksums match the release
4. Merge the PR to trigger Docker image builds

**Note**: Docker repo updates require a `BOT_TOKEN` secret with write access to `liquibase/docker`.

---

## Important Notes

### Artifacts Are Built AFTER Publishing

**This is by design** for several reasons:

1. **Simplicity**: One trigger (publish) for everything
2. **Standard pattern**: Most projects using Release Drafter work this way
3. **No cross-workflow issues**: Avoids GitHub Actions token limitations
4. **Faster process**: No waiting for builds before you can publish

**What this means:**
- The draft release has **no artifacts** - only changelog
- When you publish, artifacts build automatically (~5-10 minutes)
- Users downloading immediately after publish may need to wait for artifacts

**If you need to review artifacts before publishing:**
- Manually trigger the "Build and Attach Artifacts" workflow on the draft
- Review the artifacts
- Then publish when ready

---

## PR Labeling for Better Changelogs

Release Drafter automatically categorizes PRs based on labels. Use these labels on your PRs:

| Label | Category | Version Bump |
|-------|----------|-------------|
| `feature`, `enhancement`, `minor` | 🚀 Features | `minor` |
| `bug`, `bugfix`, `fix`, `patch` | 🐛 Bug Fixes | `patch` |
| `documentation`, `docs` | 📚 Documentation | - |
| `dependencies` | 📦 Dependencies | `patch` |
| `security` | 🔒 Security | `patch` |
| `breaking`, `major` | - | `major` |
| `maintenance`, `refactor`, `chore` | 🔧 Maintenance | - |
| `performance` | ⚡ Performance | - |
| `test`, `tests` | 🧪 Testing | - |
| `build`, `ci`, `github-actions` | 🏗️ Build & CI | - |
| `skip-changelog` | (excluded) | - |

**Auto-labeling** is configured for:
- PRs modifying `go.mod`/`go.sum` → `dependencies`
- PRs from dependabot branches → `dependencies`
- PRs with "fix" or "bug" in title → `bug`
- PRs with "feat" or "feature" in title → `feature`
- PRs modifying `*.md` files → `documentation`

---

## Workflow Files Reference

| Workflow | Trigger | Purpose | Duration |
|----------|---------|---------|----------|
| `create-release.yml` | Push to master | Create/update draft release via Release Drafter | ~5 sec |
| `attach-artifact-release.yml` | Release published | Build and upload release artifacts | ~5-10 min |
| `publish-release.yml` | Release published | Sync VERSION file after publish | ~10 sec |
| `update-docker-repo.yml` | Release published | Update LPM in docker repository | ~30 sec |

---

## Manual Workflow Dispatch (If Needed)

While the process is fully automated, you can manually trigger workflows if needed:

### Manually Build and Attach Artifacts Before Publishing

If you want to review artifacts before publishing:

1. Go to **Actions** → **Build and Attach Artifacts**
2. Click **Run workflow**
3. Leave version **empty** (auto-detects from latest draft)
4. Leave release-id **empty** (uses latest draft)
5. Click **Run workflow**
6. Wait for artifacts to build (~5-10 minutes)
7. Review artifacts on the draft release
8. Publish when ready

### Manually Update Docker Repository

If the automatic Docker update fails:

1. Go to **Actions** → **Update Docker Repository**
2. Click **Run workflow**
3. Enter the LPM version (e.g., `0.2.15` - without 'v' prefix)
4. Leave checksums empty (will fetch from release)
5. Click **Run workflow**

---

## Troubleshooting

### Artifacts not appearing after publish

**Problem:** Published release but no artifacts showing up.

**Solution:**
1. Go to **Actions** tab
2. Check the "Build and Attach Artifacts" workflow run
3. If failed, check logs for errors
4. Common issues:
   - Go module download failures → Retry the workflow
   - Build errors → Fix code, create new release
5. If successful, artifacts should appear within 5-10 minutes

### VERSION file out of sync

**Problem:** VERSION file doesn't match the latest release.

**Solution:**
1. Check the "Sync VERSION on Release Publish" workflow logs
2. If failed, manually trigger it or run:
   ```bash
   LATEST_TAG=$(git describe --tags --abbrev=0)
   VERSION=${LATEST_TAG#v}
   echo "$VERSION" > VERSION
   git add VERSION
   git commit -m "Sync VERSION to $VERSION"
   git push origin master
   ```

### Draft release not created

**Problem:** No draft release after pushing to master.

**Solution:**
1. Check the "Create Release" workflow logs
2. Ensure you have proper permissions
3. Verify `release-drafter.yml` configuration is valid
4. Check that PRs have appropriate labels

### Docker PR not created

**Problem:** No PR created in docker repository after release.

**Solution:**
1. Verify `BOT_TOKEN` secret is configured with write access to `liquibase/docker`
2. Check "Update Docker Repository" workflow logs
3. Manually trigger the workflow if needed

---

## Emergency Manual Release

If all automation fails (extremely rare):

1. **Create tag:**
   ```bash
   git tag -a vX.Y.Z -m "Release vX.Y.Z"
   git push origin --tags
   ```

2. **Create release on GitHub:**
   - Go to Releases → Draft a new release
   - Select the tag you just created
   - Write release notes
   - Click **Publish release**

3. **Manually trigger artifact build:**
   - Go to Actions → "Build and Attach Artifacts"
   - Run workflow with the version number

4. **Update VERSION file:**
   ```bash
   echo "X.Y.Z" > VERSION
   git add VERSION
   git commit -m "Sync VERSION to X.Y.Z"
   git push origin master
   ```

---

## Best Practices

✅ **DO:**
- Label PRs appropriately for accurate changelogs and versioning
- Review draft releases before publishing
- Wait for all 3 post-publish workflows to complete (~10 minutes total)
- Test releases in non-production environments first
- Keep PR titles clear and descriptive

❌ **DON'T:**
- Manually edit VERSION file (it syncs automatically on publish)
- Publish releases without reviewing the changelog
- Skip labeling PRs (affects versioning and categorization)
- Delete tags without good reason
- Force push to master

---

## Timeline: From Merge to Complete Release

| Step | Action | Duration | Manual? |
|------|--------|----------|---------|
| 1. | Merge PR to master | Instant | ✋ Manual |
| 2. | Release Drafter updates draft | ~5 sec | ✅ Auto |
| 3. | Review and publish release | Variable | ✋ Manual |
| 4. | Build artifacts (parallel) | ~5-10 min | ✅ Auto |
| 5. | Sync VERSION file (parallel) | ~10 sec | ✅ Auto |
| 6. | Create Docker PR (parallel) | ~30 sec | ✅ Auto |
| 7. | Review and merge Docker PR | Variable | ✋ Manual (optional) |

**Total automated time after publishing**: ~5-10 minutes

**Total manual steps**: 2 (publish release, merge Docker PR)

---

## Comparison: Old vs New Process

| Aspect | Old Manual Process | New Automated Process |
|--------|-------------------|---------------------|
| Manual Steps | ~5-7 steps | 1 step (publish) |
| Version Management | Manual bump | Automatic via PR labels |
| Artifact Building | Manual trigger before publish | Automatic after publish |
| VERSION File Sync | Manual | Automatic |
| Docker Updates | Manual | Automatic PR |
| Time to Release | ~20-30 minutes | ~5-10 minutes |
| Error Prone | High | Low |
| Can Review Artifacts | Yes (before publish) | No (after publish)* |
| Follows Best Practices | Partial | ✅ Fully compliant |

\* *You can manually trigger artifact build on draft if needed for review*

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

**2025-10-16:** Refactored to standard Release Drafter pattern (v2)
- **BREAKING**: Artifacts now build AFTER publishing (not before)
- Fixed critical workflow trigger issue (GITHUB_TOKEN limitation)
- Artifacts attach to published release automatically
- Simplified to truly single manual step: publish release
- All post-publish automation runs in parallel
- Aligned with industry-standard Release Drafter usage

**2025-10-16:** Initial refactor to Release Drafter pattern (v1)
- Removed manual version bumping (automatic via PR labels)
- Made artifact building automatic (attempted on draft update)
- Reduced manual steps from 3 to 1

**2025-10-14:** Original automation implementation
- Added automated VERSION bumping workflow
- Enhanced release drafter with categorization
- Added SHA256 checksums for all platforms
- Added automatic Docker repository updates
