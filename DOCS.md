# 5T3WsTools Documentation

This documentation provides explanations and examples for all functions in the 5T3WsTools library.

## Table of Contents

- [unique](#unique) - Utility functions for working with unique elements in slices
- [unicode](#unicode) - Extensions to the standard Go unicode package
- [time](#time) - Extensions to Go's standard time library
- [testing](#testing) - Extensions to the standard Go testing library
- [syscall](#syscall) - Additional functionality on top of the standard syscall library
- [sync](#sync) - Extensions to Go's standard sync package
- [structs](#structs) - Utility functions for working with struct types
- [strings](#strings) - Extensions to the standard strings library
- [strconv](#strconv) - Additional string conversion functionality
- [sort](#sort) - Extensions to Go's standard sort package
- [slices](#slices) - Utility functions that extend Go's standard slices package
- [runtime](#runtime) - Extensions to Go's standard runtime package
- [regexp](#regexp) - Extensions and utilities for the standard regexp package
- [reflect](#reflect) - Extensions to the standard library reflect package
- [plugin](#plugin) - Extensions and utilities for the standard plugin library
- [path](#path) - Extended functionality for path manipulation
- [os](#os) - Extensions to the standard os library
- [net](#net) - Extensions to Go's standard net and net/http packages
- [math](#math) - Additional mathematical functions
- [log](#log) - Extended functionality on top of Go's standard log package
- [iter](#iter) - Extensions to the Go standard library's iter package
- [io](#io) - Extensions to Go's standard io package
- [image](#image) - Extensions to the Go standard image library
- [html](#html) - Extensions to Go's standard html package
- [hash](#hash) - Hash utility functions
- [fmt](#fmt) - Extensions to the standard fmt package
- [flag](#flag) - Extensions to the standard flag package
- [expvar](#expvar) - Extensions to the standard expvar library
- [errors](#errors) - Extensions to the standard errors package
- [encoding](#encoding) - Additional encoding utilities
- [crypto](#crypto) - Extensions and utilities for Go's standard crypto libraries
- [context](#context) - Extensions to the standard library context package
- [cmp](#cmp) - Extended comparison functionality for Go values
- [bytes](#bytes) - Extensions to the standard library bytes package
- [builtin](#builtin) - Extensions to Go's builtin functions and types
- [bufio](#bufio) - Extensions to the standard bufio package

## unique

Package `unique` provides utility functions for working with unique elements in slices.

### Unq

```go
func Unq[T comparable](input []T) []T
```

Returns a new slice containing only the unique elements from the input slice, preserving order.

Example:
```go
import "github.com/5T3WsTools/unique"

func main() {
    numbers := []int{1, 2, 3, 2, 1, 4, 5}
    uniqueNumbers := unique.Unq(numbers)
    // uniqueNumbers is [1, 2, 3, 4, 5]
}
```

### IsUnique

```go
func IsUnique[T comparable](input []T) bool
```

Checks if all elements in the provided slice are unique.

Example:
```go
import "github.com/5T3WsTools/unique"

func main() {
    numbers1 := []int{1, 2, 3, 4, 5}
    allUnique := unique.IsUnique(numbers1) // true
    
    numbers2 := []int{1, 2, 3, 2, 5}
    allUnique = unique.IsUnique(numbers2) // false
}
```

### Intersection

```go
func Intersection[T comparable](a, b []T) []T
```

Returns a slice containing elements that exist in both input slices.

Example:
```go
import "github.com/5T3WsTools/unique"

func main() {
    slice1 := []int{1, 2, 3, 4}
    slice2 := []int{3, 4, 5, 6}
    common := unique.Intersection(slice1, slice2) // [3, 4]
}
```

### Union

```go
func Union[T comparable](a, b []T) []T
```

Returns a slice containing all unique elements from both input slices.

Example:
```go
import "github.com/5T3WsTools/unique"

func main() {
    slice1 := []int{1, 2, 3, 4}
    slice2 := []int{3, 4, 5, 6}
    allElements := unique.Union(slice1, slice2) // [1, 2, 3, 4, 5, 6]
}
```

### Count

```go
func Count[T comparable](input []T) map[T]int
```

Returns a map with counts of each element in the input slice.

Example:
```go
import "github.com/5T3WsTools/unique"

func main() {
    words := []string{"apple", "banana", "apple", "cherry", "banana", "apple"}
    counts := unique.Count(words)
    // counts is {"apple": 3, "banana": 2, "cherry": 1}
}
```

### Duplicates

```go
func Duplicates[T comparable](input []T) []T
```

Returns a slice containing only the elements that appear more than once in the input.

Example:
```go
import "github.com/5T3WsTools/unique"

func main() {
    words := []string{"apple", "banana", "apple", "cherry", "banana", "date"}
    dupes := unique.Duplicates(words) // ["apple", "banana"]
}
```

## unicode

Package `unicode` provides extensions to the standard Go unicode package.

### IsEmoji

```go
func IsEmoji(r rune) bool
```

Reports whether the rune is an emoji character.

Example:
```go
import "github.com/5T3WsTools/unicode"

func main() {
    isEmoji := unicode.IsEmoji('😀') // true
    isEmoji = unicode.IsEmoji('A')   // false
}
```

### IsPictographic

```go
func IsPictographic(r rune) bool
```

Reports whether the rune is a pictographic character.

Example:
```go
import "github.com/5T3WsTools/unicode"

func main() {
    isPic := unicode.IsPictographic('漢') // true (CJK ideograph)
    isPic = unicode.IsPictographic('A')   // false
}
```

### IsPrivateUse

```go
func IsPrivateUse(r rune) bool
```

Reports whether the rune is in the Unicode private use area.

Example:
```go
import "github.com/5T3WsTools/unicode"

func main() {
    isPrivate := unicode.IsPrivateUse('\uE000') // true
    isPrivate = unicode.IsPrivateUse('A')       // false
}
```

### ParseUnicodeData

```go
func ParseUnicodeData(path string) (map[rune]UnicodeData, error)
```

Parses a UnicodeData.txt file from the Unicode Character Database.

Example:
```go
import "github.com/5T3WsTools/unicode"

func main() {
    data, err := unicode.ParseUnicodeData("/path/to/UnicodeData.txt")
    if err != nil {
        // handle error
    }
    
    // Access information for a specific character
    charInfo := data['A']
    fmt.Println(charInfo.Name) // "LATIN CAPITAL LETTER A"
}
```

### HasRTL

```go
func HasRTL(s string) bool
```

Reports whether the string contains any right-to-left characters.

Example:
```go
import "github.com/5T3WsTools/unicode"

func main() {
    hasRTL := unicode.HasRTL("Hello") // false
    hasRTL = unicode.HasRTL("Hello مرحبا") // true (contains Arabic)
}
```

### IsArabic

```go
func IsArabic(r rune) bool
```

Reports whether the rune is an Arabic letter.

Example:
```go
import "github.com/5T3WsTools/unicode"

func main() {
    isArabic := unicode.IsArabic('ا') // true
    isArabic = unicode.IsArabic('A') // false
}
```

### IsHebrew

```go
func IsHebrew(r rune) bool
```

Reports whether the rune is a Hebrew letter.

Example:
```go
import "github.com/5T3WsTools/unicode"

func main() {
    isHebrew := unicode.IsHebrew('א') // true
    isHebrew = unicode.IsHebrew('A') // false
}
```

### CountUniqueScripts

```go
func CountUniqueScripts(s string) int
```

Counts how many different Unicode scripts are used in the string.

Example:
```go
import "github.com/5T3WsTools/unicode"

func main() {
    scriptCount := unicode.CountUniqueScripts("Hello") // 1
    scriptCount = unicode.CountUniqueScripts("Hello Привет") // 2 (Latin and Cyrillic)
}
```

### Truncate

```go
func Truncate(s string, maxLength int) string
```

Truncates a string to the given max length, ensuring it doesn't break a grapheme cluster.

Example:
```go
import "github.com/5T3WsTools/unicode"

func main() {
    truncated := unicode.Truncate("Hello world", 5) // "Hello"
    
    // With combining marks
    combining := "e\u0301" // "é" (e + combining acute accent)
    truncated = unicode.Truncate(combining, 1) // "" (avoids orphaning the combining mark)
}
```

## time

Package `time` extends the functionality of Go's standard time library.

### FormatRelative

```go
func FormatRelative(t stdtime.Time) string
```

Returns a string describing the time relative to now, like "5 minutes ago" or "in 2 days".

Example:
```go
import (
    "time"
    tm "github.com/5T3WsTools/time"
)

func main() {
    past := time.Now().Add(-10 * time.Minute)
    rel := tm.FormatRelative(past) // "10 minutes ago"
    
    future := time.Now().Add(24 * time.Hour)
    rel = tm.FormatRelative(future) // "tomorrow"
}
```

### ParseMultipleFormats

```go
func ParseMultipleFormats(str string, formats ...string) (stdtime.Time, error)
```

Attempts to parse a time string using multiple formats.

Example:
```go
import (
    "time"
    tm "github.com/5T3WsTools/time"
)

func main() {
    dateStr := "2023-01-15"
    formats := []string{
        "2006-01-02",
        "01/02/2006",
    }
    
    t, err := tm.ParseMultipleFormats(dateStr, formats...)
    if err != nil {
        // handle error
    }
    // t is 2023-01-15 00:00:00 +0000 UTC
}
```

### IsBusinessDay

```go
func IsBusinessDay(t stdtime.Time) bool
```

Returns true if the given time falls on a business day (Monday-Friday).

Example:
```go
import (
    "time"
    tm "github.com/5T3WsTools/time"
)

func main() {
    monday := time.Date(2023, 1, 16, 0, 0, 0, 0, time.UTC)
    isBusiness := tm.IsBusinessDay(monday) // true
    
    saturday := time.Date(2023, 1, 14, 0, 0, 0, 0, time.UTC)
    isBusiness = tm.IsBusinessDay(saturday) // false
}
```

### NextBusinessDay

```go
func NextBusinessDay(t stdtime.Time) stdtime.Time
```

Returns the next business day after the given time.

Example:
```go
import (
    "time"
    tm "github.com/5T3WsTools/time"
)

func main() {
    friday := time.Date(2023, 1, 13, 0, 0, 0, 0, time.UTC)
    next := tm.NextBusinessDay(friday) // 2023-01-16 (Monday)
    
    saturday := time.Date(2023, 1, 14, 0, 0, 0, 0, time.UTC)
    next = tm.NextBusinessDay(saturday) // 2023-01-16 (Monday)
}
```

### IsBetween

```go
func IsBetween(t, start, end stdtime.Time) bool
```

Returns true if the time is between start and end, inclusive.

Example:
```go
import (
    "time"
    tm "github.com/5T3WsTools/time"
)

func main() {
    start := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
    end := time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC)
    
    t := time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC)
    isBetween := tm.IsBetween(t, start, end) // true
    
    outside := time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)
    isBetween = tm.IsBetween(outside, start, end) // false
}
```

### StartOfMonth

```go
func StartOfMonth(t stdtime.Time) stdtime.Time
```

Returns the first day of the month containing the given time.

Example:
```go
import (
    "time"
    tm "github.com/5T3WsTools/time"
)

func main() {
    date := time.Date(2023, 1, 15, 10, 30, 0, 0, time.UTC)
    firstDay := tm.StartOfMonth(date) // 2023-01-01 00:00:00 +0000 UTC
}
```

### EndOfMonth

```go
func EndOfMonth(t stdtime.Time) stdtime.Time
```

Returns the last day of the month containing the given time.

Example:
```go
import (
    "time"
    tm "github.com/5T3WsTools/time"
)

func main() {
    date := time.Date(2023, 1, 15, 10, 30, 0, 0, time.UTC)
    lastDay := tm.EndOfMonth(date) // 2023-01-31 00:00:00 +0000 UTC
}
```

### Quarter

```go
func Quarter(t stdtime.Time) int
```

Returns the quarter (1-4) for the given time.

Example:
```go
import (
    "time"
    tm "github.com/5T3WsTools/time"
)

func main() {
    q1 := time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC)
    quarter := tm.Quarter(q1) // 1
    
    q3 := time.Date(2023, 8, 15, 0, 0, 0, 0, time.UTC)
    quarter = tm.Quarter(q3) // 3
}
```

### FormatDuration

```go
func FormatDuration(d stdtime.Duration) string
```

Formats a duration in a more human-readable way than the default.

Example:
```go
import (
    "time"
    tm "github.com/5T3WsTools/time"
)

func main() {
    d1 := 90 * time.Minute
    formatted := tm.FormatDuration(d1) // "1 hours 30 minutes"
    
    d2 := 48 * time.Hour
    formatted = tm.FormatDuration(d2) // "2 days"
}
```