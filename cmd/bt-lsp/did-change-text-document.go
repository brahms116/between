package main

func (s *Server) handleDidChangeTextDocument(params DidChangeTextDocumentParams) {
	if _, ok := s.state.documents[params.TextDocument.URI]; !ok {
		s.logger.Println("Document not found:", params.TextDocument.URI)
		return
	}

	for _, change := range params.ContentChanges {
		s.state.documents[params.TextDocument.URI] = change.Text
	}
}
