// Package maps extends the functionality of Go's built-in maps package
package mapsExt

// Merge combines multiple maps into a new map. If keys overlap, later maps take precedence.
func Merge[K comparable, V any](maps ...map[K]V) map[K]V {
	result := make(map[K]V)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// GetOrDefault retrieves a value from a map by key, returning the default value if the key doesn't exist.
func GetOrDefault[K comparable, V any](m map[K]V, key K, defaultValue V) V {
	if val, ok := m[key]; ok {
		return val
	}
	return defaultValue
}

// Keys extracts all keys from a map into a slice.
func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Values extracts all values from a map into a slice.
func Values[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

// Filter creates a new map containing only the key-value pairs that satisfy the predicate function.
func Filter[K comparable, V any](m map[K]V, predicate func(K, V) bool) map[K]V {
	result := make(map[K]V)
	for k, v := range m {
		if predicate(k, v) {
			result[k] = v
		}
	}
	return result
}

// MapValues transforms each value in the source map using the transformer function and returns a new map with the same keys.
func MapValues[K comparable, V1 any, V2 any](m map[K]V1, transformer func(V1) V2) map[K]V2 {
	result := make(map[K]V2)
	for k, v := range m {
		result[k] = transformer(v)
	}
	return result
}

// Transform creates a new map by applying a transformation function to each key-value pair.
func Transform[K1 comparable, V1 any, K2 comparable, V2 any](m map[K1]V1, transformer func(K1, V1) (K2, V2)) map[K2]V2 {
	result := make(map[K2]V2)
	for k, v := range m {
		newK, newV := transformer(k, v)
		result[newK] = newV
	}
	return result
}

// Copy creates a shallow copy of a map.
func Copy[K comparable, V any](m map[K]V) map[K]V {
	result := make(map[K]V, len(m))
	for k, v := range m {
		result[k] = v
	}
	return result
}

// Equal checks if two maps have the same key-value pairs.
func Equal[K comparable, V comparable](m1, m2 map[K]V) bool {
	if len(m1) != len(m2) {
		return false
	}
	for k, v1 := range m1 {
		if v2, ok := m2[k]; !ok || v1 != v2 {
			return false
		}
	}
	return true
}

// Difference returns two maps:
// - added: keys in m2 that are not in m1 or have different values
// - removed: keys in m1 that are not in m2
func Difference[K comparable, V comparable](m1, m2 map[K]V) (added map[K]V, removed map[K]V) {
	added = make(map[K]V)
	removed = make(map[K]V)

	for k, v1 := range m1 {
		if v2, ok := m2[k]; !ok {
			removed[k] = v1
		} else if v1 != v2 {
			added[k] = v2
		}
	}

	for k, v2 := range m2 {
		if _, ok := m1[k]; !ok {
			added[k] = v2
		}
	}

	return added, removed
}

// HasKey checks if a map contains a given key.
func HasKey[K comparable, V any](m map[K]V, key K) bool {
	_, ok := m[key]
	return ok
}

// HasValue checks if a map contains a given value.
func HasValue[K comparable, V comparable](m map[K]V, value V) bool {
	for _, v := range m {
		if v == value {
			return true
		}
	}
	return false
}

// Invert creates a new map with keys and values swapped.
// Note: If multiple keys map to the same value, one will be chosen arbitrarily.
func Invert[K comparable, V comparable](m map[K]V) map[V]K {
	result := make(map[V]K, len(m))
	for k, v := range m {
		result[v] = k
	}
	return result
}

// KeyValuePair represents a key-value pair from a map
type KeyValuePair[K comparable, V any] struct {
	Key   K
	Value V
}

// ToSlice converts a map to a slice of key-value pairs.
func ToSlice[K comparable, V any](m map[K]V) []KeyValuePair[K, V] {
	result := make([]KeyValuePair[K, V], 0, len(m))
	for k, v := range m {
		result = append(result, KeyValuePair[K, V]{k, v})
	}
	return result
}

// FromSlice creates a map from a slice of key-value pairs.
func FromSlice[K comparable, V any](pairs []KeyValuePair[K, V]) map[K]V {
	result := make(map[K]V, len(pairs))
	for _, pair := range pairs {
		result[pair.Key] = pair.Value
	}
	return result
}

// DeleteKeys removes specified keys from a map and returns the map.
func DeleteKeys[K comparable, V any](m map[K]V, keys ...K) map[K]V {
	for _, key := range keys {
		delete(m, key)
	}
	return m
}
