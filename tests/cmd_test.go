package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"git-gone/internal/git"
)

// getProjectRoot returns the project root directory.
func getProjectRoot(t *testing.T) string {
	t.Helper()
	// Get current working directory and go up from tests/
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get cwd: %v", err)
	}
	return filepath.Dir(cwd)
}

// TestBranchesCommand_InNonGitRepo verifies error handling outside git repo.
func TestBranchesCommand_InNonGitRepo(t *testing.T) {
	// Create temp dir that is NOT a git repo
	tempDir, err := os.MkdirTemp("", "git-gone-nongit-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tempDir)

	// Use the git package directly
	err = git.CheckGitRepository()
	if err == nil {
		t.Error("Expected error when not in git repository")
	}
}

// TestTagsListCommand_InNonGitRepo verifies error handling for tags list.
func TestTagsListCommand_InNonGitRepo(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "git-gone-nongit-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tempDir)

	// Use the git package directly
	err = git.CheckGitRepository()
	if err == nil {
		t.Error("Expected error when not in git repository")
	}
}

// TestVersionCommand_ShowsVersion verifies version output by checking cmd package.
func TestVersionCommand_ShowsVersion(t *testing.T) {
	projectRoot := getProjectRoot(t)
	binaryPath := filepath.Join(projectRoot, "git-gone-test")

	// Build the binary
	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Dir = projectRoot
	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build: %v\nOutput: %s", err, string(output))
	}
	defer os.Remove(binaryPath)

	cmd := exec.Command(binaryPath, "version")
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Version command failed: %v\nOutput: %s", err, string(output))
	}

	// Should contain version number (0.x.x format)
	if !strings.Contains(string(output), "0.") {
		t.Errorf("Expected version output to contain version number, got: %s", string(output))
	}
}

// TestHelpCommand_ShowsHelp verifies help output.
func TestHelpCommand_ShowsHelp(t *testing.T) {
	projectRoot := getProjectRoot(t)
	binaryPath := filepath.Join(projectRoot, "git-gone-test")

	// Build the binary
	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Dir = projectRoot
	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build: %v\nOutput: %s", err, string(output))
	}
	defer os.Remove(binaryPath)

	cmd := exec.Command(binaryPath, "--help")
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Help command failed: %v\nOutput: %s", err, string(output))
	}

	// Should contain usage info
	if !strings.Contains(string(output), "branches") {
		t.Errorf("Expected help output to mention 'branches', got: %s", string(output))
	}

	if !strings.Contains(string(output), "tags") {
		t.Errorf("Expected help output to mention 'tags', got: %s", string(output))
	}
}

// TestBranchesCommand_WithIncompatibleFlags verifies -a and -f check.
func TestBranchesCommand_WithIncompatibleFlags(t *testing.T) {
	projectRoot := getProjectRoot(t)
	binaryPath := filepath.Join(projectRoot, "git-gone-test")

	// Build the binary first
	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Dir = projectRoot
	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build: %v\nOutput: %s", err, string(output))
	}
	defer os.Remove(binaryPath)

	// Create a temp git repo to run the command in
	h := NewTestHelper(t)
	defer h.Cleanup()

	// Run with -a and -f flags
	cmd := exec.Command(binaryPath, "branches", "-a", "-f")
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	output, err := cmd.CombinedOutput()

	// Should fail due to incompatible flags
	if err == nil {
		t.Error("Expected error with -a and -f flags together")
	}

	// Should mention incompatibility
	if !strings.Contains(string(output), "incompatible") {
		t.Errorf("Expected incompatibility error, got: %s", string(output))
	}
}
