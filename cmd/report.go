package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// BranchAnalysis contains detailed information about a single branch
type BranchAnalysis struct {
	Name         string `json:"name"`
	Status       string `json:"status"`        // safe_to_delete, local_only, unmerged, protected
	DeleteMethod string `json:"delete_method"` // merged, gone_remote, force, or empty for protected
	Reason       string `json:"reason"`        // Human-readable explanation
	RemoteStatus string `json:"remote_status"` // exists, gone, local_only
	LastCommit   string `json:"last_commit"`   // Date of last commit
}

// ReportSummary contains aggregated counts for the report
type ReportSummary struct {
	SafeCount       int `json:"safe_to_delete_count"`
	LocalOnlyCount  int `json:"local_only_count"`
	UnmergedCount   int `json:"unmerged_count"`
	ProtectedCount  int `json:"protected_count"`
	MergedCount     int `json:"merged_count"`
	GoneRemoteCount int `json:"gone_remote_count"`
}

// AnalysisReport contains the complete branch analysis
type AnalysisReport struct {
	Repository    string           `json:"repository"`
	AnalysisDate  string           `json:"analysis_date"`
	DefaultBranch string           `json:"default_branch"`
	CurrentBranch string           `json:"current_branch"`
	TotalBranches int              `json:"total_branches"`
	SafeToDelete  []BranchAnalysis `json:"safe_to_delete"`
	LocalOnly     []BranchAnalysis `json:"local_only"`
	Unmerged      []BranchAnalysis `json:"unmerged"`
	Protected     []BranchAnalysis `json:"protected"`
	Summary       ReportSummary    `json:"summary"`
}

// getBranchRemoteStatus determines if a branch has a remote: exists, gone, or local_only
func getBranchRemoteStatus(branch string) string {
	// Check if branch has an upstream configured
	cmd := exec.Command("git", "config", "--get", fmt.Sprintf("branch.%s.remote", branch))
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	output, err := cmd.Output()

	if err != nil || strings.TrimSpace(string(output)) == "" {
		// No upstream configured - local only branch
		return "local_only"
	}

	// Has upstream, check if it's gone
	cmd = exec.Command("git", "branch", "--format", "%(upstream:track)", "--list", branch)
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	output, err = cmd.Output()

	if err != nil {
		return "local_only"
	}

	trackStatus := strings.TrimSpace(string(output))
	if strings.Contains(trackStatus, "[gone]") {
		return "gone"
	}

	return "exists"
}

// getBranchLastCommit returns the date of the last commit on a branch
func getBranchLastCommit(branch string) string {
	cmd := exec.Command("git", "log", "-1", "--format=%ci", branch)
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	output, err := cmd.Output()

	if err != nil {
		return "unknown"
	}

	// Parse and format the date nicely
	dateStr := strings.TrimSpace(string(output))
	if len(dateStr) >= 10 {
		return dateStr[:10] // Return just YYYY-MM-DD
	}

	return dateStr
}

// getRepositoryPath returns the root path of the current git repository
func getRepositoryPath() string {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	output, err := cmd.Output()

	if err != nil {
		return "unknown"
	}

	return strings.TrimSpace(string(output))
}

