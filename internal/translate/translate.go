package translate

import (
	"fmt"
	"strings"

	"github.com/brahms116/between/internal/ast"
	"github.com/brahms116/between/internal/lex"
	"github.com/brahms116/between/internal/st"
)

var PrimitiveTypes = map[string]struct{}{
	"Float":  {},
	"Str":    {},
	"Bool":   {},
	"Int":    {},
	"Any":    {},
	"Object": {},
	"Date":   {},
}

type TypeError struct {
	Message  string
	Location lex.Location
}

func (e TypeError) LspMessage() string {
	return e.Message
}

func (e TypeError) Error() string {
	return fmt.Sprintf("Type error at %s: %s", e.Location.Start.String(), e.Message)
}

func newTypeError(message string, loc lex.Location) TypeError {
	return TypeError{
		Message:  message,
		Location: loc,
	}
}

type translate struct {
	st                 []st.Definition
	usedPrimitiveTypes map[string]struct{}
	errors             []error
}

func newTranslate(st []st.Definition) *translate {
	return &translate{
		st:                 st,
		usedPrimitiveTypes: make(map[string]struct{}),
	}
}

func Translate(st []st.Definition) ([]ast.Definition, map[string]struct{}, []error) {
	t := newTranslate(st)
	return t.translate()
}

func (t *translate) addError(message string, location lex.Location) {
	t.errors = append(t.errors, newTypeError(message, location))
}

func (t *translate) translate() ([]ast.Definition, map[string]struct{}, []error) {
	var res []ast.Definition
	for _, d := range t.st {
		def:= t.translateDefinition(d)
		res = append(res, def)
	}
	return res, t.usedPrimitiveTypes, t.errors
}

func (t *translate) translateDefinition(d st.Definition) ast.Definition {
	if d.Product != nil {
		prod := t.translateProduct(*d.Product)
		return ast.Definition{Product: &prod}
	}
	if d.Sum != nil {
		sum := t.translateSum(*d.Sum)
		return ast.Definition{Sum: &sum}
	}
	if d.SumStr != nil {
		sumStr := t.translateSumStr(*d.SumStr)
		return ast.Definition{SumStr: &sumStr}
	}
	panic("unreachable")
}

func (t *translate) translateType(ty st.Type) ast.Type {
	if ty.List != nil {
		list := t.translateList(*ty.List)
		return ast.Type{List: &list}
	}
	if ty.TypeIdent != nil {
		ti := t.translateTypeIdent(*ty.TypeIdent)
		return ast.Type{TypeIdent: &ti}
	}
	panic("unreachable")
}

func (t *translate) translateTypeIdent(ti st.TypeIdent) ast.TypeIdent {
	if _, ok := PrimitiveTypes[ti.Id.Value]; ok {
		t.usedPrimitiveTypes[ti.Id.Value] = struct{}{}
	}

	return ast.TypeIdent{
		Id:       ti.Id.Value,
		Nullable: ti.Nullable != nil,
	}
}

func (t *translate) translateList(l st.List) ast.List {
	ty := t.translateType(l.Type)
	return ast.List{
		Nullable: l.Nullable != nil,
		Type:     ty,
	}
}

func (t *translate) translateField(f st.Field) ast.Field {
	if f.FieldFull != nil {
		var jsonName *string
		if f.FieldFull.JsonName != nil {
			jsonName = &f.FieldFull.JsonName.Value
		}

		ty := t.translateType(f.FieldFull.Type)

		return ast.Field{
			Id:       f.FieldFull.Id.Value,
			JsonName: jsonName,
			Type:     ty,
		}
	}
	if f.FieldShort != nil {
		id := lowerCaseFirstLetter(f.FieldShort.Id.Value)
		ty := ast.Type{
			TypeIdent: &ast.TypeIdent{
				Id:       f.FieldShort.Id.Value,
				Nullable: f.FieldShort.Nullable != nil,
			},
		}
		return ast.Field{
			Id:       id,
			JsonName: nil,
			Type:     ty,
		}
	}
	panic("unreachable")
}

func (t *translate) translateProduct(p st.Product) ast.Product {
	var fields []ast.Field
	for _, f := range p.Fields {
		field := t.translateField(f)
		fields = append(fields, field)
	}
	return ast.Product{
		Id:     p.Id.Value,
		Fields: fields,
	}
}

func (t *translate) translateSum(s st.Sum) ast.Sum {
	var variants []ast.Field
	for _, v := range s.Variants {
		variant := t.translateField(v)
		variants = append(variants, variant)
	}
	return ast.Sum{
		Id:       s.Id.Value,
		Variants: variants,
	}
}

func (t *translate) translateSumStr(ss st.SumStr) ast.SumStr {
	var variants []ast.SumStrVariant
	for _, v := range ss.Variants {
		variant := t.translateSumStrVariant(v)
		variants = append(variants, variant)
	}
	return ast.SumStr{
		Id:       ss.Id.Value,
		Variants: variants,
	}
}

func (t *translate) translateSumStrVariant(ssv st.SumStrVariant) ast.SumStrVariant {
	var jsonName *string
	if ssv.JsonName != nil {
		jsonName = &ssv.JsonName.Value
	}
	return ast.SumStrVariant{
		Id:       ssv.Id.Value,
		JsonName: jsonName,
	}
}

func lowerCaseFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToLower(string(s[0])) + s[1:]
}
