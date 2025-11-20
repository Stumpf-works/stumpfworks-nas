package cli

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	// Color functions
	Success = color.New(color.FgGreen).SprintFunc()
	Error   = color.New(color.FgRed).SprintFunc()
	Warning = color.New(color.FgYellow).SprintFunc()
	Info    = color.New(color.FgCyan).SprintFunc()
	Bold    = color.New(color.Bold).SprintFunc()

	// Symbols
	CheckMark = Success("✓")
	Cross     = Error("✗")
	Bullet    = Info("●")
	Arrow     = Info("→")
)

// PrintSuccess prints a success message
func PrintSuccess(format string, args ...interface{}) {
	fmt.Printf("%s %s\n", CheckMark, fmt.Sprintf(format, args...))
}

// PrintError prints an error message
func PrintError(format string, args ...interface{}) {
	fmt.Printf("%s %s\n", Cross, fmt.Sprintf(format, args...))
}

// PrintWarning prints a warning message
func PrintWarning(format string, args ...interface{}) {
	fmt.Printf("%s %s\n", Warning("⚠"), fmt.Sprintf(format, args...))
}

// PrintInfo prints an info message
func PrintInfo(format string, args ...interface{}) {
	fmt.Printf("%s %s\n", Bullet, fmt.Sprintf(format, args...))
}

// PrintHeader prints a formatted header
func PrintHeader(title string) {
	width := 60
	padding := (width - len(title) - 2) / 2
	fmt.Println()
	fmt.Println("╔" + repeat("═", width) + "╗")
	fmt.Printf("║%s %s %s║\n", repeat(" ", padding), Bold(title), repeat(" ", width-padding-len(title)-2))
	fmt.Println("╚" + repeat("═", width) + "╝")
	fmt.Println()
}

// PrintSeparator prints a separator line
func PrintSeparator() {
	fmt.Println(repeat("─", 60))
}

// repeat repeats a string n times
func repeat(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}
