package main

import (
	"bufio"
	"io"
	"log"
)

type transport struct {
	reader *bufio.Scanner
	writer io.Writer
	logger *log.Logger
}

func newTransport(reader io.Reader, writer io.Writer, logger *log.Logger) *transport {
	scanner := bufio.NewScanner(reader)
	scanner.Split(lspMsgSplit)
	return &transport{
		reader: scanner,
		writer: writer,
		logger: logger,
	}
}

func (t *transport) next() (string, []byte, bool) {
	if !t.reader.Scan() {
		return "", nil, false
	}
	msg := t.reader.Bytes()
	method, err := methodFromBody(msg)
	if err != nil {
		t.logger.Println("Error parsing method:", err)
		return t.next()
	}
	return method, msg, true
}

func (t *transport) write(msg any) {
  t.writer.Write([]byte(encodeMessage(msg)))
}
