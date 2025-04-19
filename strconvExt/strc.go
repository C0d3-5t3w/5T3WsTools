// Package strconv provides additional string conversion functionality
// beyond what is available in the standard library's strconv package.
package strconvExt

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseIntOrDefault attempts to parse a string into an integer.
// If parsing fails, it returns the provided default value.
func ParseIntOrDefault(s string, defaultVal int) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return val
}

// ParseFloatOrDefault attempts to parse a string into a float64.
// If parsing fails, it returns the provided default value.
func ParseFloatOrDefault(s string, bitSize int, defaultVal float64) float64 {
	val, err := strconv.ParseFloat(s, bitSize)
	if err != nil {
		return defaultVal
	}
	return val
}

// ParseBoolExtended parses a string to a boolean value with extended format support.
// Beyond the standard formats, it also accepts:
// - "yes", "y", "on" as true
// - "no", "n", "off" as false
// Case insensitive.
func ParseBoolExtended(s string) (bool, error) {
	// First try standard parsing
	b, err := strconv.ParseBool(s)
	if err == nil {
		return b, nil
	}

	// Try extended formats
	switch strings.ToLower(s) {
	case "yes", "y", "on":
		return true, nil
	case "no", "n", "off":
		return false, nil
	default:
		return false, fmt.Errorf("couldn't parse %q as bool", s)
	}
}

// FormatIntWithCommas formats an integer with thousand separators.
func FormatIntWithCommas(n int64) string {
	str := strconv.FormatInt(n, 10)

	// If the number is small, return it as is
	if len(str) <= 3 {
		return str
	}

	var result strings.Builder
	remainder := len(str) % 3

	// Handle first group which may be less than 3 digits
	if remainder > 0 {
		result.WriteString(str[:remainder])
		result.WriteByte(',')
	}

	// Add the remaining groups of 3
	for i := remainder; i < len(str); i += 3 {
		if i > 0 && i != remainder {
			result.WriteByte(',')
		}
		result.WriteString(str[i : i+3])
	}

	return result.String()
}

// TruncateString truncates a string to maxLength, adding ellipsis if specified.
func TruncateString(s string, maxLength int, withEllipsis bool) string {
	if len(s) <= maxLength {
		return s
	}

	if withEllipsis && maxLength > 3 {
		return s[:maxLength-3] + "..."
	}

	return s[:maxLength]
}

// ToStringOrDefault attempts to convert various types to string.
// If the conversion is not supported, it returns the default value.
func ToStringOrDefault(v interface{}, defaultVal string) string {
	if v == nil {
		return defaultVal
	}

	switch val := v.(type) {
	case string:
		return val
	case int:
		return strconv.Itoa(val)
	case int64:
		return strconv.FormatInt(val, 10)
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(val)
	default:
		return fmt.Sprintf("%v", v)
	}
}
