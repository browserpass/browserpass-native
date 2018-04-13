package errors

import (
	"os"
)

// Code exit code
type Code int

const (
	// CodeParseRequestLength error parsing request length
	CodeParseRequestLength Code = iota + 10

	// CodeParseRequest error parsing request
	CodeParseRequest
)

// ExitWithCode exit with error code
func ExitWithCode(code Code) {
	os.Exit(int(code))
}
