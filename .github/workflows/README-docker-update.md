# Docker Repository Auto-Update — Consumer Fanout

When a new LPM release is published, this workflow automatically propagates the updated version
and checksums to every consumer Docker repository in one run. Each consumer repo gets its own
isolated job, so a failure in one repo never blocks the others.

## Quick path

1. Publish a new LPM release (or dispatch manually).
2. The `prepare` job resolves the version and validates both SHA256 checksums.
3. The `update-repo` matrix job fans out — one job per consumer repo.
4. Each job opens a PR in its target repo (or records "no-op" if already up to date).
5. The `summary` job renders a per-repo outcome table in the run's Step Summary.

---

## Workflow structure

```
Trigger (workflow_run success | workflow_dispatch)
  │
  ├─ prepare          (runs once — resolves version + checksums)
  │
  ├─ update-repo      (needs: prepare, fail-fast: false)
  │    ├─ liquibase/liquibase job
  │    └─ liquibase/liquibase-pro job   ← see P2 note below
  │
  └─ summary          (needs: [prepare, update-repo], if: always())
```

The GitHub App token is **re-minted inside each matrix job**. It is never passed
via a job output — job outputs are not treated as secrets and can leak in logs.

---

## Consumer matrix

| Repo | Files updated |
|------|---------------|
| `liquibase/liquibase` | `docker/Dockerfile` `docker/Dockerfile.alpine` |
| `liquibase/liquibase-pro` | `docker/Dockerfile` *(`core/docker/*` is legacy — intentionally NOT updated)* |

Each matrix entry has two fields: `repo` and `files` (space-delimited list of paths
relative to the repo root). No other workflow logic needs to change when the list changes.

---

## How to add a consumer repo

Add one entry to the `strategy.matrix.include` list in `update-docker-repo.yml`:

```yaml
- repo: org/repo-name
  files: "path/to/Dockerfile path/to/Dockerfile.alpine"
```

That's the only change required (R9).

---

## How to remove a consumer repo

Delete its `include` entry from the matrix. No other YAML changes are needed.

---

## Workflow inputs (workflow_dispatch)

| Input | Required | Default | Description |
|-------|----------|---------|-------------|
| `lpm-version` | Yes | — | Version without `v` prefix (e.g. `0.3.5`) |
| `sha256-amd64` | No | *(auto-fetched)* | Override SHA256 for `linux-amd64`; leave empty to fetch from release |
| `sha256-arm64` | No | *(auto-fetched)* | Override SHA256 for `linux-arm64`; leave empty to fetch from release |
| `dry-run` | No | `false` | Show diff and skip PR creation — use this to test without side effects |

### Dry-run usage

The `dry-run` input lets you validate the full pipeline (checkout, sed, git diff)
against the real target repos without opening any PRs:

1. Go to **Actions → Update Docker Repository with New LPM Version → Run workflow**.
2. Set `lpm-version` to a known released version (e.g. `0.3.5`).
3. Set `dry-run: true`.
4. Each matrix job will check out the repo, run the sed updates, and print the diff
   (or "already up to date") — but will not create a PR.
5. The Summary tab shows `dry-run-would-create` or `noop` per repo.

A dry-run that returns a 403 on a checkout confirms the GitHub App is not installed
on that repo (**P2 blocker**).

---

## Partial-failure semantics

The matrix runs with `fail-fast: false`. This means:

- If `liquibase/liquibase-pro` fails (e.g. a 403 because the App is not installed),
  the `liquibase/liquibase` job still runs and opens its PR.
- The overall workflow run is marked **failed** if ANY matrix job fails. This is
  **intentional alerting** — not a bug. Check the Summary tab for the per-repo breakdown.

Do not mistake a partial run for a complete success just because some repos got PRs.

---

## Preconditions before first real release

### P2 (HIGH — blocks all private-repo jobs)

The Liquibase GitHub App (`LIQUIBASE_GITHUB_APP_ID`) must be installed on every target
repo with `contents: write` and `pull-requests: write`. This is required on
`liquibase/liquibase-pro` (private) in particular.

Fastest confirmation: run a `dry-run` dispatch. A 403 on checkout means the App is
not yet installed.

### liquibase-pro live file set (confirmed)

`liquibase/liquibase-pro` updates **only `docker/Dockerfile`**. The `core/docker/`
Dockerfiles are **legacy and must NOT be added** to the matrix — they are intentionally
left unmanaged.

The workflow fails loudly on a missing file (`::error::` annotation + exit 1), so a
wrong path breaks only the `liquibase-pro` job — other repos are unaffected.

---

## Out-of-scope repos

The following repos pin an older LPM version intentionally and are **not managed by this
workflow**. Do not add them to the matrix.

| Repo | Pinned LPM version | Reason |
|------|--------------------|--------|
| `liquibase/liquibase-test-harness` | 0.2.3 | Test infra; version pinned by test team |
| `liquibase/mongodemo` | 0.2.0 | Demo repo; version pinned intentionally |
| `liquibase/flow-demo` | 0.1.7 | Demo repo; version pinned intentionally |
| `liquibase/devops-misc` | 0.1.3 | Infra repo; version pinned intentionally |

---

## What gets updated in each Dockerfile

Three `ARG` lines are updated per file:

```dockerfile
ARG LPM_VERSION=<version>
ARG LPM_SHA256=<sha256 for linux-amd64>
ARG LPM_SHA256_ARM=<sha256 for linux-arm64>
```

The `^ARG LPM_SHA256=` sed pattern matches the amd64 line only. The trailing `=` is
a boundary that cannot match `LPM_SHA256_ARM=` — the anchor is load-bearing and must
not be shortened.

---

## Troubleshooting

| Symptom | Cause | Fix |
|---------|-------|-----|
| Job fails with 403 on checkout | GitHub App not installed on target repo | Confirm P2: install App with `contents:write` + `pull-requests:write` |
| `Listed file not found: path/to/Dockerfile` | Matrix `files` has a wrong path | Update the matrix entry to the actual live file path |
| `Invalid or missing SHA256` | checksums.txt entry not found, or value not 64 hex chars | Check release assets contain `checksums.txt`; or supply overrides via dispatch inputs |
| Overall run failed, some repos got PRs | `fail-fast:false` — one job failed, others continued | Check Summary tab for the failed repo; the partial result is the intended behavior |
| "No changes detected" for a repo | That repo is already at the target version | Expected no-op; no action needed |
| PR already exists on branch `update-lpm-<version>` | Workflow re-triggered for same version | `peter-evans/create-pull-request` updates the existing PR; no duplicate |

---

## Related workflows

- `attach-artifact-release.yml` — generates `checksums.txt` (must succeed before this workflow runs)
- `publish-release.yml` — syncs the internal VERSION file after a release

---

## Future options (deferred)

If the consumer list grows large or needs to be shared across multiple workflows,
the inline matrix can be replaced with an external `consumers.json` file and a
`fromJSON(steps.load.outputs.matrix)` dynamic matrix. This is not needed at the
current scale of 3–4 repos.
