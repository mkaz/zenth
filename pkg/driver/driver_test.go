package driver

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildAndRun(t *testing.T) {
	tests := []struct {
		file     string
		contains []string // expected substrings in output
	}{
		{
			file:     "hello.zn",
			contains: []string{"Hello, World!"},
		},
		{
			file:     "fibonacci.zn",
			contains: []string{"0", "1", "1", "2", "3", "5", "8", "13"},
		},
		{
			file:     "variables.zn",
			contains: []string{"Language: Zenth", "Counter: 2", "Sum 1..10: 55"},
		},
		{
			file:     "structs.zn",
			contains: []string{"Origin:", "Corner:", "Area: 50"},
		},
		{
			file:     "match.zn",
			contains: []string{"Fizz", "Buzz", "FizzBuzz", "Wednesday"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			// Find testdata relative to this test file
			input := filepath.Join("..", "..", "testdata", tt.file)
			if _, err := os.Stat(input); err != nil {
				t.Skipf("testdata not found: %s", input)
			}

			// Build to temp file
			tmpFile, err := os.CreateTemp("", "zenth-test-*")
			if err != nil {
				t.Fatal(err)
			}
			tmpFile.Close()
			defer os.Remove(tmpFile.Name())

			err = Build(Options{Input: input, Output: tmpFile.Name()})
			if err != nil {
				t.Fatalf("build failed: %v", err)
			}

			// Run the binary
			out, err := exec.Command(tmpFile.Name()).CombinedOutput()
			if err != nil {
				t.Fatalf("run failed: %v\noutput: %s", err, out)
			}

			output := string(out)
			for _, want := range tt.contains {
				if !strings.Contains(output, want) {
					t.Errorf("output missing %q\nfull output:\n%s", want, output)
				}
			}
		})
	}
}

func TestCompileErrors(t *testing.T) {
	tests := []struct {
		file     string
		errMsg   string
	}{
		{
			file:   "errors/type_mismatch.zn",
			errMsg: "type mismatch",
		},
		{
			file:   "errors/mutability.zn",
			errMsg: "cannot assign to immutable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			input := filepath.Join("..", "..", "testdata", tt.file)
			if _, err := os.Stat(input); err != nil {
				t.Skipf("testdata not found: %s", input)
			}

			err := Build(Options{Input: input, Output: "/dev/null"})
			if err == nil {
				t.Fatal("expected compile error, got none")
			}
			if !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("expected error containing %q, got: %v", tt.errMsg, err)
			}
		})
	}
}
