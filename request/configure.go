package request

import (
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"github.com/browserpass/browserpass-native/errors"
	"github.com/browserpass/browserpass-native/response"
	log "github.com/sirupsen/logrus"
)

func configure(request request) {
	responseData := response.MakeConfigureResponse()

	// Check that each and every store in the settings exists and is accessible.
	// Then read the default configuration for these stores (if available).
	for _, store := range request.Settings.Stores {
		normalizedStorePath, err := normalizePasswordStorePath(store.Path)
		if err != nil {
			log.Errorf(
				"Inaccessible path '%v' of the user-configured password store '%v': %+v",
				store.Path, store.Name, err)
			response.SendError(
				errors.CodeInaccessiblePasswordStore,
				"The password store is not accessible",
				&map[string]string{"error": err.Error(), "name": store.Name, "path": store.Path})
			errors.ExitWithCode(errors.CodeInaccessiblePasswordStore)
		}

		store.Path = normalizedStorePath

		responseData.StoreSettings[store.Name], err = readDefaultSettings(store.Path)
		if err != nil {
			log.Errorf(
				"Unable to read the default settings of the user-configured password store '%v' in '%v': %+v",
				store.Name, store.Path, err)
			response.SendError(
				errors.CodeUnreadablePasswordStoreDefaultSettings,
				"Unable to read the default settings of the password store",
				&map[string]string{"error": err.Error(), "name": store.Name, "path": store.Path})
			errors.ExitWithCode(errors.CodeUnreadablePasswordStoreDefaultSettings)
		}
	}

	// Check whether a store in the default location exists and is accessible.
	// If there is at least one store in the settings, it is expected that there might be no store in the default location => do not return errors.
	// However, if there are no stores in the settings, user expects to use the default password store => return an error if it is not accessible.
	possibleDefaultStorePath, err := getDefaultPasswordStorePath()
	if err != nil {
		if len(request.Settings.Stores) == 0 {
			log.Error("Unable to determine the location of the default password store: ", err)
			response.SendError(
				errors.CodeUnknownDefaultPasswordStoreLocation,
				"Unable to determine the location of the default password store",
				&map[string]string{"error": err.Error()})
			errors.ExitWithCode(errors.CodeUnknownDefaultPasswordStoreLocation)
		}
	} else {
		responseData.DefaultStore.Path, err = normalizePasswordStorePath(possibleDefaultStorePath)
		if err != nil {
			if len(request.Settings.Stores) == 0 {
				log.Errorf(
					"The path '%v' of the default password store is not accessible: %+v",
					possibleDefaultStorePath, err)
				response.SendError(
					errors.CodeInaccessibleDefaultPasswordStore,
					"The default password store is not accessible",
					&map[string]string{"error": err.Error(), "path": possibleDefaultStorePath})
				errors.ExitWithCode(errors.CodeInaccessibleDefaultPasswordStore)
			}
		}
	}

	if responseData.DefaultStore.Path != "" {
		responseData.DefaultStore.Settings, err = readDefaultSettings(responseData.DefaultStore.Path)
		if err != nil {
			log.Errorf(
				"Unable to read the default settings of the default password store in '%v': %+v",
				responseData.DefaultStore.Path, err)
			response.SendError(
				errors.CodeUnreadableDefaultPasswordStoreDefaultSettings,
				"Unable to read the default settings of the default password store",
				&map[string]string{"error": err.Error(), "path": responseData.DefaultStore.Path})
			errors.ExitWithCode(errors.CodeUnreadableDefaultPasswordStoreDefaultSettings)
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
		return "", nil
	}
	return "", err
}
