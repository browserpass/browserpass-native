package request

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	bpErrors "github.com/browserpass/browserpass-native/errors"
	"github.com/browserpass/browserpass-native/response"
	log "github.com/sirupsen/logrus"
)

func normalizePasswordStorePath(storePath string) (string, error) {
	if storePath == "" {
		return "", errors.New("The store path cannot be empty")
	}

	if strings.HasPrefix(storePath, "~/") {
		storePath = filepath.Join("$HOME", storePath[2:])
	}
	storePath = os.ExpandEnv(storePath)

	directStorePath, err := filepath.EvalSymlinks(storePath)
	if err != nil {
		return "", err
	}
	storePath = directStorePath

	stat, err := os.Stat(storePath)
	if err != nil {
		return "", err
	}
	if !stat.IsDir() {
		return "", errors.New("The specified path exists, but is not a directory")
	}
	return storePath, nil
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
	return "", fmt.Errorf("Unable to detect the location of the gpg binary to use")
}

func getGpgPath(request *request) (string, error) {
	var gpgPath string
	var err error
	if request.Settings.GpgPath != "" {
		gpgPath = request.Settings.GpgPath
		err = validateGpgBinary(gpgPath)
		if err != nil {
			log.Errorf(
				"The provided gpg binary path '%v' is invalid: %+v",
				gpgPath, err,
			)
			response.SendErrorAndExit(
				bpErrors.CodeInvalidGpgPath,
				&map[bpErrors.Field]string{
					bpErrors.FieldMessage: "The provided gpg binary path is invalid",
					bpErrors.FieldAction:  request.Action,
					bpErrors.FieldError:   err.Error(),
					bpErrors.FieldGpgPath: gpgPath,
				},
			)
			return "", err
		}
	} else {
		gpgPath, err = detectGpgBinary()
		if err != nil {
			log.Error("Unable to detect the location of the gpg binary: ", err)
			response.SendErrorAndExit(
				bpErrors.CodeUnableToDetectGpgPath,
				&map[bpErrors.Field]string{
					bpErrors.FieldMessage: "Unable to detect the location of the gpg binary",
					bpErrors.FieldAction:  request.Action,
					bpErrors.FieldError:   err.Error(),
				},
			)
			return "", err
		}
	}

	return gpgPath, nil
}

func readGPGIDs(storePath string) []string {
	IDs := make([]string, 0)
	IDFilePath := filepath.Join(storePath, ".gpg-id")
	IDFile, err := os.Open(IDFilePath)
	if err != nil {
		return IDs
	}
	defer IDFile.Close()

	scanner := bufio.NewScanner(IDFile)
	for scanner.Scan() {
		IDs = append(IDs, scanner.Text())
	}

	return IDs
}
