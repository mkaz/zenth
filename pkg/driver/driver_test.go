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
		cleanup  []string // files to remove after test
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
			file:     "matchexpr.zn",
			contains: []string{"two", "default", "Monday", "Wednesday", "Other", "2"},
		},
		{
			file: "range.zn",
			contains: []string{
				"Range(0,5): 5",
				"0 1 2 3 4",
				"Rangei(1,5): 5",
				"1 2 3 4 5",
				"Range(0,10,2): 5",
				"0 2 4 6 8",
				"Rangei(0,10,3): 4",
				"0 3 6 9",
				"i=0 v=10",
				"i=1 v=11",
				"i=2 v=12",
			},
		},
		{
			file: "range_contains.zn",
			contains: []string{
				"true",
				"false",
				"42 in range",
				"point in bounds",
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
			file:     "array_insert_remove.zn",
			contains: []string{"after insert: 5", "at[2]=3", "at[3]=4", "after remove: 4", "at[0]=1", "at[1]=3"},
		},
		{
			file:     "str_is_digit.zn",
			contains: []string{"5=true", "0=true", "9=true", "a=false", "[=false", ",=false", "digits=3"},
		},
		{
			file:     "grouped_expr.zn",
			contains: []string{"a=3", "b=24", "c=5"},
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
			file: "file_sections.zn",
			contains: []string{
				"count=3",
				"first=first section line 1\nfirst section line 2",
				"second=second section line 1",
				"third=third section line 1\nthird section line 2\nthird section line 3",
			},
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
				"to_float_len=3",
				"to_float_first=1.5",
				"to_str=10,20,30",
			},
		},
		{
			file: "base_conversion.zn",
			contains: []string{
				"hex=255",
				"bin=10",
				"oct=63",
				"dec=42",
				"to_hex=ff",
				"to_bin=1010",
				"to_oct=77",
				"to_dec=42",
				"rt=deadbeef",
			},
		},
		{
			file: "numeric_funcs.zn",
			contains: []string{
				"abs_i=5",
				"abs_f=2.5",
				"min_i=3",
				"max_i=7",
				"clamp_i=10",
				"min_f=1.5",
				"max_f=2.5",
				"clamp_f=10",
				"round=4",
				"floor=3",
				"ceil=4",
				"pow=8",
				"sqrt=3",
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
			file:     "hashmap.zn",
			contains: []string{"val=#", "len=1"},
		},
		{
			file:     "hashmap_default.zn",
			contains: []string{"foo=3", "bar=0", "baz=42", "len=2", "x=2"},
		},
		{
			file:     "hashmap_iter.zn",
			contains: []string{"sum=6", "count=3", "nkeys=3", "vsum=6"},
		},
		{
			file:     "env.zn",
			contains: []string{"HOME=", "missing=fallback", "empty="},
		},
		{
			file:     "flag.zn",
			contains: []string{"debug=false", "times=5", "msg=Hello", "args=0"},
		},
		{
			file:     "str_contains.zn",
			contains: []string{"found world", "no xyz"},
		},
		{
			file:     "strslice.zn",
			contains: []string{"bcde", "ab", "bc", "3", "30", "3", "30", "3", "20"},
		},
		{
			file:     "string_methods.zn",
			contains: []string{"u=HELLO", "l=hello", "sw=true", "ew=true", "strip=--hi--", "stripc=--hi--", "find=2", "count=2", "repall=x-b-x", "repn=x-b-a"},
		},
		{
			file:     "tuple.zn",
			contains: []string{"a", "1", "1-x-true", "left:right", "forward:10"},
		},
		{
			file:     "named_tuple.zn",
			contains: []string{"Alice", "30", "x", "42", "Bob:95", "Carol:88"},
		},
		{
			file:     "commands_tuple.zn",
			contains: []string{"Found: 3", "forward=1"},
		},
		{
			file:     "nested_generics.zn",
			contains: []string{"2", "b", "c", "100"},
		},
		{
			file: "set.zn",
			contains: []string{
				"3",
				"found start",
				"unknown not found",
				"6",
			},
		},
		{
			file: "closures.zn",
			contains: []string{
				"doubled: 2,4,6,8,10",
				"evens: 2,4",
				"chained: 30,40,50",
				"lengths: 5,3,7",
				"processed: 11,21,31",
				"typed: 3,6,9",
				"big: 20,15",
			},
		},
		{
			file: "strip_prefix.zn",
			contains: []string{
				"spec=x=5",
				"unchanged=hello",
				"trimmed=photo",
				"same=photo.png",
				"dots=.....",
				"row=ababab",
				"empty=",
			},
		},
		{
			file: "array_maxmin.zn",
			contains: []string{
				"max=9",
				"min=1",
				"single_max=42",
				"single_min=42",
				"fmax=3.14",
				"fmin=1.41",
				"neg_max=-1",
				"neg_min=-10",
			},
		},
		{
			file: "set_tuple.zn",
			contains: []string{
				"size=4",
				"after_dup=4",
				"found 6,10",
				"1,1 not found",
				"after_remove=3",
				"length=3",
				"pairs=2",
			},
		},
		{
			file: "discard.zn",
			contains: []string{
				"count=5",
				"abc",
				"012",
			},
		},
		{
			file: "hashmap_exists.zn",
			contains: []string{
				"alice exists",
				"charlie missing",
				"charlie added",
				"has0=true",
				"has1=false",
			},
		},
		{
			file: "array_sum_sorted.zn",
			contains: []string{
				"sum=15",
				"fsum=7",
				"esum=0",
				"orig=3",
				"sorted=1,1,2,3,4,5,6,9",
				"words=apple,banana,cherry",
				"negs=-5,-1,0,3,7",
				"desc=9,6,5,4,3,2,1,1",
				"asc=1,1,2,3,4,5,6,9",
				"wdesc=cherry,banana,apple",
			},
		},
		{
			file: "reduce.zn",
			contains: []string{
				"15",
				"120",
				"5",
				"115",
				"7",
			},
		},
		{
			file: "int_base.zn",
			contains: []string{
				"10",
				"255",
				"63",
				"42",
			},
		},
		{
			file:     "assert.zn",
			contains: []string{"all assertions passed"},
		},
		{
			file:     "array_destruct.zn",
			contains: []string{"10", "20", "30", "100", "200", "1", "3", "99", "6", "4", "12"},
		},
		{
			file:     "variadic_minmax.zn",
			contains: []string{"3", "7", "2", "8", "1", "9", "0.5", "2.5"},
		},
		{
			file:     "multi_return.zn",
			contains: []string{"1 9", "3 2", "world hello"},
		},
		{
			file: "enum.zn",
			contains: []string{
				"Red",
				"Green",
				"Blue",
				"is red",
				"not blue",
				"matched green",
				"#FF0000",
				"#00FF00",
			},
		},
		{
			file: "int_limits.zn",
			contains: []string{
				"9223372036854775807",
				"-9223372036854775808",
				"9223372036854775806",
				"INT_MAX is positive",
				"INT_MIN is negative",
			},
		},
		{
			file:     "array_join.zn",
			contains: []string{"hello world", "a,b,c", "xyz", "solo", "2026-03-15"},
		},
		{
			file:     "hashmap_tuple_key.zn",
			contains: []string{"origin", "point", "true", "false", "3", "yes", "1", "2"},
		},
		{
			file:     "array_extend.zn",
			contains: []string{"6", "1 2 3 4 5 6", "2", "hello world !"},
		},
		{
			file:     "zip.zn",
			contains: []string{"Alice=95", "Bob=87", "Charlie=92", "1+10", "2+20", "1,4,7", "2,5,8", "3,6,9"},
		},
		{
			file:     "array_repeat.zn",
			contains: []string{"5", "0 0 0 0 0", "---", "1 2 1 2 1 2", "3"},
		},
		{
			file:     "modules_basic.zn",
			contains: []string{"7", "42", "20", "16"},
		},
		{
			file:     "modules_alias.zn",
			contains: []string{"300", "99"},
		},
		{
			file:     "modules_dir.zn",
			contains: []string{"15", "16", "50."},
		},
		{
			file:     "modules_qualified_type.zn",
			contains: []string{"apple: 3", "banana: 5", "2"},
		},
		{
			file: "func_type.zn",
			contains: []string{
				"apply double: 10",
				"apply closure: 15",
				"even count: 3",
				"combine add: 7",
				"var call: 14",
				"map double: 2 4 6",
				"filter even: 2 4 6",
				"reduce add: 21",
				"each: 10",
				"each: 20",
			},
		},
		{
			file: "array_get.zn",
			contains: []string{
				"get0: a", "get2: c", "get5: default", "neg: none",
				"num0: 10", "num9: -1", "chain: ''", "field3: missing",
			},
		},
		{
			file: "multiline_string.zn",
			contains: []string{
				"Hello\nWorld",
				"Welcome to Zenth!\nEnjoy coding.",
				"line one\n\nline three",
				"Usage:\n  run",
			},
		},
		{
			file: "string_pad.zn",
			contains: []string{
				"'hello     '",
				"'     hello'",
				"00042",
				"hi....",
				"toolong",
				"apples      1.50",
				"bread       2.99",
				"milk        3.25",
			},
		},
		{
			file: "enumerate.zn",
			contains: []string{
				"0: apple", "1: banana", "2: cherry",
				"idx=2 val=30", "idx=3 val=40", "idx=4 val=50",
				"0:a", "1:b", "2:c",
				"pos=2 val=25", "pos=3 val=35", "pos=4 val=45",
			},
		},
		{
			file: "tuple_destruct_for.zn",
			contains: []string{
				"Alice: 95", "Bob: 87", "Charlie: 92",
				"0=apple", "1=banana", "2=cherry",
				"x", "y", "z",
				"a 1 true", "b 2 false",
			},
		},
		{
			file:     "file_write.zn",
			contains: []string{"write=hello\n", "append=hello\nworld\n", "overwrite=replaced\n"},
			cleanup:  []string{"_write_test.txt"},
		},
		{
			file: "date.zn",
			contains: []string{
				"today:",
				"epoch: 1970-01-01",
				"jan13: 2007-01-13",
				"tomorrow: 1970-01-02",
				"next_month: 1970-02-01",
				"next_year: 1971-01-01",
				"yesterday: 1970-01-01",
				"prev_month: 1970-01-01",
				"chained: 2001-09-25",
				"formatted: 13/01/2007",
				"str:",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			for _, f := range tt.cleanup {
				defer os.Remove(f)
			}
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
			errMsg: "argument 1 to takes_int has type Str, expected Int",
		},
		{
			file:   "errors/undefined_identifier.zn",
			errMsg: "undefined identifier: y",
		},
		{
			file:   "errors/unknown_method.zn",
			errMsg: "type Int has no method 'foo'",
		},
		{
			file:   "errors/split_bad_arg.zn",
			errMsg: "split() separator must be Str, got Int",
		},
		{
			file:   "errors/print_bad_flag.zn",
			errMsg: "second argument to Println must be Bool, got Int",
		},
		{
			file:   "errors/print_too_many_args.zn",
			errMsg: "Println() expects 1 or 2 arguments, got 3",
		},
		{
			file:   "errors/for_range_count_type.zn",
			errMsg: "for-range count must be Int",
		},
		{
			file:   "errors/hashmap_type_arg.zn",
			errMsg: "type arguments must be type names",
		},
		{
			file:   "errors/string_length_args.zn",
			errMsg: "length() takes no arguments, got 1",
		},
		{
			file:   "errors/string_replace_bad_count.zn",
			errMsg: "replace() argument 3 must be Int",
		},
		{
			file:   "errors/numeric_min_type_mismatch.zn",
			errMsg: "Min() arguments must all be Int or all be Float",
		},
		{
			file:   "errors/env_bad_arg.zn",
			errMsg: "expected Str",
		},
		{
			file:   "errors/flag_no_default.zn",
			errMsg: "Args.flag() requires a 'default' argument",
		},
		{
			file:   "errors/in_not_operator.zn",
			errMsg: "expected {, got in",
		},
		{
			file:   "errors/closure_filter_bool.zn",
			errMsg: "filter() closure must return Bool",
		},
		{
			file:   "errors/tuple_field_name.zn",
			errMsg: "tuple field must be numeric index",
		},
		{
			file:   "errors/tuple_type_mismatch.zn",
			errMsg: "type mismatch",
		},
		{
			file:   "errors/tuple_destructure_non_tuple.zn",
			errMsg: "tuple destructuring requires tuple value",
		},
		{
			file:   "errors/tuple_destructure_arity.zn",
			errMsg: "tuple destructuring arity mismatch",
		},
		{
			file:   "errors/named_tuple_mixed.zn",
			errMsg: "cannot mix named and positional tuple fields",
		},
		{
			file:   "errors/named_tuple_dup_field.zn",
			errMsg: "duplicate field name",
		},
		{
			file:   "errors/named_tuple_bad_field.zn",
			errMsg: "named tuple has no field",
		},
		{
			file:   "errors/named_tuple_type_mismatch.zn",
			errMsg: "type mismatch",
		},
		{
			file:   "errors/split_once_bad_arg.zn",
			errMsg: "split_once() separator must be Str",
		},
		{
			file:   "errors/to_int_bad_base.zn",
			errMsg: "to_int() base must be Int, got Str",
		},
		{
			file:   "errors/to_base_bad_arg.zn",
			errMsg: "to_base() argument must be Int, got Str",
		},
		{
			file:   "errors/hashmap_default_type.zn",
			errMsg: "hashmap default type mismatch",
		},
		{
			file:   "errors/array_max_non_numeric.zn",
			errMsg: "max() requires a numeric array",
		},
		{
			file:   "errors/repeat_bad_arg.zn",
			errMsg: "repeat() argument must be Int",
		},
		{
			file:   "errors/sum_non_numeric.zn",
			errMsg: "sum() requires a numeric array",
		},
		{
			file:   "errors/sorted_non_sortable.zn",
			errMsg: "sorted() requires a numeric or string array",
		},
		{
			file:   "errors/sorted_bad_arg.zn",
			errMsg: "sorted() argument must be \"asc\" or \"desc\"",
		},
		{
			file:   "errors/multi_return_type_mismatch.zn",
			errMsg: "return type mismatch",
		},
		{
			file:   "errors/assign_int_max.zn",
			errMsg: "cannot assign to constant",
		},
		{
			file:   "errors/enum_bad_variant.zn",
			errMsg: "has no variant",
		},
		{
			file:   "errors/reduce_not_closure.zn",
			errMsg: "reduce() first argument must be a closure",
		},
		{
			file:   "errors/reduce_wrong_params.zn",
			errMsg: "reduce() closure must take exactly 2 parameters",
		},
		{
			file:   "errors/int_base_not_str.zn",
			errMsg: "Int() with base requires first argument to be Str",
		},
		{
			file:   "errors/array_destruct_non_array.zn",
			errMsg: "array destructuring requires array value",
		},
		{
			file:   "errors/array_join_non_string.zn",
			errMsg: "join() requires Array(Str)",
		},
		{
			file:   "errors/modules_no_main.zn",
			errMsg: "cannot define fn main",
		},
		{
			file:   "errors/modules_unknown_fn.zn",
			errMsg: "has no function or type 'bogus'",
		},
		{
			file:   "errors/file_write_no_arg.zn",
			errMsg: "write() takes exactly 1 argument",
		},
		{
			file:   "errors/file_write_bad_type.zn",
			errMsg: "write() argument must be Str",
		},
		{
			file:   "errors/file_append_no_arg.zn",
			errMsg: "append() takes exactly 1 argument",
		},
		{
			file:   "errors/func_type_mismatch.zn",
			errMsg: "expected Fn(...) -> Int",
		},
		{
			file:   "errors/get_bad_args.zn",
			errMsg: "get() index must be Int, got Str",
		},
		{
			file:   "errors/pad_bad_args.zn",
			errMsg: "pad_left() width must be Int, got Str",
		},
		{
			file:   "errors/date_bad_format_arg.zn",
			errMsg: "format() argument must be Str, got Int",
		},
		{
			file:   "errors/date_bad_add_arg.zn",
			errMsg: "add() value must be Int, got Str",
		},
		{
			file:   "errors/date_bad_method.zn",
			errMsg: "Date has no method 'invalid'",
		},
		{
			file:   "errors/date_today_args.zn",
			errMsg: "Date.today() takes no arguments",
		},
		{
			file:   "errors/import_missing_prefix.zn",
			errMsg: "missing './' prefix",
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
