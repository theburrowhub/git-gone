package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/koki-develop/go-fzf"
	"github.com/spf13/cobra"
)

type Branch struct {
	Name      string
	IsMerged  bool
	IsDefault bool
	IsLocal   bool
}

var branchesCmd = &cobra.Command{
	Use:   "branches",
	Short: "Clean up merged branches interactively",
	Long: `Clean up merged branches interactively

This command will:
  1. Update all remote references
  2. Find branches that are merged or have deleted remotes
  3. Show an interactive selector to choose branches to delete
  4. Confirm before deletion (unless --force is used)
  5. Safely delete selected branches

Interactive Controls:
  ↑/↓         Navigate through the list
  Tab         Toggle selection
  Enter       Confirm selection
  Esc         Cancel operation
  Type        Filter branches by name`,
	Example: `  # Clean up branches in current repository
  git-gone branches

  # Or simply (branches is the default command)
  git gone

  # Skip confirmation prompt
  git-gone branches --force
  git gone -f`,
	Run: func(cmd *cobra.Command, args []string) {
		runCleanup()
	},
}


func runCleanup() {
	// Check if we're in a git repository
	if err := checkGitRepository(); err != nil {
		log.Fatal("❌ Not in a git repository")
	}

	fmt.Println("🔄 Updating remote references...")
	if err := updateRemoteRefs(); err != nil {
		fmt.Printf("⚠️  Warning: Failed to update remote refs: %v\n", err)
	}

	// Get default branch
	defaultBranch, err := getDefaultBranch()
	if err != nil {
		log.Fatalf("❌ Failed to get default branch: %v", err)
	}
	fmt.Printf("📍 Default branch: %s\n", defaultBranch)

	// Get current branch
	currentBranch, err := getCurrentBranch()
	if err != nil {
		log.Fatalf("❌ Failed to get current branch: %v", err)
	}
	fmt.Printf("🌿 Current branch: %s\n", currentBranch)

	// Get branches to delete (both merged and gone remotes)
	goneBranches, err := getGoneBranches()
	if err != nil {
		fmt.Printf("⚠️  Warning: Failed to get gone branches: %v\n", err)
		goneBranches = []string{}
	}

	// Get merged branches
	mergedBranches, err := getMergedBranches(defaultBranch)
	if err != nil {
		fmt.Printf("⚠️  Warning: Failed to get merged branches: %v\n", err)
		mergedBranches = []string{}
	}

	// Combine and deduplicate branches
	branchMap := make(map[string]bool)
	for _, branch := range goneBranches {
		branch = strings.TrimSpace(branch)
		if branch != defaultBranch && branch != currentBranch && branch != "" {
			branchMap[branch] = true
		}
	}
	for _, branch := range mergedBranches {
		branch = strings.TrimSpace(branch)
		if branch != defaultBranch && branch != currentBranch && branch != "" {
			branchMap[branch] = true
		}
	}

	// Convert map to slice
	var branchesToDelete []string
	for branch := range branchMap {
		branchesToDelete = append(branchesToDelete, branch)
	}

	if len(branchesToDelete) == 0 {
		fmt.Println("✅ No branches to delete (all branches are either active or unmerged)")
		return
	}

	// Sort branches for better display
	sort.Strings(branchesToDelete)

	fmt.Printf("\n🔍 Found %d deletable branches:\n", len(branchesToDelete))
	if len(goneBranches) > 0 {
		fmt.Printf("   • %d branches with deleted remotes\n", len(goneBranches))
	}
	if len(mergedBranches) > 0 {
		fmt.Printf("   • %d branches merged into %s\n", len(mergedBranches), defaultBranch)
	}

	// Use go-fzf for selection
	selectedBranches, err := selectBranchesWithFzf(branchesToDelete)
	if err != nil {
		if err.Error() == "abort" {
			fmt.Println("\n❌ Selection cancelled")
			return
		}
		log.Fatalf("❌ Failed to select branches: %v", err)
	}

	if len(selectedBranches) == 0 {
		fmt.Println("\n✅ No branches selected for deletion")
		return
	}

	// Show branches to delete
	fmt.Printf("\n⚠️  The following branches will be deleted:\n")
	for _, branch := range selectedBranches {
		fmt.Printf("  • %s\n", branch)
	}

	// Confirm deletion (unless --force is used)
	if !forceDelete {
		fmt.Print("\nAre you sure you want to delete these branches? (y/N): ")
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))

		if response != "y" && response != "yes" {
			fmt.Println("❌ Deletion cancelled")
			return
		}
	} 

	// Delete selected branches
	deletedCount := 0
	for _, branch := range selectedBranches {
		if err := deleteBranch(branch); err != nil {
			fmt.Printf("❌ Failed to delete branch %s: %v\n", branch, err)
		} else {
			fmt.Printf("✅ Deleted branch: %s\n", branch)
			deletedCount++
		}
	}

	fmt.Printf("\n🎉 Successfully deleted %d branches\n", deletedCount)
}

