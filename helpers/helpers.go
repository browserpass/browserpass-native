package helpers

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	passwordFile.Close()
	passDir := fmt.Sprintf("%s/", filepath.Dir(passwordFile.Name()))

	tmpFile, err := ioutil.TempFile("",
		fmt.Sprintf(
			"bp_%s_", strings.TrimPrefix(passwordFile.Name(), passDir),
		),
	)

	if err != nil {
		return "", fmt.Errorf("error: %s", err.Error())
	}

	// Ensure tmpFile with plain text secret is removed
	defer os.Remove(tmpFile.Name())

	var stderr bytes.Buffer
	gpgOptions := []string{
		"--decrypt", "--yes", "--quiet", "--batch", "--output", tmpFile.Name(), passwordFile.Name(),
	}

	cmd := exec.Command(gpgPath, gpgOptions...)
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("error: %s, Stderr: %s", err.Error(), stderr.String())
	}

	scanner := bufio.NewScanner(tmpFile)
	out := []string{}
	for scanner.Scan() {
		out = append(out, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error: %s", err.Error())
	}

	return strings.Join(out, "\n"), nil
}

func GpgEncryptFile(filePath string, contents string, recipients []string, gpgPath string) error {
	err := os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		return fmt.Errorf("Unable to create directory structure: %s", err.Error())
	}

	var stdout, stderr bytes.Buffer
	gpgOptions := []string{"--encrypt", "--yes", "--quiet", "--batch", "--output", filePath}
	for _, recipient := range recipients {
		gpgOptions = append(gpgOptions, "--recipient", recipient)
	}

	cmd := exec.Command(gpgPath, gpgOptions...)
	cmd.Stdin = strings.NewReader(contents)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err = cmd.Run(); err != nil {
		return fmt.Errorf("error %s, stderr: %s", err.Error(), stderr.String())
	}

	return nil
}

func DetectGpgRecipients(filePath string) ([]string, error) {
	dir := filepath.Dir(filePath)
	for {
		file, err := ioutil.ReadFile(filepath.Join(dir, ".gpg-id"))
		if err == nil {
			return strings.Split(strings.TrimSpace(string(file)), "\n"), nil
		}

		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("Unable to open `.gpg-id` file: %s", err.Error())
		}

		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			return nil, fmt.Errorf("Unable to find '.gpg-id' file")
		}

		dir = parentDir
	}
}

func IsDirectoryEmpty(dirPath string) (bool, error) {
	f, err := os.Open(dirPath)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}

	return false, err
}
