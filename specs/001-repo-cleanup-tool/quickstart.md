# Quickstart: git-gone

## Installation

### Quick Install (Recommended)

```bash
curl -sSL https://raw.githubusercontent.com/theburrowhub/git-gone/main/install.sh | bash
```

### From Source

```bash
git clone https://github.com/theburrowhub/git-gone.git
cd git-gone
go build -o git-gone
sudo mv git-gone /usr/local/bin/
```

## Basic Usage

### Clean Merged Branches

Navigate to any git repository and run:

```bash
# As git plugin
git gone

# Or standalone
git-gone
git-gone branches
```

**What happens**:
1. Updates all remote references
2. Identifies branches safe to delete (merged or with deleted remotes)
3. Shows interactive selector - use arrow keys and Tab to select
4. Press Enter to confirm selection
5. Confirm deletion with y/N prompt
6. Branches are deleted with a summary

### Clean Stale Tags

```bash
# List stale tags
git-gone tags list

# Clean stale tags interactively
git-gone tags clean
```

## Common Workflows

### Daily Cleanup

After merging a PR, clean up your local branches:

```bash
git checkout main
git pull
git gone
```

### Batch Cleanup (Skip Selection)

Delete all merged branches without interactive selection:

```bash
git gone -a
```

### Force Mode (Skip Confirmation)

For scripting or when you're confident:

```bash
git gone -f
```

### Include Unmerged Branches

When you want to clean abandoned work (dangerous!):

```bash
git gone -u
```

Note: Unmerged branches are marked with `(!)` and require typing `DELETE` to confirm.

## Interactive Controls

| Key | Action |
|-----|--------|
| ↑/↓ | Navigate through list |
| Tab | Toggle selection |
| Enter | Confirm selection |
| Esc | Cancel operation |
| Type | Filter branches by name |

## Command Reference

```bash
# Branch cleanup (default)
git gone                    # Interactive cleanup
git gone -a                 # Select all, still confirm
git gone -f                 # Skip confirmation (after selection)
git gone -u                 # Include unmerged branches

# Tag management
git gone tags list          # Show stale tags
git gone tags clean         # Interactive tag cleanup
git gone tags clean -a -f   # Delete all stale tags

# Utility
git gone version            # Show version
git gone self-update        # Update to latest
git gone self-update -c     # Check for updates only
git gone help               # Show help
```

## Safety Features

- **Never deletes default branch** (main/master)
- **Never deletes current branch**
- **Confirmation required** before any deletion
- **Extra confirmation** for unmerged branches (type "DELETE")
- **Recoverable** via `git reflog` for 30 days

## Troubleshooting

### "Not in a git repository"

Make sure you're inside a git repository:

```bash
cd your-project
git status  # Should work
git gone    # Now this works
```

### "Failed to update remote refs"

Check your network connection. The tool will continue with local data but may miss recently deleted remotes.

### "No branches to delete"

Your repository is clean! All branches are either:
- The default branch
- The current branch
- Not merged and don't have deleted remotes

## Next Steps

- Run `git gone help` for complete command reference
- Run `git gone completion bash >> ~/.bashrc` for shell completion
- Check for updates with `git gone self-update -c`

