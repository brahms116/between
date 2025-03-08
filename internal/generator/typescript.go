package generator

import (
	"fmt"

	"github.com/brahms116/between/internal/ast"
)

var TS_PRIMITIVES map[string]string = map[string]string{
	"Float":  "number",
	"Str":    "string",
	"Bool":   "boolean",
	"Int":    "number",
	"Any":    "unknown",
	"Object": "Record<string, unknown>",
}

func printTsDefinitions(ds []ast.Definition) string {
	var definitionString string
	for _, d := range ds {
		definitionString += printTsDefinition(d)
	}
	return definitionString
}

func printTsDefinition(d ast.Definition) string {
	if d.SumStr != nil {
		return printTsSumStr(*d.SumStr)
	}
	if d.Sum != nil {
		return printTsSum(*d.Sum)
	}
	if d.Product != nil {
		return printTsProduct(*d.Product)
	}
	panic("Invalid definition")
}

func printTsSumStr(s ast.SumStr) string {
	var variantsString string
	for _, variant := range s.Variants {
		variantsString += fmt.Sprintf(`| "%s" `, variant)
	}
	return fmt.Sprintf(`type %s = %s; `, s.Id, variantsString)
}

func printTsSum(s ast.Sum) string {
	var variantsString string
	for _, variant := range s.Variants {
		variantsString += fmt.Sprintf(`| { _type:"%s"; %s} `, variant.Id, printTsField(variant))
	}
	return fmt.Sprintf(`type %s = %s; `, s.Id, variantsString)
}

func printTsProduct(p ast.Product) string {
	var fieldsString string
	for _, field := range p.Fields {
		fieldsString += printTsField(field) + " "
	}
	return fmt.Sprintf(`interface %s { %s}; `, p.Id, fieldsString)
}

func printTsField(f ast.Field) string {
	nullable, typeString := printTsType(f.Type)
	var nullableString string
	if nullable {
		nullableString = "?"
	}
	return fmt.Sprintf(`%s%s: %s;`, f.Id, nullableString, typeString)
}

func printTsType(t ast.Type) (bool, string) {
	return t.IsNullable(), printTsTypeTail(t, true)
}

func printTsTypeTail(t ast.Type, isTopLevel bool) string {
	if t.List != nil {
		if t.List.Nullable && !isTopLevel {
			return fmt.Sprintf(`(%s[]|undefined)`, printTsTypeTail(t.List.Type, false))
		}
		return fmt.Sprintf(`%s[]`, printTsTypeTail(t.List.Type, false))
	}
	typeString, ok := TS_PRIMITIVES[t.TypeIdent.Id]
	if !ok {
		typeString = t.TypeIdent.Id
	}
	if t.TypeIdent.Nullable && !isTopLevel {
		return fmt.Sprintf(`(%s|undefined)`, typeString)
	}
	return typeString
}
