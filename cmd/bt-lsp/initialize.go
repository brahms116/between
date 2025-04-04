package main

import "encoding/json"

type InitializeRequest struct {
	Request
	Params InitializeRequestParams `json:"params"`
}

type InitializeRequestParams struct {
	ClientInfo *ClientInfo `json:"clientInfo"`
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type InitializeResponse struct {
	Response
	Result InitializeResult `json:"result"`
}

type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   ServerInfo         `json:"serverInfo"`
}

type ServerCapabilities struct {
	TextDocumentSync       int                    `json:"textDocumentSync"`
	SemanticTokensProvider *SemanticTokensOptions `json:"semanticTokensProvider,omitempty"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type SemanticTokensLegend struct {
	TokenTypes     []string `json:"tokenTypes"`
	TokenModifiers []string `json:"tokenModifiers"`
}

type SemanticTokensOptions struct {
	Legend SemanticTokensLegend `json:"legend"`
	Full   *bool                `json:"full,omitempty"`
}

func newInitializeResponse(id int) InitializeResponse {
	semanticTokensSyncFull := true
	return InitializeResponse{
		Response: Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: InitializeResult{
			Capabilities: ServerCapabilities{
				TextDocumentSync: 1,
				SemanticTokensProvider: &SemanticTokensOptions{
					Legend: SemanticTokensLegend{},
					Full:   &semanticTokensSyncFull,
				},
			},
			ServerInfo: ServerInfo{
				Name:    "between-lsp",
				Version: "0.0.0.0.0.0-beta1.final",
			},
		},
	}
}

func (s *Server) handleInitialize(body []byte) {
	var request InitializeRequest
	if err := json.Unmarshal(body, &request); err != nil {
		s.logger.Println("Error unmarshalling initialize request:", err)
		return
	}
	response := newInitializeResponse(request.ID)
	s.transport.write(response)
}