// analyzeBranches collects and classifies all branches in the repository
func analyzeBranches(includeUnmerged bool) *AnalysisReport {
	report := &AnalysisReport{
		Repository:   getRepositoryPath(),
		AnalysisDate: time.Now().Format("2006-01-02 15:04:05"),
		SafeToDelete: []BranchAnalysis{},
		LocalOnly:    []BranchAnalysis{},
		Unmerged:     []BranchAnalysis{},
		Protected:    []BranchAnalysis{},
	}

	// Get default branch
	defaultBranch, err := getDefaultBranch()
	if err != nil {
		defaultBranch = "main"
	}
	report.DefaultBranch = defaultBranch

	// Get current branch
	currentBranch, err := getCurrentBranch()
	if err != nil {
		currentBranch = ""
	}
	report.CurrentBranch = currentBranch

	// Get all local branches
	allBranches, err := getAllLocalBranches()
	if err != nil {
		return report
	}
	report.TotalBranches = len(allBranches)

	// Get gone branches (remote deleted)
	goneBranches, _ := getGoneBranches()
	goneBranchesMap := make(map[string]bool)
	for _, b := range goneBranches {
		goneBranchesMap[strings.TrimSpace(b)] = true
	}

	// Get merged branches
	mergedBranches, _ := getMergedBranches(defaultBranch)
	mergedBranchesMap := make(map[string]bool)
	for _, b := range mergedBranches {
		mergedBranchesMap[strings.TrimSpace(b)] = true
	}

	// Classify each branch
	for _, branch := range allBranches {
		branch = strings.TrimSpace(branch)
		if branch == "" {
			continue
		}

		analysis := BranchAnalysis{
			Name:         branch,
			RemoteStatus: getBranchRemoteStatus(branch),
			LastCommit:   getBranchLastCommit(branch),
		}

		// Check if protected (default or current branch)
		if branch == defaultBranch {
			analysis.Status = "protected"
			analysis.DeleteMethod = ""
			analysis.Reason = "Default branch"
			report.Protected = append(report.Protected, analysis)
			report.Summary.ProtectedCount++
			continue
		}

		if branch == currentBranch {
			analysis.Status = "protected"
			analysis.DeleteMethod = ""
			analysis.Reason = "Currently checked out"
			report.Protected = append(report.Protected, analysis)
			report.Summary.ProtectedCount++
			continue
		}

		// Check if gone (remote deleted)
		if goneBranchesMap[branch] {
			analysis.Status = "safe_to_delete"
			analysis.DeleteMethod = "gone_remote"
			analysis.Reason = "Remote tracking branch deleted"
			analysis.RemoteStatus = "gone"
			report.SafeToDelete = append(report.SafeToDelete, analysis)
			report.Summary.SafeCount++
			report.Summary.GoneRemoteCount++
			continue
		}

		// Check if merged
		if mergedBranchesMap[branch] {
			if analysis.RemoteStatus == "local_only" {
				// Merged but never pushed - separate category
				analysis.Status = "local_only"
				analysis.DeleteMethod = "merged"
				analysis.Reason = "Merged but never pushed to remote (local-only)"
				report.LocalOnly = append(report.LocalOnly, analysis)
				report.Summary.LocalOnlyCount++
				report.Summary.MergedCount++
			} else {
				// Merged with remote
				analysis.Status = "safe_to_delete"
				analysis.DeleteMethod = "merged"
				analysis.Reason = fmt.Sprintf("Merged into %s", defaultBranch)
				report.SafeToDelete = append(report.SafeToDelete, analysis)
				report.Summary.SafeCount++
				report.Summary.MergedCount++
			}
			continue
		}

		// Not merged - only include if --unmerged flag is set
		if includeUnmerged {
			analysis.Status = "unmerged"
			analysis.DeleteMethod = "force"
			analysis.Reason = "Not merged, requires force delete"
			report.Unmerged = append(report.Unmerged, analysis)
			report.Summary.UnmergedCount++
		}
	}

	return report
}

