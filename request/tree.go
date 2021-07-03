package request

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/browserpass/browserpass-native/errors"
	"github.com/browserpass/browserpass-native/response"
	"github.com/mattn/go-zglob/fastwalk"
	log "github.com/sirupsen/logrus"
)

func listDirectories(request *request) {
	responseData := response.MakeTreeResponse()

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
					errors.FieldAction:    "tree",
					errors.FieldError:     err.Error(),
					errors.FieldStoreID:   store.ID,
					errors.FieldStoreName: store.Name,
					errors.FieldStorePath: store.Path,
				},
			)
		}

		store.Path = normalizedStorePath

		var mu sync.Mutex
		directories := []string{}
		err = fastwalk.FastWalk(store.Path, func(path string, typ os.FileMode) error {
			if typ == os.ModeSymlink {
				followedPath, err := filepath.EvalSymlinks(path)
				if err == nil {
					fi, err := os.Lstat(followedPath)
					if err == nil && fi.IsDir() {
						return fastwalk.TraverseLink
					}
				}
			}

			if typ.IsDir() && path != store.Path {
				if filepath.Base(path) == ".git" {
					return filepath.SkipDir
				}
				mu.Lock()
				directories = append(directories, path)
				mu.Unlock()
			}

			return nil
		})

		if err != nil {
			log.Errorf(
				"Unable to list the directory tree in the password store '%+v' at its location: %+v",
				store, err,
			)
			response.SendErrorAndExit(
				errors.CodeUnableToListDirectoriesInPasswordStore,
				&map[errors.Field]string{
					errors.FieldMessage:   "Unable to list the directory tree in the password store",
					errors.FieldAction:    "tree",
					errors.FieldError:     err.Error(),
					errors.FieldStoreID:   store.ID,
					errors.FieldStoreName: store.Name,
					errors.FieldStorePath: store.Path,
				},
			)
		}

		for i, directory := range directories {
			relativePath, err := filepath.Rel(store.Path, directory)
			if err != nil {
				log.Errorf(
					"Unable to determine the relative path for a file '%v' in the password store '%+v': %+v",
					directory, store, err,
				)
				response.SendErrorAndExit(
					errors.CodeUnableToDetermineRelativeDirectoryPathInPasswordStore,
					&map[errors.Field]string{
						errors.FieldMessage:   "Unable to determine the relative path for a directory in the password store",
						errors.FieldAction:    "tree",
						errors.FieldError:     err.Error(),
						errors.FieldDirectory: directory,
						errors.FieldStoreID:   store.ID,
						errors.FieldStoreName: store.Name,
						errors.FieldStorePath: store.Path,
					},
				)
			}
			directories[i] = strings.Replace(relativePath, "\\", "/", -1) // normalize Windows paths
		}

		sort.Strings(directories)
		responseData.Directories[store.ID] = directories
	}

	response.SendOk(responseData)
}
