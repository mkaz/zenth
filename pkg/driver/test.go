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

// TestResult holds the result of running a single test file.
type TestResult struct {
	File  string
	Tests []SingleTest
	Error string // non-empty if the file failed to compile
}

// SingleTest holds the result of a single test function.
type SingleTest struct {
	Name   string
	Passed bool
	Error  string // failure message (empty if passed)
}

// RunTests discovers and runs test files, returning results.
func RunTests(path string, verbose bool) ([]TestResult, error) {
	files, err := discoverTestFiles(path)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no test files found in %s", path)
	}

	var results []TestResult
	for _, file := range files {
		result := runTestFile(file, verbose)
		results = append(results, result)
	}
	return results, nil
}

// discoverTestFiles finds all *_test.zn files in the given path.
// If path is a single file, returns just that file.
func discoverTestFiles(path string) ([]string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("cannot access %s: %w", path, err)
	}

	if !info.IsDir() {
		// Single file
		if !strings.HasSuffix(path, "_test.zn") {
			return nil, fmt.Errorf("test files must end with _test.zn: %s", path)
		}
		return []string{path}, nil
	}

	// Directory: find all *_test.zn files
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read directory %s: %w", path, err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), "_test.zn") {
			files = append(files, filepath.Join(path, entry.Name()))
		}
	}
	return files, nil
}

// findTestFunctions parses a Zenth source and returns the names of fn test_*() functions.
func findTestFunctions(filename, source string) ([]string, error) {
	l := lexer.New(filename, source)
	tokens, err := l.Tokenize()
	if err != nil {
		return nil, fmt.Errorf("lex error: %w", err)
	}

	p := parser.New(tokens)
	prog, err := p.Parse()
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	var testNames []string
	for _, stmt := range prog.Stmts {
		if fn, ok := stmt.(*ast.FnDecl); ok {
			if strings.HasPrefix(fn.Name, "test_") {
				testNames = append(testNames, fn.Name)
			}
		}
	}
	return testNames, nil
}

// runTestFile compiles and runs a single test file, returning results.
func runTestFile(file string, verbose bool) TestResult {
	result := TestResult{File: file}

	// Read source
	src, err := os.ReadFile(file)
	if err != nil {
		result.Error = fmt.Sprintf("cannot read %s: %v", file, err)
		return result
	}

	// Find test function names
	testNames, err := findTestFunctions(file, string(src))
	if err != nil {
		result.Error = fmt.Sprintf("parse error: %v", err)
		return result
	}

	if len(testNames) == 0 {
		result.Error = "no test functions found (functions must start with test_)"
		return result
	}

	// Compile through the normal pipeline
	l := lexer.New(file, string(src))
	tokens, err := l.Tokenize()
	if err != nil {
		result.Error = fmt.Sprintf("lex error: %v", err)
		return result
	}

	p := parser.New(tokens)
	prog, err := p.Parse()
	if err != nil {
		result.Error = fmt.Sprintf("%v", err)
		return result
	}

	sourceDir := filepath.Dir(file)

	// Resolve local imports: set GoPackagePath and collect module entries
	var modules []moduleEntry
	for _, stmt := range prog.Stmts {
		imp, ok := stmt.(*ast.ImportDecl)
		if !ok || !imp.IsLocal {
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
			result.Error = fmt.Sprintf("cannot resolve module %s: %v", imp.Path, err)
			return result
		}
		modules = append(modules, moduleEntry{
			name:      modName,
			goRelPath: goRelPath,
			znFiles:   znFiles,
		})
	}

	ch := checker.NewWithDir(sourceDir)
	if err := ch.Check(prog); err != nil {
		result.Error = fmt.Sprintf("%v", err)
		return result
	}

	gen := codegen.New()
	goSrc := gen.Generate(prog)

	// Inject a test harness main() into the generated Go source.
	// The generated code won't have a main() since the .zn file doesn't have one.
	goSrc = injectTestMain(goSrc, testNames)

	// Build to temp binary
	tmpDir, err := os.MkdirTemp("", "zenth-test-*")
	if err != nil {
		result.Error = fmt.Sprintf("cannot create temp dir: %v", err)
		return result
	}
	defer os.RemoveAll(tmpDir)

	goMod := "module zenth_output\n\ngo 1.26\n"
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		result.Error = fmt.Sprintf("cannot write go.mod: %v", err)
		return result
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(goSrc), 0644); err != nil {
		result.Error = fmt.Sprintf("cannot write main.go: %v", err)
		return result
	}

	// Build and write each local module as a Go package
	for _, mod := range modules {
		modGoSrc, err := buildModuleGoSrc(mod, sourceDir)
		if err != nil {
			result.Error = fmt.Sprintf("building module %s: %v", mod.name, err)
			return result
		}
		modDir := filepath.Join(tmpDir, filepath.FromSlash(mod.goRelPath))
		if err := os.MkdirAll(modDir, 0755); err != nil {
			result.Error = fmt.Sprintf("cannot create module dir: %v", err)
			return result
		}
		modFile := filepath.Join(modDir, mod.name+".go")
		if err := os.WriteFile(modFile, []byte(modGoSrc), 0644); err != nil {
			result.Error = fmt.Sprintf("cannot write module file: %v", err)
			return result
		}
	}

	tmpBin := filepath.Join(tmpDir, "testbin")
	cmd := exec.Command("go", "build", "-o", tmpBin, ".")
	cmd.Dir = tmpDir
	buildOut, err := cmd.CombinedOutput()
	if err != nil {
		result.Error = fmt.Sprintf("build failed: %s", string(buildOut))
		return result
	}

	// Run the test binary
	runCmd := exec.Command(tmpBin)
	out, _ := runCmd.CombinedOutput()
	// We don't check the error because the test binary exits 1 on failure

	// Parse structured output
	result.Tests = parseTestOutput(string(out), testNames)
	return result
}

