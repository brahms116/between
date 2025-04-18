package main

func newInitializeResponse() InitializeResult {
	semanticTokensSyncFull := true
	return InitializeResult{
		Capabilities: ServerCapabilities{
			TextDocumentSync: 1,
			SemanticTokensProvider: &SemanticTokensOptions{
				Legend: SemanticTokensLegend{
					TokenTypes: []string{
						SEMTOK_CLASS,
						SEMTOK_KEYWORD,
            SEMTOK_PROPERTY,
            SEMTOK_STRING,
            SEMTOK_ENUM_MEMBER,
					},
				},
				Full: &semanticTokensSyncFull,
			},
		},
		ServerInfo: ServerInfo{
			Name:    "between-lsp",
			Version: "0.0.0.0.0.0-beta1.final",
		},
	}
}

func (s *Server) handleInitialize(_ InitializeParams) InitializeResult {
	return newInitializeResponse()
}
