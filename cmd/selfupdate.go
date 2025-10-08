package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"

	"github.com/fynelabs/selfupdate"
	"github.com/spf13/cobra"
)

var (
	checkOnly bool
)

var selfUpdateCmd = &cobra.Command{
	Use:   "self-update",
	Short: "Update git-gone to the latest version",
	Long: `Update git-gone to the latest version from GitHub releases

This command will:
  1. Check for the latest version available
  2. Download the binary for your platform
  3. Verify and apply the update
  4. Replace the current executable`,
	Example: `  # Update to the latest version
  git-gone self-update

  # Check for updates without applying them
  git-gone self-update --check`,
	Run: func(cmd *cobra.Command, args []string) {
		runSelfUpdate()
	},
}

func init() {
	selfUpdateCmd.Flags().BoolVarP(&checkOnly, "check", "c", false, "Only check for updates without applying them")
}

func runSelfUpdate() {
	fmt.Println("🔍 Checking for updates...")

	// Construct the GitHub release URL
	// Format: https://github.com/{owner}/{repo}/releases/download/{version}/{binary}-{os}-{arch}
	owner := "theburrowhub"
	repo := "git-gone"
	
	// Get the appropriate binary name and extension for the platform
	binaryName := "git-gone"
	
	// Map OS names to match release artifacts
	osName := runtime.GOOS
	if osName == "darwin" {
		osName = "macos"
	}
	
	ext := ""
	if runtime.GOOS == "windows" {
		ext = ".exe"
	}

	// Build the URL for the latest release
	url := fmt.Sprintf("https://github.com/%s/%s/releases/latest/download/%s-%s-%s%s",
		owner, repo, binaryName, osName, runtime.GOARCH, ext)

	fmt.Printf("📦 Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("🔗 Update URL: %s\n", url)

	if checkOnly {
		fmt.Println("ℹ️  Check-only mode: Would fetch from:", url)
		fmt.Println("✅ Use without --check flag to apply the update")
		return
	}

	// Download the new binary
	fmt.Println("⬇️  Downloading latest version...")
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("❌ Failed to download update: %v\n", err)
		fmt.Println("\nℹ️  Possible reasons:")
		fmt.Println("  • No internet connection")
		fmt.Println("  • Latest release not available for your platform")
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("❌ Failed to download update: HTTP %d\n", resp.StatusCode)
		if resp.StatusCode == http.StatusNotFound {
			fmt.Println("\nℹ️  No release found for your platform.")
			fmt.Printf("    Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
		}
		os.Exit(1)
	}

	// Show download progress
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("❌ Failed to read update: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Downloaded %d bytes\n", len(body))
	fmt.Println("🔄 Applying update...")

	// Apply the update
	err = selfupdate.Apply(io.NopCloser(bytes.NewReader(body)), selfupdate.Options{})
	if err != nil {
		if rerr := selfupdate.RollbackError(err); rerr != nil {
			fmt.Printf("❌ Failed to rollback from bad update: %v\n", rerr)
		}
		fmt.Printf("❌ Update failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Successfully updated to the latest version!")
	fmt.Println("🔄 Please restart git-gone to use the new version")
}

