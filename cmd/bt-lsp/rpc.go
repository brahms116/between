package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

func encodeMessage(item any) string {
	content, err := json.Marshal(item)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(content), content)
}

func methodFromBody(msg []byte) (string, error) {
	var base struct {
		Method string `json:"method"`
	}
	err := json.Unmarshal(msg, &base)
	if err != nil {
		return "", fmt.Errorf("failed to message to json: %w", err)
	}
	return base.Method, nil
}

func lspMsgSplit(msg []byte, _ bool) (advance int, token []byte, err error) {
	header, body, found := bytes.Cut(msg, []byte("\r\n\r\n"))
	if !found {
		return 0, nil, nil
	}

	contentLengthBytes := header[len("Content-Length: "):]
	contentLength, err := strconv.Atoi(string(contentLengthBytes))
	if err != nil {
		return 0, nil, fmt.Errorf("failed to parse length from content-length: %w", err)
	}
	if contentLength > len(body) {
		return 0, nil, nil
	}
	totalLength := contentLength + len(header) + len("\r\n\r\n")
	return totalLength, body[:contentLength], nil
}
