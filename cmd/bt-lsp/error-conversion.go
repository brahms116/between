package main

import "github.com/brahms116/between/internal/parser"

func lspErrorsToDiagnostics(errs []error) []Diagnostic {
	diagnostics := make([]Diagnostic, 0, len(errs))
	for _, err := range errs {
		if d := lspErrorToDiagnostic(err); d != nil {
			diagnostics = append(diagnostics, *d)
		}
	}
	return diagnostics
}

func lspErrorToDiagnostic(err error) *Diagnostic {
	switch e := err.(type) {
	case parser.UnexpectedTokenError:
		return &Diagnostic{
			Range:    lexLocationToLspRange(e.Actual.Loc),
			Severity: &DiagnosticSeverityError,
			Message:  e.LspMessage(),
		}
	default:
		return nil
	}
}
