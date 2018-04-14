package errors

import (
	"os"
)

// Code exit code
type Code int

// Error codes that are sent to the browser extension and used as exit codes in the app.
// DO NOT MODIFY THE VALUES, always append new error codes to the bottom.
const (
	// CodeParseRequestLength error parsing request length
	CodeParseRequestLength Code = 10

	// CodeParseRequest error parsing request
	CodeParseRequest Code = 11

	// CodeInvalidRequestAction error parsing request action
	CodeInvalidRequestAction = 12
)

// ExitWithCode exit with error code
func ExitWithCode(code Code) {
	os.Exit(int(code))
}
