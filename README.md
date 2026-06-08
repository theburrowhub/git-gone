# git-gone

A Git plugin for cleaning up git repositories interactively. Removes merged branches, stale tags, and other cleanup tasks.

Works as a native Git extension: `git gone`

## Features

- 🔄 Updates all local references with origin
- 🎯 Detects branches ready for deletion using multiple methods:
  - Branches merged into the default branch (traditional merge)
  - Branches with deleted remotes (squash/rebase merges)
  - Platform-independent operation (works regardless of system language)
- 🏷️ Detects stale tags (local tags not on remote)
- 🔍 Interactive multi-selection using fuzzy finder
- ✅ Safe deletion with confirmation prompt
- ⚠️ Extra safety for dangerous operations (unmerged branches require typing "DELETE")
- 📊 Clear status indicators throughout the process
- 📋 Generate detailed branch analysis reports (text/JSON/CSV)

## Installation

### Homebrew (macOS/Linux)

```bash
brew install theburrowhub/tap/git-gone
```

Or first tap the repository and then install:

```bash
brew tap theburrowhub/tap
brew install git-gone
```

### Quick install

```bash
# Download and run the installation script
# Will download pre-built binary if available, or build from source automatically
curl -sSL https://raw.githubusercontent.com/theburrowhub/git-gone/main/install.sh | bash
```

**Note**: The installation script will automatically:
- Try to download pre-built releases if available
- Fall back to building from source if no releases exist (requires Git and Go 1.24+)
- Install the binary to `~/.local/bin` by default

### From source

```bash
# Clone the repository
git clone https://github.com/theburrowhub/git-gone.git
cd git-gone

# Use the installation script for local build
./install.sh

# Or build manually
go build -o git-gone
sudo mv git-gone /usr/local/bin/
```

### Download binary

