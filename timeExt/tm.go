// Package time extends the functionality of Go's standard time library
package timeExt

import (
	"fmt"
	stdtime "time"
)

// FormatRelative returns a string describing the time relative to now
// like "5 minutes ago" or "in 2 days"
func FormatRelative(t stdtime.Time) string {
	now := stdtime.Now()
	diff := now.Sub(t)

	if diff > 0 {
		// Past
		switch {
		case diff < stdtime.Minute:
			return "just now"
		case diff < stdtime.Hour:
			return fmt.Sprintf("%d minutes ago", int(diff.Minutes()))
		case diff < 24*stdtime.Hour:
			return fmt.Sprintf("%d hours ago", int(diff.Hours()))
		case diff < 48*stdtime.Hour:
			return "yesterday"
		case diff < 7*24*stdtime.Hour:
			return fmt.Sprintf("%d days ago", int(diff.Hours()/24))
		default:
			return t.Format("Jan 2, 2006")
		}
	} else {
		// Future
		diff = -diff
		switch {
		case diff < stdtime.Minute:
			return "in a moment"
		case diff < stdtime.Hour:
			return fmt.Sprintf("in %d minutes", int(diff.Minutes()))
		case diff < 24*stdtime.Hour:
			return fmt.Sprintf("in %d hours", int(diff.Hours()))
		case diff < 48*stdtime.Hour:
			return "tomorrow"
		case diff < 7*24*stdtime.Hour:
			return fmt.Sprintf("in %d days", int(diff.Hours()/24))
		default:
			return t.Format("Jan 2, 2006")
		}
	}
}

// ParseMultipleFormats attempts to parse a time string using multiple formats
func ParseMultipleFormats(str string, formats ...string) (stdtime.Time, error) {
	var firstErr error

	for _, format := range formats {
		t, err := stdtime.Parse(format, str)
		if err == nil {
			return t, nil
		}
		if firstErr == nil {
			firstErr = err
		}
	}

	if firstErr == nil {
		return stdtime.Time{}, fmt.Errorf("no formats provided to parse %q", str)
	}
	return stdtime.Time{}, fmt.Errorf("could not parse %q with any of the provided formats", str)
}

// IsBusinessDay returns true if the given time falls on a business day (Monday-Friday)
func IsBusinessDay(t stdtime.Time) bool {
	weekday := t.Weekday()
	return weekday != stdtime.Saturday && weekday != stdtime.Sunday
}

// NextBusinessDay returns the next business day after the given time
func NextBusinessDay(t stdtime.Time) stdtime.Time {
	t = t.AddDate(0, 0, 1)
	for !IsBusinessDay(t) {
		t = t.AddDate(0, 0, 1)
	}
	return t
}

// IsBetween returns true if the time is between start and end, inclusive
func IsBetween(t, start, end stdtime.Time) bool {
	return (t.Equal(start) || t.After(start)) && (t.Equal(end) || t.Before(end))
}

// StartOfMonth returns the first day of the month containing the given time
func StartOfMonth(t stdtime.Time) stdtime.Time {
	return stdtime.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// EndOfMonth returns the last day of the month containing the given time
func EndOfMonth(t stdtime.Time) stdtime.Time {
	return StartOfMonth(t).AddDate(0, 1, -1)
}

// Quarter returns the quarter (1-4) for the given time
func Quarter(t stdtime.Time) int {
	return int(t.Month()-1)/3 + 1
}

// FormatDuration formats a duration in a more human-readable way than the default
func FormatDuration(d stdtime.Duration) string {
	if d < stdtime.Minute {
		return fmt.Sprintf("%d seconds", int(d.Seconds()))
	} else if d < stdtime.Hour {
		return fmt.Sprintf("%d minutes", int(d.Minutes()))
	} else if d < 24*stdtime.Hour {
		h := int(d.Hours())
		m := int(d.Minutes()) % 60
		if m == 0 {
			return fmt.Sprintf("%d hours", h)
		}
		return fmt.Sprintf("%d hours %d minutes", h, m)
	} else {
		days := int(d.Hours()) / 24
		h := int(d.Hours()) % 24
		if h == 0 {
			return fmt.Sprintf("%d days", days)
		}
		return fmt.Sprintf("%d days %d hours", days, h)
	}
}
