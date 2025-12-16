package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SpinnerModel is a Bubbletea model for displaying a spinner during operations.
type SpinnerModel struct {
	spinner  spinner.Model
	message  string
	quitting bool
	done     bool
	err      error
}

// NewSpinner creates a new spinner with a message.
func NewSpinner(message string) SpinnerModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(ColorInfo)
	return SpinnerModel{
		spinner: s,
		message: message,
	}
}

// Init implements tea.Model.
func (m SpinnerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

// Update implements tea.Model.
func (m SpinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case spinnerDoneMsg:
		m.done = true
		m.err = msg.err
		return m, tea.Quit
	}
	return m, nil
}

// View implements tea.Model.
func (m SpinnerModel) View() string {
	if m.done {
		if m.err != nil {
			return fmt.Sprintf("%s %s\n", EmojiError, m.err.Error())
		}
		return ""
	}
	return fmt.Sprintf("%s %s\n", m.spinner.View(), m.message)
}

type spinnerDoneMsg struct {
	err error
}

// RunWithSpinner runs a function while showing a spinner.
func RunWithSpinner(message string, fn func() error) error {
	done := make(chan error, 1)

	go func() {
		done <- fn()
	}()

	// Simple spinner without Bubbletea for simpler integration
	chars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	i := 0
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case err := <-done:
			// Clear the spinner line
			fmt.Printf("\r%s\r", "                                                  ")
			return err
		case <-ticker.C:
			fmt.Printf("\r%s %s", chars[i%len(chars)], message)
			i++
		}
	}
}
