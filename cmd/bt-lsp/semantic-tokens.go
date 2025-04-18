package main

import (
	"github.com/brahms116/between/internal/parser"
)

func (s *Server) handleSemanticTokensFull(params SemanticTokensParams) SemanticTokens {

	tree, err := parser.LexAndParse(s.state.documents[params.TextDocument.URI])
	if err != nil {
		return SemanticTokens{
			Data: make([]int, 0),
		}
	}
  tokens := convertTreeToSemanticTokens(tree)
  s.logger.Printf("tokens: %v \n", tokens)

	return SemanticTokens{
		Data: tokens,
  }
}
