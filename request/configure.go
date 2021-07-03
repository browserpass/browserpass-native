package request

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"github.com/browserpass/browserpass-native/errors"
	"github.com/browserpass/browserpass-native/helpers"
	"github.com/browserpass/browserpass-native/response"
	log "github.com/sirupsen/logrus"
)

func configure(request *request) {
	responseData := response.MakeConfigureResponse()

	// User configured gpgPath in the browser, check if it is a valid binary to use
	if request.Settings.GpgPath != "" {
		err := helpers.ValidateGpgBinary(request.Settings.GpgPath)
		if err != nil {
			log.Errorf(
				"The provided gpg binary path '%v' is invalid: %+v",
				request.Settings.GpgPath, err,
			)
			response.SendErrorAndExit(
				errors.CodeInvalidGpgPath,
				&map[errors.Field]string{
					errors.FieldMessage: "The provided gpg binary path is invalid",
					errors.FieldAction:  "configure",
					errors.FieldError:   err.Error(),
					errors.FieldGpgPath: request.Settings.GpgPath,
				},
			)
		}
	}

	// Check that each and every store in the settings exists and is accessible.
	// Then read the default configuration for these stores (if available).
	for _, store := range request.Settings.Stores {
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
					errors.FieldAction:    "configure",
					errors.FieldError:     err.Error(),
					errors.FieldStoreID:   store.ID,
					errors.FieldStoreName: store.Name,
					errors.FieldStorePath: store.Path,
				},
			)
		}

		store.Path = normalizedStorePath

		responseData.StoreSettings[store.ID], err = readDefaultSettings(store.Path)
		if err == nil {
			var storeSettings StoreSettings
			err = json.Unmarshal([]byte(responseData.StoreSettings[store.ID]), &storeSettings)
		}
		if err != nil {
			log.Errorf(
				"Unable to read .browserpass.json of the user-configured password store '%+v': %+v",
				store, err,
			)
			response.SendErrorAndExit(
				errors.CodeUnreadablePasswordStoreDefaultSettings,
				&map[errors.Field]string{
					errors.FieldMessage:   "Unable to read .browserpass.json of the password store",
					errors.FieldAction:    "configure",
					errors.FieldError:     err.Error(),
					errors.FieldStoreID:   store.ID,
					errors.FieldStoreName: store.Name,
					errors.FieldStorePath: store.Path,
				},
			)
		}
	}

	// Check whether a store in the default location exists and is accessible.
	// If there is at least one store in the settings, user will not use the default store => skip its validation.
	// However, if there are no stores in the settings, user expects to use the default password store => return an error if it is not accessible.
	if len(request.Settings.Stores) == 0 {
		possibleDefaultStorePath, err := getDefaultPasswordStorePath()
		if err != nil {
			log.Error("Unable to determine the location of the default password store: ", err)
			response.SendErrorAndExit(
				errors.CodeUnknownDefaultPasswordStoreLocation,
				&map[errors.Field]string{
					errors.FieldMessage: "Unable to determine the location of the default password store",
					errors.FieldAction:  "configure",
					errors.FieldError:   err.Error(),
				},
			)
		} else {
			responseData.DefaultStore.Path, err = normalizePasswordStorePath(possibleDefaultStorePath)
			if err != nil {
				log.Errorf(
					"The default password store is not accessible at the location '%v': %+v",
					possibleDefaultStorePath, err,
				)
				response.SendErrorAndExit(
					errors.CodeInaccessibleDefaultPasswordStore,
					&map[errors.Field]string{
						errors.FieldMessage:   "The default password store is not accessible",
						errors.FieldAction:    "configure",
						errors.FieldError:     err.Error(),
						errors.FieldStorePath: possibleDefaultStorePath,
					},
				)
			}
		}

		responseData.DefaultStore.Settings, err = readDefaultSettings(responseData.DefaultStore.Path)
		if err == nil {
			var storeSettings StoreSettings
			err = json.Unmarshal([]byte(responseData.DefaultStore.Settings), &storeSettings)
		}
		if err != nil {
			log.Errorf(
				"Unable to read .browserpass.json of the default password store in '%v': %+v",
				responseData.DefaultStore.Path, err,
			)
			response.SendErrorAndExit(
				errors.CodeUnreadableDefaultPasswordStoreDefaultSettings,
				&map[errors.Field]string{
					errors.FieldMessage:   "Unable to read .browserpass.json of the default password store",
					errors.FieldAction:    "configure",
					errors.FieldError:     err.Error(),
					errors.FieldStorePath: responseData.DefaultStore.Path,
				},
			)
		}
	}

	response.SendOk(responseData)
}

func getDefaultPasswordStorePath() (string, error) {
	path := os.Getenv("PASSWORD_STORE_DIR")
	if path != "" {
		return path, nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	path = filepath.Join(usr.HomeDir, ".password-store")
	return path, nil
}

func readDefaultSettings(storePath string) (string, error) {
	content, err := ioutil.ReadFile(filepath.Join(storePath, ".browserpass.json"))
	if err == nil {
		return string(content), nil
	}
	if os.IsNotExist(err) {
		return "{}", nil
	}
	return "", err
}
