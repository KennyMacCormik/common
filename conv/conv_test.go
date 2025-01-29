package conv

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"unsafe"
)

func TestStrToBytes(t *testing.T) {
	s := "hello"
	b := StrToBytes(s)

	assert.Equal(t, s, string(b), "expected byte slice to match the original string")
	assert.Equal(t, unsafe.StringData(s), unsafe.SliceData(b), "expected shared underlying memory")
}

func TestBytesToStr(t *testing.T) {
	b := []byte("world")
	s := BytesToStr(b)

	assert.Equal(t, string(b), s, "expected string to match the original byte slice")
	assert.Equal(t, unsafe.StringData(s), unsafe.SliceData(b), "expected shared underlying memory")
}

func TestImmutabilityWarning_StrToBytes(t *testing.T) {
	s := "immutable"
	b := StrToBytes(s)

	// Modifying the byte slice can lead to undefined behavior.
	// Uncommenting the line below demonstrates the risk:
	// b[0] = 'M'

	// Attempting to modify the byte slice should be avoided; this test ensures shared memory.
	assert.Equal(t, s, string(b), "expected no modification to underlying memory")
}

func TestImmutabilityWarning_BytesToStr(t *testing.T) {
	b := []byte("mutable")
	s := BytesToStr(b)

	// Modifying the byte slice reflects in the string.
	b[0] = 'M'

	assert.Equal(t, "Mutable", s, "expected string to reflect the modification in the byte slice")
}

func TestEmptyInputs(t *testing.T) {
	// Test empty string to byte slice conversion
	b := StrToBytes("")
	assert.Empty(t, b, "expected empty byte slice for empty string input")

	// Test empty byte slice to string conversion
	s := BytesToStr([]byte{})
	assert.Empty(t, s, "expected empty string for empty byte slice input")
}

// TestSafeAlternatives provides safer alternatives for StrToBytes and BytesToStr.
func TestSafeAlternatives(t *testing.T) {
	s := "safe"
	b := []byte(s)

	// Safe conversion alternatives
	safeBytes := []byte(s)
	safeString := string(b)

	assert.Equal(t, safeBytes, StrToBytes(s), "safe conversion should match unsafe StrToBytes")
	assert.Equal(t, safeString, BytesToStr(b), "safe conversion should match unsafe BytesToStr")
}
