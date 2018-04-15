package errors

import (
	"os"
)

// Code exit code
type Code int

// Error codes that are sent to the browser extension and used as exit codes in the app.
// DO NOT MODIFY THE VALUES, always append new error codes to the bottom.
const (
	// CodeParseRequestLength error parsing a request length
	CodeParseRequestLength Code = 10

	// CodeParseRequest error parsing a request
	CodeParseRequest Code = 11

	// CodeInvalidRequestAction error parsing a request action
	CodeInvalidRequestAction = 12

	// CodeInaccessiblePasswordStore error accessing a user-configured password store
	CodeInaccessiblePasswordStore = 13

	// CodeInaccessibleDefaultPasswordStore error accessing the default password store
	CodeInaccessibleDefaultPasswordStore = 14

	// CodeUnknownDefaultPasswordStoreLocation error determining the location of the default password store
	CodeUnknownDefaultPasswordStoreLocation = 15

	// CodeUnreadablePasswordStoreDefaultSettings error reading the default settings of a user-configured password store
	CodeUnreadablePasswordStoreDefaultSettings = 16

	// CodeUnreadableDefaultPasswordStoreDefaultSettings error reading the default settings of the default password store
	CodeUnreadableDefaultPasswordStoreDefaultSettings = 17
)

// ExitWithCode exit with error code
func ExitWithCode(code Code) {
	os.Exit(int(code))
}
