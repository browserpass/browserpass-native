package errors

import (
	"os"
)

// Code exit code
type Code int

// Error codes that are sent to the browser extension and used as exit codes in the app.
// DO NOT MODIFY THE VALUES, always append new error codes to the bottom.
const (
	CodeParseRequestLength                                    Code = 10
	CodeParseRequest                                          Code = 11
	CodeInvalidRequestAction                                  Code = 12
	CodeInaccessiblePasswordStore                             Code = 13
	CodeInaccessibleDefaultPasswordStore                      Code = 14
	CodeUnknownDefaultPasswordStoreLocation                   Code = 15
	CodeUnreadablePasswordStoreDefaultSettings                Code = 16
	CodeUnreadableDefaultPasswordStoreDefaultSettings         Code = 17
	CodeUnableToListFilesInPasswordStore                      Code = 18
	CodeUnableToDetermineRelativeFilePathInPasswordStore      Code = 19
	CodeInvalidPasswordStore                                  Code = 20
	CodeInvalidGpgPath                                        Code = 21
	CodeUnableToDetectGpgPath                                 Code = 22
	CodeInvalidPasswordFileExtension                          Code = 23
	CodeUnableToDecryptPasswordFile                           Code = 24
	CodeUnableToListDirectoriesInPasswordStore                Code = 25
	CodeUnableToDetermineRelativeDirectoryPathInPasswordStore Code = 26
	CodeEmptyContents                                         Code = 27
	CodeUnableToDetermineGpgRecipients                        Code = 28
	CodeUnableToEncryptPasswordFile                           Code = 29
	CodeUnableToDeletePasswordFile                            Code = 30
	CodeUnableToDetermineIsDirectoryEmpty                     Code = 31
	CodeUnableToDeleteEmptyDirectory                          Code = 32
)

// Field extra field in the error response params
type Field string

// Extra fields that can be sent to the browser extension as part of an error response.
// FieldMessage is always present, others are optional.
const (
	FieldMessage   Field = "message"
	FieldAction    Field = "action"
	FieldError     Field = "error"
	FieldStoreID   Field = "storeId"
	FieldStoreName Field = "storeName"
	FieldStorePath Field = "storePath"
	FieldFile      Field = "file"
	FieldDirectory Field = "directory"
	FieldGpgPath   Field = "gpgPath"
)

// ExitWithCode exit with error code
func ExitWithCode(code Code) {
	os.Exit(int(code))
}
