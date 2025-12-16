package tui

import (
	"fmt"

	"github.com/koki-develop/go-fzf"
)

// SelectItems presents an interactive selector for items and returns selected indices.
func SelectItems(items []string, prompt string) ([]string, error) {
	if len(items) == 0 {
		return []string{}, nil
	}

	f, err := fzf.New(
		fzf.WithLimit(len(items)),
		fzf.WithNoLimit(true),
		fzf.WithPrompt(prompt),
		fzf.WithCursor("> "),
		fzf.WithSelectedPrefix("[✓] "),
		fzf.WithUnselectedPrefix("[ ] "),
		fzf.WithInputPlaceholder("Type to filter, Tab to select/deselect, Enter to confirm, Esc to cancel"),
	)
	if err != nil {
		return nil, err
	}

	indices, err := f.Find(items, func(i int) string {
		return items[i]
	})

	if err != nil {
		return nil, fmt.Errorf("abort")
	}

	selected := make([]string, len(indices))
	for i, idx := range indices {
		selected[i] = items[idx]
	}

	return selected, nil
}

// SelectBranches presents an interactive selector for branches.
func SelectBranches(branches []string) ([]string, error) {
	return SelectItems(branches, "Select branches to delete > ")
}

// SelectTags presents an interactive selector for tags.
func SelectTags(tags []string) ([]string, error) {
	return SelectItems(tags, "Select tags to delete > ")
}
