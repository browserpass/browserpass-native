package request

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/browserpass/browserpass-native/errors"
	"github.com/browserpass/browserpass-native/response"
	log "github.com/sirupsen/logrus"
)

func createFile(request *request) {
	store, ok := request.Settings.Stores[request.StoreID]
	if !ok {
		log.Errorf(
			"The password store with ID '%v' is not present in the list of stores '%+v'",
			request.StoreID, request.Settings.Stores,
		)
		response.SendErrorAndExit(
			errors.CodeInvalidPasswordStore,
			&map[errors.Field]string{
				errors.FieldMessage: "The password store is not present in the list of stores",
				errors.FieldAction:  "create",
				errors.FieldStoreID: request.StoreID,
			},
		)
	}
	storePath, err := normalizePasswordStorePath(store.Path)
	if err != nil {
		return // TODO
	}

	gpgPath, err := getGpgPath(request)
	if err != nil {
		return // TODO
	}

	credentials := request.Credentials
	fileString := fmt.Sprintf("%s\nlogin: %s\n", credentials.Password, credentials.Login)

	if credentials.Email != "" {
		fileString = fileString + fmt.Sprintf("email: %s\n", credentials.Email)
	}

	err = encryptContent(storePath, request.File, fileString, gpgPath)
	if err != nil {
		response.SendErrorAndExit(
			errors.CodeUnableToEncryptPasswordFile,
			&map[errors.Field]string{
				errors.FieldMessage: "Could not encrypt new password file",
				errors.FieldAction:  "create",
				errors.FieldError:   err.Error(),
				errors.FieldStoreID: request.StoreID,
			},
		)
	}

	err = gitAddAndCommit(storePath, request.File)
	if err != nil {
		response.SendErrorAndExit(
			errors.CodeUnableToGitCommit,
			&map[errors.Field]string{
				errors.FieldMessage: "Could not commit file to git repository",
				errors.FieldAction:  "create",
				errors.FieldError:   err.Error(),
				errors.FieldStoreID: request.StoreID,
				errors.FieldFile:    request.File,
			},
		)
	}

	response.SendOk(nil)
}

func encryptContent(storePath, file, content, gpgPath string) error {
	IDs := readGPGIDs(storePath)
	passwordFilePath := filepath.Join(storePath, file)
	err := os.MkdirAll(filepath.Dir(passwordFilePath), os.ModePerm)
	if err != nil {
		return err
	}

	var stderr bytes.Buffer
	gpgOptions := []string{"--encrypt", "-o", passwordFilePath, "--quiet"}

	for _, id := range IDs {
		gpgOptions = append(gpgOptions, "-r")
		gpgOptions = append(gpgOptions, id)
	}

	contentReader := strings.NewReader(content)
	cmd := exec.Command(gpgPath, gpgOptions...)
	cmd.Stdin = contentReader
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Error: %s, Stderr: %s", err.Error(), stderr.String())
	}

	return nil
}

func gitAddAndCommit(storePath, file string) error {
	gitBaseOptions := []string{"-C", storePath}

	var stderr bytes.Buffer
	cmd := exec.Command("git", append(gitBaseOptions, []string{"add", file}...)...)
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Error: %s, Stderr: %s", err.Error(), stderr.String())
	}

	commitMessage := fmt.Sprintf("Add password %s from browserpass", file)
	cmd = exec.Command("git", append(gitBaseOptions, []string{"commit", "-m", commitMessage}...)...)
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Error: %s, Stderr: %s", err.Error(), stderr.String())
	}

	return nil
}
