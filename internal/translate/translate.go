package translate

import (
	"strings"

	"github.com/brahms116/between/internal/ast"
	"github.com/brahms116/between/internal/st"
)

type translate struct {
	st []st.Definition
}

func Translate(st []st.Definition) ([]ast.Definition, error) {
	t := &translate{st: st}
	return t.translate()
}

func (t *translate) translate() ([]ast.Definition, error) {
	var res []ast.Definition
	for _, d := range t.st {
		def, err := t.translateDefinition(d)
		if err != nil {
			return nil, err
		}
		res = append(res, def)
	}
	return res, nil
}

func (t *translate) translateDefinition(d st.Definition) (ast.Definition, error) {
	if d.Product != nil {
		prod, err := t.translateProduct(*d.Product)
		if err != nil {
			return ast.Definition{}, err
		}
		return ast.Definition{Product: &prod}, nil
	}
	if d.Sum != nil {
		sum, err := t.translateSum(*d.Sum)
		if err != nil {
			return ast.Definition{}, err
		}
		return ast.Definition{Sum: &sum}, nil
	}
	if d.SumStr != nil {
		sumStr, err := t.translateSumStr(*d.SumStr)
		if err != nil {
			return ast.Definition{}, err
		}
		return ast.Definition{SumStr: &sumStr}, nil
	}
	panic("unreachable")
}

func (t *translate) translateType(ty st.Type) (ast.Type, error) {
	if ty.List != nil {
		list, err := t.translateList(*ty.List)
		if err != nil {
			return ast.Type{}, err
		}
		return ast.Type{List: &list}, nil
	}
	if ty.TypeIdent != nil {
		ti, err := t.translateTypeIdent(*ty.TypeIdent)
		if err != nil {
			return ast.Type{}, err
		}
		return ast.Type{TypeIdent: &ti}, nil
	}
	panic("unreachable")
}

func (t *translate) translateTypeIdent(ti st.TypeIdent) (ast.TypeIdent, error) {
	return ast.TypeIdent{
		Id:       ti.Id.Value,
		Nullable: ti.Nullable != nil,
	}, nil
}

func (t *translate) translateList(l st.List) (ast.List, error) {
	ty, err := t.translateType(l.Type)
	if err != nil {
		return ast.List{}, err
	}
	return ast.List{
		Nullable: l.Nullable != nil,
		Type:     ty,
	}, nil
}

func (t *translate) translateField(f st.Field) (ast.Field, error) {
	if f.FieldFull != nil {
		var jsonName *string
		if f.FieldFull.JsonName != nil {
			jsonName = &f.FieldFull.JsonName.Value
		}

		ty, err := t.translateType(f.FieldFull.Type)
		if err != nil {
			return ast.Field{}, err
		}

		return ast.Field{
			Id:       f.FieldFull.Id.Value,
			JsonName: jsonName,
			Type:     ty,
		}, nil
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
		}, nil
	}
	panic("unreachable")
}

func (t *translate) translateProduct(p st.Product) (ast.Product, error) {
	var fields []ast.Field
	for _, f := range p.Fields {
		field, err := t.translateField(f)
		if err != nil {
			return ast.Product{}, err
		}
		fields = append(fields, field)
	}
	return ast.Product{
		Id:     p.Id.Value,
		Fields: fields,
	}, nil
}

func (t *translate) translateSum(s st.Sum) (ast.Sum, error) {
	var variants []ast.Field
	for _, v := range s.Variants {
		variant, err := t.translateField(v)
		if err != nil {
			return ast.Sum{}, err
		}
		variants = append(variants, variant)
	}
	return ast.Sum{
		Id:       s.Id.Value,
		Variants: variants,
	}, nil
}

func (t *translate) translateSumStr(ss st.SumStr) (ast.SumStr, error) {
	var variants []ast.SumStrVariant
	for _, v := range ss.Variants {
		variant, err := t.translateSumStrVariant(v)
		if err != nil {
			return ast.SumStr{}, err
		}
		variants = append(variants, variant)
	}
	return ast.SumStr{
		Id:       ss.Id.Value,
		Variants: variants,
	}, nil
}

func (t *translate) translateSumStrVariant(ssv st.SumStrVariant) (ast.SumStrVariant, error) {
	var jsonName *string
	if ssv.JsonName != nil {
		jsonName = &ssv.JsonName.Value
	}
	return ast.SumStrVariant{
		Id:       ssv.Id.Value,
		JsonName: jsonName,
	}, nil
}

func lowerCaseFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToLower(string(s[0])) + s[1:]
}
