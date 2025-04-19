// Package iter provides extensions to the Go standard library's iter package.
// It offers additional utility functions for working with sequences and iterators.
package iterExt

import (
	"cmp"
	"iter"
)

// Map transforms each element from the sequence using the provided function.
func Map[T, R any](seq iter.Seq[T], fn func(T) R) iter.Seq[R] {
	return func(yield func(R) bool) {
		seq(func(v T) bool {
			return yield(fn(v))
		})
	}
}

// Filter returns a sequence containing only elements matching the predicate.
func Filter[T any](seq iter.Seq[T], predicate func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		seq(func(v T) bool {
			if predicate(v) {
				if !yield(v) {
					return false
				}
			}
			return true
		})
	}
}

// Reduce combines all elements in the sequence into a single value.
func Reduce[T, R any](seq iter.Seq[T], initial R, reducer func(R, T) R) R {
	result := initial
	seq(func(v T) bool {
		result = reducer(result, v)
		return true
	})
	return result
}

// ForEach applies the given function to each element in the sequence.
func ForEach[T any](seq iter.Seq[T], fn func(T)) {
	seq(func(v T) bool {
		fn(v)
		return true
	})
}

// Count returns the number of elements in the sequence.
func Count[T any](seq iter.Seq[T]) int {
	count := 0
	seq(func(v T) bool {
		count++
		return true
	})
	return count
}

// Any returns true if any element satisfies the predicate.
func Any[T any](seq iter.Seq[T], predicate func(T) bool) bool {
	result := false
	seq(func(v T) bool {
		if predicate(v) {
			result = true
			return false // Stop iteration
		}
		return true
	})
	return result
}

// All returns true if all elements satisfy the predicate.
func All[T any](seq iter.Seq[T], predicate func(T) bool) bool {
	result := true
	seq(func(v T) bool {
		if !predicate(v) {
			result = false
			return false // Stop iteration
		}
		return true
	})
	return result
}

// Find returns the first element that satisfies the predicate, along with a
// boolean indicating whether such an element was found.
func Find[T any](seq iter.Seq[T], predicate func(T) bool) (T, bool) {
	var result T
	found := false

	seq(func(v T) bool {
		if predicate(v) {
			result = v
			found = true
			return false // Stop iteration
		}
		return true
	})

	return result, found
}

// Take returns a sequence containing at most n elements from the original sequence.
func Take[T any](seq iter.Seq[T], n int) iter.Seq[T] {
	return func(yield func(T) bool) {
		count := 0
		seq(func(v T) bool {
			if count >= n {
				return false
			}
			count++
			return yield(v)
		})
	}
}

// Skip returns a sequence without the first n elements of the original sequence.
func Skip[T any](seq iter.Seq[T], n int) iter.Seq[T] {
	return func(yield func(T) bool) {
		count := 0
		seq(func(v T) bool {
			if count < n {
				count++
				return true
			}
			return yield(v)
		})
	}
}

// Min returns the minimum element in the sequence according to the provided less function.
func Min[T any](seq iter.Seq[T], less func(a, b T) bool) (T, bool) {
	var min T
	found := false

	seq(func(v T) bool {
		if !found || less(v, min) {
			min = v
			found = true
		}
		return true
	})

	return min, found
}

// Max returns the maximum element in the sequence according to the provided less function.
func Max[T any](seq iter.Seq[T], less func(a, b T) bool) (T, bool) {
	var max T
	found := false

	seq(func(v T) bool {
		if !found || less(max, v) {
			max = v
			found = true
		}
		return true
	})

	return max, found
}

// OrderedMin returns the minimum element in a sequence of comparable elements.
func OrderedMin[T cmp.Ordered](seq iter.Seq[T]) (T, bool) {
	return Min(seq, func(a, b T) bool { return a < b })
}

// OrderedMax returns the maximum element in a sequence of comparable elements.
func OrderedMax[T cmp.Ordered](seq iter.Seq[T]) (T, bool) {
	return Max(seq, func(a, b T) bool { return a < b })
}
