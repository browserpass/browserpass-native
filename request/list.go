package request

import (
	"path/filepath"
	"sort"

	"github.com/browserpass/browserpass-native/errors"
	"github.com/browserpass/browserpass-native/response"
	"github.com/mattn/go-zglob"
	log "github.com/sirupsen/logrus"
)

func listFiles(request request) {
	responseData := response.MakeListResponse()

	for _, store := range request.Settings.Stores {
		normalizedStorePath, err := normalizePasswordStorePath(store.Path)
		if err != nil {
			log.Errorf(
				"The password store '%v' is not accessible at the location '%v': %+v",
				store.Name, store.Path, err,
			)
			response.SendErrorAndExit(
				errors.CodeInaccessiblePasswordStore,
				"The password store is not accessible",
				&map[string]string{"action": "list", "error": err.Error(), "name": store.Name, "path": store.Path},
			)
		}

		store.Path = normalizedStorePath

		files, err := zglob.GlobFollowSymlinks(filepath.Join(store.Path, "/**/*.gpg"))
		if err != nil {
			log.Errorf(
				"Unable to list the files in the password store '%v' at the location '%v': %+v",
				store.Name, store.Path, err,
			)
			response.SendErrorAndExit(
				errors.CodeUnableToListFilesInPasswordStore,
				"Unable to list the files in the password store",
				&map[string]string{"action": "list", "error": err.Error(), "name": store.Name, "path": store.Path},
			)
		}

		for i, file := range files {
			relativePath, err := filepath.Rel(store.Path, file)
			if err != nil {
				log.Errorf(
					"Unable to determine the relative path for a file '%v' in the password store '%v' at the location '%v': %+v",
					file, store.Name, store.Path, err,
				)
				response.SendErrorAndExit(
					errors.CodeUnableToDetermineRelativeFilePathInPasswordStore,
					"Unable to determine the relative path for a file in the password store",
					&map[string]string{"action": "list", "error": err.Error(), "file": file, "name": store.Name, "path": store.Path},
				)
			}
			files[i] = relativePath
		}

		sort.Strings(files)
		responseData.Files[store.Name] = files
	}

	response.SendOk(responseData)
}
