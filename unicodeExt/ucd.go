// Package unicode provides extensions to the standard Go unicode package.
package unicodeExt

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// UnicodeData represents a parsed entry from the Unicode Character Database.
type UnicodeData struct {
	CodePoint       rune
	Name            string
	GeneralCategory string
	CanonicalClass  int
	BidiClass       string
	DecompType      string
	DecompMapping   []rune
	NumericType     string
	NumericValue    float64
	BidiMirrored    bool
	OldName         string
	SimpleUpperCase rune
	SimpleLowerCase rune
	SimpleTitleCase rune
}

// IsEmoji reports whether the rune is an emoji character.
func IsEmoji(r rune) bool {
	return unicode.In(r, &unicode.RangeTable{
		R32: []unicode.Range32{
			{0x1F600, 0x1F64F, 1}, // Emoticons
			{0x1F300, 0x1F5FF, 1}, // Misc Symbols and Pictographs
			{0x1F680, 0x1F6FF, 1}, // Transport and Map
			{0x1F900, 0x1F9FF, 1}, // Supplemental Symbols and Pictographs
		},
		LatinOffset: 0,
	})
}

// IsPictographic reports whether the rune is a pictographic character.
func IsPictographic(r rune) bool {
	return (r >= 0x4E00 && r <= 0x9FFF) || // CJK Unified Ideographs
		(r >= 0x3400 && r <= 0x4DBF) || // CJK Unified Ideographs Extension A
		(r >= 0x1F000 && r <= 0x1F9FF) // Various emoji blocks
}

// IsPrivateUse reports whether the rune is in the Unicode private use area.
func IsPrivateUse(r rune) bool {
	return (r >= 0xE000 && r <= 0xF8FF) || // Private Use Area
		(r >= 0xF0000 && r <= 0xFFFFD) || // Supplementary Private Use Area-A
		(r >= 0x100000 && r <= 0x10FFFD) // Supplementary Private Use Area-B
}

