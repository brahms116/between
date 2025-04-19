package main

import (
	"github.com/brahms116/between/internal/parser"
)

func (s *Server) handleSemanticTokensFull(params SemanticTokensParams) SemanticTokens {
	tree, _ := parser.LexAndParse(s.state.documents[params.TextDocument.URI])
	tokens := convertTreeToSemanticTokens(tree)
	s.logger.Printf("tokens: %v \n", tokens)
	return SemanticTokens{
		Data: tokens,
	}
}
