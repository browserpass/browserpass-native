package request

import (
	"bytes"
	goerrors "errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/browserpass/browserpass-native/errors"
	"github.com/browserpass/browserpass-native/response"
	log "github.com/sirupsen/logrus"
)

func fetchDecryptedContents(request request) {
	responseData := response.MakeFetchResponse()

	if !strings.HasSuffix(request.File, ".gpg") {
		log.Errorf("The requested password file '%v' does not have the expected '.gpg' extension", request.File)
		response.SendErrorAndExit(
			errors.CodeInvalidPasswordFile,
			"The requested password file does not have the expected '.gpg' extension",
			&map[string]string{"action": "fetch", "file": request.File},
		)
	}

	store, ok := request.Settings.Stores[request.Store]
	if !ok {
		log.Errorf(
			"The password store '%v' is not present in the list of stores '%v'",
			request.Store, request.Settings.Stores,
		)
		response.SendErrorAndExit(
			errors.CodeInvalidPasswordStore,
			"The password store is not present in the list of stores",
			&map[string]string{"action": "fetch", "name": request.Store},
		)
	}

	normalizedStorePath, err := normalizePasswordStorePath(store.Path)
	if err != nil {
		log.Errorf(
			"The password store '%v' is not accessible at the location '%v': %+v",
			store.Name, store.Path, err,
		)
		response.SendErrorAndExit(
			errors.CodeInaccessiblePasswordStore,
			"The password store is not accessible",
			&map[string]string{"action": "fetch", "error": err.Error(), "name": store.Name, "path": store.Path},
		)
	}
	store.Path = normalizedStorePath

	gpgPath := request.Settings.GpgPath
	if gpgPath != "" {
		err = validateGpgBinary(gpgPath)
		if err != nil {
			log.Errorf(
				"The provided gpg binary path '%v' is invalid: %+v",
				gpgPath, err,
			)
			response.SendErrorAndExit(
				errors.CodeInvalidGpgPath,
				"The provided gpg binary path is invalid",
				&map[string]string{"action": "fetch", "error": err.Error(), "gpgPath": gpgPath},
			)
		}
	} else {
		gpgPath, err = detectGpgBinary()
		if err != nil {
			log.Error("Unable to detect the location of the gpg binary: ", err)
			response.SendErrorAndExit(
				errors.CodeUnableToDetectGpgPath,
				"Unable to detect the location of the gpg binary",
				&map[string]string{"action": "fetch", "error": err.Error()},
			)
		}
	}

	responseData.Contents, err = decryptFile(store, request.File, gpgPath)
	if err != nil {
		log.Errorf(
			"Unable to decrypt the password file '%v' in the password store '%v' located in '%v': %+v",
			request.File, store.Name, store.Path, err,
		)
		response.SendErrorAndExit(
			errors.CodeUnableToDecryptPasswordFile,
			"Unable to decrypt the password file",
			&map[string]string{"action": "fetch", "error": err.Error(), "file": request.File, "name": store.Name, "path": store.Path},
		)
	}

	response.SendOk(responseData)
}

func detectGpgBinary() (string, error) {
	// Look in $PATH first, then check common locations - the first successful result wins
	gpgBinaryPriorityList := []string{
		"gpg2", "gpg",
		"/bin/gpg2", "/usr/bin/gpg2", "/usr/local/bin/gpg2",
		"/bin/gpg", "/usr/bin/gpg", "/usr/local/bin/gpg",
	}

	for _, binary := range gpgBinaryPriorityList {
		err := validateGpgBinary(binary)
		if err == nil {
			return binary, nil
		}
	}

	return "", goerrors.New("Unable to detect the location of the gpg binary to use")
}

func validateGpgBinary(gpgPath string) error {
	return exec.Command(gpgPath, "--version").Run()
}

func decryptFile(store store, file string, gpgPath string) (string, error) {
	passwordFilePath := filepath.Join(store.Path, file)
	passwordFile, err := os.Open(passwordFilePath)
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
		return "", goerrors.New(fmt.Sprintf("Error: %s, Stderr: %s", err.Error(), stderr.String()))
	}

	return stdout.String(), nil
}
