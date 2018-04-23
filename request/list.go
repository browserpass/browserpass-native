package request

import (
	"path/filepath"
	"sort"

	"github.com/browserpass/browserpass-native/errors"
	"github.com/browserpass/browserpass-native/response"
	"github.com/mattn/go-zglob"
	log "github.com/sirupsen/logrus"
)

func listFiles(request *request) {
	responseData := response.MakeListResponse()

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
					errors.FieldAction:    "list",
					errors.FieldError:     err.Error(),
					errors.FieldStoreID:   store.ID,
					errors.FieldStoreName: store.Name,
					errors.FieldStorePath: store.Path,
				},
			)
		}

		store.Path = normalizedStorePath

		files, err := zglob.GlobFollowSymlinks(filepath.Join(store.Path, "/**/*.gpg"))
		if err != nil {
			log.Errorf(
				"Unable to list the files in the password store '%+v' at its location: %+v",
				store, err,
			)
			response.SendErrorAndExit(
				errors.CodeUnableToListFilesInPasswordStore,
				&map[errors.Field]string{
					errors.FieldMessage:   "Unable to list the files in the password store",
					errors.FieldAction:    "list",
					errors.FieldError:     err.Error(),
					errors.FieldStoreID:   store.ID,
					errors.FieldStoreName: store.Name,
					errors.FieldStorePath: store.Path,
				},
			)
		}

		for i, file := range files {
			relativePath, err := filepath.Rel(store.Path, file)
			if err != nil {
				log.Errorf(
					"Unable to determine the relative path for a file '%v' in the password store '%+v': %+v",
					file, store, err,
				)
				response.SendErrorAndExit(
					errors.CodeUnableToDetermineRelativeFilePathInPasswordStore,
					&map[errors.Field]string{
						errors.FieldMessage:   "Unable to determine the relative path for a file in the password store",
						errors.FieldAction:    "list",
						errors.FieldError:     err.Error(),
						errors.FieldFile:      file,
						errors.FieldStoreID:   store.ID,
						errors.FieldStoreName: store.Name,
						errors.FieldStorePath: store.Path,
					},
				)
			}
			files[i] = relativePath
		}

		sort.Strings(files)
		responseData.Files[store.ID] = files
	}

	response.SendOk(responseData)
}
