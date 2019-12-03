package main

import "errors"

const usage = `"ct-template" is a command-line tool for ct(container-linux-config-transpiler).
Usage:
    ct-template <template dir> [-p <param file>]

Options:
    -p <param file>, --param <param file>

Other options:
    -?, -h, --help      display this help and exit
    -v, --version       output program version and exit

GitHub repository URL: https://github.com/masa213f/ct-template
`

type cmdlineOption struct {
	tmplDir     string
	paramFile   string
	showUsage   bool
	showVersion bool
}

func defaultOption() *cmdlineOption {
	return &cmdlineOption{
		tmplDir:     "",
		paramFile:   "",
		showUsage:   false,
		showVersion: false,
	}
}

func parseOptions(args []string) (*cmdlineOption, error) {
	opt := defaultOption()
	argsLen := len(args)
	for i := 0; i < argsLen; i++ {
		o := args[i]
		switch o {
		case "-p", "--param":
			i++
			if i >= argsLen {
				return nil, errors.New("opt error")
			}
			opt.paramFile = args[i]
		case "-?", "-h", "--help":
			opt.showUsage = true
		case "-v", "--version":
			opt.showVersion = true
		default:
			opt.tmplDir = args[i]
		}
	}

	if opt.showUsage || opt.showVersion {
		return opt, nil
	}

	if len(opt.tmplDir) == 0 {
		opt.tmplDir = "."
	}

	return opt, nil
}
