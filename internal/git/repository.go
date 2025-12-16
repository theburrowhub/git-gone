// Package git provides operations for interacting with git repositories.
package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Repository represents a git repository context.
type Repository struct {
	Path          string
	DefaultBranch string
	CurrentBranch string
	HasRemote     bool
	RemoteURL     string
}

// CheckGitRepository verifies that the current directory is a git repository.
func CheckGitRepository() error {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("not a git repository")
	}
	return nil
}

// GetDefaultBranch returns the default branch name (main, master, etc.).
func GetDefaultBranch() (string, error) {
	// Try to get the default branch from remote
	cmd := exec.Command("git", "symbolic-ref", "refs/remotes/origin/HEAD")
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	output, err := cmd.Output()
	if err == nil {
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

	return "main", nil
}

// GetCurrentBranch returns the currently checked out branch name.
func GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// NewRepository creates a Repository instance for the current directory.
func NewRepository() (*Repository, error) {
	if err := CheckGitRepository(); err != nil {
		return nil, err
	}

	defaultBranch, err := GetDefaultBranch()
	if err != nil {
		return nil, fmt.Errorf("failed to get default branch: %w", err)
	}

	currentBranch, err := GetCurrentBranch()
	if err != nil {
		return nil, fmt.Errorf("failed to get current branch: %w", err)
	}

	repo := &Repository{
		DefaultBranch: defaultBranch,
		CurrentBranch: currentBranch,
	}

	// Check for remote
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	output, err := cmd.Output()
	if err == nil {
		repo.HasRemote = true
		repo.RemoteURL = strings.TrimSpace(string(output))
	}

	return repo, nil
}
