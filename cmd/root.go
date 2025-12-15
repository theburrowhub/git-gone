package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Version information - set during build with ldflags
var (
	Version    = "0.2.0"
	CommitHash = "unknown"
	BuildTime  = "unknown"
)

// Global flags
var (
	forceDelete     bool
	selectAll       bool
	includeUnmerged bool
)

var rootCmd = &cobra.Command{
	Use:   "git-gone",
	Short: "Clean up merged git branches interactively",
	Long: `git-gone - Clean up merged git branches interactively

git-gone helps you clean up local git branches that have been merged
or whose remote tracking branches have been deleted.`,
	Run: func(cmd *cobra.Command, args []string) {
		// By default, run the branches command when no subcommand is provided
		branchesCmd.Run(cmd, args)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Add persistent flags that are available to root and all subcommands
	rootCmd.PersistentFlags().BoolVarP(&forceDelete, "force", "f", false, "Skip confirmation prompt and delete selected branches immediately")
	rootCmd.PersistentFlags().BoolVarP(&selectAll, "all", "a", false, "Select all candidate branches without interactive selection (incompatible with -f)")
	rootCmd.PersistentFlags().BoolVarP(&includeUnmerged, "unmerged", "u", false, "Include unmerged branches in the list (marked with ⚠️, always requires confirmation)")

	// Add subcommands
	rootCmd.AddCommand(branchesCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(selfUpdateCmd)
}
