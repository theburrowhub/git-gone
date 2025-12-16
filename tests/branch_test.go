package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"git-gone/internal/git"
)

// TestHelper provides utilities for git tests.
type TestHelper struct {
	t       *testing.T
	tempDir string
	origDir string
}

// NewTestHelper creates a new test helper with a temp git repo.
func NewTestHelper(t *testing.T) *TestHelper {
	t.Helper()

	// Save original directory
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// Create temp directory
	tempDir, err := os.MkdirTemp("", "git-gone-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Change to temp directory
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp dir: %v", err)
	}

	// Initialize git repo
	runGitCmd(t, "init")
	runGitCmd(t, "config", "user.email", "test@example.com")
	runGitCmd(t, "config", "user.name", "Test User")

	// Create initial commit on main
	createFile(t, "README.md", "# Test Repo")
	runGitCmd(t, "add", ".")
	runGitCmd(t, "commit", "-m", "Initial commit")
	runGitCmd(t, "branch", "-M", "main")

	return &TestHelper{
		t:       t,
		tempDir: tempDir,
		origDir: origDir,
	}
}

// Cleanup removes the temp directory and restores working dir.
func (h *TestHelper) Cleanup() {
	os.Chdir(h.origDir)
	os.RemoveAll(h.tempDir)
}

// CreateBranch creates a new branch with a commit.
func (h *TestHelper) CreateBranch(name string) {
	h.t.Helper()
	runGitCmd(h.t, "checkout", "-b", name)
	createFile(h.t, name+".txt", "Content for "+name)
	runGitCmd(h.t, "add", ".")
	runGitCmd(h.t, "commit", "-m", "Commit on "+name)
}

// MergeBranch merges a branch into main.
func (h *TestHelper) MergeBranch(name string) {
	h.t.Helper()
	runGitCmd(h.t, "checkout", "main")
	runGitCmd(h.t, "merge", name, "--no-ff", "-m", "Merge "+name)
}

// CheckoutMain switches to main branch.
func (h *TestHelper) CheckoutMain() {
	h.t.Helper()
	runGitCmd(h.t, "checkout", "main")
}

func runGitCmd(t *testing.T, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s failed: %v\nOutput: %s", strings.Join(args, " "), err, string(output))
	}
	return string(output)
}

func createFile(t *testing.T, name, content string) {
	t.Helper()
	dir := filepath.Dir(name)
	if dir != "." && dir != "" {
		os.MkdirAll(dir, 0755)
	}
	if err := os.WriteFile(name, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create file %s: %v", name, err)
	}
}

func TestGetMergedBranches_ReturnsMergedBranches(t *testing.T) {
	h := NewTestHelper(t)
	defer h.Cleanup()

	// Create and merge a branch
	h.CreateBranch("feature-merged")
	h.CheckoutMain()
	h.MergeBranch("feature-merged")

	// Get merged branches
	merged, err := git.GetMergedBranches("main")
	if err != nil {
		t.Fatalf("GetMergedBranches failed: %v", err)
	}

	// Should contain feature-merged and main
	found := false
	for _, branch := range merged {
		if branch == "feature-merged" {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected 'feature-merged' in merged branches, got: %v", merged)
	}
}

func TestGetMergedBranches_ExcludesUnmergedBranches(t *testing.T) {
	h := NewTestHelper(t)
	defer h.Cleanup()

	// Create a branch but don't merge it
	h.CreateBranch("feature-unmerged")
	h.CheckoutMain()

	// Get merged branches
	merged, err := git.GetMergedBranches("main")
	if err != nil {
		t.Fatalf("GetMergedBranches failed: %v", err)
	}

	// Should NOT contain feature-unmerged
	for _, branch := range merged {
		if branch == "feature-unmerged" {
			t.Errorf("Expected 'feature-unmerged' NOT in merged branches, but it was found")
		}
	}
}

func TestFilterProtectedBranches_ExcludesCurrentAndDefault(t *testing.T) {
	branches := []string{"main", "develop", "feature-1", "feature-2", "current"}

	filtered := git.FilterProtectedBranches(branches, "main", "current")

	// Should only have develop, feature-1, feature-2
	expected := []string{"develop", "feature-1", "feature-2"}
	if len(filtered) != len(expected) {
		t.Errorf("Expected %d branches, got %d: %v", len(expected), len(filtered), filtered)
	}

	for _, exp := range expected {
		found := false
		for _, f := range filtered {
			if f == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected %s in filtered branches", exp)
		}
	}

	// Should not contain main or current
	for _, f := range filtered {
		if f == "main" || f == "current" {
			t.Errorf("Protected branch %s should not be in filtered list", f)
		}
	}
}

func TestFilterProtectedBranches_HandlesEmptyList(t *testing.T) {
	filtered := git.FilterProtectedBranches([]string{}, "main", "current")

	if filtered != nil && len(filtered) != 0 {
		t.Errorf("Expected empty result for empty input, got: %v", filtered)
	}
}

func TestGetAllLocalBranches_ReturnsAllBranches(t *testing.T) {
	h := NewTestHelper(t)
	defer h.Cleanup()

	// Create multiple branches
	h.CreateBranch("branch-a")
	h.CheckoutMain()
	h.CreateBranch("branch-b")
	h.CheckoutMain()

	branches, err := git.GetAllLocalBranches()
	if err != nil {
		t.Fatalf("GetAllLocalBranches failed: %v", err)
	}

	// Should contain main, branch-a, branch-b
	expectedBranches := []string{"main", "branch-a", "branch-b"}
	for _, exp := range expectedBranches {
		found := false
		for _, b := range branches {
			if b == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected branch '%s' in result, got: %v", exp, branches)
		}
	}
}
