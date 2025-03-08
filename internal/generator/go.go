package generator

import (
	"fmt"
	"strings"

	"github.com/brahms116/between/internal/ast"
)

var GO_PRIMITIVES map[string]string = map[string]string{
	"Float":  "float32",
	"Str":    "string",
	"Bool":   "bool",
	"Int":    "int",
	"Any":    "any",
	"Object": "map[string]any",
}

func printGoDefinitions(ds []ast.Definition) string {
	var definitionString string
	for _, d := range ds {
		definitionString += printGoDefinition(d)
	}
	return "package main;" + definitionString
}

func printGoDefinition(d ast.Definition) string {
	if d.SumStr != nil {
		return printGoSumStr(*d.SumStr)
	}
	if d.Sum != nil {
		return printGoSum(*d.Sum)
	}
	if d.Product != nil {
		return printGoProduct(*d.Product)
	}
	panic("Invalid definition")
}

func printGoSumStr(s ast.SumStr) string {
	typeDec := fmt.Sprintf(`type %s string;`, s.Id)
	var variantsString string
	for _, variant := range s.Variants {
		// Variants can't be optional, yet?
		variantName := s.Id + "_" + variant
		variantsString += fmt.Sprintf(`const %s = "%s";`, variantName, variant)
	}
	return typeDec + variantsString
}

func printGoSum(s ast.Sum) string {
	var variantsString string
	for _, variant := range s.Variants {
		// Variants can't be optional, yet?
		variantsString += fmt.Sprintf(`%s`, printGoField(variant, true))
	}
	return fmt.Sprintf(`type %s struct { %s};`, s.Id, variantsString)
}

func printGoProduct(p ast.Product) string {
	var fieldsString string
	for _, field := range p.Fields {
		fieldsString += printGoField(field, false) + " "
	}
	return fmt.Sprintf(`type %s struct { %s};`, p.Id, fieldsString)
}

func printGoField(f ast.Field, forcePointer bool) string {
	fieldName := capitalizeHead(f.Id)
	var omitEmptyTag string
	if f.Type.IsNullable() || forcePointer {
		omitEmptyTag = ",omitEmpty"
	}
	jsonTag := fmt.Sprintf("`json:\"%s%s\"`", f.Id, omitEmptyTag)

	return fmt.Sprintf(`%s %s %s;`, fieldName, printGoType(f.Type, forcePointer), jsonTag)
}

func printGoType(t ast.Type, forcePointer bool) string {
	var optionalString string
	if t.IsNullable() || forcePointer {
		optionalString = "*"
	}
	if t.List != nil {
		return fmt.Sprintf(`%s[]%s`, optionalString, printGoType(t.List.Type, false))
	}
	typeString, ok := GO_PRIMITIVES[t.TypeIdent.Id]
	if !ok {
		typeString = t.TypeIdent.Id
	}
	return optionalString + typeString
}

func capitalizeHead(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}
