package version

import "fmt"

const major = 3
const minor = 1
const patch = 2

// Code version as integer
const Code = major*1000000 + minor*1000 + patch

// String version as string
func String() string {
	return fmt.Sprintf("%d.%d.%d", major, minor, patch)
}
