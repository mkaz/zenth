package main

import (
	"fmt"
	"os"

	"github.com/mkaz/zenth/pkg/driver"
)

const version = "0.7.4"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "build":
		handleBuild()
	case "run":
		handleRun()
	case "test":
		handleTest()
	case "version":
		fmt.Printf("zenth %s\n", version)
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func handleBuild() {
	opts := driver.Options{}

	args := os.Args[2:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-o":
			if i+1 >= len(args) {
				fatal("-o requires an argument")
			}
			i++
			opts.Output = args[i]
		case "--emit-go":
			opts.EmitGo = true
		case "-v", "--verbose":
			opts.Verbose = true
		default:
			if args[i][0] == '-' {
				fatal("unknown flag: %s", args[i])
			}
			opts.Input = args[i]
		}
	}

	if opts.Input == "" {
		fatal("no input file specified")
	}

	if err := driver.Build(opts); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func handleRun() {
	// Build to temp location, then run
	opts := driver.Options{}

	args := os.Args[2:]
	var remainingArgs []string
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-v", "--verbose":
			opts.Verbose = true
		default:
			if opts.Input == "" {
				if args[i][0] == '-' {
					fatal("unknown flag: %s", args[i])
				}
				opts.Input = args[i]
			} else {
				// Everything after the .zn file is passed to the program
				remainingArgs = args[i:]
				i = len(args) // break out of loop
			}
		}
	}

	if opts.Input == "" {
		fatal("no input file specified")
	}

	// Build to temp file
	tmpFile, err := os.CreateTemp("", "zenth-run-*")
	if err != nil {
		fatal("cannot create temp file: %v", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	opts.Output = tmpFile.Name()
	if err := driver.Build(opts); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// Run the binary, forwarding remaining args
	cmd := execCommand(tmpFile.Name(), remainingArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}
}

func handleTest() {
	var verbose bool
	path := "tests" // default test directory

	args := os.Args[2:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-v", "--verbose":
			verbose = true
		default:
			if args[i][0] == '-' {
				fatal("unknown flag: %s", args[i])
			}
			path = args[i]
		}
	}

	results, err := driver.RunTests(path, verbose)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	output, allPassed := driver.FormatResults(results)
	fmt.Print(output)
	if !allPassed {
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`zenth - The Zenth Programming Language

Usage:
  zenth build [flags] <file.zn>    Compile a Zenth program to a binary
  zenth run [flags] <file.zn>      Compile and run a Zenth program
  zenth test [path]                Run tests (default: tests/ directory)
  zenth version                    Print version
  zenth help                       Show this help

Build flags:
  -o <name>       Output binary name (default: input filename without .zn)
  --emit-go       Print generated Go source instead of compiling
  -v, --verbose   Verbose output

Test flags:
  -v, --verbose   Verbose output`)
}

func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
	os.Exit(1)
}
