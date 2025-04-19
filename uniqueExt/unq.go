// Package unique provides utility functions for working with unique elements in slices.
package uniqueExt

// Unq returns a new slice containing only the unique comparable elements from the input slice.
// The order of elements is preserved (first occurrence kept).
func Unq[T comparable](input []T) []T {
	seen := make(map[T]struct{})
	result := []T{}

	for _, item := range input {
		if _, ok := seen[item]; !ok {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}

// IsUnique checks if all elements in the provided slice are unique.
// Returns true if all elements are unique, false otherwise.
func IsUnique[T comparable](input []T) bool {
	seen := make(map[T]struct{})

	for _, item := range input {
		if _, ok := seen[item]; ok {
			return false
		}
		seen[item] = struct{}{}
	}

	return true
}

// Intersection returns a slice containing elements that exist in both input slices.
// Only unique elements are returned.
func Intersection[T comparable](a, b []T) []T {
	setA := make(map[T]struct{})
	result := []T{}

	for _, item := range a {
		setA[item] = struct{}{}
	}

	seen := make(map[T]struct{})
	for _, item := range b {
		if _, ok := setA[item]; ok {
			if _, alreadySeen := seen[item]; !alreadySeen {
				seen[item] = struct{}{}
				result = append(result, item)
			}
		}
	}

	return result
}

// Union returns a slice containing all unique elements from both input slices.
func Union[T comparable](a, b []T) []T {
	seen := make(map[T]struct{})
	result := []T{}

	// Process slice a
	for _, item := range a {
		if _, ok := seen[item]; !ok {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}

	// Process slice b
	for _, item := range b {
		if _, ok := seen[item]; !ok {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}

// Count returns a map with counts of each element in the input slice.
func Count[T comparable](input []T) map[T]int {
	counts := make(map[T]int)
	for _, item := range input {
		counts[item]++
	}
	return counts
}

// Duplicates returns a slice containing only the elements that appear more than once in the input.
// Each duplicate is included only once in the result.
func Duplicates[T comparable](input []T) []T {
	counts := Count(input)
	result := []T{}
	seen := make(map[T]struct{})

	for item, count := range counts {
		if count > 1 {
			if _, ok := seen[item]; !ok {
				seen[item] = struct{}{}
				result = append(result, item)
			}
		}
	}

	return result
}
