package main

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

var (
	version string
)

func showUsage(o io.Writer) {
	fmt.Fprintf(o, usage)
}

func showVersion(o io.Writer) {
	fmt.Fprintln(o, version)
}

func main() {
	cmdOpt, err := parseOptions(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(2)
	}

	if cmdOpt.showVersion {
		showVersion(os.Stdout)
		os.Exit(0)
	}

	if cmdOpt.showUsage {
		showUsage(os.Stdout)
		os.Exit(0)
	}

	input, err := readTemplateDir(cmdOpt.tmplDir)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	output := convert(input)
	yamlOutput, err := yaml.Marshal(output)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println(string(yamlOutput))
}
