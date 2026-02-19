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
			contains: []string{"Language: Zenth", "Counter: 2", "Sum 1..10: 45"},
		},
		{
			file: "objects.zn",
			contains: []string{
				"Origin:", "Corner:", "Area: 50",
				"a1: 0 2.5",
				"a2: 100 2.5",
				"a3: 1000 5",
			},
		},
		{
			file:     "match.zn",
			contains: []string{"Fizz", "Buzz", "FizzBuzz", "Wednesday"},
		},
		{
			file: "range.zn",
			contains: []string{
				"range(0,5): 5",
				"0 1 2 3 4",
				"rangei(1,5): 5",
				"1 2 3 4 5",
				"range(0,10,2): 5",
				"0 2 4 6 8",
				"rangei(0,10,3): 4",
				"0 3 6 9",
				"i=0 v=10",
				"i=1 v=11",
				"i=2 v=12",
			},
		},
		{
			file: "interpolation.zn",
			contains: []string{
				"Hello World!",
				"3 + 4 = 7",
				"No interpolation here",
				"Raw {name} string",
				"Use {braces} literally",
				"Result: 30",
				"Area: 50",
				"Say: hi",
			},
		},
		{
			file:     "ifexpr.zn",
			contains: []string{"big", "not huge", "A", "42"},
		},
		{
			file:     "defaults.zn",
			contains: []string{"Hello, Alice!", "Hi, Bob!", "3", "6"},
		},
		{
			file:     "multiassign.zn",
			contains: []string{"2 1", "5"},
		},
		{
			file:     "pop.zn",
			contains: []string{"last=6", "four=4", "len=5"},
		},
		{
			file:     "addpush.zn",
			contains: []string{"after add: 4", "last=5", "after push: 5", "first=1", "second=2"},
		},
		{
			file:     "length.zn",
			contains: []string{"len=2", "empty=0", "after=3"},
		},
		{
			file:     "in_array.zn",
			contains: []string{"found 20", "no 99", "found bob", "no dave", "found true"},
		},
		{
			file:     "exit.zn",
			contains: []string{"before exit"},
		},
		{
			file:     "loop.zn",
			contains: []string{"hello", "count=5"},
		},
		{
			file:     "file.zn",
			contains: []string{"exists=true", "name=file.zn", "ext=.zn", "read=ok", "lines=ok", "nope=false"},
		},
		{
			file: "conversions.zn",
			contains: []string{
				"int_str=42",
				"int_f64=3",
				"int_true=1",
				"int_false=0",
				"f64_str=3.14",
				"f64_int=42",
				"to_int_len=3",
				"to_int_first=1",
				"to_int_last=3",
				"to_f64_len=3",
				"to_f64_first=1.5",
				"to_str=10,20,30",
			},
		},
		{
			file:     "named_args_func.zn",
			contains: []string{"1,2,3", "5,6,7"},
		},
		{
			file:     "method_defaults.zn",
			contains: []string{"Hello, world", "Hello, Alice"},
		},
		{
			file:     "split.zn",
			contains: []string{"n1=3", "this|a|dog", "n2=4", "my|dog|has|fleas"},
		},
		{
			file:     "string_index.zn",
			contains: []string{"ch=b", "u1=é"},
		},
		{
			file:     "string_length.zn",
			contains: []string{"m=5", "f=5"},
		},
		{
			file:     "debug_print.zn",
			contains: []string{"This is debug code", "This is standard"},
		},
		{
			file:     "object_array.zn",
			contains: []string{"len=1", "x=1 y=2", "cleared=0"},
		},
		{
			file:     "object_print.zn",
			contains: []string{"Point(x=1, y=2, val=\"#\")", "Label<demo>"},
		},
		{
			file:     "map.zn",
			contains: []string{"val=#", "len=1"},
		},
		{
			file:     "flag.zn",
			contains: []string{"debug=false", "times=5", "msg=Hello"},
		},
		{
			file:     "str_contains.zn",
			contains: []string{"found world", "no xyz"},
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
		file   string
		errMsg string
	}{
		{
			file:   "errors/type_mismatch.zn",
			errMsg: "type mismatch",
		},
		{
			file:   "errors/mutability.zn",
			errMsg: "cannot assign to immutable",
		},
		{
			file:   "errors/pop_too_many_args.zn",
			errMsg: "pop() takes 0 or 1 arguments",
		},
		{
			file:   "errors/arg_type_mismatch.zn",
			errMsg: "argument 1 to takes_int has type str, expected int",
		},
		{
			file:   "errors/undefined_identifier.zn",
			errMsg: "undefined identifier: y",
		},
		{
			file:   "errors/unknown_method.zn",
			errMsg: "type int has no method 'foo'",
		},
		{
			file:   "errors/split_bad_arg.zn",
			errMsg: "split() separator must be str, got int",
		},
		{
			file:   "errors/print_bad_flag.zn",
			errMsg: "second argument to println must be bool, got int",
		},
		{
			file:   "errors/print_too_many_args.zn",
			errMsg: "println() expects 1 or 2 arguments, got 3",
		},
		{
			file:   "errors/for_range_count_type.zn",
			errMsg: "for-range count must be int",
		},
		{
			file:   "errors/map_type_arg.zn",
			errMsg: "map() type arguments must be type names",
		},
		{
			file:   "errors/string_length_args.zn",
			errMsg: "length() takes no arguments, got 1",
		},
		{
			file:   "errors/flag_no_default.zn",
			errMsg: "flag() requires a 'default' argument",
		},
		{
			file:   "errors/in_not_operator.zn",
			errMsg: "expected {, got in",
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
