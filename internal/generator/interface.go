package generator

import (
	"fmt"

	"github.com/brahms116/between/internal/ast"
)

type OutputFormat string

const (
	TypescriptOut OutputFormat = "Typescript"
)

func GenerateOutput(definitions []ast.Definition, outputFormat OutputFormat) string {
	if outputFormat == TypescriptOut {
		return printTsDefinitions(definitions)
	}
	panic(fmt.Sprintf("Unknown output format %s", outputFormat))
}
