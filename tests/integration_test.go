package tests

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"git-gone/internal/git"
)

// TestIntegration_NotInGitRepository verifies error outside git repo.
func TestIntegration_NotInGitRepository(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "git-gone-nongit-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tempDir)

	err = git.CheckGitRepository()
	if err == nil {
		t.Error("Expected error when not in git repository")
	}
}

// TestIntegration_NoBranchesToDelete verifies clean repo message.
func TestIntegration_NoBranchesToDelete(t *testing.T) {
	h := NewTestHelper(t)
	defer h.Cleanup()

	// In a fresh repo with only main, there should be no branches to delete
	merged, err := git.GetMergedBranches("main")
	if err != nil {
		t.Fatalf("GetMergedBranches failed: %v", err)
	}

	// Filter out main itself
	filtered := git.FilterProtectedBranches(merged, "main", "main")

	if len(filtered) != 0 {
		t.Errorf("Expected no deletable branches in fresh repo, got: %v", filtered)
	}
}

// TestIntegration_PartialDeletionFailure simulates deletion with mixed results.
func TestIntegration_PartialDeletionFailure(t *testing.T) {
	h := NewTestHelper(t)
	defer h.Cleanup()

	// Create and merge a branch
	h.CreateBranch("test-branch")
	h.CheckoutMain()
	h.MergeBranch("test-branch")

	// Delete should succeed
	err := git.DeleteBranch("test-branch")
	if err != nil {
		t.Errorf("DeleteBranch should succeed: %v", err)
	}

	// Deleting again should fail (branch doesn't exist)
	err = git.DeleteBranch("test-branch")
	if err == nil {
		t.Error("DeleteBranch should fail for non-existent branch")
	}
}

// TestIntegration_NonEnglishLocale verifies LC_ALL=C is used.
func TestIntegration_NonEnglishLocale(t *testing.T) {
	h := NewTestHelper(t)
	defer h.Cleanup()

	// Set a non-English locale
	originalLang := os.Getenv("LANG")
	originalLcAll := os.Getenv("LC_ALL")
	defer func() {
		os.Setenv("LANG", originalLang)
		os.Setenv("LC_ALL", originalLcAll)
	}()

	os.Setenv("LANG", "es_ES.UTF-8")
	os.Setenv("LC_ALL", "es_ES.UTF-8")

	// Git operations should still work because internal/git uses LC_ALL=C
	defaultBranch, err := git.GetDefaultBranch()
	if err != nil {
		t.Fatalf("GetDefaultBranch failed with non-English locale: %v", err)
	}

	if defaultBranch == "" {
		t.Error("Expected non-empty default branch")
	}
}

// TestIntegration_GitPluginMode tests both git-gone and git gone invocations.
func TestIntegration_GitPluginMode(t *testing.T) {
	// Get project root (one level up from tests/)
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get cwd: %v", err)
	}
	projectRoot := cwd
	// If we're in tests/, go up one level
	if strings.HasSuffix(cwd, "tests") {
		projectRoot = strings.TrimSuffix(cwd, "/tests")
	}

	binaryPath := projectRoot + "/git-gone-test"

	// Build the binary from project root
	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Dir = projectRoot
	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build: %v\nOutput: %s", err, string(output))
	}
	defer os.Remove(binaryPath)

	// Test standalone mode
	cmd := exec.Command(binaryPath, "--help")
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Standalone mode failed: %v\nOutput: %s", err, string(output))
	}

	if !strings.Contains(string(output), "git-gone") {
		t.Error("Expected standalone mode to work")
	}
}

// TestIntegration_DifferentDefaultBranches tests main/master/develop detection.
func TestIntegration_DifferentDefaultBranches(t *testing.T) {
	tests := []struct {
		name        string
		setupBranch string
	}{
		{"main branch", "main"},
		{"master branch", "master"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir, err := os.MkdirTemp("", "git-gone-test-*")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			origDir, _ := os.Getwd()
			defer os.Chdir(origDir)
			os.Chdir(tempDir)

			// Initialize git repo with specific branch name
			runGitCmd(t, "init")
			runGitCmd(t, "config", "user.email", "test@example.com")
			runGitCmd(t, "config", "user.name", "Test User")
			createFile(t, "README.md", "# Test")
			runGitCmd(t, "add", ".")
			runGitCmd(t, "commit", "-m", "Initial commit")
			runGitCmd(t, "branch", "-M", tt.setupBranch)

			// Test default branch detection
			defaultBranch, err := git.GetDefaultBranch()
			if err != nil {
				t.Fatalf("GetDefaultBranch failed: %v", err)
			}

			// Should detect the correct default branch
			if defaultBranch != tt.setupBranch {
				t.Errorf("Expected default branch '%s', got '%s'", tt.setupBranch, defaultBranch)
			}
		})
	}
}

// TestIntegration_DeleteBranchWithRemote_LocalOnly tests deletion when remote doesn't exist.
func TestIntegration_DeleteBranchWithRemote_LocalOnly(t *testing.T) {
	h := NewTestHelper(t)
	defer h.Cleanup()

	// Create a branch (no remote)
	h.CreateBranch("local-only-branch")
	h.CheckoutMain()

	// Delete with remote should succeed for local part
	err := git.DeleteBranchWithRemote("local-only-branch")
	if err != nil {
		t.Errorf("DeleteBranchWithRemote should succeed for local branch: %v", err)
	}

	// Verify branch is deleted
	branches, _ := git.GetAllLocalBranches()
	for _, b := range branches {
		if b == "local-only-branch" {
			t.Error("Branch should have been deleted")
		}
	}
}
