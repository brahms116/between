package main

func (s *Server) handleDidOpenTextDocument(params DidOpenTextDocumentParams) {
  s.state.addDocument(params.TextDocument.URI, params.TextDocument.Text)
}

