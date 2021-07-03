package request

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/browserpass/browserpass-native/errors"
	"github.com/browserpass/browserpass-native/helpers"
	"github.com/browserpass/browserpass-native/response"
	log "github.com/sirupsen/logrus"
)

func deleteFile(request *request) {
	responseData := response.MakeDeleteResponse()

	if !strings.HasSuffix(request.File, ".gpg") {
		log.Errorf("The requested password file '%v' does not have the expected '.gpg' extension", request.File)
		response.SendErrorAndExit(
			errors.CodeInvalidPasswordFileExtension,
			&map[errors.Field]string{
				errors.FieldMessage: "The requested password file does not have the expected '.gpg' extension",
				errors.FieldAction:  "delete",
				errors.FieldFile:    request.File,
			},
		)
	}

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
				errors.FieldAction:  "delete",
				errors.FieldStoreID: request.StoreID,
			},
		)
	}

	normalizedStorePath, err := normalizePasswordStorePath(store.Path)
	if err != nil {
		log.Errorf(
			"The password store '%+v' is not accessible at its location: %+v",
			store, err,
		)
		response.SendErrorAndExit(
			errors.CodeInaccessiblePasswordStore,
			&map[errors.Field]string{
				errors.FieldMessage:   "The password store is not accessible",
				errors.FieldAction:    "delete",
				errors.FieldError:     err.Error(),
				errors.FieldStoreID:   store.ID,
				errors.FieldStoreName: store.Name,
				errors.FieldStorePath: store.Path,
			},
		)
	}
	store.Path = normalizedStorePath

	filePath := filepath.Join(store.Path, request.File)

	err = os.Remove(filePath)
	if err != nil {
		log.Error("Unable to delete the password file: ", err)
		response.SendErrorAndExit(
			errors.CodeUnableToDeletePasswordFile,
			&map[errors.Field]string{
				errors.FieldMessage:   "Unable to delete the password file",
				errors.FieldAction:    "delete",
				errors.FieldError:     err.Error(),
				errors.FieldFile:      request.File,
				errors.FieldStoreID:   store.ID,
				errors.FieldStoreName: store.Name,
				errors.FieldStorePath: store.Path,
			},
		)
	}

	parentDir := filepath.Dir(filePath)
	for {
		if parentDir == store.Path {
			break
		}

		isEmpty, err := helpers.IsDirectoryEmpty(parentDir)
		if err != nil {
			log.Error("Unable to determine if directory is empty and can be deleted: ", err)
			response.SendErrorAndExit(
				errors.CodeUnableToDetermineIsDirectoryEmpty,
				&map[errors.Field]string{
					errors.FieldMessage:   "Unable to determine if directory is empty and can be deleted",
					errors.FieldAction:    "delete",
					errors.FieldError:     err.Error(),
					errors.FieldDirectory: parentDir,
					errors.FieldStoreID:   store.ID,
					errors.FieldStoreName: store.Name,
					errors.FieldStorePath: store.Path,
				},
			)
		}

		if !isEmpty {
			break
		}

		err = os.Remove(parentDir)
		if err != nil {
			log.Error("Unable to delete the empty directory: ", err)
			response.SendErrorAndExit(
				errors.CodeUnableToDeleteEmptyDirectory,
				&map[errors.Field]string{
					errors.FieldMessage:   "Unable to delete the empty directory",
					errors.FieldAction:    "delete",
					errors.FieldError:     err.Error(),
					errors.FieldDirectory: parentDir,
					errors.FieldStoreID:   store.ID,
					errors.FieldStoreName: store.Name,
					errors.FieldStorePath: store.Path,
				},
			)
		}

		parentDir = filepath.Dir(parentDir)
	}

	response.SendOk(responseData)
}
