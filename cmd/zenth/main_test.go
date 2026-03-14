package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var zenthBin string

func TestMain(m *testing.M) {
	tmp, err := os.CreateTemp("", "zenth-test-*")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	tmp.Close()
	zenthBin = tmp.Name()

	cmd := exec.Command("go", "build", "-o", zenthBin, ".")
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "build failed: %s\n%s", err, out)
		os.Exit(1)
	}

	code := m.Run()
	os.Remove(zenthBin)
	os.Exit(code)
}

// runZenth runs the test binary with the given args and returns stdout, stderr, and any error.
func runZenth(args ...string) (stdout, stderr string, err error) {
	cmd := exec.Command(zenthBin, args...)
	var outBuf, errBuf strings.Builder
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err = cmd.Run()
	return outBuf.String(), errBuf.String(), err
}

// writeZenthFile creates a .zn file in dir with the given content and returns its path.
func writeZenthFile(t *testing.T, dir, content string) string {
	t.Helper()
	path := filepath.Join(dir, "test.zn")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestVersion(t *testing.T) {
	stdout, _, err := runZenth("version")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(stdout, "0.4.8") {
		t.Errorf("expected version output to contain '0.4.8', got: %q", stdout)
	}
}

func TestHelp(t *testing.T) {
	for _, arg := range []string{"help", "--help", "-h"} {
		t.Run(arg, func(t *testing.T) {
			stdout, _, err := runZenth(arg)
			if err != nil {
				t.Fatalf("unexpected error for %q: %v", arg, err)
			}
			if !strings.Contains(stdout, "Usage:") {
				t.Errorf("expected 'Usage:' in output for %q, got: %q", arg, stdout)
			}
			for _, sub := range []string{"build", "run", "test", "version"} {
				if !strings.Contains(stdout, sub) {
					t.Errorf("expected %q in help output for %q", sub, arg)
				}
			}
		})
	}
}

func TestNoArgs(t *testing.T) {
	stdout, _, err := runZenth()
	if err == nil {
		t.Fatal("expected non-zero exit code when no args provided")
	}
	// Usage is printed to stdout even on error
	if !strings.Contains(stdout, "Usage:") {
		t.Errorf("expected usage output, got stdout=%q", stdout)
	}
}

func TestUnknownCommand(t *testing.T) {
	_, stderr, err := runZenth("badcmd")
	if err == nil {
		t.Fatal("expected non-zero exit code for unknown command")
	}
	if !strings.Contains(stderr, "unknown command") {
		t.Errorf("expected 'unknown command' in stderr, got: %q", stderr)
	}
}

func TestBuildMissingInput(t *testing.T) {
	_, stderr, err := runZenth("build")
	if err == nil {
		t.Fatal("expected non-zero exit code when no input file")
	}
	if !strings.Contains(stderr, "no input file") {
		t.Errorf("expected 'no input file' in stderr, got: %q", stderr)
	}
}

func TestBuildUnknownFlag(t *testing.T) {
	_, stderr, err := runZenth("build", "--badopt")
	if err == nil {
		t.Fatal("expected non-zero exit code for unknown flag")
	}
	if !strings.Contains(stderr, "unknown flag") {
		t.Errorf("expected 'unknown flag' in stderr, got: %q", stderr)
	}
}

func TestBuildAndRunHello(t *testing.T) {
	dir := t.TempDir()
	src := writeZenthFile(t, dir, `fn main() { println("hello from test"); }`)

	stdout, stderr, err := runZenth("run", src)
	if err != nil {
		t.Fatalf("unexpected error: %v\nstderr: %s", err, stderr)
	}
	if got := strings.TrimSpace(stdout); got != "hello from test" {
		t.Errorf("expected 'hello from test', got: %q", got)
	}
}

func TestBuildOutputFlag(t *testing.T) {
	dir := t.TempDir()
	src := writeZenthFile(t, dir, `fn main() { println("output flag test"); }`)
	binPath := filepath.Join(dir, "mybin")

	_, stderr, err := runZenth("build", "-o", binPath, src)
	if err != nil {
		t.Fatalf("build failed: %v\nstderr: %s", err, stderr)
	}

	// Verify binary exists
	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		t.Fatal("expected output binary to exist")
	}

	// Run the binary and check output
	cmd := exec.Command(binPath)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("running built binary failed: %v", err)
	}
	if got := strings.TrimSpace(string(out)); got != "output flag test" {
		t.Errorf("expected 'output flag test', got: %q", got)
	}
}

func TestEmitGo(t *testing.T) {
	dir := t.TempDir()
	src := writeZenthFile(t, dir, `fn main() { println("emit go"); }`)

	stdout, stderr, err := runZenth("build", "--emit-go", src)
	if err != nil {
		t.Fatalf("unexpected error: %v\nstderr: %s", err, stderr)
	}
	if !strings.Contains(stdout, "package main") {
		t.Errorf("expected 'package main' in emitted Go source, got: %q", stdout)
	}
}

func TestRunWithArgs(t *testing.T) {
	dir := t.TempDir()
	src := writeZenthFile(t, dir, `fn main() {
    let msg = flag(default="default_val");
    println(msg);
}`)

	stdout, stderr, err := runZenth("run", src, "--msg", "custom_val")
	if err != nil {
		t.Fatalf("unexpected error: %v\nstderr: %s", err, stderr)
	}
	if got := strings.TrimSpace(stdout); got != "custom_val" {
		t.Errorf("expected 'custom_val', got: %q", got)
	}
}

func TestBuildInvalidSource(t *testing.T) {
	dir := t.TempDir()
	src := writeZenthFile(t, dir, `fn main() { let x: int = "not an int"; }`)

	_, _, err := runZenth("build", src)
	if err == nil {
		t.Fatal("expected non-zero exit code for invalid source")
	}
}
