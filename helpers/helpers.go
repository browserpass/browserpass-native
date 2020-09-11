package helpers

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func DetectGpgBinary() (string, error) {
	// Look in $PATH first, then check common locations - the first successful result wins
	gpgBinaryPriorityList := []string{
		"gpg2", "gpg",
		"/bin/gpg2", "/usr/bin/gpg2", "/usr/local/bin/gpg2",
		"/bin/gpg", "/usr/bin/gpg", "/usr/local/bin/gpg",
	}

	for _, binary := range gpgBinaryPriorityList {
		err := ValidateGpgBinary(binary)
		if err == nil {
			return binary, nil
		}
	}
	return "", fmt.Errorf("Unable to detect the location of the gpg binary to use")
}

func ValidateGpgBinary(gpgPath string) error {
	return exec.Command(gpgPath, "--version").Run()
}

func GpgDecryptFile(filePath string, gpgPath string) (string, error) {
	passwordFile, err := os.Open(filePath)
	if err != nil {
		return "", err
	}

	var stdout, stderr bytes.Buffer
	gpgOptions := []string{"--decrypt", "--yes", "--quiet", "--batch", "-"}

	cmd := exec.Command(gpgPath, gpgOptions...)
	cmd.Stdin = passwordFile
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("Error: %s, Stderr: %s", err.Error(), stderr.String())
	}

	return stdout.String(), nil
}
