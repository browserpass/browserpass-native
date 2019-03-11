package request

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
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
