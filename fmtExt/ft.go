// Package ft extends the functionality of the standard fmt package
// with additional formatting utilities and helper functions.
package fmtExt

import (
	"fmt"
	"strings"
	"time"
)

// PrintWithTimestamp prints the provided arguments with a timestamp prefix
func PrintWithTimestamp(a ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] ", timestamp)
	fmt.Print(a...)
}

// PrintlnWithTimestamp prints the provided arguments with a timestamp prefix
// followed by a newline
func PrintlnWithTimestamp(a ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] ", timestamp)
	fmt.Println(a...)
}

// PrintfWithTimestamp prints a formatted string with a timestamp prefix
func PrintfWithTimestamp(format string, a ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] ", timestamp)
	fmt.Printf(format, a...)
}

// BoxPrint prints text inside a box made of ASCII characters
func BoxPrint(text string) {
	lines := strings.Split(text, "\n")

	// Find the longest line
	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	// Top border
	fmt.Println("+" + strings.Repeat("-", maxLen+2) + "+")

	// Content
	for _, line := range lines {
		fmt.Printf("| %-*s |\n", maxLen, line)
	}

	// Bottom border
	fmt.Println("+" + strings.Repeat("-", maxLen+2) + "+")
}

// Common ANSI color codes
const (
	ColorRed    = "31"
	ColorGreen  = "32"
	ColorYellow = "33"
	ColorBlue   = "34"
	ColorPurple = "35"
	ColorCyan   = "36"
)

// PrintColored prints text with ANSI color codes
func PrintColored(text string, colorCode string) {
	fmt.Printf("\033[%sm%s\033[0m", colorCode, text)
}

// PrintColoredln prints text with ANSI color codes followed by newline
func PrintColoredln(text string, colorCode string) {
	PrintColored(text, colorCode)
	fmt.Println()
}

// PrintRed prints text in red
func PrintRed(text string) {
	PrintColored(text, ColorRed)
}

// PrintGreen prints text in green
func PrintGreen(text string) {
	PrintColored(text, ColorGreen)
}

// PrintYellow prints text in yellow
func PrintYellow(text string) {
	PrintColored(text, ColorYellow)
}

// PrintIndented prints text with a specified indentation level
func PrintIndented(text string, indentLevel int) {
	indentation := strings.Repeat("  ", indentLevel)
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		fmt.Println(indentation + line)
	}
}

// PrettyJSON formats and prints a JSON string with proper indentation
func PrettyJSON(jsonStr string) {
	var result strings.Builder
	indentLevel := 0
	inQuotes := false

	for i := 0; i < len(jsonStr); i++ {
		char := jsonStr[i]

		switch char {
		case '"':
			// Check for escaped quotes
			if i > 0 && jsonStr[i-1] != '\\' {
				inQuotes = !inQuotes
			}
			result.WriteByte(char)
		case '{', '[':
			result.WriteByte(char)
			if !inQuotes {
				indentLevel++
				result.WriteString("\n" + strings.Repeat("  ", indentLevel))
			}
		case '}', ']':
			if !inQuotes {
				indentLevel--
				result.WriteString("\n" + strings.Repeat("  ", indentLevel))
			}
			result.WriteByte(char)
		case ',':
			result.WriteByte(char)
			if !inQuotes {
				result.WriteString("\n" + strings.Repeat("  ", indentLevel))
			}
		case ':':
			result.WriteByte(char)
			if !inQuotes {
				result.WriteByte(' ')
			}
		default:
			result.WriteByte(char)
		}
	}

	fmt.Println(result.String())
}

// Repeat prints a string repeated n times
func Repeat(s string, n int) {
	fmt.Print(strings.Repeat(s, n))
}
