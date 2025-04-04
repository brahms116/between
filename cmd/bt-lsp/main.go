package main

import "os"

func main() {
	logger := newLogger("server.log")
	transport := newTransport(os.Stdin, os.Stdout, logger)
	server := newServer(transport, logger)
	server.Start()
}
