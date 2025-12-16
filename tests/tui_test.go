package tests

import (
	"testing"

	"git-gone/internal/tui"
)

func TestEmojiConstants_AreDefined(t *testing.T) {
	// Verify emoji constants are not empty
	emojis := map[string]string{
		"EmojiSuccess":   tui.EmojiSuccess,
		"EmojiError":     tui.EmojiError,
		"EmojiWarning":   tui.EmojiWarning,
		"EmojiRefresh":   tui.EmojiRefresh,
		"EmojiLocation":  tui.EmojiLocation,
		"EmojiBranch":    tui.EmojiBranch,
		"EmojiSearch":    tui.EmojiSearch,
		"EmojiCelebrate": tui.EmojiCelebrate,
		"EmojiDanger":    tui.EmojiDanger,
		"EmojiTag":       tui.EmojiTag,
	}

	for name, emoji := range emojis {
		if emoji == "" {
			t.Errorf("%s should not be empty", name)
		}
	}
}

func TestSelectBranches_WithEmptyList_ReturnsEmpty(t *testing.T) {
	// SelectBranches should return empty slice for empty input
	result, err := tui.SelectBranches([]string{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("Expected empty result, got: %v", result)
	}
}

func TestSelectTags_WithEmptyList_ReturnsEmpty(t *testing.T) {
	// SelectTags should return empty slice for empty input
	result, err := tui.SelectTags([]string{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("Expected empty result, got: %v", result)
	}
}
