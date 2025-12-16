package tui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ConfirmDeletion prompts for a simple y/N confirmation.
func ConfirmDeletion(message string) bool {
	fmt.Printf("\n%s (y/N): ", message)
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}

// TypedConfirmation requires the user to type an exact string to confirm.
func TypedConfirmation(prompt, expected string) bool {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(response)
	return response == expected
}

// ConfirmDangerousOperation shows a warning and requires typing DELETE to confirm.
func ConfirmDangerousOperation(items []string, itemType string) bool {
	fmt.Printf("\n%s WARNING: You are about to delete %d %s:\n", EmojiDanger, len(items), itemType)
	for _, item := range items {
		fmt.Printf("   • %s\n", item)
	}
	fmt.Printf("\n%s  This action cannot be undone! Type 'DELETE' to confirm: ", EmojiWarning)
	return TypedConfirmation("", "DELETE")
}
