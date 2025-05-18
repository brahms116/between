package main

import "log"

func (s *Server) addDocument(uri string, text string) {
	s.state.addDocument(uri, text)
	s.publishDocumentDiagnostics(uri)
}

func (s *Server) updateDocument(uri string, text string) {
	s.state.updateDocument(uri, text)
	s.publishDocumentDiagnostics(uri)
}

func (s *Server) onDocumentChanged(uri string) {
	s.publishDocumentDiagnostics(uri)
}

type lspState struct {
	logger    *log.Logger
	documents map[string]documentState
}

func newLspState(logger *log.Logger) *lspState {
	return &lspState{
		logger:    logger,
		documents: make(map[string]documentState),
	}
}

func (s *lspState) addDocument(uri string, text string) {
	ds := newDocumentState(text, s.logger)
	s.documents[uri] = ds
}

func (s *lspState) updateDocument(uri string, text string) {
	if ds, ok := s.documents[uri]; ok {
		ds.updateText(text)
		s.documents[uri] = ds
	} else {
		panic("Document not found")
	}
}
