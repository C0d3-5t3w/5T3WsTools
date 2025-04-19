// Package rxp provides extensions and utilities for the standard regexp package.
package rxp

import (
	"regexp"
	"strings"
)

// Matcher wraps a compiled regexp for extended functionality
type Matcher struct {
	*regexp.Regexp
}

// New creates a new Matcher from a pattern string
func New(pattern string) (*Matcher, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return &Matcher{re}, nil
}

// MustNew creates a new Matcher from a pattern string and panics if compilation fails
func MustNew(pattern string) *Matcher {
	re := regexp.MustCompile(pattern)
	return &Matcher{re}
}

// MatchAll returns all non-overlapping matches of the regexp in the input string
// along with their start and end positions
func (m *Matcher) MatchAll(s string) []Match {
	matches := m.FindAllStringSubmatchIndex(s, -1)
	result := make([]Match, 0, len(matches))

	for _, match := range matches {
		if len(match) >= 2 {
			result = append(result, Match{
				Text:  s[match[0]:match[1]],
				Start: match[0],
				End:   match[1],
			})
		}
	}

	return result
}

// Match represents a regexp match with position information
type Match struct {
	Text  string // The matched text
	Start int    // Start position in the original string
	End   int    // End position in the original string
}

// MatchFull returns true only if the entire string matches the regexp
func (m *Matcher) MatchFull(s string) bool {
	matches := m.FindStringIndex(s)
	if matches == nil {
		return false
	}
	return matches[0] == 0 && matches[1] == len(s)
}

// ExtractGroups returns a map of named capture groups and their values
func (m *Matcher) ExtractGroups(s string) map[string]string {
	match := m.FindStringSubmatch(s)
	if match == nil {
		return nil
	}

	result := make(map[string]string)
	for i, name := range m.SubexpNames() {
		if i != 0 && name != "" && i < len(match) {
			result[name] = match[i]
		}
	}

	return result
}

// MatchAny returns true if the string matches any of the provided patterns
func MatchAny(s string, patterns []string) bool {
	for _, pattern := range patterns {
		if matched, _ := regexp.MatchString(pattern, s); matched {
			return true
		}
	}
	return false
}

// CountMatches returns the number of matches in the string
func (m *Matcher) CountMatches(s string) int {
	return len(m.FindAllString(s, -1))
}

// Common predefined patterns
var (
	EmailPattern   = `(?i)[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}`
	URLPattern     = `https?://(?:[-\w.]|(?:%[\da-fA-F]{2}))+(?::\d+)?(?:/[^?\s]*)?(?:\?[^#\s]*)?(?:#[^\s]*)?`
	IPV4Pattern    = `\b(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\b`
	DateISOPattern = `\d{4}-\d{2}-\d{2}`
	TimePattern    = `\d{2}:\d{2}(:\d{2})?`
)

// Email returns a matcher for email addresses
func Email() *Matcher {
	return MustNew(EmailPattern)
}

// URL returns a matcher for URLs
func URL() *Matcher {
	return MustNew(URLPattern)
}

// IPV4 returns a matcher for IPv4 addresses
func IPV4() *Matcher {
	return MustNew(IPV4Pattern)
}

// Replace replaces all matches of the regexp with a replacement string
// while providing access to the match information in the callback
func (m *Matcher) Replace(s string, replacer func(Match) string) string {
	matches := m.MatchAll(s)
	if len(matches) == 0 {
		return s
	}

	var result strings.Builder
	lastEnd := 0

	for _, match := range matches {
		// Add text before the match
		result.WriteString(s[lastEnd:match.Start])

		// Add replacement
		result.WriteString(replacer(match))

		lastEnd = match.End
	}

	// Add text after the last match
	result.WriteString(s[lastEnd:])

	return result.String()
}
