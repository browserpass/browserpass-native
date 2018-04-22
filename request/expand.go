package request

import (
	"path/filepath"
	"strings"

	"github.com/browserpass/browserpass-native/errors"
	"github.com/browserpass/browserpass-native/response"
	zglob "github.com/mattn/go-zglob"
	log "github.com/sirupsen/logrus"
)

func expand(request *request) {
	responseData := response.MakeExpandResponse()

	for _, container := range request.Settings.Stores {
		normalizedStorePath, err := normalizePasswordStorePath(container.Path)
		if err != nil {
			log.Errorf(
				"The password stores container '%+v' is not accessible at its location: %+v",
				container, err,
			)
			response.SendErrorAndExit(
				errors.CodeInaccessiblePasswordStoresContainer,
				&map[errors.Field]string{
					errors.FieldMessage:   "The password stores container is not accessible",
					errors.FieldAction:    "expand",
					errors.FieldError:     err.Error(),
					errors.FieldStoreID:   container.ID,
					errors.FieldStoreName: container.Name,
					errors.FieldStorePath: container.Path,
				},
			)
		}

		container.Path = normalizedStorePath

		files, err := zglob.GlobFollowSymlinks(filepath.Join(container.Path, "/**/*.gpg"))
		if err != nil {
			log.Errorf(
				"Unable to list the files in the password stores container '%+v' at its location: %+v",
				container, err,
			)
			response.SendErrorAndExit(
				errors.CodeUnableToListFilesInPasswordStoresContainer,
				&map[errors.Field]string{
					errors.FieldMessage:   "Unable to list the files in the password stores container",
					errors.FieldAction:    "expand",
					errors.FieldError:     err.Error(),
					errors.FieldStoreID:   container.ID,
					errors.FieldStoreName: container.Name,
					errors.FieldStorePath: container.Path,
				},
			)
		}

		stores := make(map[string]string)
		for _, file := range files {
			relativePath, err := filepath.Rel(container.Path, file)
			if err != nil {
				log.Errorf(
					"Unable to determine the relative path for a file '%v' in the password stores container '%+v': %+v",
					file, container, err,
				)
				response.SendErrorAndExit(
					errors.CodeUnableToDetermineRelativeFilePathInPasswordStoresContainer,
					&map[errors.Field]string{
						errors.FieldMessage:   "Unable to determine the relative path for a file in the password stores container",
						errors.FieldAction:    "expand",
						errors.FieldError:     err.Error(),
						errors.FieldFile:      file,
						errors.FieldStoreID:   container.ID,
						errors.FieldStoreName: container.Name,
						errors.FieldStorePath: container.Path,
					},
				)
			}
			parts := strings.SplitN(relativePath, "/", 2)
			if len(parts) > 1 {
				storeName := parts[0]
				if stores[storeName] == "" {
					stores[storeName] = filepath.Join(container.Path, storeName)
				}
			}
		}

		storesList := make([]*response.ExpandedStore, len(stores))
		i := 0
		for storeName, storePath := range stores {
			storesList[i] = response.MakeExpandedStore(storeName, storePath)
			i++
		}
		responseData.Stores[container.ID] = storesList
	}

	response.SendOk(responseData)
}
