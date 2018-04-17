package request

import (
	"bytes"
	"io"
	"testing"
)

func Test_ParseRequestLength_ConsidersFirstFourBytes(t *testing.T) {
	// Arrange
	var expected uint32
	expected = 201334791 // 0x0c002007

	// First 4 bytes represent the value of `expected` in Little Endian format,
	// the rest should be completely ignored during the parsing.
	input := bytes.NewReader([]byte{7, 32, 0, 12, 13, 13, 13})

	// Act
	actual, err := parseRequestLength(input)

	// Assert
	if err != nil {
		t.Fatalf("Error parsing request length: %v", err)
	}

	if expected != actual {
		t.Fatalf("The actual length '%v' does not match the expected value of '%v'", actual, expected)
	}
}

func Test_ParseRequestLength_ConnectionAborted(t *testing.T) {
	// Arrange
	expectedErr := io.ErrUnexpectedEOF
	input := bytes.NewReader([]byte{7})

	// Act
	_, err := parseRequestLength(input)

	// Assert
	if expectedErr != err {
		t.Fatalf("The expected error is '%v', but got '%v'", expectedErr, err)
	}
}