Download the latest release from the [releases page](https://github.com/theburrowhub/git-gone/releases) for your platform.

## Usage

### Branch Cleanup

Navigate to any git repository and run:

```bash
git gone
# or
git-gone
# or explicitly
git-gone branches
```

#### Branch Command Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--force` | `-f` | Skip confirmation prompt (doesn't apply to unmerged branches) |
| `--all` | `-a` | Select all candidate branches without interactive selection |
| `--unmerged` | `-u` | Include unmerged branches in the list (marked with `(!)`) |

**Note**: `-a` and `-f` are incompatible. The `-a` flag is designed for review before deletion.

#### Examples

```bash
# Interactive cleanup (default)
git gone

# Select all candidates, then confirm
git gone -a

# Force delete safe branches (merged/gone), skip confirmation
git gone -f

# Include unmerged branches (requires typing DELETE to confirm)
git gone -u

# Include all branches including unmerged, review before confirmation
git gone -a -u
```

### Tag Cleanup

```bash
# List stale tags (local tags not on remote)
git-gone tags list

# List ALL local tags (not just stale)
git-gone tags list --no-stale
git-gone tags list -n

# Clean stale tags interactively
git-gone tags clean

# Clean all stale tags with confirmation
git-gone tags clean --all

# Clean all stale tags without confirmation
git-gone tags clean --all --force

# Clean ANY local tag (not just stale)
git-gone tags clean --no-stale
git-gone tags clean -n
```

### Other Commands

```bash
# Show version information
git-gone version

# Self-update to latest release
git-gone self-update

# Show help
git-gone --help
git-gone branches --help
git-gone tags --help
```

## Report Mode

Generate a detailed analysis report without deleting any branches:

### Basic Report

```bash
git-gone report
```

### Output Formats

```bash
# Text format (default)
git-gone report --output text

# JSON format (for scripting/automation)
git-gone report --output json
git-gone report -o json

# CSV format (for spreadsheets)
git-gone report --output csv
```

### Save to File

```bash
git-gone report --output json --file report.json
```

### Include Unmerged Branches

```bash
git-gone report --unmerged
git-gone report -u
```

### Report Categories

The report classifies branches into:

- **Safe to Delete**: Merged branches or branches with deleted remotes
- **Local-only**: Merged but never pushed to remote (review recommended)
- **Unmerged**: Not merged, requires `--unmerged` flag to include
- **Protected**: Default branch or currently checked out

### Example Report Output

```
============================================================
              GIT-GONE BRANCH ANALYSIS REPORT
============================================================
Repository: /path/to/repo
Date: 2026-01-09 14:30:00
Default Branch: main
Current Branch: feature/ui

------------------------------------------------------------
SAFE TO DELETE (2 branches)
------------------------------------------------------------
  * feature/old-login
    Method: merged | Reason: Merged into main
    Remote: gone | Last commit: 2025-12-15

------------------------------------------------------------
LOCAL-ONLY (1 branch) - Merged but never pushed
------------------------------------------------------------
  * temp/local-experiment
    Method: merged | Reason: Merged but never pushed to remote
    Remote: local_only | Last commit: 2025-12-10

------------------------------------------------------------
PROTECTED (2 branches)
------------------------------------------------------------
  * main
    Reason: Default branch

  * feature/ui
    Reason: Currently checked out

============================================================
SUMMARY: 2 safe | 1 local-only | 0 unmerged | 2 protected
============================================================
```

## Command Structure

git-gone uses a subcommand structure powered by Cobra:

```
git-gone
├── branches              # Default command (branch cleanup)
│   ├── --all, -a        # Select all candidates
│   ├── --force, -f      # Skip confirmation
│   └── --unmerged, -u   # Include unmerged branches
├── tags                  # Tag management
│   ├── list             # List stale tags
│   │   └── --no-stale, -n  # List ALL local tags
│   └── clean            # Clean stale tags
│       ├── --all, -a    # Select all stale tags
│       ├── --force, -f  # Skip confirmation
│       └── --no-stale, -n  # Include ALL local tags
├── report               # Generate analysis report (no deletion)
│   ├── --output, -o     # Output format (text/json/csv)
│   ├── --file           # Save report to file
│   └── --unmerged, -u   # Include unmerged branches
├── version              # Show version info
├── self-update          # Update to latest release
└── help                 # Auto-generated help
```

## Interactive Controls

- **↑/↓**: Navigate through the list
- **Tab/Space**: Toggle selection of current item
- **Enter**: Confirm selection and proceed
- **Esc**: Cancel operation
- **Type**: Filter items by name

## Branch Categories

| Indicator | Description | Risk Level |
|-----------|-------------|------------|
| (none) | Merged branch or gone remote | Safe |
| `(!)` | Unmerged branch | Dangerous - requires typing "DELETE" |

When using `-u` flag, a legend is shown:
```
   (!) Unmerged
```

## Safety Features

- Never deletes the default branch (main/master/develop)
- Never deletes the currently checked out branch
- Merged branches: Simple y/N confirmation
- Unmerged branches (`-u` flag): Requires typing "DELETE" to confirm
- Shows per-item deletion success/failure
- Attempts safe deletion first, falls back to force only if needed
- Remote deletion only for unmerged branches (when applicable)

## Requirements

- Go 1.24 or higher (for building from source)
- Git installed and configured
- Terminal with UTF-8 support (for emoji indicators)

## Dependencies

- [github.com/spf13/cobra](https://github.com/spf13/cobra) - CLI framework
- [github.com/koki-develop/go-fzf](https://github.com/koki-develop/go-fzf) - Fuzzy finder
- [github.com/charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [github.com/charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss) - TUI styling

## How it Works as a Git Plugin

When you install `git-gone`, the binary is placed in your `$PATH`. Git automatically recognizes executables named `git-<command>` as git subcommands, allowing you to use it as `git gone`.

## Development

### Running tests

```bash
go test ./tests/... -v
```

### Creating a new release

Releases are automatically created when you push a tag starting with `v`:

```bash
# Create a new tag
git tag v0.3.0 -m "Release version 0.3.0"

# Push the tag to GitHub
git push origin v0.3.0
```

GitHub Actions will automatically:
1. Build binaries for all platforms (Linux, macOS, Windows)
2. Create a GitHub release with the binaries
3. Generate checksums for all files

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
