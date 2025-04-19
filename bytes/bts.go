// Package bytes provides utilities that extend the Go standard library's bytes package
package bytes

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"
)

// ErrInsufficientBytes is returned when a byte slice doesn't have enough bytes for conversion
var ErrInsufficientBytes = errors.New("insufficient bytes for conversion")

// ToUint16 converts a byte slice to uint16 using the specified byte order
func ToUint16(b []byte, order binary.ByteOrder) (uint16, error) {
	if len(b) < 2 {
		return 0, ErrInsufficientBytes
	}
	return order.Uint16(b), nil
}

// ToUint32 converts a byte slice to uint32 using the specified byte order
func ToUint32(b []byte, order binary.ByteOrder) (uint32, error) {
	if len(b) < 4 {
		return 0, ErrInsufficientBytes
	}
	return order.Uint32(b), nil
}

// ToUint64 converts a byte slice to uint64 using the specified byte order
func ToUint64(b []byte, order binary.ByteOrder) (uint64, error) {
	if len(b) < 8 {
		return 0, ErrInsufficientBytes
	}
	return order.Uint64(b), nil
}

// ToFloat32 converts a byte slice to float32 using the specified byte order
func ToFloat32(b []byte, order binary.ByteOrder) (float32, error) {
	if len(b) < 4 {
		return 0, ErrInsufficientBytes
	}
	return math.Float32frombits(order.Uint32(b)), nil
}

// ToFloat64 converts a byte slice to float64 using the specified byte order
func ToFloat64(b []byte, order binary.ByteOrder) (float64, error) {
	if len(b) < 8 {
		return 0, ErrInsufficientBytes
	}
	return math.Float64frombits(order.Uint64(b)), nil
}

// FromUint16 converts uint16 to a byte slice using the specified byte order
func FromUint16(v uint16, order binary.ByteOrder) []byte {
	b := make([]byte, 2)
	order.PutUint16(b, v)
	return b
}

// FromUint32 converts uint32 to a byte slice using the specified byte order
func FromUint32(v uint32, order binary.ByteOrder) []byte {
	b := make([]byte, 4)
	order.PutUint32(b, v)
	return b
}

// FromUint64 converts uint64 to a byte slice using the specified byte order
func FromUint64(v uint64, order binary.ByteOrder) []byte {
	b := make([]byte, 8)
	order.PutUint64(b, v)
	return b
}

// FromFloat32 converts float32 to a byte slice using the specified byte order
func FromFloat32(v float32, order binary.ByteOrder) []byte {
	b := make([]byte, 4)
	order.PutUint32(b, math.Float32bits(v))
	return b
}

// FromFloat64 converts float64 to a byte slice using the specified byte order
func FromFloat64(v float64, order binary.ByteOrder) []byte {
	b := make([]byte, 8)
	order.PutUint64(b, math.Float64bits(v))
	return b
}

// SafeSlice returns a slice of the specified length from the input bytes,
// or a smaller slice if the input is not long enough
func SafeSlice(b []byte, start, length int) []byte {
	if start < 0 {
		start = 0
	}
	if start >= len(b) {
		return []byte{}
	}
	end := start + length
	if end > len(b) {
		end = len(b)
	}
	return b[start:end]
}

// Repeat creates a new byte slice by repeating the input slice n times
func Repeat(b []byte, n int) []byte {
	if n <= 0 {
		return []byte{}
	}
	result := make([]byte, len(b)*n)
	copy(result, b)
	for i := 1; i < n; i++ {
		copy(result[i*len(b):], b)
	}
	return result
}

// ReverseInPlace reverses a byte slice in place
func ReverseInPlace(b []byte) {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
}

// Reverse returns a new byte slice with the elements in reverse order
func Reverse(b []byte) []byte {
	result := make([]byte, len(b))
	for i, j := 0, len(b)-1; j >= 0; i, j = i+1, j-1 {
		result[i] = b[j]
	}
	return result
}

// IsASCII checks if all bytes in the slice are valid ASCII characters
func IsASCII(b []byte) bool {
	for _, c := range b {
		if c > 127 {
			return false
		}
	}
	return true
}

// ContainsAny returns true if the byte slice contains any of the bytes in chars
func ContainsAny(b []byte, chars []byte) bool {
	for _, c := range chars {
		if bytes.IndexByte(b, c) >= 0 {
			return true
		}
	}
	return false
}

// RemoveAll returns a copy with all occurrences of the given bytes removed
func RemoveAll(s, chars []byte) []byte {
	result := make([]byte, 0, len(s))
	for _, c := range s {
		if !bytes.Contains(chars, []byte{c}) {
			result = append(result, c)
		}
	}
	return result
}
