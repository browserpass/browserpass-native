package request

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"os"

	"github.com/browserpass/browserpass-native/errors"
	"github.com/browserpass/browserpass-native/response"
	log "github.com/sirupsen/logrus"
)

type StoreSettings struct {
	GpgPath string `json:"gpgPath"`
}

type store struct {
	ID       string        `json:"id"`
	Name     string        `json:"name"`
	Path     string        `json:"path"`
	Settings StoreSettings `json:"settings"`
}

type settings struct {
	GpgPath string           `json:"gpgPath"`
	Stores  map[string]store `json:"stores"`
}

type request struct {
	Action       string      `json:"action"`
	Settings     settings    `json:"settings"`
	File         string      `json:"file"`
	Contents     string      `json:"contents"`
	StoreID      string      `json:"storeId"`
	EchoResponse interface{} `json:"echoResponse"`
}

// Process handles browser request
func Process() {
	requestLength, err := parseRequestLength(os.Stdin)
	if err != nil {
		log.Error("Unable to parse the length of the browser request: ", err)
		response.SendErrorAndExit(
			errors.CodeParseRequestLength,
			&map[errors.Field]string{
				errors.FieldMessage: "Unable to parse the length of the browser request",
				errors.FieldError:   err.Error(),
			},
		)
	}

	request, err := parseRequest(requestLength, os.Stdin)
	if err != nil {
		log.Error("Unable to parse the browser request: ", err)
		response.SendErrorAndExit(
			errors.CodeParseRequest,
			&map[errors.Field]string{
				errors.FieldMessage: "Unable to parse the browser request",
				errors.FieldError:   err.Error(),
			},
		)
	}

	switch request.Action {
	case "configure":
		configure(request)
	case "list":
		listFiles(request)
	case "tree":
		listDirectories(request)
	case "fetch":
		fetchDecryptedContents(request)
	case "save":
		saveEncryptedContents(request)
	case "delete":
		deleteFile(request)
	case "echo":
		response.SendRaw(request.EchoResponse)
	default:
		log.Errorf("Received a browser request with an unknown action: %+v", request)
		response.SendErrorAndExit(
			errors.CodeInvalidRequestAction,
			&map[errors.Field]string{
				errors.FieldMessage: "Invalid request action",
				errors.FieldAction:  request.Action,
			},
		)
	}
}

// Request length is the first 4 bytes in LittleEndian encoding
func parseRequestLength(input io.Reader) (uint32, error) {
	var length uint32
	if err := binary.Read(input, binary.LittleEndian, &length); err != nil {
		return 0, err
	}
	return length, nil
}

// Request is a json with a predefined structure
func parseRequest(messageLength uint32, input io.Reader) (*request, error) {
	var parsed request
	reader := &io.LimitedReader{R: input, N: int64(messageLength)}
	if err := json.NewDecoder(reader).Decode(&parsed); err != nil {
		return nil, err
	}
	return &parsed, nil
}
