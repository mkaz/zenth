package driver

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mkaz/zenth/pkg/checker"
	"github.com/mkaz/zenth/pkg/codegen"
	"github.com/mkaz/zenth/pkg/lexer"
	"github.com/mkaz/zenth/pkg/parser"
)

// Options configures the build.
type Options struct {
	Input    string // input .zn file
	Output   string // output binary name
	EmitGo   bool   // if true, print generated Go and stop
	Verbose  bool
}

// Build compiles a Zenth source file to a native binary.
func Build(opts Options) error {
	// Read source
	src, err := os.ReadFile(opts.Input)
	if err != nil {
		return fmt.Errorf("cannot read %s: %w", opts.Input, err)
	}

	// Lex
	l := lexer.New(opts.Input, string(src))
	tokens, err := l.Tokenize()
	if err != nil {
		return fmt.Errorf("lex error: %w", err)
	}

	if opts.Verbose {
		fmt.Fprintf(os.Stderr, "lexed %d tokens\n", len(tokens))
	}

	// Parse
	p := parser.New(tokens)
	prog, err := p.Parse()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if opts.Verbose {
		fmt.Fprintf(os.Stderr, "parsed %d top-level declarations\n", len(prog.Stmts))
	}

	// Type check
	ch := checker.New()
	if err := ch.Check(prog); err != nil {
		return fmt.Errorf("%w", err)
	}

	if opts.Verbose {
		fmt.Fprintln(os.Stderr, "type check passed")
	}

	// Generate Go code
	gen := codegen.New()
	goSrc := gen.Generate(prog)

	if opts.EmitGo {
		fmt.Print(goSrc)
		return nil
	}

	// Write to temp directory and build
	tmpDir, err := os.MkdirTemp("", "zenth-build-*")
	if err != nil {
		return fmt.Errorf("cannot create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Write go.mod
	goMod := "module zenth_output\n\ngo 1.21\n"
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		return fmt.Errorf("cannot write go.mod: %w", err)
	}

	// Write main.go
	if err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(goSrc), 0644); err != nil {
		return fmt.Errorf("cannot write main.go: %w", err)
	}

	if opts.Verbose {
		fmt.Fprintln(os.Stderr, "generated Go source:")
		fmt.Fprintln(os.Stderr, goSrc)
	}

	// Determine output path
	output := opts.Output
	if output == "" {
		base := filepath.Base(opts.Input)
		output = strings.TrimSuffix(base, filepath.Ext(base))
	}
	absOutput, err := filepath.Abs(output)
	if err != nil {
		return fmt.Errorf("cannot resolve output path: %w", err)
	}

	// Run go build
	cmd := exec.Command("go", "build", "-o", absOutput, ".")
	cmd.Dir = tmpDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go build failed: %w", err)
	}

	if opts.Verbose {
		fmt.Fprintf(os.Stderr, "built: %s\n", absOutput)
	}

	return nil
}
