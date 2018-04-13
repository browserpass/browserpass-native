package main

import (
	"fmt"
)

const major = 3
const minor = 0
const patch = 0

const versionCode = major*1000000 + minor*1000 + patch

func versionString() string {
	return fmt.Sprintf("%d.%d.%d", major, minor, patch)
}
