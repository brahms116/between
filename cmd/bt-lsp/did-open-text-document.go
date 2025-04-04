package main

func (s *Server) handleDidOpenTextDocument(params DidOpenTextDocumentParams) {
	s.state.documents[params.TextDocument.URI] = params.TextDocument.Text
}

