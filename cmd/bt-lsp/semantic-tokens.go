package main

import ()

func (s *Server) handleSemanticTokensFull(params SemanticTokensParams) SemanticTokens {
	tokens := convertTreeToSemanticTokens(
		s.state.documents[params.TextDocument.URI].syntaxTree,
  )
	s.logger.Printf("tokens: %v \n", tokens)
	return SemanticTokens{
		Data: tokens,
	}
}
