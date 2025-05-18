package main

func (s *Server) publishDocumentDiagnostics(documentUri string) {
	ds, ok := s.state.documents[documentUri]
	if !ok {
		return
	}

	params := PublishDiagnosticsParams{
		DocumentURI: documentUri,
		Diagnostics: ds.diagnostics,
	}
	s.sendNotification(PublishDiagnosticsMethod, params)
}
