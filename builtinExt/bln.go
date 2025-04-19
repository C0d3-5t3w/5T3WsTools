// Package builtin provides extensions to Go's builtin functions and types.
// It offers helpers for common operations on slices, maps, and basic types
// that are not included in the standard library.
package builtinExt

import (
	"fmt"
	"reflect"
	"strings"
)

// SliceOperations provides utility functions for slices
type SliceOperations[T comparable] struct{}

// MapOperations provides utility functions for maps
type MapOperations[K comparable, V any] struct{}

// StringOperations provides utility functions for strings
type StringOperations struct{}

// Slice returns a singleton instance of SliceOperations
func Slice[T comparable]() SliceOperations[T] {
	return SliceOperations[T]{}
}

// Map returns a singleton instance of MapOperations
func Map[K comparable, V any]() MapOperations[K, V] {
	return MapOperations[K, V]{}
}

// String returns a singleton instance of StringOperations
func String() StringOperations {
	return StringOperations{}
}

// Contains checks if a slice contains a specific element
func (SliceOperations[T]) Contains(slice []T, element T) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}

// Unique returns a new slice with duplicate elements removed
func (SliceOperations[T]) Unique(slice []T) []T {
	if len(slice) == 0 {
		return slice
	}

	seen := make(map[T]struct{}, len(slice))
	result := make([]T, 0, len(slice))

	for _, item := range slice {
		if _, ok := seen[item]; !ok {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}

// Filter returns a new slice containing only elements that satisfy the predicate
func (SliceOperations[T]) Filter(slice []T, predicate func(T) bool) []T {
	result := make([]T, 0, len(slice))
	for _, item := range slice {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

// Map applies a function to each element in a slice and returns a new slice
// Converting from method with type parameters to standalone function
func SliceMap[T any, R any](slice []T, mapper func(T) R) []R {
	result := make([]R, len(slice))
	for i, item := range slice {
		result[i] = mapper(item)
	}
	return result
}

// Keys returns all keys from a map as a slice
func (MapOperations[K, V]) Keys(m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Values returns all values from a map as a slice
func (MapOperations[K, V]) Values(m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

// Merge combines two maps into a new one. In case of duplicate keys, values from the second map take precedence
func (MapOperations[K, V]) Merge(m1, m2 map[K]V) map[K]V {
	result := make(map[K]V, len(m1)+len(m2))

	for k, v := range m1 {
		result[k] = v
	}

	for k, v := range m2 {
		result[k] = v
	}

	return result
}

// SplitAndTrim splits a string and trims whitespace from each part
func (StringOperations) SplitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}
	return parts
}

// ToMap converts a slice of key-value pairs into a map
func (StringOperations) ToMap(pairs []string, sep string) map[string]string {
	result := make(map[string]string, len(pairs))

	for _, pair := range pairs {
		parts := strings.SplitN(pair, sep, 2)
		if len(parts) == 2 {
			result[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	return result
}

// IsEmptyOrWhitespace checks if a string is empty or contains only whitespace
func (StringOperations) IsEmptyOrWhitespace(s string) bool {
	return strings.TrimSpace(s) == ""
}

// Must panics if err is non-nil, otherwise returns value
func Must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}

// DefaultIfZero returns the default value if the given value is the zero value for its type
func DefaultIfZero[T comparable](value, defaultValue T) T {
	var zero T
	if value == zero {
		return defaultValue
	}
	return value
}

// IsZero checks if a value is the zero value for its type
func IsZero(v interface{}) bool {
	return reflect.ValueOf(v).IsZero()
}

// ToString converts various types to their string representation
func ToString(v interface{}) string {
	return fmt.Sprintf("%v", v)
}
