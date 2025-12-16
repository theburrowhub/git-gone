package tests

import (
	"testing"

	"git-gone/internal/tui"
)

func TestEmojiPrefixes_AreConsistent(t *testing.T) {
	// All emoji constants should be consistent format (emoji + optional space)
	emojis := []struct {
		name  string
		emoji string
	}{
		{"Success", tui.EmojiSuccess},
		{"Error", tui.EmojiError},
		{"Warning", tui.EmojiWarning},
		{"Refresh", tui.EmojiRefresh},
		{"Location", tui.EmojiLocation},
		{"Branch", tui.EmojiBranch},
		{"Search", tui.EmojiSearch},
		{"Celebrate", tui.EmojiCelebrate},
		{"Danger", tui.EmojiDanger},
		{"Tag", tui.EmojiTag},
	}

	for _, e := range emojis {
		t.Run(e.name, func(t *testing.T) {
			// Emoji should not be empty
			if e.emoji == "" {
				t.Errorf("%s emoji should not be empty", e.name)
			}

			// Emoji should not start with space
			if len(e.emoji) > 0 && e.emoji[0] == ' ' {
				t.Errorf("%s emoji should not start with space", e.name)
			}
		})
	}
}

func TestTUIStyles_AreDefined(t *testing.T) {
	// Verify styles are usable (won't panic)
	_ = tui.TitleStyle.Render("Test")
	_ = tui.SuccessStyle.Render("Test")
	_ = tui.ErrorStyle.Render("Test")
	_ = tui.WarningStyle.Render("Test")
	_ = tui.DangerStyle.Render("Test")
	_ = tui.MutedStyle.Render("Test")
	_ = tui.BranchStyle.Render("Test")
	_ = tui.TagStyle.Render("Test")
}
