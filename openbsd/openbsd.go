// +build openbsd

package openbsd

import "golang.org/x/sys/unix"

// Pledge allowed system calls
func Pledge(promises string) {
	unix.Pledge(promises, nil)
}
