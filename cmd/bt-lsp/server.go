package main

import (
	"encoding/json"
	"log"
)

const JSON_RPC = "2.0"

type Server struct {
	logger    *log.Logger
	transport *transport
	state     *lspState
}

func newServer(transport *transport, logger *log.Logger) *Server {
	return &Server{
		logger:    logger,
		transport: transport,
		state:     newLspState(logger),
	}
}

func handleNotification[T any](s *Server, body []byte, handler func(body T)) {
	var notification RpcNotification[T]
	if err := json.Unmarshal(body, &notification); err != nil {
		s.logger.Println("Error unmarshalling notification:", err)
		return
	}
	handler(notification.Params)
}

func handleRequest[T any, K any](s *Server, body []byte, handler func(body T) K) {
	var request RpcRequest[T]
	if err := json.Unmarshal(body, &request); err != nil {
		s.logger.Println("Error unmarshalling request:", err)
		return
	}
	response := handler(request.Params)
	s.transport.write(RpcResponse[K]{
		RPC:    JSON_RPC,
		ID:     request.ID,
		Result: &response,
	})
}

func (s *Server) sendNotification(method string, params any) {
	notification := RpcNotification[any]{
		RPC:    JSON_RPC,
		Method: method,
		Params: params,
	}
	s.transport.write(notification)
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
			handleRequest(s, msg, s.handleInitialize)
		case "initialized":
			// Handle initialized notification
		case "textDocument/didOpen":
			handleNotification(s, msg, s.handleDidOpenTextDocument)
		case "textDocument/didChange":
			handleNotification(s, msg, s.handleDidChangeTextDocument)
		case "textDocument/semanticTokens/full":
			handleRequest(s, msg, s.handleSemanticTokensFull)
		default:
			s.logger.Println("Unknown method:", method)
		}
	}
}
