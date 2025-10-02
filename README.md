# GitCleaner

A Go application for cleaning up merged git branches interactively using fuzzy finder.

## Features

- 🔄 Updates all local references with origin
- 🎯 Detects branches ready for deletion using multiple methods:
  - Branches merged into the default branch (traditional merge)
  - Branches with deleted remotes (squash/rebase merges)
  - Platform-independent operation (works regardless of system language)
- 🔍 Interactive multi-selection using go-fzf
- ✅ Safe deletion with confirmation prompt
- 📊 Clear status indicators throughout the process

## Installation

### Quick install (recommended)

```bash
# Download and run the installation script
curl -sSL https://raw.githubusercontent.com/YOUR_USERNAME/gitcleaner/main/install.sh | bash
```

### From source

```bash
# Clone the repository
git clone https://github.com/YOUR_USERNAME/gitcleaner.git
cd gitcleaner

# Use the installation script for local build
./install.sh

# Or build manually
go build -o gitcleaner
sudo mv gitcleaner /usr/local/bin/
```

### Download binary

Download the latest release from the [releases page](https://github.com/YOUR_USERNAME/gitcleaner/releases) for your platform.

## Usage

Navigate to any git repository and run:

```bash
gitcleaner
```

The application will:
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

## License

MIT
