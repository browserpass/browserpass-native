package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"
)

const (
	errorCodeReadRequestLength = iota + 1
	errorCodeReadRequest
)

type okResponse struct {
	Status  string      `json:"status"`
	Version int         `json:"version"`
	Data    interface{} `json:"data"`
}

type errorResponse struct {
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Version int         `json:"version"`
	Params  interface{} `json:"params"`
}

func sendError(errorCode int, errorMsg string) {
	send(&errorResponse{
		Status:  "error",
		Code:    errorCode,
		Version: versionCode,
		Params: map[string]string{
			"message": errorMsg,
		},
	})
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
