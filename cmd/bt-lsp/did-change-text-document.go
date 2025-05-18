package main

func (s *Server) handleDidChangeTextDocument(params DidChangeTextDocumentParams) {
  s.updateDocument(params.TextDocument.URI, params.ContentChanges[0].Text)
}
