// Package strings extends the functionality of the standard strings library
package strings

import (
	"strings"
	"unicode"
)

// IsEmpty returns true if the string is empty
func IsEmpty(s string) bool {
	return len(s) == 0
}

// IsBlank returns true if the string is empty or only contains whitespace
func IsBlank(s string) bool {
	for _, r := range s {
		if !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

// DefaultIfEmpty returns the default value if the string is empty
func DefaultIfEmpty(s string, defaultValue string) string {
	if IsEmpty(s) {
		return defaultValue
	}
	return s
}

// DefaultIfBlank returns the default value if the string is blank
func DefaultIfBlank(s string, defaultValue string) string {
	if IsBlank(s) {
		return defaultValue
	}
	return s
}

// Reverse returns the string with its characters in reverse order
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Capitalize returns the string with the first character converted to uppercase
func Capitalize(s string) string {
	if s == "" {
		return ""
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

// Uncapitalize returns the string with the first character converted to lowercase
func Uncapitalize(s string) string {
	if s == "" {
		return ""
	}
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}

// SwapCase returns the string with uppercase changed to lowercase and vice versa
func SwapCase(s string) string {
	r := []rune(s)
	for i := 0; i < len(r); i++ {
		if unicode.IsUpper(r[i]) {
			r[i] = unicode.ToLower(r[i])
		} else if unicode.IsLower(r[i]) {
			r[i] = unicode.ToUpper(r[i])
		}
	}
	return string(r)
}

// TruncateWithSuffix truncates a string to the specified length and appends the suffix
func TruncateWithSuffix(s string, maxLength int, suffix string) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength] + suffix
}

// LeftPad pads the string on the left with the specified char to the given width
func LeftPad(s string, width int, char rune) string {
	if len(s) >= width {
		return s
	}
	return strings.Repeat(string(char), width-len(s)) + s
}

// RightPad pads the string on the right with the specified char to the given width
func RightPad(s string, width int, char rune) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(string(char), width-len(s))
}

// ContainsAny returns true if the string contains any of the specified substrings
func ContainsAny(s string, substrings ...string) bool {
	for _, sub := range substrings {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

// ContainsAll returns true if the string contains all of the specified substrings
func ContainsAll(s string, substrings ...string) bool {
	for _, sub := range substrings {
		if !strings.Contains(s, sub) {
			return false
		}
	}
	return true
}

// IsAlpha returns true if the string contains only Unicode letters
func IsAlpha(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return s != ""
}

// IsAlphanumeric returns true if the string contains only Unicode letters or digits
func IsAlphanumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return s != ""
}

// IsNumeric returns true if the string contains only Unicode digits
func IsNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return s != ""
}
