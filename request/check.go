package request

import (
	"os"
	"path/filepath"

	"github.com/browserpass/browserpass-native/response"
)

type existsResponse struct {
	Exists bool `json:"exists"`
}

func checkFile(request *request) {
	exists := false

	for _, store := range request.Settings.Stores {
		normalizedStorePath, err := normalizePasswordStorePath(store.Path)
		if err != nil {
			continue // TODO: should respond with error
		}

		absoluteFilePath := filepath.Join(normalizedStorePath, request.File)
		_, err = os.Stat(absoluteFilePath)

		if err == nil {
			exists = true
		}
	}

	response.SendOk(&existsResponse{exists})
}
