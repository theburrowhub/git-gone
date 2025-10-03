## [Unreleased]

### Added
- GitHub Actions workflows for CI/CD
- Automatic release generation with correct binary names (git-gone)
- Support for Windows builds in CI

### Fixed
- Installation script now handles both legacy (gitcleaner) and new (git-gone) binary names
- Improved error handling in installation script with automatic fallback to source compilation

## [0.2.0] - 2025-10-03

- **BREAKING CHANGE**: Renamed project from `gitcleaner` to `git-gone`
- feat: Now works as a Git plugin - can be invoked with `git gone`
- feat: Added support for short flags `-h` and `-v` (required for Git plugin usage)
- refactor: Updated all references throughout the codebase
- docs: Updated documentation to reflect new name and usage
- fix: Use `-h` flag with `git gone` command as Git intercepts `--help`

## [0.1.1] - 2025-10-02

- Merge branch 'main' of github.com:theburrowhub/gitcleaner
- fix: correct Go version in go.mod and format code

## [0.1.0] - 2025-10-02

- Initial commit: GitCleaner application
- feat: add support for detecting branches with deleted remotes (gone/desaparecido)
- docs: update README to document support for squash/rebase merge detection
- feat: add LC_ALL=C to all git commands for platform-independent output
- docs: update README to reflect platform-independent operation
- feat: add installation script and GitHub Actions workflows

