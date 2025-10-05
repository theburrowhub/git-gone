package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	verboseVersion bool
	formatVersion  string
)

type VersionInfo struct {
	Version    string `json:"version"`
	CommitHash string `json:"commit,omitempty"`
	BuildTime  string `json:"build_time,omitempty"`
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  `Display the version of git-gone. Use --verbose for additional build information.`,
	Example: `  # Show version number only
  git-gone version

  # Show detailed version information
  git-gone version --verbose
  git-gone version -v

  # Show version in JSON format
  git-gone version --format json

  # Show detailed version in JSON format
  git-gone version -v --format json`,
	Run: func(cmd *cobra.Command, args []string) {
		info := VersionInfo{
			Version: Version,
		}

		if verboseVersion {
			info.CommitHash = CommitHash
			info.BuildTime = BuildTime
		}

		switch formatVersion {
		case "json":
			output, err := json.MarshalIndent(info, "", "  ")
			if err != nil {
				fmt.Printf("Error formatting JSON: %v\n", err)
				return
			}
			fmt.Println(string(output))
		default:
			// Plain text format
			fmt.Println(info.Version)
			if verboseVersion {
				fmt.Printf("Commit: %s\n", info.CommitHash)
				fmt.Printf("Built: %s\n", info.BuildTime)
			}
		}
	},
}

func init() {
	versionCmd.Flags().BoolVarP(&verboseVersion, "verbose", "v", false, "Show detailed version information including commit and build time")
	versionCmd.Flags().StringVar(&formatVersion, "format", "text", "Output format (text or json)")
}
