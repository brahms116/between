package main

func (s *Server) handleDidChangeTextDocument(params DidChangeTextDocumentParams) {
  s.state.updateDocument(params.TextDocument.URI, params.ContentChanges[0].Text)
}
