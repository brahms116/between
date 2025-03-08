package main

import (
	"flag"
	"fmt"
)

type flags struct {
	inputFileLocation  string
	outputFileLocation string
}

func newFlags() (flags, error) {
	f := flags{}

	flag.StringVar(&f.inputFileLocation, "input", "", "path to the input file")
	flag.StringVar(&f.outputFileLocation, "output", "", "path to the output file")
	flag.Parse()

	if f.outputFileLocation == "" {
		return f, fmt.Errorf("--output is required")
	}
	if f.inputFileLocation == "" {
		return f, fmt.Errorf("--input is required")
	}
	return f, nil
}
