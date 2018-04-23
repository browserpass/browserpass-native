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
	CodeInvalidRequestAction Code = 12

	// CodeInaccessiblePasswordStore error accessing a user-configured password store
	CodeInaccessiblePasswordStore Code = 13

	// CodeInaccessibleDefaultPasswordStore error accessing the default password store
	CodeInaccessibleDefaultPasswordStore Code = 14

	// CodeUnknownDefaultPasswordStoreLocation error determining the location of the default password store
	CodeUnknownDefaultPasswordStoreLocation Code = 15

	// CodeUnreadablePasswordStoreDefaultSettings error reading the default settings of a user-configured password store
	CodeUnreadablePasswordStoreDefaultSettings Code = 16

	// CodeUnreadableDefaultPasswordStoreDefaultSettings error reading the default settings of the default password store
	CodeUnreadableDefaultPasswordStoreDefaultSettings Code = 17

	// CodeUnableToListFilesInPasswordStore error listing files in a password store
	CodeUnableToListFilesInPasswordStore Code = 18

	// CodeUnableToDetermineRelativeFilePathInPasswordStore error determining a relative path for a file in a password store
	CodeUnableToDetermineRelativeFilePathInPasswordStore Code = 19

	// CodeInvalidPasswordStore error looking for a password store with the given ID
	CodeInvalidPasswordStore Code = 20

	// CodeInvalidGpgPath error looking for a gpg binary at the given path
	CodeInvalidGpgPath Code = 21

	// CodeUnableToDetectGpgPath error detecting the location of the gpg binary
	CodeUnableToDetectGpgPath Code = 22

	// CodeInvalidPasswordFileExtension error unexpected password file extension
	CodeInvalidPasswordFileExtension Code = 23

	// CodeUnableToDecryptPasswordFile error decrypting a password file
	CodeUnableToDecryptPasswordFile Code = 24

	// CodeInaccessiblePasswordStoresContainer error accessing a user-configured password stores container
	CodeInaccessiblePasswordStoresContainer Code = 25

	// CodeUnableToListFilesInPasswordStoresContainer error listing files in a password stores container
	CodeUnableToListFilesInPasswordStoresContainer Code = 26

	// CodeUnableToDetermineRelativeFilePathInPasswordStoresContainer error determining a relative path for a file in a password stores container
	CodeUnableToDetermineRelativeFilePathInPasswordStoresContainer Code = 27
)

// Field extra field in the error response params
type Field string

const (
	// FieldMessage a user-friendly error message, always present
	FieldMessage Field = "message"

	// FieldAction a browser request action that resulted in a failure
	FieldAction Field = "action"

	// FieldError an error message returned from an external system
	FieldError Field = "error"

	// FieldStoreID a password store id
	FieldStoreID Field = "storeId"

	// FieldStoreName a password store name
	FieldStoreName Field = "storeName"

	// FieldStorePath a password store path
	FieldStorePath Field = "storePath"

	// FieldFile a password file
	FieldFile Field = "file"

	// FieldGpgPath a path to the gpg binary
	FieldGpgPath Field = "gpgPath"
)

// ExitWithCode exit with error code
func ExitWithCode(code Code) {
	os.Exit(int(code))
}
