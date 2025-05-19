package main

import (
	"github.com/brahms116/between/internal/lex"
	"github.com/brahms116/between/internal/parser"
)

func parseAndLexErrorsToDiagnostics(errs []error) []Diagnostic {
	diagnostics := make([]Diagnostic, 0, len(errs))
	for _, err := range errs {
		if d := parseAndLexErrorToDiagnostic(err); d != nil {
			diagnostics = append(diagnostics, *d)
		}
	}
	return diagnostics
}

func parseAndLexErrorToDiagnostic(err error) *Diagnostic {
	switch e := err.(type) {
	case parser.UnexpectedTokenError:
		return &Diagnostic{
			Range:    lexLocationToLspRange(e.Actual.Loc),
			Severity: &DiagnosticSeverityError,
			Message:  e.LspMessage(),
		}
  case lex.UnexpectedCharError:
    pt := lexPointToLspPosition(e.Point)
    return &Diagnostic{
      Severity: &DiagnosticSeverityError,
      Message:  e.LspMessage(),
      Range: Range{
        Start: pt,
        End:   pt,
      },
    }
	default:
		return nil
	}
}
