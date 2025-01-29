// Package conv provides utilities for efficient and unsafe conversions between
// strings and byte slices in Go. These conversions are designed for performance-critical
// scenarios where avoiding data copying is essential.
//
// WARNING: The methods in this package rely on unsafe operations and should be used
// with caution. Modifying the underlying data after conversion can result in undefined behavior.
//
// Note: For safer alternatives that ensure immutability, consider using the standard
// string-to-byte-slice and byte-slice-to-string conversions provided by Go.
package conv

import "unsafe"

// StrToBytes converts a string to a byte slice without copying data.
// The returned []byte shares the same underlying memory as the input string.
// WARNING: Modifying the []byte can lead to undefined behavior as strings are immutable in Go.
// Use only in performance-critical scenarios where immutability can be guaranteed.
func StrToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// BytesToStr converts a byte slice to a string without copying data.
// The returned string shares the same underlying memory as the input []byte.
// WARNING: The input []byte mustn't be modified after the conversion, as strings are immutable.
// Use only when you can ensure the byte slice's immutability after conversion.
func BytesToStr(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