// injectTestMain appends a Go func main() that calls each test function
// with panic recovery and prints structured PASS/FAIL output.
func injectTestMain(goSrc string, testNames []string) string {
	// Ensure "os" is imported (needed for os.Exit in the test harness).
	// The codegen may not have included it if the .zn file didn't use os.
	if !strings.Contains(goSrc, `"os"`) {
		goSrc = strings.Replace(goSrc, "import (\n", "import (\n\t\"os\"\n", 1)
	}

	var b strings.Builder
	b.WriteString(goSrc)
	b.WriteString("\n")
	b.WriteString("func main() {\n")
	b.WriteString("\tpassed := 0\n")
	b.WriteString("\tfailed := 0\n")

	for _, name := range testNames {
		// Each test is called in a closure with deferred recover.
		// Top-level functions in codegen keep their original name (no export capitalization).
		b.WriteString(fmt.Sprintf("\tfunc() {\n"))
		b.WriteString(fmt.Sprintf("\t\tdefer func() {\n"))
		b.WriteString(fmt.Sprintf("\t\t\tif r := recover(); r != nil {\n"))
		b.WriteString(fmt.Sprintf("\t\t\t\tfmt.Printf(\"FAIL %s %%v\\n\", r)\n", name))
		b.WriteString(fmt.Sprintf("\t\t\t\tfailed++\n"))
		b.WriteString(fmt.Sprintf("\t\t\t} else {\n"))
		b.WriteString(fmt.Sprintf("\t\t\t\tfmt.Println(\"PASS %s\")\n", name))
		b.WriteString(fmt.Sprintf("\t\t\t\tpassed++\n"))
		b.WriteString(fmt.Sprintf("\t\t\t}\n"))
		b.WriteString(fmt.Sprintf("\t\t}()\n"))
		b.WriteString(fmt.Sprintf("\t\t%s()\n", name))
		b.WriteString(fmt.Sprintf("\t}()\n"))
	}

	b.WriteString("\tfmt.Printf(\"%%d passed, %%d failed\\n\", passed, failed)\n")
	b.WriteString("\tif failed > 0 { os.Exit(1) }\n")
	b.WriteString("}\n")
	return b.String()
}

// parseTestOutput reads the structured output from the test binary.
func parseTestOutput(output string, testNames []string) []SingleTest {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	lineMap := make(map[string]string) // testname -> full line
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "PASS ") || strings.HasPrefix(line, "FAIL ") {
			parts := strings.SplitN(line, " ", 3)
			if len(parts) >= 2 {
				lineMap[parts[1]] = line
			}
		}
	}

	var tests []SingleTest
	for _, name := range testNames {
		line, ok := lineMap[name]
		if !ok {
			tests = append(tests, SingleTest{Name: name, Passed: false, Error: "no output (possible crash)"})
			continue
		}
		if strings.HasPrefix(line, "PASS") {
			tests = append(tests, SingleTest{Name: name, Passed: true})
		} else {
			msg := strings.TrimPrefix(line, "FAIL "+name+" ")
			tests = append(tests, SingleTest{Name: name, Passed: false, Error: msg})
		}
	}
	return tests
}

// FormatResults formats test results for display.
func FormatResults(results []TestResult) (output string, allPassed bool) {
	var b strings.Builder
	totalPassed := 0
	totalFailed := 0

	for _, r := range results {
		b.WriteString(r.File + "\n")
		if r.Error != "" {
			b.WriteString(fmt.Sprintf("  ERROR: %s\n", r.Error))
			totalFailed++
			continue
		}
		for _, t := range r.Tests {
			if t.Passed {
				b.WriteString(fmt.Sprintf("  %s ... PASS\n", t.Name))
				totalPassed++
			} else {
				b.WriteString(fmt.Sprintf("  %s ... FAIL\n", t.Name))
				if t.Error != "" {
					// Indent each line of the error
					for _, line := range strings.Split(t.Error, "\n") {
						b.WriteString(fmt.Sprintf("    %s\n", line))
					}
				}
				totalFailed++
			}
		}
	}

	b.WriteString(fmt.Sprintf("\n%d passed, %d failed\n", totalPassed, totalFailed))
	return b.String(), totalFailed == 0
}
