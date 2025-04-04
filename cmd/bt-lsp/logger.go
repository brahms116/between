package main

import (
	"log"
	"os"
)

func newLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("Failed to open log file: " + err.Error())
	}

	return log.New(logfile, "[between-lsp] ", log.Ldate|log.Ltime|log.Lshortfile)
}
