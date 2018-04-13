package main

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

type request struct {
	Action   string      `json:"action"`
	Settings interface{} `json:"settings"`
}

func process() {
	requestLength := parseRequestLength()
	request := parseRequest(requestLength)

	switch request.Action {
	case "configure":
		break
	case "list":
		break
	case "fetch":
		break
	default:
		log.Errorf("Received a browser request with an unknown action: %+v", request)
	}
}

// Request length is the first 4 bytes in LittleEndian encoding on stdin
func parseRequestLength() uint32 {
	var length uint32
	if err := binary.Read(os.Stdin, binary.LittleEndian, &length); err != nil {
		// TODO: Original browserpass ignores EOF as if it is expected, is it true?
		// if err == io.EOF {
		// 	return
		// }
		log.Error("Unable to read the length of the browser request: ", err)
		sendError(errorCodeReadRequestLength, "Unable to read the length of the browser request")
		os.Exit(errorCodeReadRequestLength)
	}
	return length
}

// Request is a json with a predefined structure
func parseRequest(messageLength uint32) request {
	var parsed request
	reader := &io.LimitedReader{R: os.Stdin, N: int64(messageLength)}
	if err := json.NewDecoder(reader).Decode(&parsed); err != nil {
		log.Error("Unable to read the browser request: ", err)
		sendError(errorCodeReadRequest, "Unable to read the browser request")
		os.Exit(errorCodeReadRequest)
	}
	return parsed
}
