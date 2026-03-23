package driver

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mkaz/zenth/pkg/ast"
	"github.com/mkaz/zenth/pkg/checker"
	"github.com/mkaz/zenth/pkg/codegen"
	"github.com/mkaz/zenth/pkg/lexer"
	"github.com/mkaz/zenth/pkg/parser"
)

// Options configures the build.
type Options struct {
	Input  string // input .zn file
	Output string // output binary name
	EmitGo bool   // if true, print generated Go and stop
	Verbose bool
}

// moduleEntry describes a local Zenth module to build as a Go package.
type moduleEntry struct {
	name      string   // Go package name (e.g., "utils")
	goRelPath string   // relative path in temp dir (e.g., "utils" or "geo/vector")
	znFiles   []string // source .zn file paths
}

// Build compiles a Zenth source file to a native binary.
func Build(opts Options) error {
	sourceDir := filepath.Dir(opts.Input)

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

	// Resolve local imports: set GoPackagePath and collect module entries
	// Also collect external Go module paths for import_go statements
	var modules []moduleEntry
	var goExternalPaths []string
	for _, stmt := range prog.Stmts {
		imp, ok := stmt.(*ast.ImportDecl)
		if !ok {
			continue
		}
		if imp.IsGoExternal {
			goExternalPaths = append(goExternalPaths, imp.Path)
			continue
		}
		if !imp.IsLocal {
			continue
		}
		rel := strings.TrimPrefix(imp.Path, "./")
		rel = strings.TrimPrefix(rel, "../")
		goRelPath := filepath.ToSlash(rel)
		imp.GoPackagePath = "zenth_output/" + goRelPath

		modName := imp.Alias
		if modName == "" {
			modName = filepath.Base(rel)
		}

		znFiles, err := checker.ResolveModuleFiles(sourceDir, imp.Path)
		if err != nil {
			return fmt.Errorf("cannot resolve module %s: %w", imp.Path, err)
		}
		modules = append(modules, moduleEntry{
			name:      modName,
			goRelPath: goRelPath,
			znFiles:   znFiles,
		})
	}

	// Type check with source dir so local imports can be resolved
	ch := checker.NewWithDir(sourceDir)
	if err := ch.Check(prog); err != nil {
		return fmt.Errorf("%w", err)
	}

	if opts.Verbose {
		fmt.Fprintln(os.Stderr, "type check passed")
	}

	// Generate Go code for the main program
	gen := codegen.New()
	goSrc := gen.Generate(prog)

	// Collect external deps needed by codegen helpers
	goExternalPaths = append(goExternalPaths, gen.ExternalDeps()...)

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
	goMod := "module zenth_output\n\ngo 1.26\n"
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		return fmt.Errorf("cannot write go.mod: %w", err)
	}

	// Fetch external Go modules declared via import_go
	for _, modPath := range goExternalPaths {
		if opts.Verbose {
			fmt.Fprintf(os.Stderr, "go get %s\n", modPath)
		}
		getCmd := exec.Command("go", "get", modPath)
		getCmd.Dir = tmpDir
		getCmd.Stdout = os.Stdout
		getCmd.Stderr = os.Stderr
		if err := getCmd.Run(); err != nil {
			return fmt.Errorf("go get %s failed: %w", modPath, err)
		}
	}

	// Write main.go
	if err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(goSrc), 0644); err != nil {
		return fmt.Errorf("cannot write main.go: %w", err)
	}

	if opts.Verbose {
		fmt.Fprintln(os.Stderr, "generated Go source:")
		fmt.Fprintln(os.Stderr, goSrc)
	}

	// Build and write each local module as a Go package
	for _, mod := range modules {
		modGoSrc, err := buildModuleGoSrc(mod, sourceDir)
		if err != nil {
			return fmt.Errorf("building module %s: %w", mod.name, err)
		}
		modDir := filepath.Join(tmpDir, filepath.FromSlash(mod.goRelPath))
		if err := os.MkdirAll(modDir, 0755); err != nil {
			return fmt.Errorf("cannot create module dir: %w", err)
		}
		modFile := filepath.Join(modDir, mod.name+".go")
		if err := os.WriteFile(modFile, []byte(modGoSrc), 0644); err != nil {
			return fmt.Errorf("cannot write module file: %w", err)
		}
		if opts.Verbose {
			fmt.Fprintf(os.Stderr, "generated module: %s\n", modFile)
		}
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

// buildModuleGoSrc parses, type-checks, and generates Go source for a local module.
func buildModuleGoSrc(mod moduleEntry, sourceDir string) (string, error) {
	// Combine all .zn files in the module into one program
	combined := &ast.Program{}
	for _, znFile := range mod.znFiles {
		src, err := os.ReadFile(znFile)
		if err != nil {
			return "", fmt.Errorf("cannot read %s: %w", znFile, err)
		}
		l := lexer.New(znFile, string(src))
		tokens, err := l.Tokenize()
		if err != nil {
			return "", fmt.Errorf("lex error in %s: %w", znFile, err)
		}
		p := parser.New(tokens)
		prog, err := p.Parse()
		if err != nil {
			return "", fmt.Errorf("parse error in %s: %w", znFile, err)
		}
		combined.Stmts = append(combined.Stmts, prog.Stmts...)
	}

	// Type-check the module (sourceDir used for any nested local imports)
	ch := checker.NewWithDir(sourceDir)
	if err := ch.Check(combined); err != nil {
		return "", fmt.Errorf("type error in module %s: %w", mod.name, err)
	}

	// Generate Go source with the module's package name
	gen := codegen.New()
	gen.PackageName = mod.name
	return gen.Generate(combined), nil
}
