package errors

import "os"

type Code int

const (
	CodeReadRequestLength Code = iota + 1
	CodeReadRequest
)

func ExitWithCode(code Code) {
	os.Exit(int(code))
}
