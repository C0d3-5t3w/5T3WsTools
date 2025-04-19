// Package slices provides utility functions that extend Go's standard slices package.
package slicesExt

import (
	"slices"
)

// Filter returns a new slice containing only the elements of s for which keep returns true.
func Filter[E any](s []E, keep func(E) bool) []E {
	var result []E
	for _, v := range s {
		if keep(v) {
			result = append(result, v)
		}
	}
	return result
}

// Map returns a new slice containing the results of applying the function f to each element of s.
func Map[E, T any](s []E, f func(E) T) []T {
	result := make([]T, len(s))
	for i, v := range s {
		result[i] = f(v)
	}
	return result
}

// Reduce applies the function f to each element of s, accumulating a result.
func Reduce[E, T any](s []E, initial T, f func(T, E) T) T {
	result := initial
	for _, v := range s {
		result = f(result, v)
	}
	return result
}

// Chunk splits a slice into chunks of the specified size.
// The last chunk may contain fewer elements than size.
func Chunk[E any](s []E, size int) [][]E {
	if size <= 0 {
		panic("Chunk size must be positive")
	}

	result := make([][]E, 0, (len(s)+size-1)/size)
	for i := 0; i < len(s); i += size {
		end := i + size
		if end > len(s) {
			end = len(s)
		}
		result = append(result, s[i:end])
	}
	return result
}

// Unique returns a new slice with duplicate elements removed.
// The order of elements is preserved.
func Unique[E comparable](s []E) []E {
	if len(s) <= 1 {
		return slices.Clone(s)
	}

	seen := make(map[E]struct{}, len(s))
	result := make([]E, 0, len(s))

	for _, v := range s {
		if _, exists := seen[v]; !exists {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

// Intersect returns a slice of elements that appear in all provided slices.
func Intersect[E comparable](slices ...[]E) []E {
	if len(slices) == 0 {
		return nil
	}

	counts := make(map[E]int)
	for _, slice := range slices {
		seen := make(map[E]struct{})
		for _, v := range slice {
			if _, exists := seen[v]; !exists {
				seen[v] = struct{}{}
				counts[v]++
			}
		}
	}

	var result []E
	target := len(slices)
	for v, count := range counts {
		if count == target {
			result = append(result, v)
		}
	}

	return result
}

// Union returns a slice containing all unique elements from all provided slices.
func Union[E comparable](slices ...[]E) []E {
	var combined []E
	for _, s := range slices {
		combined = append(combined, s...)
	}
	return Unique(combined)
}

// Difference returns a slice of elements in s1 that don't appear in s2.
func Difference[E comparable](s1, s2 []E) []E {
	exclude := make(map[E]struct{})
	for _, v := range s2 {
		exclude[v] = struct{}{}
	}

	var result []E
	for _, v := range s1 {
		if _, exists := exclude[v]; !exists {
			result = append(result, v)
		}
	}
	return result
}

// All returns true if the predicate returns true for all elements in the slice.
func All[E any](s []E, predicate func(E) bool) bool {
	for _, v := range s {
		if !predicate(v) {
			return false
		}
	}
	return true
}

// Any returns true if the predicate returns true for any element in the slice.
func Any[E any](s []E, predicate func(E) bool) bool {
	for _, v := range s {
		if predicate(v) {
			return true
		}
	}
	return false
}

// GroupBy groups elements in a slice by keys returned by the keyFunc.
func GroupBy[E any, K comparable](s []E, keyFunc func(E) K) map[K][]E {
	result := make(map[K][]E)
	for _, v := range s {
		key := keyFunc(v)
		result[key] = append(result[key], v)
	}
	return result
}

// ForEach applies the function f to each element of s.
func ForEach[E any](s []E, f func(E)) {
	for _, v := range s {
		f(v)
	}
}
