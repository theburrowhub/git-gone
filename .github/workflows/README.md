# GitHub Actions Workflows

This directory contains the CI/CD workflows for git-gone.

## Workflows

### 1. Pull Request Checks (`pr.yml`)
**Trigger**: Pull requests to `main` branch

**Actions performed**:
- Go format check (`gofmt`)
- Go vet analysis
- Linting with `golangci-lint`
- Run all tests with race detection
- Coverage reporting to Codecov
- Check `go mod tidy` status

**Purpose**: Validate code quality before merging PRs.

---

### 2. Release and Build (`release.yml`)
**Trigger**: Push to `main` branch

**Actions performed**:
1. **Version Bump**:
   - Calculate next semantic version using conventional commits
   - Update version in code
   - Update `CHANGELOG.md` with new release notes
   - Commit version bump changes
   - Create and push git tag

2. **Build Binaries**:
   - Build for multiple platforms:
     - Linux (amd64, arm64)
     - macOS (amd64, arm64)
   - Inject version information with ldflags
   - Create compressed archives

3. **Create Release**:
   - Generate SHA256 checksums
   - Create GitHub release with all artifacts
   - Extract changelog for release notes

**Requirements**:
- Uses conventional commit messages for version calculation
- Requires `svu` (Semantic Version Util)
- Automatically skips CI if commit message contains `[skip ci]`

---

## Conventional Commits

The version bumping uses conventional commits to determine the next version:

- `fix:` → patch version bump (1.0.0 → 1.0.1)
- `feat:` → minor version bump (1.0.0 → 1.1.0)
- `BREAKING CHANGE:` → major version bump (1.0.0 → 2.0.0)

## Setup Requirements

1. **Codecov Integration** (optional):
   - Add `CODECOV_TOKEN` to repository secrets for private repos

2. **Permissions**:
   - Workflows need write permissions for:
     - Contents (to push commits and tags)
     - Pull requests (to comment on PRs)
     - Actions (to upload artifacts)

3. **Branch Protection** (recommended):
   - Protect `main` branch
   - Require PR reviews
   - Require status checks to pass

## Workflow Summary

```
Pull Request → pr.yml (tests + lint)
     ↓
   Merge to main
     ↓
release.yml (version bump + build + release)
```

This ensures every PR is tested, and every merge to main automatically creates a release.