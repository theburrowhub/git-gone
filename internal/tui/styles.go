// Package tui provides terminal user interface components.
package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// Emoji prefixes for consistent messaging.
const (
	EmojiSuccess   = "✅"
	EmojiError     = "❌"
	EmojiWarning   = "⚠️"
	EmojiRefresh   = "🔄"
	EmojiLocation  = "📍"
	EmojiBranch    = "🌿"
	EmojiSearch    = "🔍"
	EmojiCelebrate = "🎉"
	EmojiDanger    = "🚨"
	EmojiTag       = "🏷️"
)

// Colors for TUI elements.
var (
	ColorSuccess = lipgloss.Color("#00FF00")
	ColorError   = lipgloss.Color("#FF0000")
	ColorWarning = lipgloss.Color("#FFFF00")
	ColorInfo    = lipgloss.Color("#00BFFF")
	ColorMuted   = lipgloss.Color("#888888")
	ColorDanger  = lipgloss.Color("#FF4500")
)

// Styles for TUI elements.
var (
	// TitleStyle for headers and titles.
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorInfo)

	// SuccessStyle for success messages.
	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorSuccess)

	// ErrorStyle for error messages.
	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorError)

	// WarningStyle for warning messages.
	WarningStyle = lipgloss.NewStyle().
			Foreground(ColorWarning)

	// DangerStyle for dangerous operations.
	DangerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorDanger)

	// MutedStyle for less important info.
	MutedStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	// BranchStyle for branch names.
	BranchStyle = lipgloss.NewStyle().
			Foreground(ColorInfo)

	// TagStyle for tag names.
	TagStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFA500"))
)
