package main

func (s *Server) handleDidOpenTextDocument(params DidOpenTextDocumentParams) {
  s.addDocument(params.TextDocument.URI, params.TextDocument.Text)
}

