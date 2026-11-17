# Release Process

Remotes: `gh` = GitHub (releases, CI), `origin` = GitLab (mirror, no pipelines).
Releases are published by GitHub Actions (`.github/workflows/release.yml`) when a `v*` tag reaches GitHub.
Release notes come from `CHANGELOG.yaml` via [chlog](https://github.com/ariel-frischer/chlog).

## Steps

```bash
# 1. Promote unreleased changelog entries and commit
chlog release 0.4.0 && chlog sync && chlog check
git add CHANGELOG.yaml CHANGELOG.md
git commit -m "chore(release): v0.4.0"
git push gh main && git push origin main

# 2. Tag; pushing the tag to GitHub triggers the release workflow
git tag -a v0.4.0 -m "Release v0.4.0"
git push gh v0.4.0 && git push origin v0.4.0

# 3. Verify
gh run watch "$(gh run list -L 1 --workflow Release --json databaseId -q '.[0].databaseId')" --exit-status
gh release view v0.4.0 --json body,assets -q '"assets=\(.assets|length) body=\(.body|length)"'
```

A good release has 7 assets (`cclean_<ver>_{linux,darwin,windows}_{amd64,arm64}` archives, zip on
Windows, plus `checksums.txt`) and a non-empty body matching `chlog extract <ver>`.

## Workflow

The release job checks out the tag, runs tests, installs chlog (pinned, `cmd/chlog`), runs
`chlog extract <ver> > /tmp/release-notes.md`, then `goreleaser release --release-notes`.
CI (`ci.yml`) runs `chlog validate`, `chlog check` and `goreleaser check` on every push, so release
tooling breakage shows up before tagging.

## When a release run fails

Fix the cause on `main`, then re-run the workflow for the existing tag:

```bash
gh workflow run Release --ref main -f tag=v0.4.0
```

The re-run uses the workflow file from `main` but the source and `.goreleaser.yaml` from the tag.
`release.replace_existing_artifacts` lets it overwrite assets of a partially published release.
Don't move or re-create a pushed tag: the Go module proxy may already have cached it.

If notes are missing on an otherwise good release:
`chlog extract 0.4.0 > /tmp/n.md && gh release edit v0.4.0 --notes-file /tmp/n.md`.

## Gotchas

- `chlog` is installed from `github.com/ariel-frischer/chlog/cmd/chlog`; the module root has no main package.
- Never set `changelog.disable: true` in `.goreleaser.yaml`: it skips goreleaser's changelog step,
  which also discards `--release-notes`, and the release publishes with an empty body (v0.3.0).
- `make patch` / `make release` tag and push to `origin` (GitLab) and run goreleaser locally. Use them only
  as a fallback when Actions is unavailable, and never together with a tag pushed to `gh` (both would publish).
