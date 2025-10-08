# GitHub Actions Workflows

This directory contains the CI/CD workflows for git-gone.

## Workflows

### `pr.yml` - Pull Request Checks

Runs on every pull request to `main` branch.

**What it does:**
- Format check (`gofmt`)
- Go vet
- Linting (`golangci-lint`)
- Unit tests with race detection
- Coverage report (uploaded to Codecov)
- Module tidiness check

### `release.yml` - Release and Build

Triggers on tag pushes (e.g., `v1.2.3`).

**What it does:**
- Uses [GoReleaser](https://goreleaser.com) to:
  - Build binaries for multiple platforms (Linux, macOS on amd64 and arm64)
  - Create archives (tar.gz and zip)
  - Generate checksums
  - Create GitHub release with changelog
  - Upload all artifacts

## Making a Release

To create a new release:

1. Make sure you're on the `main` branch and it's up to date:
   ```bash
   git checkout main
   git pull
   ```

2. Create and push a new tag following semantic versioning:
   ```bash
   # For a patch release (bug fixes)
   git tag -a v0.4.2 -m "Release v0.4.2"
   
   # For a minor release (new features, backwards compatible)
   git tag -a v0.5.0 -m "Release v0.5.0"
   
   # For a major release (breaking changes)
   git tag -a v1.0.0 -m "Release v1.0.0"
   
   git push origin v0.4.2  # Replace with your version
   ```

3. The GitHub Action will automatically:
   - Build binaries for all platforms
   - Create a GitHub release
   - Upload all artifacts
   - Generate and attach checksums

4. The release will be available at: https://github.com/theburrowhub/gitcleaner/releases

## Testing Locally

You can test the release process locally without pushing:

```bash
# Test building for all platforms (snapshot mode, no publishing)
goreleaser release --snapshot --clean

# Test building only for your current platform
goreleaser build --snapshot --clean --single-target

# Validate configuration
goreleaser check
```

## Configuration

The release process is configured in `.goreleaser.yaml` at the root of the repository.

Key features:
- Multi-platform builds (Linux, macOS)
- Multi-architecture (amd64, arm64)
- Automatic changelog generation from commits
- Checksum generation for security
- Archive creation (tar.gz and zip)

## Conventional Commits

For better changelogs, use conventional commit messages:

- `feat:` for new features
- `fix:` for bug fixes
- `docs:` for documentation changes
- `refactor:` for code refactoring
- `test:` for test changes
- `ci:` for CI/CD changes
- `chore:` for other changes

Example:
```bash
git commit -m "feat: add interactive mode for branch selection"
git commit -m "fix: handle empty git repositories correctly"
```

These will be automatically categorized in the release notes.
