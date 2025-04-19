// Package cmp provides extended comparison functionality for Go values.
package cmpExt

import (
	"math"
	"reflect"
	"sort"
)

// Equal performs a deep equality check between two values.
func Equal(x, y interface{}) bool {
	return reflect.DeepEqual(x, y)
}

// FloatEqual compares two float64 values with a specified tolerance.
func FloatEqual(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}

// IntSliceEqual checks if two int slices contain the same elements,
// regardless of their order.
func IntSliceEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	// Create copies to avoid modifying the originals
	aCopy := make([]int, len(a))
	bCopy := make([]int, len(b))

	copy(aCopy, a)
	copy(bCopy, b)

	// Sort both slices
	sort.Ints(aCopy)
	sort.Ints(bCopy)

	// Compare sorted slices
	for i := range aCopy {
		if aCopy[i] != bCopy[i] {
			return false
		}
	}
	return true
}

// StringSliceEqual checks if two string slices contain the same elements,
// regardless of their order.
func StringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	// Create copies to avoid modifying the originals
	aCopy := make([]string, len(a))
	bCopy := make([]string, len(b))

	copy(aCopy, a)
	copy(bCopy, b)

	// Sort both slices
	sort.Strings(aCopy)
	sort.Strings(bCopy)

	// Compare sorted slices
	for i := range aCopy {
		if aCopy[i] != bCopy[i] {
			return false
		}
	}
	return true
}

// EqualWithComparator compares two values using a custom comparison function.
func EqualWithComparator(a, b interface{}, comparator func(a, b interface{}) bool) bool {
	return comparator(a, b)
}

// StructFieldEqual compares two structs based on specific field names.
func StructFieldEqual(a, b interface{}, fieldNames ...string) bool {
	aVal := reflect.ValueOf(a)
	bVal := reflect.ValueOf(b)

	// Get the proper value if it's a pointer
	if aVal.Kind() == reflect.Ptr {
		aVal = aVal.Elem()
	}
	if bVal.Kind() == reflect.Ptr {
		bVal = bVal.Elem()
	}

	// Check if both are structs
	if aVal.Kind() != reflect.Struct || bVal.Kind() != reflect.Struct {
		return false
	}

	// Check each specified field
	for _, fieldName := range fieldNames {
		aField := aVal.FieldByName(fieldName)
		bField := bVal.FieldByName(fieldName)

		// Check if fields exist
		if !aField.IsValid() || !bField.IsValid() {
			return false
		}

		// Compare field values
		if !reflect.DeepEqual(aField.Interface(), bField.Interface()) {
			return false
		}
	}

	return true
}

// MapEqual checks if two maps contain the same key-value pairs.
func MapEqual(a, b map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}

	for k, aVal := range a {
		bVal, exists := b[k]
		if !exists {
			return false
		}
		if !reflect.DeepEqual(aVal, bVal) {
			return false
		}
	}

	return true
}

// Contains checks if slice contains a specific value.
func Contains(slice interface{}, val interface{}) bool {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return false
	}

	for i := 0; i < s.Len(); i++ {
		if reflect.DeepEqual(s.Index(i).Interface(), val) {
			return true
		}
	}
	return false
}
