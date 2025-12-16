package cmd

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"git-gone/internal/git"
	"git-gone/internal/tui"

	"github.com/koki-develop/go-fzf"
	"github.com/spf13/cobra"
)

// Tag-specific flags
var includeNonStale bool

var tagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "Manage and clean up tags",
	Long: `Manage and clean up tags in the repository.

This command provides subcommands to list and clean stale tags
(tags that exist locally but not on the remote).

Use --no-stale (-n) to include ALL local tags, not just stale ones.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Default behavior: show help
		cmd.Help()
	},
}

var tagsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List stale tags (local tags not on remote)",
	Long: `List stale tags that exist locally but not on the remote.

These tags may have been deleted from the remote or were never pushed.

Use --no-stale (-n) to list ALL local tags instead.`,
	Example: `  # List stale tags only
  git-gone tags list

  # List ALL local tags
  git-gone tags list --no-stale
  git-gone tags list -n`,
	Run: func(cmd *cobra.Command, args []string) {
		runTagsList()
	},
}

var tagsCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean up stale tags interactively",
	Long: `Clean up stale tags interactively.

This command will:
  1. Find tags that exist locally but not on the remote
  2. Show an interactive selector to choose tags to delete
  3. Confirm before deletion (unless --force is used)
  4. Safely delete selected tags

Use --no-stale (-n) to include ALL local tags, not just stale ones.

