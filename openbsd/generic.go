// +build !openbsd

package openbsd

// Pledge allowed system calls, available only on OpenBSD systems
func Pledge(promises string) {
}
