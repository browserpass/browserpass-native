package request

import (
	"path/filepath"
	"strings"

	"github.com/browserpass/browserpass-native/errors"
	"github.com/browserpass/browserpass-native/helpers"
	"github.com/browserpass/browserpass-native/response"
	log "github.com/sirupsen/logrus"
)

func saveEncryptedContents(request *request) {
	responseData := response.MakeSaveResponse()

	if !strings.HasSuffix(request.File, ".gpg") {
		log.Errorf("The requested password file '%v' does not have the expected '.gpg' extension", request.File)
		response.SendErrorAndExit(
			errors.CodeInvalidPasswordFileExtension,
			&map[errors.Field]string{
				errors.FieldMessage: "The requested password file does not have the expected '.gpg' extension",
				errors.FieldAction:  "save",
				errors.FieldFile:    request.File,
			},
		)
	}

	if request.Contents == "" {
		log.Errorf("The entry contents is missing")
		response.SendErrorAndExit(
			errors.CodeEmptyContents,
			&map[errors.Field]string{
				errors.FieldMessage: "The entry contents is missing",
				errors.FieldAction:  "save",
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
				errors.FieldAction:  "save",
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
				errors.FieldAction:    "save",
				errors.FieldError:     err.Error(),
				errors.FieldStoreID:   store.ID,
				errors.FieldStoreName: store.Name,
				errors.FieldStorePath: store.Path,
			},
		)
	}
	store.Path = normalizedStorePath

	var gpgPath string
	if request.Settings.GpgPath != "" || store.Settings.GpgPath != "" {
		if request.Settings.GpgPath != "" {
			gpgPath = request.Settings.GpgPath
		} else {
			gpgPath = store.Settings.GpgPath
		}
		err = helpers.ValidateGpgBinary(gpgPath)
		if err != nil {
			log.Errorf(
				"The provided gpg binary path '%v' is invalid: %+v",
				gpgPath, err,
			)
			response.SendErrorAndExit(
				errors.CodeInvalidGpgPath,
				&map[errors.Field]string{
					errors.FieldMessage: "The provided gpg binary path is invalid",
					errors.FieldAction:  "save",
					errors.FieldError:   err.Error(),
					errors.FieldGpgPath: gpgPath,
				},
			)
		}
	} else {
		gpgPath, err = helpers.DetectGpgBinary()
		if err != nil {
			log.Error("Unable to detect the location of the gpg binary: ", err)
			response.SendErrorAndExit(
				errors.CodeUnableToDetectGpgPath,
				&map[errors.Field]string{
					errors.FieldMessage: "Unable to detect the location of the gpg binary",
					errors.FieldAction:  "save",
					errors.FieldError:   err.Error(),
				},
			)
		}
	}

	filePath := filepath.Join(store.Path, request.File)

	recipients, err := helpers.DetectGpgRecipients(filePath)
	if err != nil {
		log.Error("Unable to determine recipients for the gpg encryption: ", err)
		response.SendErrorAndExit(
			errors.CodeUnableToDetermineGpgRecipients,
			&map[errors.Field]string{
				errors.FieldMessage:   "Unable to determine recipients for the gpg encryption",
				errors.FieldAction:    "save",
				errors.FieldError:     err.Error(),
				errors.FieldFile:      request.File,
				errors.FieldStoreID:   store.ID,
				errors.FieldStoreName: store.Name,
				errors.FieldStorePath: store.Path,
			},
		)
	}

	err = helpers.GpgEncryptFile(filePath, request.Contents, recipients, gpgPath)
	if err != nil {
		log.Errorf(
			"Unable to encrypt the password file '%v' in the password store '%+v': %+v",
			request.File, store, err,
		)
		response.SendErrorAndExit(
			errors.CodeUnableToEncryptPasswordFile,
			&map[errors.Field]string{
				errors.FieldMessage:   "Unable to encrypt the password file",
				errors.FieldAction:    "save",
				errors.FieldError:     err.Error(),
				errors.FieldFile:      request.File,
				errors.FieldStoreID:   store.ID,
				errors.FieldStoreName: store.Name,
				errors.FieldStorePath: store.Path,
			},
		)
	}

	response.SendOk(responseData)
}
