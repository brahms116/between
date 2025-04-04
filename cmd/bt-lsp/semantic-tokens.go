package main

func (s *Server) handleSemanticTokensFull(params SemanticTokensParams) SemanticTokens {
	return SemanticTokens{
		Data: []int{
			0, 0, 4, 1, 0,
		},
	}
}
