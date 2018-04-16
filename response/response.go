package response

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"os"

	"github.com/browserpass/browserpass-native/errors"
	"github.com/browserpass/browserpass-native/version"
	log "github.com/sirupsen/logrus"
)

type okResponse struct {
	Status  string      `json:"status"`
	Version int         `json:"version"`
	Data    interface{} `json:"data"`
}

type errorResponse struct {
	Status  string      `json:"status"`
	Code    errors.Code `json:"code"`
	Version int         `json:"version"`
	Params  interface{} `json:"params"`
}

// ConfigureResponse a response format for the "configure" request
type ConfigureResponse struct {
	DefaultStore struct {
		Path     string `json:"path"`
		Settings string `json:"settings"`
	} `json:"defaultStore"`
	StoreSettings map[string]string `json:"storeSettings"`
}

// MakeConfigureResponse initializes an empty configure response
func MakeConfigureResponse() *ConfigureResponse {
	return &ConfigureResponse{
		StoreSettings: make(map[string]string),
	}
}

// ListResponse a response format for the "list" request
type ListResponse struct {
	Files map[string][]string `json:"files"`
}

// MakeListResponse initializes an empty list response
func MakeListResponse() *ListResponse {
	return &ListResponse{
		Files: make(map[string][]string),
	}
}

// FetchResponse a response format for the "fetch" request
type FetchResponse struct {
	Contents string `json:"contents"`
}

// MakeFetchResponse initializes an empty fetch response
func MakeFetchResponse() *FetchResponse {
	return &FetchResponse{}
}

// SendOk sends a success response to the browser extension in the predefined json format
func SendOk(data interface{}) {
	send(&okResponse{
		Status:  "ok",
		Version: version.Code,
		Data:    data,
	})
}

// SendErrorAndExit sends an error response to the browser extension in the predefined json format and exits with the specified exit code
func SendErrorAndExit(errorCode errors.Code, params *map[errors.Field]string) {
	send(&errorResponse{
		Status:  "error",
		Code:    errorCode,
		Version: version.Code,
		Params:  params,
	})

	errors.ExitWithCode(errorCode)
}

func send(data interface{}) {
	switch data.(type) {
	case *okResponse:
	case *errorResponse:
		break
	default:
		log.Fatalf("Only data of type OkResponse and ErrorResponse is allowed to be sent to the browser extension, attempted to send: %+v", data)
	}

	var bytesBuffer bytes.Buffer
	if err := json.NewEncoder(&bytesBuffer).Encode(data); err != nil {
		log.Fatal("Unable to encode data for sending: ", err)
	}

	if err := binary.Write(os.Stdout, binary.LittleEndian, uint32(bytesBuffer.Len())); err != nil {
		log.Fatal("Unable to send the length of the response data: ", err)
	}
	if _, err := bytesBuffer.WriteTo(os.Stdout); err != nil {
		log.Fatal("Unable to send the response data: ", err)
	}
}
