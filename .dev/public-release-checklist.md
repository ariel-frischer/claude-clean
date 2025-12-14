# Public Release Checklist

Checklist for preparing this repository for public release.

## Documentation

- [ ] **README.md** - Review for clarity, accuracy, and completeness
  - [ ] Installation instructions work on all platforms
  - [ ] Usage examples are up-to-date
  - [ ] Screenshots/GIFs demonstrating output (optional but nice)
  - [ ] Badges (build status, Go version, license, etc.)
- [ ] **CONTRIBUTING.md** - Review contribution guidelines
- [ ] **CHANGELOG.md** - Ensure all notable changes are documented
- [ ] **SECURITY.md** - Verify security policy and contact info
- [ ] **LICENSE** - Confirm license choice (currently MIT)
- [ ] **CLAUDE.md** - Decide if this should be public or moved to .dev/

## Code Quality

- [ ] All tests passing (`make test`)
- [ ] Code formatted (`make fmt`)
- [ ] No linting issues (`make vet`)
- [ ] No TODO/FIXME comments that shouldn't be public
- [ ] Remove any debug/test code
- [ ] Review error messages for clarity

## Security & Secrets

- [ ] No hardcoded credentials, API keys, or tokens
- [ ] No personal paths or usernames in code
- [ ] No sensitive data in mocks/ directory
- [ ] .gitignore covers all sensitive files
- [ ] Review git history for accidentally committed secrets
  - Run: `git log --all --full-history -S "password" -S "secret" -S "api_key" -S "token"`
  - Run: `git log --all --full-history -S "ANTHROPIC" -S "OPENAI" -S "Bearer"`
  - Consider using tools like `trufflehog` or `gitleaks` for thorough scanning

## Repository Hygiene

- [ ] **Git history audit**
  - [ ] No embarrassing/unprofessional commit messages
    - Run: `git log --oneline --all` and review
  - [ ] No WIP/fixup/squash commits that should have been cleaned up
  - [ ] No "oops" or "fix typo" chains that could be squashed
  - [ ] No commits with personal rants or inappropriate language
- [ ] **No dev artifacts in history**
  - [ ] Check for accidentally committed binaries: `git log --all --diff-filter=A --name-only | grep -E '\.(exe|bin|o|so|dylib)$'`
  - [ ] Check for large files: `git rev-list --objects --all | git cat-file --batch-check='%(objecttype) %(objectname) %(objectsize) %(rest)' | awk '/^blob/ {print $3, $4}' | sort -rn | head -20`
  - [ ] No .env files or config with real values in history
  - [ ] No node_modules, vendor dirs, or other dependency dumps
  - [ ] No personal notes, scratch files, or test outputs
- [ ] **If history needs cleaning**
  - Consider `git filter-repo` for removing sensitive files from all history
  - Or start fresh with a squashed initial commit (nuclear option)
- [ ] Consistent branch naming (main vs master)
- [ ] Remove unused branches
- [ ] Tags are semantic versioned (v1.0.0, v1.1.0, etc.)
- [ ] .gitattributes configured properly

## GitHub Setup

- [ ] **Repository settings**
  - [ ] Description and topics/tags set
  - [ ] Website URL (if applicable)
  - [ ] Social preview image (1280x640px recommended)

- [ ] **.github/ directory** (create if missing)
  - [ ] `ISSUE_TEMPLATE/bug_report.md`
  - [ ] `ISSUE_TEMPLATE/feature_request.md`
  - [ ] `PULL_REQUEST_TEMPLATE.md`
  - [ ] `CODEOWNERS` (optional)
  - [ ] `FUNDING.yml` (optional, for sponsors)

- [ ] **GitHub Actions CI/CD**
  - [ ] Build workflow (test on push/PR)
  - [ ] Release workflow (build binaries on tag)
  - [ ] Dependabot for dependency updates

## Community Standards

- [ ] **CODE_OF_CONDUCT.md** - Add if not present
- [ ] Issue labels configured (bug, enhancement, good first issue, etc.)
- [ ] Branch protection rules (require PR reviews, status checks)
- [ ] Enable Discussions (optional)

## Release Preparation

- [ ] Version number set correctly in code/Makefile
- [ ] Create initial release with binaries
- [ ] Verify binary checksums work
- [ ] Test installation instructions from scratch

## Pre-Release Final Checks

- [ ] Clone repo fresh and verify build works
- [ ] Run full test suite one more time
- [ ] Check all links in documentation work
- [ ] Review README one final time as a new user would
- [ ] Verify license headers in source files (if required)

## Post-Release

- [ ] Announce release (if applicable)
- [ ] Monitor issues for early feedback
- [ ] Set up notifications for new issues/PRs

---

## Current Status

### Already Complete
- [x] README.md exists with comprehensive documentation
- [x] LICENSE (MIT)
- [x] CONTRIBUTING.md
- [x] SECURITY.md
- [x] CHANGELOG.md
- [x] .gitignore
- [x] Makefile with build/test/fmt/vet targets
- [x] Tests exist and pass
- [x] .github/ directory with issue templates, PR template, workflows
- [x] CODE_OF_CONDUCT.md
- [x] CI/CD GitHub Actions (ci.yml, release.yml)
- [x] README badges (CI, Release, Go version, License)
- [x] Reviewed mocks/ - no sensitive content found
- [x] Moved CLAUDE.md to .dev/
- [x] dependabot.yml for automated dependency updates
- [x] Removed binary from git history (3MB claude-clean-output)

### Needs Attention
- [ ] Force push to remote (history was rewritten)
- [ ] Run full test suite one more time before release
- [ ] Clone repo fresh and verify build works