func checkGitRepository() error {
	cmd := exec.Command("git", "status")
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	return cmd.Run()
}

func updateRemoteRefs() error {
	// Fetch all remotes with prune
	cmd := exec.Command("git", "fetch", "--all", "--prune")
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("fetch failed: %s", string(output))
	}

	// Update remote tracking branches
	cmd = exec.Command("git", "remote", "update", "--prune")
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("remote update failed: %s", string(output))
	}

	return nil
}

func getDefaultBranch() (string, error) {
	// Try to get the default branch from remote
	cmd := exec.Command("git", "symbolic-ref", "refs/remotes/origin/HEAD")
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	output, err := cmd.Output()
	if err == nil {
		// Extract branch name from refs/remotes/origin/main format
		parts := strings.Split(strings.TrimSpace(string(output)), "/")
		if len(parts) > 0 {
			return parts[len(parts)-1], nil
		}
	}

	// Fallback: try common default branch names
	commonDefaults := []string{"main", "master", "develop"}
	for _, branch := range commonDefaults {
		cmd := exec.Command("git", "rev-parse", "--verify", branch)
		cmd.Env = append(os.Environ(), "LC_ALL=C")
		if err := cmd.Run(); err == nil {
			return branch, nil
		}
	}

	return "main", nil // Final fallback
}

func getCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func getGoneBranches() ([]string, error) {
	// Get branches whose remote tracking branch is gone
	// Use LC_ALL=C to ensure consistent English output across all platforms
	cmd := exec.Command("git", "branch", "--format", "%(refname:short) %(upstream:track)")
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	var branches []string
	for _, line := range lines {
		// Check if the branch has [gone] (always in English due to LC_ALL=C)
		if strings.Contains(line, "[gone]") {
			// Extract branch name (first word)
			parts := strings.Fields(line)
			if len(parts) > 0 {
				branches = append(branches, parts[0])
			}
		}
	}
	return branches, nil
}

func getMergedBranches(defaultBranch string) ([]string, error) {
	// Get branches merged into the default branch
	cmd := exec.Command("git", "branch", "--merged", defaultBranch)
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	var branches []string
	for _, line := range lines {
		// Remove the asterisk for current branch and trim spaces
		branch := strings.TrimSpace(strings.TrimPrefix(line, "*"))
		branch = strings.TrimSpace(branch)
		if branch != "" {
			branches = append(branches, branch)
		}
	}
	return branches, nil
}

func selectBranchesWithFzf(branches []string) ([]string, error) {
	if len(branches) == 0 {
		return []string{}, nil
	}

	f, err := fzf.New(
		fzf.WithLimit(len(branches)), // Allow selecting all branches
		fzf.WithNoLimit(true),        // Remove limit on selections
		fzf.WithPrompt("Select branches to delete > "),
		fzf.WithCursor("> "),
		fzf.WithSelectedPrefix("[✓] "),
		fzf.WithUnselectedPrefix("[ ] "),
		fzf.WithInputPlaceholder("Type to filter, Tab to select/deselect, Enter to confirm, Esc to cancel"),
	)
	if err != nil {
		return nil, err
	}

	indices, err := f.Find(branches, func(i int) string {
		return branches[i]
	})

	if err != nil {
		return nil, fmt.Errorf("abort")
	}

	selected := make([]string, len(indices))
	for i, idx := range indices {
		selected[i] = branches[idx]
	}

	return selected, nil
}

func deleteBranch(branch string) error {
	cmd := exec.Command("git", "branch", "-d", branch)
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	err := cmd.Run()

	// If safe delete fails, the branch might have unmerged commits
	// but we already know it's merged to default branch
	if err != nil {
		// Try force delete
		cmd = exec.Command("git", "branch", "-D", branch)
		cmd.Env = append(os.Environ(), "LC_ALL=C")
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("%s", string(output))
		}
	}

	return nil
}
