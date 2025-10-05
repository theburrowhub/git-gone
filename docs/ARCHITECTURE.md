# git-gone Architecture

## Project Structure

```
git-gone/
├── cmd/                    # Command package
│   ├── root.go            # Root command setup
│   ├── version.go         # Version subcommand
│   └── branches.go        # Branches subcommand (cleanup logic)
├── main.go                # Entry point
├── go.mod                 # Go module definition
└── install.sh             # Installation script
```

## Commands Hierarchy

### Root Command
- **Usage**: `git-gone`
- **Default behavior**: Runs `branches` subcommand when no subcommand is specified

### Subcommands

#### `version`
- **Usage**: `git-gone version`
- **Purpose**: Display version information
- **File**: `cmd/version.go`

#### `branches` (default)
- **Usage**: `git-gone branches` or just `git-gone`
- **Purpose**: Interactive branch cleanup
- **File**: `cmd/branches.go`
- **Features**:
  - Updates remote references
  - Finds merged branches
  - Finds branches with deleted remotes
  - Interactive selection with fzf
  - Confirmation before deletion
  - Safe and force delete

## Adding New Subcommands

To add a new subcommand:

1. Create a new file in `cmd/` (e.g., `cmd/mycommand.go`)
2. Define your command:

```go
package cmd

import "github.com/spf13/cobra"

var myCmd = &cobra.Command{
    Use:   "mycommand",
    Short: "Short description",
    Long:  `Long description`,
    Run: func(cmd *cobra.Command, args []string) {
        // Your logic here
    },
}
```

3. Register it in `cmd/root.go` in the `init()` function:

```go
func init() {
    rootCmd.AddCommand(myCmd)
    // ... other commands
}
```

## Command Execution Flow

1. `main.go` calls `cmd.Execute()`
2. Cobra parses arguments
3. If no subcommand is specified:
   - Root command's `Run` function executes
   - Which calls `branchesCmd.Run()` (default behavior)
4. If subcommand is specified:
   - Cobra routes to the appropriate subcommand

## Build Information

Version information is injected at build time using ldflags:

```bash
go build -ldflags "\
  -X git-gone/cmd.Version=1.0.0 \
  -X git-gone/cmd.CommitHash=abc123 \
  -X git-gone/cmd.BuildTime=2025-01-01T00:00:00Z"
```

This information is shared across all commands through the `cmd` package.
