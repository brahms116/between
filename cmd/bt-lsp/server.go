package main

import (
	"log"
)

type Server struct {
	logger    *log.Logger
	transport *transport
}

func newServer(transport *transport, logger *log.Logger) *Server {
	return &Server{
		logger:    logger,
		transport: transport,
	}
}

func (s *Server) Start() {
	s.logger.Println("Starting server")
	for {
		method, msg, ok := s.transport.next()
		if !ok {
			s.logger.Println("No more messages to process")
			break
		}
		s.logger.Printf("Method %s, msg: %s", method, msg)
		switch method {
		case "initialize":
			s.handleInitialize(msg)
		case "initialized":
			// Handle initialized notification
		default:
			s.logger.Println("Unknown method:", method)
		}
	}
}
