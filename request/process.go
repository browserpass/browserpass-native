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

type store struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type settings struct {
	GpgPath string           `json:"gpgPath"`
	Stores  map[string]store `json:"stores"`
}

type request struct {
	Action   string   `json:"action"`
	Settings settings `json:"settings"`
	File     string   `json:"file"`
	Store    string   `json:"store"`
}

// Process handles browser request
func Process() {
	requestLength := parseRequestLength()
	request := parseRequest(requestLength)

	switch request.Action {
	case "configure":
		configure(request)
	case "list":
		listFiles(request)
	case "fetch":
		fetchDecryptedContents(request)
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

// Request length is the first 4 bytes in LittleEndian encoding on stdin
func parseRequestLength() uint32 {
	var length uint32
	if err := binary.Read(os.Stdin, binary.LittleEndian, &length); err != nil {
		log.Error("Unable to parse the length of the browser request: ", err)
		response.SendErrorAndExit(
			errors.CodeParseRequestLength,
			&map[errors.Field]string{
				errors.FieldMessage: "Unable to parse the length of the browser request",
				errors.FieldError:   err.Error(),
			},
		)
	}
	return length
}

// Request is a json with a predefined structure
func parseRequest(messageLength uint32) request {
	var parsed request
	reader := &io.LimitedReader{R: os.Stdin, N: int64(messageLength)}
	if err := json.NewDecoder(reader).Decode(&parsed); err != nil {
		log.Error("Unable to parse the browser request: ", err)
		response.SendErrorAndExit(
			errors.CodeParseRequest,
			&map[errors.Field]string{
				errors.FieldMessage: "Unable to parse the browser request",
				errors.FieldError:   err.Error(),
			},
		)
	}
	return parsed
}