Interactive Controls:
  ↑/↓         Navigate through the list
  Tab         Toggle selection
  Enter       Confirm selection
  Esc         Cancel operation
  Type        Filter tags by name`,
	Example: `  # Clean up stale tags interactively
  git-gone tags clean

  # Select all stale tags and confirm
  git-gone tags clean --all

  # Delete all stale tags without confirmation
  git-gone tags clean --all --force

  # Clean ANY local tag (not just stale)
  git-gone tags clean --no-stale
  git-gone tags clean -n`,
	Run: func(cmd *cobra.Command, args []string) {
		runTagsClean()
	},
}

func init() {
	// Add --no-stale flag to tags subcommands
	tagsListCmd.Flags().BoolVarP(&includeNonStale, "no-stale", "n", false, "Include ALL local tags, not just stale ones")
	tagsCleanCmd.Flags().BoolVarP(&includeNonStale, "no-stale", "n", false, "Include ALL local tags, not just stale ones")

	tagsCmd.AddCommand(tagsListCmd)
	tagsCmd.AddCommand(tagsCleanCmd)
}

func runTagsList() {
	// Check if we're in a git repository
	if err := git.CheckGitRepository(); err != nil {
		fmt.Printf("%s Not in a git repository\n", tui.EmojiError)
		os.Exit(1)
	}

	var tags []string
	var err error
	var listType string

	if includeNonStale {
		// List ALL local tags
		tags, err = git.GetLocalTags()
		if err != nil {
			fmt.Printf("%s Failed to get local tags: %v\n", tui.EmojiError, err)
			os.Exit(1)
		}
		listType = "local"
	} else {
		// Check if remote exists for stale detection
		if !git.HasRemote() {
			fmt.Printf("%s No remote 'origin' configured. Cannot determine stale tags.\n", tui.EmojiWarning)
			fmt.Printf("   Use --no-stale (-n) to list all local tags instead.\n")
			return
		}

		fmt.Printf("%s Fetching remote tags...\n", tui.EmojiRefresh)

		tags, err = git.GetStaleTags()
		if err != nil {
			fmt.Printf("%s Failed to get stale tags: %v\n", tui.EmojiError, err)
			os.Exit(1)
		}
		listType = "stale"
	}

	if len(tags) == 0 {
		if includeNonStale {
			fmt.Printf("%s No local tags found.\n", tui.EmojiSuccess)
		} else {
			fmt.Printf("%s No stale tags found. All local tags exist on remote.\n", tui.EmojiSuccess)
		}
		return
	}

	sort.Strings(tags)

	if includeNonStale {
		fmt.Printf("\n%s  Found %d local tag(s):\n", tui.EmojiTag, len(tags))
	} else {
		fmt.Printf("\n%s  Found %d %s tag(s) (local only, not on remote):\n", tui.EmojiTag, len(tags), listType)
	}
	for _, tag := range tags {
		fmt.Printf("   • %s\n", tag)
	}
}

func runTagsClean() {
	// Validate incompatible flags
	if selectAll && forceDelete {
		fmt.Printf("%s Options -a (--all) and -f (--force) are incompatible\n", tui.EmojiError)
		os.Exit(1)
	}

	// Check if we're in a git repository
	if err := git.CheckGitRepository(); err != nil {
		fmt.Printf("%s Not in a git repository\n", tui.EmojiError)
		os.Exit(1)
	}

	var tags []string
	var err error

	if includeNonStale {
		// Get ALL local tags
		tags, err = git.GetLocalTags()
		if err != nil {
			fmt.Printf("%s Failed to get local tags: %v\n", tui.EmojiError, err)
			os.Exit(1)
		}
	} else {
		// Check if remote exists for stale detection
		if !git.HasRemote() {
			fmt.Printf("%s No remote 'origin' configured. Cannot determine stale tags.\n", tui.EmojiWarning)
			fmt.Printf("   Use --no-stale (-n) to manage all local tags instead.\n")
			return
		}

		fmt.Printf("%s Fetching remote tags...\n", tui.EmojiRefresh)

		tags, err = git.GetStaleTags()
		if err != nil {
			fmt.Printf("%s Failed to get stale tags: %v\n", tui.EmojiError, err)
			os.Exit(1)
		}
	}

	if len(tags) == 0 {
		if includeNonStale {
			fmt.Printf("%s No local tags found.\n", tui.EmojiSuccess)
		} else {
			fmt.Printf("%s No stale tags found. All local tags exist on remote.\n", tui.EmojiSuccess)
		}
		return
	}

	sort.Strings(tags)

	if includeNonStale {
		fmt.Printf("\n%s  Found %d local tag(s):\n", tui.EmojiTag, len(tags))
	} else {
		fmt.Printf("\n%s  Found %d stale tag(s) (local only, not on remote):\n", tui.EmojiTag, len(tags))
	}

	// Select tags: use all if -a flag is set, otherwise use interactive fzf
	var selectedTags []string
	if selectAll {
		selectedTags = tags
	} else {
		selectedTags, err = selectTagsWithFzf(tags)
		if err != nil {
			if err.Error() == "abort" {
				fmt.Printf("\n%s Selection cancelled\n", tui.EmojiError)
				return
			}
			fmt.Printf("%s Failed to select tags: %v\n", tui.EmojiError, err)
			os.Exit(1)
		}
	}

	if len(selectedTags) == 0 {
		fmt.Printf("\n%s No tags selected for deletion\n", tui.EmojiSuccess)
		return
	}

	// Show tags to delete
	fmt.Printf("\n%s The following tags will be deleted:\n", tui.EmojiWarning)
	for _, tag := range selectedTags {
		fmt.Printf("  • %s\n", tag)
	}

	// Confirm deletion (unless --force is used)
	if !forceDelete {
		fmt.Print("\nAre you sure you want to delete these tags? (y/N): ")
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))

		if response != "y" && response != "yes" {
			fmt.Printf("%s Deletion cancelled\n", tui.EmojiError)
			return
		}
	}

	// Delete tags
	deletedCount := 0
	for _, tag := range selectedTags {
		if err := git.DeleteTag(tag); err != nil {
			fmt.Printf("%s Failed to delete tag %s: %v\n", tui.EmojiError, tag, err)
		} else {
			fmt.Printf("%s Deleted tag: %s\n", tui.EmojiSuccess, tag)
			deletedCount++
		}
	}

	fmt.Printf("\n%s Successfully deleted %d tag(s)\n", tui.EmojiCelebrate, deletedCount)
}

func selectTagsWithFzf(tags []string) ([]string, error) {
	if len(tags) == 0 {
		return []string{}, nil
	}

	f, err := fzf.New(
		fzf.WithLimit(len(tags)),
		fzf.WithNoLimit(true),
		fzf.WithPrompt("Select tags to delete > "),
		fzf.WithCursor("> "),
		fzf.WithSelectedPrefix("[✓] "),
		fzf.WithUnselectedPrefix("[ ] "),
		fzf.WithInputPlaceholder("Type to filter, Tab to select/deselect, Enter to confirm, Esc to cancel"),
	)
	if err != nil {
		return nil, err
	}

	indices, err := f.Find(tags, func(i int) string {
		return tags[i]
	})

	if err != nil {
		return nil, fmt.Errorf("abort")
	}

	selected := make([]string, len(indices))
	for i, idx := range indices {
		selected[i] = tags[idx]
	}

	return selected, nil
}
