package main

import (
	"log"

	"github.com/brahms116/between/internal/parser"
	"github.com/brahms116/between/internal/st"
)

type documentState struct {
	logger      *log.Logger
	text        string
	syntaxTree  []st.Definition
	diagnostics []Diagnostic
}

func newDocumentState(text string, logger *log.Logger) documentState {
	ds := documentState{}
	ds.logger = logger
	ds.updateText(text)
	return ds
}

func (ds *documentState) updateText(text string) {
	ds.text = text
	tree, errs := parser.LexAndParse(text)
	ds.syntaxTree = tree
	ds.diagnostics = parseAndLexErrorsToDiagnostics(errs)
}