// generateTextReport creates a human-readable text report
func generateTextReport(report *AnalysisReport) string {
	var sb strings.Builder

	sb.WriteString("============================================================\n")
	sb.WriteString("              GIT-GONE BRANCH ANALYSIS REPORT\n")
	sb.WriteString("============================================================\n")
	sb.WriteString(fmt.Sprintf("Repository: %s\n", report.Repository))
	sb.WriteString(fmt.Sprintf("Date: %s\n", report.AnalysisDate))
	sb.WriteString(fmt.Sprintf("Default Branch: %s\n", report.DefaultBranch))
	sb.WriteString(fmt.Sprintf("Current Branch: %s\n", report.CurrentBranch))
	sb.WriteString("\n")

	// Safe to Delete
	if len(report.SafeToDelete) > 0 {
		sb.WriteString("------------------------------------------------------------\n")
		sb.WriteString(fmt.Sprintf("SAFE TO DELETE (%d branches)\n", len(report.SafeToDelete)))
		sb.WriteString("------------------------------------------------------------\n")
		for _, branch := range report.SafeToDelete {
			sb.WriteString(fmt.Sprintf("  * %s\n", branch.Name))
			sb.WriteString(fmt.Sprintf("    Method: %s | Reason: %s\n", branch.DeleteMethod, branch.Reason))
			sb.WriteString(fmt.Sprintf("    Remote: %s | Last commit: %s\n", branch.RemoteStatus, branch.LastCommit))
			sb.WriteString("\n")
		}
	}

	// Local Only
	if len(report.LocalOnly) > 0 {
		sb.WriteString("------------------------------------------------------------\n")
		sb.WriteString(fmt.Sprintf("LOCAL-ONLY (%d branches) - Merged but never pushed\n", len(report.LocalOnly)))
		sb.WriteString("------------------------------------------------------------\n")
		for _, branch := range report.LocalOnly {
			sb.WriteString(fmt.Sprintf("  * %s\n", branch.Name))
			sb.WriteString(fmt.Sprintf("    Method: %s | Reason: %s\n", branch.DeleteMethod, branch.Reason))
			sb.WriteString(fmt.Sprintf("    Remote: %s | Last commit: %s\n", branch.RemoteStatus, branch.LastCommit))
			sb.WriteString("\n")
		}
	}

	// Unmerged
	if len(report.Unmerged) > 0 {
		sb.WriteString("------------------------------------------------------------\n")
		sb.WriteString(fmt.Sprintf("UNMERGED (%d branches)\n", len(report.Unmerged)))
		sb.WriteString("------------------------------------------------------------\n")
		for _, branch := range report.Unmerged {
			sb.WriteString(fmt.Sprintf("  * %s\n", branch.Name))
			sb.WriteString(fmt.Sprintf("    Method: %s | Reason: %s\n", branch.DeleteMethod, branch.Reason))
			sb.WriteString(fmt.Sprintf("    Remote: %s | Last commit: %s\n", branch.RemoteStatus, branch.LastCommit))
			sb.WriteString("\n")
		}
	}

	// Protected
	if len(report.Protected) > 0 {
		sb.WriteString("------------------------------------------------------------\n")
		sb.WriteString(fmt.Sprintf("PROTECTED (%d branches)\n", len(report.Protected)))
		sb.WriteString("------------------------------------------------------------\n")
		for _, branch := range report.Protected {
			sb.WriteString(fmt.Sprintf("  * %s\n", branch.Name))
			sb.WriteString(fmt.Sprintf("    Reason: %s\n", branch.Reason))
			sb.WriteString("\n")
		}
	}

	// Summary
	sb.WriteString("============================================================\n")
	sb.WriteString(fmt.Sprintf("SUMMARY: %d safe | %d local-only | %d unmerged | %d protected\n",
		report.Summary.SafeCount,
		report.Summary.LocalOnlyCount,
		report.Summary.UnmergedCount,
		report.Summary.ProtectedCount))
	sb.WriteString("============================================================\n")

	return sb.String()
}

// generateJSONReport creates a JSON formatted report
func generateJSONReport(report *AnalysisReport) string {
	output, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}
	return string(output)
}

// generateCSVReport creates a CSV formatted report
func generateCSVReport(report *AnalysisReport) string {
	var sb strings.Builder
	writer := csv.NewWriter(&sb)

	// Header
	_ = writer.Write([]string{"Name", "Status", "Delete Method", "Reason", "Remote Status", "Last Commit"})

	// All branches
	allBranches := append(report.SafeToDelete, report.LocalOnly...)
	allBranches = append(allBranches, report.Unmerged...)
	allBranches = append(allBranches, report.Protected...)

	for _, branch := range allBranches {
		_ = writer.Write([]string{
			branch.Name,
			branch.Status,
			branch.DeleteMethod,
			branch.Reason,
			branch.RemoteStatus,
			branch.LastCommit,
		})
	}

	writer.Flush()
	return sb.String()
}

// outputReport writes the report to stdout or a file
func outputReport(report *AnalysisReport, format string, filePath string) {
	var output string

	switch format {
	case "json":
		output = generateJSONReport(report)
	case "csv":
		output = generateCSVReport(report)
	default:
		output = generateTextReport(report)
	}

	if filePath != "" {
		err := os.WriteFile(filePath, []byte(output), 0644)
		if err != nil {
			fmt.Printf("❌ Failed to write report to file: %v\n", err)
			fmt.Println(output)
			return
		}
		fmt.Printf("✅ Report saved to: %s\n", filePath)
	} else {
		fmt.Println(output)
	}
}
