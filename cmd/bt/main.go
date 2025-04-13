package main

import (
	"log"
	"os"
	"strings"

	"github.com/brahms116/between/internal/generator"
	"github.com/brahms116/between/internal/parser"
	"github.com/brahms116/between/internal/translate"
)

type OutputFormat string

const (
	TypescriptOut OutputFormat = "Typescript"
	GolangOut     OutputFormat = "Golang"
)

var extentionOutputMap map[string]OutputFormat = map[string]OutputFormat{
	"ts": TypescriptOut,
	"go": GolangOut,
}

func parseOutputFileDetails(outputFileLocation string) (filename string, format OutputFormat) {
	parts := strings.Split(outputFileLocation, "/")
	fileName := parts[len(parts)-1]
	parts = strings.Split(fileName, ".")
	extension := parts[len(parts)-1]
	fileName = parts[0]

	outputFormat, ok := extentionOutputMap[extension]
	if !ok {
		log.Panic("Unsupported output format")
	}

	return fileName, outputFormat
}

func main() {
	args, err := newFlags()
	if err != nil {
		log.Panic(err)
	}

	fileName, outputFormat := parseOutputFileDetails(args.outputFileLocation)

	input, err := os.ReadFile(args.inputFileLocation)
	if err != nil {
		log.Panic(err)
	}

	st, err := parser.LexAndParse(string(input))
	definitions, primitives ,err := translate.Translate(st)
	if err != nil {
		log.Panic(err)
	}

	var output string
	switch outputFormat {
	case TypescriptOut:
		output = generator.PrintTsDefinitions(definitions)
	case GolangOut:
		goPackageName := args.goPackageName
		if goPackageName == "" {
			goPackageName = fileName
		}
		output = generator.PrintGoDefinitions(definitions, primitives, generator.GoGeneratorOptions{PackageName: goPackageName})
	}

	err = os.WriteFile(args.outputFileLocation, []byte(output), 0644)
	if err != nil {
		log.Panic(err)
	}
}