// ParseUnicodeData parses a UnicodeData.txt file from the Unicode Character Database.
func ParseUnicodeData(path string) (map[rune]UnicodeData, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data := make(map[rune]UnicodeData)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), ";")
		if len(fields) != 15 {
			continue // Skip malformed lines
		}

		// Parse code point
		cp, err := strconv.ParseInt(fields[0], 16, 32)
		if err != nil {
			continue
		}

		// Parse canonical class
		cc, _ := strconv.Atoi(fields[3])

		// Parse numeric value
		var numVal float64
		if fields[8] != "" {
			numVal, _ = strconv.ParseFloat(fields[8], 64)
		}

		// Parse bidi mirrored
		bidiMirrored := fields[9] == "Y"

		// Parse case mappings
		var upperCase, lowerCase, titleCase rune
		if fields[12] != "" {
			val, _ := strconv.ParseInt(fields[12], 16, 32)
			upperCase = rune(val)
		}
		if fields[13] != "" {
			val, _ := strconv.ParseInt(fields[13], 16, 32)
			lowerCase = rune(val)
		}
		if fields[14] != "" {
			val, _ := strconv.ParseInt(fields[14], 16, 32)
			titleCase = rune(val)
		}

		// Parse decomposition mapping
		var decomp []rune
		if len(fields[5]) > 0 {
			decompStr := fields[5]
			if strings.HasPrefix(decompStr, "<") {
				re := regexp.MustCompile(`[0-9A-F]+`)
				hexValues := re.FindAllString(decompStr, -1)
				for _, hex := range hexValues {
					val, _ := strconv.ParseInt(hex, 16, 32)
					decomp = append(decomp, rune(val))
				}
			} else {
				hexValues := strings.Fields(decompStr)
				for _, hex := range hexValues {
					val, _ := strconv.ParseInt(hex, 16, 32)
					decomp = append(decomp, rune(val))
				}
			}
		}

		// Store the data
		data[rune(cp)] = UnicodeData{
			CodePoint:       rune(cp),
			Name:            fields[1],
			GeneralCategory: fields[2],
			CanonicalClass:  cc,
			BidiClass:       fields[4],
			DecompType:      fields[5],
			DecompMapping:   decomp,
			NumericType:     fields[6],
			NumericValue:    numVal,
			BidiMirrored:    bidiMirrored,
			OldName:         fields[10],
			SimpleUpperCase: upperCase,
			SimpleLowerCase: lowerCase,
			SimpleTitleCase: titleCase,
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

// RTL represents right-to-left scripts like Arabic and Hebrew
var RTL = &unicode.RangeTable{
	R16: []unicode.Range16{
		{0x0590, 0x05FF, 1}, // Hebrew
		{0x0600, 0x06FF, 1}, // Arabic
		{0x0750, 0x077F, 1}, // Arabic Supplement
		{0x08A0, 0x08FF, 1}, // Arabic Extended-A
		{0xFB50, 0xFDFF, 1}, // Arabic Presentation Forms-A
		{0xFE70, 0xFEFF, 1}, // Arabic Presentation Forms-B
	},
	LatinOffset: 0,
}

// HasRTL reports whether the string contains any right-to-left characters.
func HasRTL(s string) bool {
	for _, r := range s {
		if unicode.In(r, RTL) {
			return true
		}
	}
	return false
}

// IsArabic reports whether the rune is an Arabic letter.
func IsArabic(r rune) bool {
	return r >= 0x0600 && r <= 0x06FF || r >= 0x0750 && r <= 0x077F
}

// IsHebrew reports whether the rune is a Hebrew letter.
func IsHebrew(r rune) bool {
	return r >= 0x0590 && r <= 0x05FF
}

// IsLatinExtended reports whether the rune is in the Latin Extended blocks.
func IsLatinExtended(r rune) bool {
	return (r >= 0x0100 && r <= 0x024F) || // Latin Extended-A/B
		(r >= 0x1E00 && r <= 0x1EFF) || // Latin Extended Additional
		(r >= 0x2C60 && r <= 0x2C7F) // Latin Extended-C
}

// IsThai reports whether the rune is a Thai character.
func IsThai(r rune) bool {
	return r >= 0x0E00 && r <= 0x0E7F
}

// IsHangul reports whether the rune is a Hangul character.
func IsHangul(r rune) bool {
	return (r >= 0xAC00 && r <= 0xD7A3) || // Hangul Syllables
		(r >= 0x1100 && r <= 0x11FF) // Hangul Jamo
}

// CountUniqueScripts counts how many different Unicode scripts are used in the string.
func CountUniqueScripts(s string) int {
	scripts := make(map[string]bool)

	for _, r := range s {
		switch {
		case unicode.Is(unicode.Latin, r):
			scripts["Latin"] = true
		case unicode.Is(unicode.Cyrillic, r):
			scripts["Cyrillic"] = true
		case unicode.Is(unicode.Greek, r):
			scripts["Greek"] = true
		case unicode.Is(unicode.Han, r):
			scripts["Han"] = true
		case IsHangul(r):
			scripts["Hangul"] = true
		case IsArabic(r):
			scripts["Arabic"] = true
		case IsHebrew(r):
			scripts["Hebrew"] = true
		case IsThai(r):
			scripts["Thai"] = true
		}
	}

	return len(scripts)
}

// Truncate truncates a string to the given max length, making sure not to break
// a grapheme cluster (character + combining marks).
func Truncate(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}

	// Truncate to maxLength bytes, but ensure we don't break a UTF-8 sequence
	truncated := s[:maxLength]
	for len(truncated) > 0 && !utf8.RuneStart(truncated[len(truncated)-1]) {
		truncated = truncated[:len(truncated)-1]
	}

	// Check if we need to trim a base character to avoid orphaned combining marks
	lastCharPos := 0
	for i, r := range truncated {
		if !unicode.Is(unicode.Mn, r) { // Mn = nonspacing marks
			lastCharPos = i
		}
	}

	nextIndex := len(truncated)
	if nextIndex < len(s) {
		r, _ := utf8.DecodeRuneInString(s[nextIndex:])
		if unicode.Is(unicode.Mn, r) {
			// Trim to the last non-combining character
			return truncated[:lastCharPos]
		}
	}

	return truncated
}
