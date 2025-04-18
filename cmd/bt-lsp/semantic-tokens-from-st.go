package main

import (
	"github.com/brahms116/between/internal/lex"
	"github.com/brahms116/between/internal/st"
)

func pointDelta(a, b lex.Point) lex.Point {
	rowDelta := a.Row - b.Row

	colDelta := a.Col
	if rowDelta == 0 {
		colDelta = a.Col - b.Col
	}
	return lex.Point{
		Row: rowDelta,
		Col: colDelta,
	}
}

type treeToSemanticTokens struct {
	lastStartPoint lex.Point
	sematicTokens  []int
}

func newTreeToSemanticTokens() *treeToSemanticTokens {
	return &treeToSemanticTokens{
		sematicTokens: make([]int, 0, 1000),
	}
}

func (t *treeToSemanticTokens) addSemanticToken(semanticTokenIndex int, location lex.Location) {
	delta := pointDelta(location.Start, t.lastStartPoint)
	t.sematicTokens = append(
		t.sematicTokens,
		delta.Row,
		delta.Col,
		location.ByteEnd-location.ByteStart,
		semanticTokenIndex,
		0,
	)
  t.lastStartPoint = location.Start
}

func convertTreeToSemanticTokens(definitions []st.Definition) []int {
	t := newTreeToSemanticTokens()
	t.convertDefinitions(definitions)
	return t.sematicTokens
}

func (t *treeToSemanticTokens) convertDefinitions(ds []st.Definition) {
	for _, d := range ds {
		t.convertDefinition(d)
	}
}

func (t *treeToSemanticTokens) convertDefinition(d st.Definition) {
	if d.Product != nil {
		t.convertProduct(*d.Product)
	} else if d.Sum != nil {
		t.convertSum(*d.Sum)
	} else if d.SumStr != nil {
		t.convertSumStr(*d.SumStr)
	} else {
		panic("unreachable")
	}
}

func (t *treeToSemanticTokens) convertSum(s st.Sum) {
	t.addSemanticToken(SEMTOK_KEYWORD_INDEX, s.Keyword.Loc)
	t.addSemanticToken(SEMTOK_CLASS_INDEX, s.Id.Loc)
	for _, v := range s.Variants {
		t.convertField(v)
	}
}

func (t *treeToSemanticTokens) convertProduct(p st.Product) {
	t.addSemanticToken(SEMTOK_KEYWORD_INDEX, p.Keyword.Loc)
	t.addSemanticToken(SEMTOK_CLASS_INDEX, p.Id.Loc)
	for _, f := range p.Fields {
		t.convertField(f)
	}
}

func (t *treeToSemanticTokens) convertField(f st.Field) {
	if f.FieldFull != nil {
		t.addSemanticToken(SEMTOK_PROPERTY_INDEX, f.FieldFull.Id.Loc)
		if f.FieldFull.JsonName != nil {
			t.addSemanticToken(SEMTOK_STRING_INDEX, f.FieldFull.JsonName.Loc)
		}
		t.convertType(f.FieldFull.Type)
	} else if f.FieldShort != nil {
		t.addSemanticToken(SEMTOK_CLASS_INDEX, f.FieldShort.Id.Loc)
	} else {
		panic("unreachable")
	}
}

func (t *treeToSemanticTokens) convertType(ty st.Type) {
	if ty.List != nil {
		t.convertType(ty.List.Type)
	} else if ty.TypeIdent != nil {
		t.addSemanticToken(SEMTOK_CLASS_INDEX, ty.TypeIdent.Id.Loc)
	}
}

func (t *treeToSemanticTokens) convertSumStr(ss st.SumStr) {
	t.addSemanticToken(SEMTOK_KEYWORD_INDEX, ss.Keyword.Loc)
	t.addSemanticToken(SEMTOK_CLASS_INDEX, ss.Id.Loc)
	for _, v := range ss.Variants {
		t.convertSumStrVariant(v)
	}
}

func (t *treeToSemanticTokens) convertSumStrVariant(sv st.SumStrVariant) {
	t.addSemanticToken(SEMTOK_ENUM_MEMBER_INDEX, sv.Id.Loc)
	if sv.JsonName != nil {
		t.addSemanticToken(SEMTOK_STRING_INDEX, sv.JsonName.Loc)
	}
}
