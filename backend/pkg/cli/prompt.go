package cli

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

// ConfirmPrompt asks the user for yes/no confirmation
func ConfirmPrompt(label string) bool {
	prompt := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
	}

	result, err := prompt.Run()
	if err != nil {
		return false
	}

	return strings.ToLower(result) == "y" || strings.ToLower(result) == "yes"
}

// TextPrompt asks the user for text input
func TextPrompt(label string, defaultValue string) (string, error) {
	prompt := promptui.Prompt{
		Label:   label,
		Default: defaultValue,
	}

	return prompt.Run()
}

// PasswordPrompt asks the user for a password
func PasswordPrompt(label string) (string, error) {
	prompt := promptui.Prompt{
		Label: label,
		Mask:  '*',
	}

	return prompt.Run()
}

// SelectPrompt asks the user to select from a list
func SelectPrompt(label string, items []string) (string, error) {
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}

	_, result, err := prompt.Run()
	return result, err
}

// MultiSelectPrompt asks the user to select multiple items
func MultiSelectPrompt(label string, items []string) ([]string, error) {
	selected := []string{}

	for {
		remaining := []string{}
		for _, item := range items {
			found := false
			for _, s := range selected {
				if s == item {
					found = true
					break
				}
			}
			if !found {
				remaining = append(remaining, item)
			}
		}

		if len(remaining) == 0 {
			break
		}

		remaining = append(remaining, "[Done - Finish Selection]")

		prompt := promptui.Select{
			Label: fmt.Sprintf("%s (Selected: %d)", label, len(selected)),
			Items: remaining,
		}

		_, result, err := prompt.Run()
		if err != nil {
			return nil, err
		}

		if result == "[Done - Finish Selection]" {
			break
		}

		selected = append(selected, result)
	}

	return selected, nil
}
