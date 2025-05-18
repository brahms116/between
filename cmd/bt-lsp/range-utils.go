package main

import "github.com/brahms116/between/internal/lex"

func lexPointToLspPosition(point lex.Point) Position {
	return Position{
		point.Row,
		point.Col,
	}
}

func lexLocationToLspRange(l lex.Location) Range {
	return Range{
		Start: lexPointToLspPosition(l.Start),
		End:   lexPointToLspPosition(l.End),
	}
}
