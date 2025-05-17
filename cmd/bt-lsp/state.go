package main

import (
	"github.com/brahms116/between/internal/parser"
	"github.com/brahms116/between/internal/st"
)

type lspState struct {
	documents map[string]documentState
}

func newLspState() *lspState {
	return &lspState{
		documents: make(map[string]documentState),
	}
}

func (s *lspState) addDocument(uri string, text string) {
	ds := newDocumentState(text)
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

type documentState struct {
	text       string
	syntaxTree []st.Definition
}

func newDocumentState(text string) documentState {
	ds := documentState{}
	ds.updateText(text)
	return ds
}

func (ds *documentState) updateText(text string) {
	ds.text = text
	tree, _ := parser.LexAndParse(text)
	ds.syntaxTree = tree
}
