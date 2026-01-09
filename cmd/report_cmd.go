package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Report command flags
var (
	reportOutputFormat    string
	reportOutputFile      string
	reportIncludeUnmerged bool
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate branch analysis report without deleting",
	Long: `Generate a detailed analysis report of all branches in the repository.

This command analyzes your repository branches and generates a comprehensive
report showing which branches are safe to delete, which are local-only,
unmerged, or protected.

The report includes:
  - Safe to delete: Merged branches or branches with deleted remotes
  - Local-only: Branches that were never pushed to remote
  - Unmerged: Branches with unmerged changes (when -u flag is used)
  - Protected: Default branch and currently checked out branch

Output formats available:
  - text: Human-readable formatted report (default)
  - json: Machine-readable JSON format
  - csv: Spreadsheet-compatible CSV format`,
	Example: `  # Generate a text report to stdout
  git-gone report

  # Generate a JSON report
  git-gone report --output json

  # Save report to a file
  git-gone report --file branches-report.txt

  # Generate CSV report including unmerged branches
  git-gone report -u --output csv --file report.csv`,
	Run: func(cmd *cobra.Command, args []string) {
		runReport()
	},
}

func init() {
	reportCmd.Flags().StringVarP(&reportOutputFormat, "output", "o", "text", "Report output format (text, json, csv)")
	reportCmd.Flags().StringVar(&reportOutputFile, "file", "", "Write report to file instead of stdout")
	reportCmd.Flags().BoolVarP(&reportIncludeUnmerged, "unmerged", "u", false, "Include unmerged branches in the report")
}

func runReport() {
	// Check if we're in a git repository
	if err := checkGitRepository(); err != nil {
		fmt.Println("❌ Not in a git repository")
		os.Exit(1)
	}

	fmt.Println("🔄 Updating remote references...")
	if err := updateRemoteRefs(); err != nil {
		fmt.Printf("⚠️  Warning: Failed to update remote refs: %v\n", err)
	}

	// Set the global includeUnmerged flag for analyzeBranches to use
	includeUnmerged = reportIncludeUnmerged

	fmt.Println("📊 Analyzing branches...")
	report := analyzeBranches()
	outputReport(report, reportOutputFormat, reportOutputFile)
}
