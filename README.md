# git-gone

A Git plugin for cleaning up merged git branches interactively using fuzzy finder.
Works as a native Git extension: `git gone`

## Features

- 🔄 Updates all local references with origin
- 🎯 Detects branches ready for deletion using multiple methods:
  - Branches merged into the default branch (traditional merge)
  - Branches with deleted remotes (squash/rebase merges)
  - Platform-independent operation (works regardless of system language)
- 🔍 Interactive multi-selection using go-fzf
- ✅ Safe deletion with confirmation prompt
- 📊 Clear status indicators throughout the process
- 📋 Generate detailed branch analysis reports (text/JSON/CSV)

## Installation

### Quick install (recommended)

```bash
# Download and run the installation script
# Will download pre-built binary if available, or build from source automatically
curl -sSL https://raw.githubusercontent.com/theburrowhub/git-gone/main/install.sh | bash
```

**Note**: The installation script will automatically:
- Try to download pre-built releases if available
- Fall back to building from source if no releases exist (requires Git and Go 1.19+)
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

### Basic Usage

Navigate to any git repository and run:

```bash
git gone
# or
git-gone
# or explicitly
git-gone branches
```

### Available Commands

- **`git-gone` or `git-gone branches`**: Clean up merged branches (default)
- **`git-gone --report`**: Generate branch analysis report without deleting
- **`git-gone version`**: Show version information
- **`git-gone help`**: Show help information

### What the tool does:
1. Update all remote references (`git fetch --all --prune`)
2. Identify the default branch (main/master)
3. Find deletable branches using two methods:
   - Branches merged into the default branch (traditional merges)
   - Branches whose remote tracking branch is gone (squash/rebase merges)
4. Present an interactive list where you can:
   - Use arrow keys to navigate
   - Press `Tab` or `Space` to select/deselect branches
   - Press `Enter` to confirm selection
   - Press `Esc` to cancel
   - Type to filter branches by name
5. Ask for confirmation before deleting selected branches
6. Delete the selected branches safely

## Report Mode

Generate a detailed analysis report without deleting any branches:

### Basic Report

```bash
git-gone --report
git-gone -r
```

### Output Formats

```bash
# Text format (default)
git-gone --report --output text

# JSON format (for scripting/automation)
git-gone --report --output json

# CSV format (for spreadsheets)
git-gone --report --output csv
```

### Save to File

```bash
git-gone --report --output json --file report.json
```

### Include Unmerged Branches

```bash
git-gone --report --unmerged
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

```bash
# Show general help
git-gone --help
git-gone -h

# Show version
git-gone version

# Clean branches (default command)
git-gone
git-gone branches

# Get help for specific command
git-gone branches --help

# Generate analysis report
git-gone --report
git-gone -r

# Report in JSON format
git-gone --report --output json
git-gone -r -o json
```

**Note**: When using as a Git plugin (`git gone`), commands work the same way:
- `git gone` - runs branch cleanup
- `git gone version` - shows version
- `git gone -h` - shows help

## Interactive Controls

- **↑/↓**: Navigate through the list
- **Tab/Space**: Toggle selection of current branch
- **Enter**: Confirm selection and proceed
- **Esc**: Cancel operation
- **Type**: Filter branches by name

## Safety Features

- Never deletes the default branch (main/master)
- Never deletes the current branch
- Shows only branches that have been merged
- Requires explicit confirmation before deletion
- Attempts safe deletion first (`git branch -d`)
- Falls back to force deletion only if necessary

## Requirements

- Go 1.19 or higher
- Git installed and configured
- Terminal with UTF-8 support (for emoji indicators)

## Dependencies

- [github.com/koki-develop/go-fzf](https://github.com/koki-develop/go-fzf) - Fuzzy finder library

## How it Works as a Git Plugin

When you install `git-gone`, the binary is placed in your `$PATH`. Git automatically recognizes executables named `git-<command>` as git subcommands, allowing you to use it as `git gone`.

## Development

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
