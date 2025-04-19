// Package srt extends the functionality of Go's standard sort package
package srt

import (
	"sort"
	"sync"
)

// SortBy sorts a slice using a custom comparator function
// The less function should return true if element i should sort before element j
func SortBy[T any](data []T, less func(i, j T) bool) {
	sort.Slice(data, func(i, j int) bool {
		return less(data[i], data[j])
	})
}

// SortByKey sorts a slice by a specific key extracted from each element
func SortByKey[T any, K any](data []T, key func(T) K, lessKey func(a, b K) bool) {
	sort.Slice(data, func(i, j int) bool {
		return lessKey(key(data[i]), key(data[j]))
	})
}

// SortByKeyOrdered sorts a slice by a specific key that implements Ordered interface
func SortByKeyOrdered[T any, K Ordered](data []T, key func(T) K) {
	sort.Slice(data, func(i, j int) bool {
		return key(data[i]) < key(data[j])
	})
}

// Ordered is a constraint that permits any ordered type: numbers, strings, etc.
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string
}

// ParallelSort performs a parallel merge sort for large slices
// It's more efficient than standard sort for large datasets
func ParallelSort[T any](data []T, less func(i, j T) bool, parallelism int) {
	if len(data) < 2 {
		return
	}

	if parallelism <= 1 || len(data) < 1000 {
		SortBy(data, less)
		return
	}

	mid := len(data) / 2
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		ParallelSort(data[:mid], less, parallelism/2)
	}()
	ParallelSort(data[mid:], less, parallelism/2)
	wg.Wait()

	// Merge the sorted halves
	merge(data, mid, less)
}

// merge combines two sorted slices into one
func merge[T any](data []T, mid int, less func(i, j T) bool) {
	temp := make([]T, len(data))
	copy(temp, data)

	i, j, k := 0, mid, 0
	for i < mid && j < len(data) {
		if less(temp[i], temp[j]) {
			data[k] = temp[i]
			i++
		} else {
			data[k] = temp[j]
			j++
		}
		k++
	}

	for i < mid {
		data[k] = temp[i]
		i++
		k++
	}
}

// IsSorted checks if a slice is sorted according to a comparator
func IsSorted[T any](data []T, less func(i, j T) bool) bool {
	for i := 0; i < len(data)-1; i++ {
		if less(data[i+1], data[i]) {
			return false
		}
	}
	return true
}

// BinarySearch performs a binary search on a sorted slice
// Returns index and true if found, insertion point and false if not found
func BinarySearch[T any](data []T, target T, less func(a, b T) bool) (int, bool) {
	low, high := 0, len(data)-1
	for low <= high {
		mid := low + (high-low)/2
		if less(data[mid], target) {
			low = mid + 1
		} else if less(target, data[mid]) {
			high = mid - 1
		} else {
			return mid, true
		}
	}
	return low, false
}

// Deduplicate removes duplicate elements from a sorted slice in-place
// Returns the new length of the slice
func Deduplicate[T comparable](data []T) int {
	if len(data) < 2 {
		return len(data)
	}

	j := 1
	for i := 1; i < len(data); i++ {
		if data[i] != data[i-1] {
			if i != j {
				data[j] = data[i]
			}
			j++
		}
	}

	return j
}
