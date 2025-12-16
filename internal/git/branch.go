package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RemoteStatus represents the tracking status of a branch.
type RemoteStatus int

const (
	// TrackingActive means the branch has an active remote tracking branch.
	TrackingActive RemoteStatus = iota
	// TrackingGone means the remote tracking branch was deleted.
	TrackingGone
	// NoTracking means the branch has no remote tracking.
	NoTracking
)

// Branch represents a local git branch with its metadata.
type Branch struct {
	Name         string
	IsCurrent    bool
	IsDefault    bool
	IsMerged     bool
	RemoteStatus RemoteStatus
	RemoteName   string
}

// IsProtected returns true if the branch cannot be deleted.
func (b *Branch) IsProtected() bool {
	return b.IsCurrent || b.IsDefault
}

// IsSafeToDelete returns true if the branch can be safely deleted.
func (b *Branch) IsSafeToDelete() bool {
	return !b.IsProtected() && (b.IsMerged || b.RemoteStatus == TrackingGone)
}

// GetMergedBranches returns branches that have been merged into the default branch.
func GetMergedBranches(defaultBranch string) ([]string, error) {
	cmd := exec.Command("git", "branch", "--merged", defaultBranch)
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	var branches []string
	for _, line := range lines {
		branch := strings.TrimSpace(strings.TrimPrefix(line, "*"))
		branch = strings.TrimSpace(branch)
		if branch != "" {
			branches = append(branches, branch)
		}
	}
	return branches, nil
}

// GetGoneBranches returns branches whose remote tracking branch has been deleted.
func GetGoneBranches() ([]string, error) {
	cmd := exec.Command("git", "branch", "--format", "%(refname:short) %(upstream:track)")
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	var branches []string
	for _, line := range lines {
		if strings.Contains(line, "[gone]") {
			parts := strings.Fields(line)
			if len(parts) > 0 {
				branches = append(branches, parts[0])
			}
		}
	}
	return branches, nil
}

// GetAllLocalBranches returns all local branch names.
func GetAllLocalBranches() ([]string, error) {
	cmd := exec.Command("git", "branch", "--format", "%(refname:short)")
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	var branches []string
	for _, line := range lines {
		branch := strings.TrimSpace(line)
		if branch != "" {
			branches = append(branches, branch)
		}
	}
	return branches, nil
}

// FilterProtectedBranches removes the default and current branches from the list.
func FilterProtectedBranches(branches []string, defaultBranch, currentBranch string) []string {
	var filtered []string
	for _, branch := range branches {
		if branch != defaultBranch && branch != currentBranch && branch != "" {
			filtered = append(filtered, branch)
		}
	}
	return filtered
}

// DeleteBranch deletes a local branch. It tries safe delete first, then force delete.
func DeleteBranch(name string) error {
	cmd := exec.Command("git", "branch", "-d", name)
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	err := cmd.Run()

	if err != nil {
		cmd = exec.Command("git", "branch", "-D", name)
		cmd.Env = append(os.Environ(), "LC_ALL=C")
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("%s", string(output))
		}
	}

	return nil
}

// DeleteBranchWithRemote deletes both local and remote branch.
func DeleteBranchWithRemote(name string) error {
	// Try to delete remote branch first
	cmd := exec.Command("git", "push", "origin", "--delete", name)
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	output, err := cmd.CombinedOutput()
	if err != nil {
		outputStr := string(output)
		if !strings.Contains(outputStr, "remote ref does not exist") {
			// Log warning but continue with local delete
			fmt.Printf("⚠️  Warning: Failed to delete remote branch %s: %s\n", name, outputStr)
		}
	}

	// Force delete local branch
	cmd = exec.Command("git", "branch", "-D", name)
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", string(output))
	}

	return nil
}
