package tui

import (
	"fmt"
	"time"
)

// RunWithSpinner runs a function while showing a lightweight terminal spinner.
func RunWithSpinner(message string, fn func() error) error {
	done := make(chan error, 1)

	go func() {
		done <- fn()
	}()

	chars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	i := 0
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case err := <-done:
			// Clear the entire spinner line regardless of message length.
			fmt.Print("\r\033[2K")
			return err
		case <-ticker.C:
			fmt.Printf("\r%s %s", chars[i%len(chars)], message)
			i++
		}
	}
}
