package main

import (
	"log"
	"os"

	"github.com/brahms116/between/internal/generator"
	"github.com/brahms116/between/internal/parser"
)

func main() {
	args, err := newFlags()
	if err != nil {
		log.Panic(err)
	}

	input, err := os.ReadFile(args.inputFileLocation)
	if err != nil {
		log.Panic(err)
	}

	definitions, err := parser.LexAndParse(string(input))
	if err != nil {
		log.Panic(err)
	}

	output := generator.GenerateOutput(definitions, generator.TypescriptOut)
	err = os.WriteFile(args.outputFileLocation, []byte(output), 0644)
	if err != nil {
		log.Panic(err)
	}
}
