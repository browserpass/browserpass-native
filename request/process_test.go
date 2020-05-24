package request

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"testing"
	"os"
	"fmt"
	"path/filepath"
)

func Test_ParseRequestLength_ConsidersFirstFourBytes(t *testing.T) {
	// Arrange
	expected := uint32(201334791) // 0x0c002007

	// The first 4 bytes represent the value of `expected` in Little Endian format,
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

func Test_ParseRequest_CanParse(t *testing.T) {
	// Arrange
	expected := &request{
		Action: "list",
		Settings: settings{
			Stores: map[string]store{
				"id1": store{
					ID:   "id1",
					Name: "default",
					Path: "~/.password-store",
				},
			},
		},
	}

	jsonBytes, err := json.Marshal(expected)
	if err != nil {
		t.Fatal("Unable to marshal the expected object to initialize the test")
	}

	inputLength := uint32(len(jsonBytes))
	input := bytes.NewReader(jsonBytes)

	// Act
	actual, err := parseRequest(inputLength, input)

	// Assert
	if err != nil {
		t.Fatalf("Error parsing request: %v", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("The request was parsed incorrectly.\nExpected: %+v\nActual:   %+v", expected, actual)
	}
}

func Test_ParseRequest_WrongLength(t *testing.T) {
	// Arrange
	expectedErr := io.ErrUnexpectedEOF

	jsonBytes, err := json.Marshal(&request{Action: "list"})
	if err != nil {
		t.Fatal("Unable to marshal the expected object to initialize the test")
	}

	wrongInputLength := uint32(len(jsonBytes)) - 1
	input := bytes.NewReader(jsonBytes)

	// Act
	_, err = parseRequest(wrongInputLength, input)

	// Assert
	if expectedErr != err {
		t.Fatalf("The expected error is '%v', but got '%v'", expectedErr, err)
	}
}

func Test_ParseRequest_InvalidJson(t *testing.T) {
	// Arrange
	jsonBytes := []byte("not_a_json")
	inputLength := uint32(len(jsonBytes))
	input := bytes.NewReader(jsonBytes)

	// Act
	_, err := parseRequest(inputLength, input)

	// Assert
	if err == nil {
		t.Fatalf("Expected a parsing error, but didn't get it")
	}
}


func Test_ParseRequest_GetDefaultPasswordStorePath(t *testing.T) {
	// Arrange
	classicPath := fmt.Sprintf("/home/%s/.password-store", os.Getenv("USER"))
	xdgDataPath := fmt.Sprintf("/home/%s/.local", os.Getenv("USER"))
	xdgPath := filepath.Join(xdgDataPath, "password-store")
	os.Unsetenv("XDG_DATA_HOME")

	// Act
	path, _ := getDefaultPasswordStorePath()

	// Assert
	if path != classicPath {
		t.Fatalf("Expected '%s', got '%s'", classicPath, path)
	}

	// Arrange
	os.Setenv("XDG_DATA_HOME", xdgDataPath)

	// Act
	path, _ = getDefaultPasswordStorePath()

	// Assert
	if path != xdgPath {
		t.Fatalf("Expected '%s', got '%s'", xdgPath, path)
	}
}
