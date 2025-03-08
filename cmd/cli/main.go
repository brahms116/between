package main

import (
	"log"
	"os"
	"strings"

	"github.com/brahms116/between/internal/generator"
	"github.com/brahms116/between/internal/parser"
)

var extentionOutputMap map[string]generator.OutputFormat = map[string]generator.OutputFormat{
	"ts": generator.TypescriptOut,
	"go": generator.GolangOut,
}

func main() {
	args, err := newFlags()
	if err != nil {
		log.Panic(err)
	}

	parts := strings.Split(args.outputFileLocation, ".")
	extension := parts[len(parts)-1]

	outputFormat, ok := extentionOutputMap[extension]
	if !ok {
		log.Panic("Unsupported output format")
	}

	input, err := os.ReadFile(args.inputFileLocation)
	if err != nil {
		log.Panic(err)
	}

	definitions, err := parser.LexAndParse(string(input))
	if err != nil {
		log.Panic(err)
	}

	output := generator.GenerateOutput(definitions, outputFormat)
	err = os.WriteFile(args.outputFileLocation, []byte(output), 0644)
	if err != nil {
		log.Panic(err)
	}
}
