package main

import (
	"flag"
	"fmt"
)

type flags struct {
	inputFileLocation  string
	outputFileLocation string
	goPackageName      string
}

func newFlags() (flags, error) {
	f := flags{}

	flag.StringVar(&f.inputFileLocation, "input", "", "path to the input file: e.g. ./input.bt")
	flag.StringVar(&f.outputFileLocation, "output", "", "path to the output file: e.g. ./output.go")
	flag.StringVar(&f.goPackageName, "go-package-name", "", "used when output is a golang file, specifies the package name for the generated go file, defaults to the name of the output file, e.g. mypackage.go will be mypackage")
	flag.Parse()

	if f.outputFileLocation == "" {
		return f, fmt.Errorf("--output is required")
	}
	if f.inputFileLocation == "" {
		return f, fmt.Errorf("--input is required")
	}
	return f, nil
}
