package codegen

import (
	"fmt"
	"strings"

	"github.com/mkaz/zenth/pkg/ast"
	"github.com/mkaz/zenth/pkg/token"
)

// Generator translates a Zenth AST to Go source code.
type Generator struct {
	buf                strings.Builder
	indent             int
	imports            map[string]string // Go import path -> alias (or empty)
	objs               map[string]*ast.ObjDecl
	funcs              map[string]*ast.FnDecl // "name" or "StructName.methodName"
	typeAliases        map[string]*ast.TypeExpr
	needsRange         bool
	needsRangei        bool
	needsPop           bool
	needsAdd           bool
	needsPush          bool
	loopCounter        int
	needsContains      bool
	needsFile          bool
	needsIntConv       bool
	needsF64Conv       bool
	needsSliceToInt    bool
	needsSliceToF64    bool
	needsSliceToStr    bool
	needsStrIndex      bool
	needsStrSlice      bool
	needsHashmapObjKey bool
	needsMap           bool
	needsFilter        bool
	needsSplitOnce     bool
	needsAbsInt        bool
	needsMinInt        bool
	needsMaxInt        bool
	needsClampInt      bool
	needsClampF64      bool
	needsIntBase       bool
	needsToBase        bool
	needsMapget          bool
	needsHashmapKeys    bool
	needsHashmapValues  bool
	needsFlag          bool
	flagDecls          []flagDecl
	tempCounter        int
}

type flagDecl struct {
	name       string // variable name
	goType     string // "bool", "int", "string"
	defaultVal string // Go expression for default value
}

// New creates a new code Generator.
func New() *Generator {
	return &Generator{
		imports:     make(map[string]string),
		objs:        make(map[string]*ast.ObjDecl),
		funcs:       make(map[string]*ast.FnDecl),
		typeAliases: make(map[string]*ast.TypeExpr),
	}
}

// Generate produces Go source code from a Zenth AST.
func (g *Generator) Generate(prog *ast.Program) string {
	// First pass: collect objects, functions, and imports
	for _, stmt := range prog.Stmts {
		switch s := stmt.(type) {
		case *ast.ObjDecl:
			g.objs[s.Name] = s
			for _, m := range s.Methods {
				g.funcs[s.Name+"."+m.Name] = m
			}
		case *ast.FnDecl:
			g.funcs[s.Name] = s
		case *ast.ImportDecl:
			g.addImport(s)
		case *ast.TypeAliasDecl:
			g.typeAliases[s.Name] = s.Type
		}
	}

	// Generate all top-level declarations into a buffer
	var body strings.Builder
	for _, stmt := range prog.Stmts {
		switch stmt.(type) {
		case *ast.ImportDecl:
			continue // handled separately
		default:
			old := g.buf
			g.buf = strings.Builder{}
			g.genNode(stmt)
			body.WriteString(g.buf.String())
			g.buf = old
		}
	}

	// Now build the final output
	g.buf.Reset()
	g.writeln("package main")
	g.writeln("")

	// Always import fmt for print/println
	g.imports["fmt"] = ""
	if g.needsFile {
		g.imports["os"] = ""
		g.imports["strings"] = ""
		g.imports["path/filepath"] = ""
	}
	if g.needsIntConv || g.needsF64Conv || g.needsSliceToInt || g.needsSliceToF64 || g.needsIntBase || g.needsToBase {
		g.imports["strconv"] = ""
	}
	if g.needsIntConv || g.needsF64Conv || g.needsSliceToInt || g.needsSliceToF64 || g.needsSliceToStr || g.needsIntBase || g.needsToBase {
		g.imports["os"] = ""
	}
	if g.needsFlag {
		g.imports["flag"] = ""
	}

	if len(g.imports) > 0 {
		g.writeln("import (")
		g.indent++
		for path, alias := range g.imports {
			if alias != "" {
				g.writef("%s %q\n", alias, path)
			} else {
				g.writef("%q\n", path)
			}
		}
		g.indent--
		g.writeln(")")
		g.writeln("")
	}

	// Emit package-level flag variable declarations
	for _, fd := range g.flagDecls {
		switch fd.goType {
		case "bool":
			g.writef("var _zflag_%s = flag.Bool(%q, %s, \"\")\n", fd.name, fd.name, fd.defaultVal)
		case "int":
			g.writef("var _zflag_%s = flag.Int(%q, %s, \"\")\n", fd.name, fd.name, fd.defaultVal)
		case "string":
			g.writef("var _zflag_%s = flag.String(%q, %s, \"\")\n", fd.name, fd.name, fd.defaultVal)
		}
	}
	if len(g.flagDecls) > 0 {
		g.writeln("")
	}

	if g.needsRange {
		g.writeln("func zenth_range(start, end, step int) []int {")
		g.writeln("\tvar r []int")
		g.writeln("\tif step > 0 {")
		g.writeln("\t\tfor i := start; i < end; i += step { r = append(r, i) }")
		g.writeln("\t} else if step < 0 {")
		g.writeln("\t\tfor i := start; i > end; i += step { r = append(r, i) }")
		g.writeln("\t}")
		g.writeln("\treturn r")
		g.writeln("}")
		g.writeln("")
	}
	if g.needsRangei {
		g.writeln("func zenth_rangei(start, end, step int) []int {")
		g.writeln("\tvar r []int")
		g.writeln("\tif step > 0 {")
		g.writeln("\t\tfor i := start; i <= end; i += step { r = append(r, i) }")
		g.writeln("\t} else if step < 0 {")
		g.writeln("\t\tfor i := start; i >= end; i += step { r = append(r, i) }")
		g.writeln("\t}")
		g.writeln("\treturn r")
		g.writeln("}")
		g.writeln("")
	}
	if g.needsPop {
		g.writeln("func zenth_pop[T any](s *[]T) T {")
		g.writeln("\tif len(*s) == 0 {")
		g.writeln("\t\tvar zero T")
		g.writeln("\t\treturn zero")
		g.writeln("\t}")
		g.writeln("\tn := len(*s) - 1")
		g.writeln("\telem := (*s)[n]")
		g.writeln("\t*s = (*s)[:n]")
		g.writeln("\treturn elem")
		g.writeln("}")
		g.writeln("")
		g.writeln("func zenth_pop_at[T any](s *[]T, i int) T {")
		g.writeln("\tif i < 0 || i >= len(*s) {")
		g.writeln("\t\tvar zero T")
		g.writeln("\t\treturn zero")
		g.writeln("\t}")
		g.writeln("\telem := (*s)[i]")
		g.writeln("\t*s = append((*s)[:i], (*s)[i+1:]...)")
		g.writeln("\treturn elem")
		g.writeln("}")
		g.writeln("")
	}
	if g.needsAdd {
		g.writeln("func zenth_add[T any](s *[]T, elem T) {")
		g.writeln("\t*s = append(*s, elem)")
		g.writeln("}")
		g.writeln("")
	}
	if g.needsPush {
		g.writeln("func zenth_push[T any](s *[]T, elem T) {")
		g.writeln("\t*s = append([]T{elem}, *s...)")
		g.writeln("}")
		g.writeln("")
	}
	if g.needsContains {
		g.writeln("func zenth_contains[T comparable](s []T, elem T) bool {")
		g.writeln("\tfor _, v := range s {")
		g.writeln("\t\tif v == elem { return true }")
		g.writeln("\t}")
		g.writeln("\treturn false")
		g.writeln("}")
		g.writeln("")
	}
	if g.needsFile {
		g.writeln("type ZenthFile struct { Path string }")
		g.writeln("")
		g.writeln("func zenth_file(path string) ZenthFile { return ZenthFile{Path: path} }")
		g.writeln("")
		g.writeln("func zenth_file_exists(f ZenthFile) bool {")
		g.writeln("\t_, err := os.Stat(f.Path)")
		g.writeln("\treturn err == nil")
		g.writeln("}")
		g.writeln("")
		g.writeln("func zenth_file_read(f ZenthFile) string {")
		g.writeln("\tdata, err := os.ReadFile(f.Path)")
		g.writeln("\tif err != nil {")
		g.writeln("\t\tfmt.Fprintf(os.Stderr, \"error: %v\\n\", err)")
		g.writeln("\t\tos.Exit(1)")
		g.writeln("\t}")
		g.writeln("\treturn string(data)")
		g.writeln("}")
		g.writeln("")
		g.writeln("func zenth_file_lines(f ZenthFile) []string {")
		g.writeln("\treturn strings.Split(strings.TrimRight(zenth_file_read(f), \"\\n\"), \"\\n\")")
		g.writeln("}")
		g.writeln("")
		g.writeln("func zenth_file_name(f ZenthFile) string { return filepath.Base(f.Path) }")
		g.writeln("")
		g.writeln("func zenth_file_ext(f ZenthFile) string { return filepath.Ext(f.Path) }")
		g.writeln("")
	}
	if g.needsIntConv {
		g.writeln("func zenth_int(v interface{}) int {")
		g.writeln("\tswitch x := v.(type) {")
		g.writeln("\tcase int:")
		g.writeln("\t\treturn x")
		g.writeln("\tcase float64:")
		g.writeln("\t\treturn int(x)")
		g.writeln("\tcase float32:")
		g.writeln("\t\treturn int(x)")
		g.writeln("\tcase bool:")
		g.writeln("\t\tif x { return 1 }")
		g.writeln("\t\treturn 0")
		g.writeln("\tcase string:")
		g.writeln("\t\tn, err := strconv.Atoi(x)")
		g.writeln("\t\tif err != nil {")
		g.writeln("\t\t\tfmt.Fprintf(os.Stderr, \"error: cannot convert %q to int\\n\", x)")
		g.writeln("\t\t\tos.Exit(1)")
		g.writeln("\t\t}")
		g.writeln("\t\treturn n")
		g.writeln("\tdefault:")
		g.writeln("\t\tfmt.Fprintf(os.Stderr, \"error: cannot convert %T to int\\n\", v)")
		g.writeln("\t\tos.Exit(1)")
		g.writeln("\t\treturn 0")
		g.writeln("\t}")
		g.writeln("}")
		g.writeln("")
	}
	if g.needsIntBase {
		g.writeln("func zenth_int_base(s string, base int) int {")
		g.writeln("\tn, err := strconv.ParseInt(s, base, 64)")
		g.writeln("\tif err != nil {")
		g.writeln("\t\tfmt.Fprintf(os.Stderr, \"error: cannot convert %q to int with base %d\\n\", s, base)")
		g.writeln("\t\tos.Exit(1)")
		g.writeln("\t}")
		g.writeln("\treturn int(n)")
		g.writeln("}")
		g.writeln("")
	}
	if g.needsToBase {
		g.writeln("func zenth_to_base(n int64, base int) string {")
		g.writeln("\treturn strconv.FormatInt(n, base)")
		g.writeln("}")
		g.writeln("")
	}
	if g.needsF64Conv {
		g.writeln("func zenth_f64(v interface{}) float64 {")
		g.writeln("\tswitch x := v.(type) {")
		g.writeln("\tcase float64:")
		g.writeln("\t\treturn x")
		g.writeln("\tcase float32:")
		g.writeln("\t\treturn float64(x)")
		g.writeln("\tcase int:")
		g.writeln("\t\treturn float64(x)")
		g.writeln("\tcase string:")
		g.writeln("\t\tn, err := strconv.ParseFloat(x, 64)")
		g.writeln("\t\tif err != nil {")
		g.writeln("\t\t\tfmt.Fprintf(os.Stderr, \"error: cannot convert %q to f64\\n\", x)")
		g.writeln("\t\t\tos.Exit(1)")
		g.writeln("\t\t}")
		g.writeln("\t\treturn n")
		g.writeln("\tdefault:")
		g.writeln("\t\tfmt.Fprintf(os.Stderr, \"error: cannot convert %T to f64\\n\", v)")
		g.writeln("\t\tos.Exit(1)")
		g.writeln("\t\treturn 0")
		g.writeln("\t}")
		g.writeln("}")
		g.writeln("")
	}
	if g.needsSliceToInt {
		g.writeln("func zenth_slice_to_int(s interface{}) []int {")
		g.writeln("\tvar result []int")
		g.writeln("\tswitch xs := s.(type) {")
		g.writeln("\tcase []interface{}:")
		g.writeln("\t\tfor _, v := range xs {")
		g.writeln("\t\t\tn, err := strconv.Atoi(fmt.Sprint(v))")
		g.writeln("\t\t\tif err != nil {")
		g.writeln("\t\t\t\tfmt.Fprintf(os.Stderr, \"error: cannot convert %q to int\\n\", fmt.Sprint(v))")
		g.writeln("\t\t\t\tos.Exit(1)")
		g.writeln("\t\t\t}")
		g.writeln("\t\t\tresult = append(result, n)")
		g.writeln("\t\t}")
		g.writeln("\tcase []string:")
		g.writeln("\t\tfor _, v := range xs {")
		g.writeln("\t\t\tn, err := strconv.Atoi(v)")
		g.writeln("\t\t\tif err != nil {")
		g.writeln("\t\t\t\tfmt.Fprintf(os.Stderr, \"error: cannot convert %q to int\\n\", v)")
		g.writeln("\t\t\t\tos.Exit(1)")
		g.writeln("\t\t\t}")
		g.writeln("\t\t\tresult = append(result, n)")
		g.writeln("\t\t}")
		g.writeln("\t}")
		g.writeln("\treturn result")
		g.writeln("}")
		g.writeln("")
	}
	if g.needsSliceToF64 {
		g.writeln("func zenth_slice_to_f64(s interface{}) []float64 {")
		g.writeln("\tvar result []float64")
		g.writeln("\tswitch xs := s.(type) {")
		g.writeln("\tcase []interface{}:")
		g.writeln("\t\tfor _, v := range xs {")
		g.writeln("\t\t\tn, err := strconv.ParseFloat(fmt.Sprint(v), 64)")
		g.writeln("\t\t\tif err != nil {")
		g.writeln("\t\t\t\tfmt.Fprintf(os.Stderr, \"error: cannot convert %q to f64\\n\", fmt.Sprint(v))")
		g.writeln("\t\t\t\tos.Exit(1)")
		g.writeln("\t\t\t}")
		g.writeln("\t\t\tresult = append(result, n)")
		g.writeln("\t\t}")
		g.writeln("\tcase []string:")
		g.writeln("\t\tfor _, v := range xs {")
		g.writeln("\t\t\tn, err := strconv.ParseFloat(v, 64)")
		g.writeln("\t\t\tif err != nil {")
		g.writeln("\t\t\t\tfmt.Fprintf(os.Stderr, \"error: cannot convert %q to f64\\n\", v)")
		g.writeln("\t\t\t\tos.Exit(1)")
		g.writeln("\t\t\t}")
		g.writeln("\t\t\tresult = append(result, n)")
		g.writeln("\t\t}")
		g.writeln("\t}")
		g.writeln("\treturn result")
		g.writeln("}")
		g.writeln("")
	}
	if g.needsSliceToStr {
		g.writeln("func zenth_slice_to_str(s interface{}) []string {")
		g.writeln("\tvar result []string")
		g.writeln("\tswitch xs := s.(type) {")
		g.writeln("\tcase []interface{}:")
		g.writeln("\t\tfor _, v := range xs {")
		g.writeln("\t\t\tresult = append(result, fmt.Sprint(v))")
		g.writeln("\t\t}")
		g.writeln("\tcase []int:")
		g.writeln("\t\tfor _, v := range xs {")
		g.writeln("\t\t\tresult = append(result, fmt.Sprint(v))")
		g.writeln("\t\t}")
		g.writeln("\tcase []float64:")
		g.writeln("\t\tfor _, v := range xs {")
		g.writeln("\t\t\tresult = append(result, fmt.Sprint(v))")
		g.writeln("\t\t}")
		g.writeln("\tdefault:")
		g.writeln("\t\tfmt.Fprintf(os.Stderr, \"error: cannot convert %T to []str\\n\", s)")
		g.writeln("\t\tos.Exit(1)")
		g.writeln("\t}")
		g.writeln("\treturn result")
		g.writeln("}")
		g.writeln("")
	}
	if g.needsStrIndex {
		g.writeln("func zenth_str_index(s string, i int) string {")
		g.writeln("\tr := []rune(s)")
		g.writeln("\treturn string(r[i])")
		g.writeln("}")
		g.writeln("")
	}
	if g.needsStrSlice {
		g.writeln("func zenth_str_slice(s string, lo, hi int) string {")
		g.writeln("\tr := []rune(s)")
		g.writeln("\tif hi < 0 { hi = len(r) }")
		g.writeln("\treturn string(r[lo:hi])")
		g.writeln("}")
		g.writeln("")
	}
	if g.needsHashmapObjKey {
		g.writeln("func zenth_hashmap_obj_key(v interface{}) string {")
		g.writeln("\treturn fmt.Sprintf(\"%#v\", v)")
		g.writeln("}")
		g.writeln("")
	}
	if g.needsMap {
		g.writeln("func zenth_map[T any, U any](s []T, f func(T) U) []U {")
		g.writeln("\tr := make([]U, len(s))")
		g.writeln("\tfor i, v := range s { r[i] = f(v) }")
		g.writeln("\treturn r")
		g.writeln("}")
		g.writeln("")
	}
	if g.needsFilter {
		g.writeln("func zenth_filter[T any](s []T, f func(T) bool) []T {")
		g.writeln("\tvar r []T")
		g.writeln("\tfor _, v := range s { if f(v) { r = append(r, v) } }")
		g.writeln("\treturn r")
		g.writeln("}")
		g.writeln("")
	}
	if g.needsSplitOnce {
		g.writeln("func zenth_split_once(s, sep string) []interface{} {")
		g.writeln("\tparts := strings.SplitN(s, sep, 2)")
		g.writeln("\tif len(parts) == 1 {")
		g.writeln("\t\treturn []interface{}{parts[0], \"\"}")
		g.writeln("\t}")
		g.writeln("\treturn []interface{}{parts[0], parts[1]}")
		g.writeln("}")
		g.writeln("")
	}
	if g.needsAbsInt {
		g.writeln("func zenth_abs_int(v int) int {")
		g.writeln("\tif v < 0 { return -v }")
		g.writeln("\treturn v")
		g.writeln("}")
		g.writeln("")
	}
	if g.needsMinInt {
		g.writeln("func zenth_min_int(a, b int) int {")
		g.writeln("\tif a < b { return a }")
		g.writeln("\treturn b")
		g.writeln("}")
		g.writeln("")
	}
	if g.needsMaxInt {
		g.writeln("func zenth_max_int(a, b int) int {")
		g.writeln("\tif a > b { return a }")
		g.writeln("\treturn b")
		g.writeln("}")
		g.writeln("")
	}
	if g.needsClampInt {
		g.writeln("func zenth_clamp_int(x, lo, hi int) int {")
		g.writeln("\tif x < lo { return lo }")
		g.writeln("\tif x > hi { return hi }")
		g.writeln("\treturn x")
		g.writeln("}")
		g.writeln("")
	}
	if g.needsClampF64 {
		g.writeln("func zenth_clamp_f64(x, lo, hi float64) float64 {")
		g.writeln("\treturn math.Max(lo, math.Min(x, hi))")
		g.writeln("}")
		g.writeln("")
	}
	if g.needsMapget {
		g.writeln("func zenth_mapget[K comparable, V any](m map[K]V, key K, def V) V {")
		g.writeln("\tif v, ok := m[key]; ok { return v }")
		g.writeln("\treturn def")
		g.writeln("}")
		g.writeln("")
	}
	if g.needsHashmapKeys {
		g.writeln("func zenth_hashmap_keys[K comparable, V any](m map[K]V) []K {")
		g.writeln("\tkeys := make([]K, 0, len(m))")
		g.writeln("\tfor k := range m { keys = append(keys, k) }")
		g.writeln("\treturn keys")
		g.writeln("}")
		g.writeln("")
	}
	if g.needsHashmapValues {
		g.writeln("func zenth_hashmap_values[K comparable, V any](m map[K]V) []V {")
		g.writeln("\tvals := make([]V, 0, len(m))")
		g.writeln("\tfor _, v := range m { vals = append(vals, v) }")
		g.writeln("\treturn vals")
		g.writeln("}")
		g.writeln("")
	}

	g.buf.WriteString(body.String())

	return g.buf.String()
}

func (g *Generator) addImport(imp *ast.ImportDecl) {
	goPath := mapImportPath(imp.Path)
	g.imports[goPath] = imp.Alias
}

func mapImportPath(zenthPath string) string {
	switch zenthPath {
	case "fmt":
		return "fmt"
	case "math":
		return "math"
	case "os":
		return "os"
	case "strings", "str":
		return "strings"
	case "io":
		return "os" // map to os for file I/O
	default:
		return zenthPath
	}
}

func (g *Generator) write(s string) {
	g.buf.WriteString(s)
}

func (g *Generator) writeln(s string) {
	g.writeIndent()
	g.buf.WriteString(s)
	g.buf.WriteString("\n")
}

func (g *Generator) writef(format string, args ...any) {
	g.writeIndent()
	fmt.Fprintf(&g.buf, format, args...)
}

func (g *Generator) writeIndent() {
	for i := 0; i < g.indent; i++ {
		g.buf.WriteString("\t")
	}
}

func (g *Generator) genNode(node ast.Node) {
	switch n := node.(type) {
	case *ast.FnDecl:
		g.genFnDecl(n)
	case *ast.ObjDecl:
		g.genObjDecl(n)
	case *ast.InterfaceDecl:
		g.genInterfaceDecl(n)
	case *ast.TypeAliasDecl:
		g.writef("type %s = %s\n\n", n.Name, g.genType(n.Type))
	case *ast.Block:
		g.genBlock(n)
	case *ast.LetStmt:
		g.genLetStmt(n)
	case *ast.VarStmt:
		g.genVarStmt(n)
	case *ast.ConstStmt:
		g.genConstStmt(n)
	case *ast.TupleDestructStmt:
		g.genTupleDestructStmt(n)
	case *ast.AssignStmt:
		g.genAssignStmt(n)
	case *ast.MultiAssignStmt:
		g.genMultiAssignStmt(n)
	case *ast.ReturnStmt:
		g.genReturnStmt(n)
	case *ast.IfStmt:
		g.genIfStmt(n)
	case *ast.ForStmt:
		g.genForStmt(n)
	case *ast.ForInStmt:
		g.genForInStmt(n)
	case *ast.LoopStmt:
		g.genLoopStmt(n)
	case *ast.MatchStmt:
		g.genMatchStmt(n)
	case *ast.BreakStmt:
		g.writeln("break")
	case *ast.ContinueStmt:
		g.writeln("continue")
	case *ast.IncDecStmt:
		g.genIncDecStmt(n)
	case *ast.ExprStmt:
		g.writeIndent()
		g.genExpr(n.Expr)
		g.write("\n")
	default:
		g.writef("// unhandled: %T\n", node)
	}
}

func (g *Generator) genFnDecl(f *ast.FnDecl) {
	g.writeIndent()
	g.write("func ")
	if f.OwnerObj != "" {
		g.write("(self *")
		g.write(f.OwnerObj)
		g.write(") ")
		g.write(exportName(f.Name))
	} else {
		g.write(f.Name)
	}
	g.write("(")
	for i, p := range f.Params {
		if i > 0 {
			g.write(", ")
		}
		g.write(p.Name)
		g.write(" ")
		g.write(g.genType(p.Type))
	}
	g.write(")")
	if f.ReturnType != nil {
		g.write(" ")
		g.write(g.genType(f.ReturnType))
	}
	g.write(" {\n")
	g.indent++
	if f.Name == "main" && f.OwnerObj == "" {
		// Generate body into temp buffer so we know if flag() was used
		oldBuf := g.buf
		g.buf = strings.Builder{}
		if f.Body != nil {
			for _, stmt := range f.Body.Stmts {
				g.genNode(stmt)
			}
		}
		bodyStr := g.buf.String()
		g.buf = oldBuf
		if g.needsFlag {
			g.writeln("flag.Parse()")
		}
		g.buf.WriteString(bodyStr)
	} else {
		if f.Body != nil {
			for _, stmt := range f.Body.Stmts {
				g.genNode(stmt)
			}
		}
	}
	g.indent--
	g.writeln("}")
	g.writeln("")
}

func (g *Generator) genObjDecl(s *ast.ObjDecl) {
	g.writef("type %s struct {\n", s.Name)
	g.indent++
	for _, f := range s.Fields {
		g.writef("%s %s\n", exportName(f.Name), g.genType(f.Type))
	}
	g.indent--
	g.writeln("}")
	g.writeln("")

	if !objHasStringMethod(s) {
		g.genDefaultObjStringMethod(s)
	}

	// Generate methods
	for _, m := range s.Methods {
		g.genFnDecl(m)
	}
}

func (g *Generator) genDefaultObjStringMethod(s *ast.ObjDecl) {
	g.writef("func (self *%s) String() string {\n", s.Name)
	g.indent++
	g.writeIndent()
	g.writef("if self == nil { return %q }\n", s.Name+"(nil)")

	if len(s.Fields) == 0 {
		g.writeIndent()
		g.writef("return %q\n", s.Name+"()")
		g.indent--
		g.writeln("}")
		g.writeln("")
		return
	}

	var format strings.Builder
	format.WriteString(s.Name + "(")
	for i, f := range s.Fields {
		if i > 0 {
			format.WriteString(", ")
		}
		format.WriteString(f.Name + "=" + objFieldFormatVerb(f.Type))
	}
	format.WriteString(")")

	g.writeIndent()
	g.writef("return fmt.Sprintf(%q", format.String())
	for _, f := range s.Fields {
		g.write(", ")
		g.write("self." + exportName(f.Name))
	}
	g.write(")\n")
	g.indent--
	g.writeln("}")
	g.writeln("")
}

func objHasStringMethod(s *ast.ObjDecl) bool {
	for _, m := range s.Methods {
		if strings.EqualFold(m.Name, "string") {
			return true
		}
	}
	return false
}

func objFieldFormatVerb(t *ast.TypeExpr) string {
	if t != nil && !t.IsSlice && t.Name == "str" {
		return "%q"
	}
	return "%v"
}

func (g *Generator) genInterfaceDecl(iface *ast.InterfaceDecl) {
	g.writef("type %s interface {\n", iface.Name)
	g.indent++
	for _, m := range iface.Methods {
		g.writeIndent()
		g.write(exportName(m.Name))
		g.write("(")
		for i, p := range m.Params {
			if i > 0 {
				g.write(", ")
			}
			g.write(p.Name)
			g.write(" ")
			g.write(g.genType(p.Type))
		}
		g.write(")")
		if m.ReturnType != nil {
			g.write(" ")
			g.write(g.genType(m.ReturnType))
		}
		g.write("\n")
	}
	g.indent--
	g.writeln("}")
	g.writeln("")
}

func (g *Generator) genBlock(b *ast.Block) {
	for _, stmt := range b.Stmts {
		g.genNode(stmt)
	}
}

func (g *Generator) genLetStmt(s *ast.LetStmt) {
	g.writeIndent()
	if s.Infer || s.Type == nil {
		g.write(s.Name + " := ")
	} else {
		g.write("var " + s.Name + " " + g.genType(s.Type) + " = ")
	}
	g.genExpr(s.Value)
	g.write("\n")
	// Suppress Go's "declared and not used" by referencing with _
	g.writef("_ = %s\n", s.Name)
}

func (g *Generator) genVarStmt(s *ast.VarStmt) {
	g.writeIndent()
	if s.Infer || s.Type == nil {
		g.write(s.Name + " := ")
	} else {
		g.write("var " + s.Name + " " + g.genType(s.Type) + " = ")
	}
	g.genExpr(s.Value)
	g.write("\n")
}

func (g *Generator) genConstStmt(s *ast.ConstStmt) {
	g.writeIndent()
	// Go const requires compile-time constant expressions
	// Use var for safety since Zenth const may reference runtime values
	g.write(s.Name + " := ")
	g.genExpr(s.Value)
	g.write("\n")
	g.writef("_ = %s\n", s.Name)
}

func (g *Generator) genTupleDestructStmt(s *ast.TupleDestructStmt) {
	tmp := fmt.Sprintf("__ztuple%d", g.tempCounter)
	g.tempCounter++
	g.writeIndent()
	g.write(tmp + " := ")
	g.genExpr(s.Value)
	g.write("\n")
	for i, name := range s.Names {
		if name == "_" {
			continue
		}
		g.writeIndent()
		if s.Kind == token.Var {
			g.write("var " + name + " = ")
		} else {
			g.write(name + " := ")
		}
		g.write(fmt.Sprintf("%s[%d]", tmp, i))
		if i < len(s.ElemGoTypes) && s.ElemGoTypes[i] != "" && s.ElemGoTypes[i] != "interface{}" {
			g.write(".(" + s.ElemGoTypes[i] + ")")
		}
		g.write("\n")
		if s.Kind != token.Var {
			g.writef("_ = %s\n", name)
		}
	}
}

func (g *Generator) genAssignStmt(s *ast.AssignStmt) {
	// For compound ops on hashmap index with default, emit init-if-missing first
	if s.Op != token.Assign {
		if idx, ok := s.Target.(*ast.IndexExpr); ok && idx.HashmapDefaultVal != nil {
			g.writeIndent()
			g.write("if _, ok := ")
			g.genExpr(idx.Object)
			g.write("[")
			if idx.HashmapObjKey {
				g.needsHashmapObjKey = true
				g.write("zenth_hashmap_obj_key(")
				g.genExpr(idx.Index)
				g.write(")")
			} else {
				g.genExpr(idx.Index)
			}
			g.write("]; !ok { ")
			g.genExpr(idx.Object)
			g.write("[")
			if idx.HashmapObjKey {
				g.write("zenth_hashmap_obj_key(")
				g.genExpr(idx.Index)
				g.write(")")
			} else {
				g.genExpr(idx.Index)
			}
			g.write("] = ")
			g.genExpr(idx.HashmapDefaultVal)
			g.write(" }\n")
		}
	}
	g.writeIndent()
	// For assignments to hashmap index (both plain and compound), use raw m[key] not mapget
	if idx, ok := s.Target.(*ast.IndexExpr); ok && idx.HashmapDefaultVal != nil {
		g.genRawIndexExpr(idx)
	} else {
		g.genExpr(s.Target)
	}
	switch s.Op {
	case token.Assign:
		g.write(" = ")
	case token.PlusAssign:
		g.write(" += ")
	case token.MinusAssign:
		g.write(" -= ")
	case token.StarAssign:
		g.write(" *= ")
	case token.SlashAssign:
		g.write(" /= ")
	}
	g.genExpr(s.Value)
	g.write("\n")
}

// genRawIndexExpr generates m[key] without mapget wrapper
func (g *Generator) genRawIndexExpr(idx *ast.IndexExpr) {
	g.genExpr(idx.Object)
	g.write("[")
	if idx.HashmapObjKey {
		g.needsHashmapObjKey = true
		g.write("zenth_hashmap_obj_key(")
		g.genExpr(idx.Index)
		g.write(")")
	} else {
		g.genExpr(idx.Index)
	}
	g.write("]")
}

func (g *Generator) genMultiAssignStmt(s *ast.MultiAssignStmt) {
	g.writeIndent()
	for i, target := range s.Targets {
		if i > 0 {
			g.write(", ")
		}
		g.genExpr(target)
	}
	g.write(" = ")
	for i, value := range s.Values {
		if i > 0 {
			g.write(", ")
		}
		g.genExpr(value)
	}
	g.write("\n")
}

func (g *Generator) genReturnStmt(s *ast.ReturnStmt) {
	g.writeIndent()
	g.write("return")
	if s.Value != nil {
		g.write(" ")
		g.genExpr(s.Value)
	}
	g.write("\n")
}

func (g *Generator) genIfStmt(s *ast.IfStmt) {
	g.writeIndent()
	g.write("if ")
	g.genExpr(s.Condition)
	g.write(" {\n")
	g.indent++
	for _, stmt := range s.Body.Stmts {
		g.genNode(stmt)
	}
	g.indent--
	if s.Else != nil {
		switch e := s.Else.(type) {
		case *ast.IfStmt:
			g.writeIndent()
			g.write("} else ")
			// Remove indent since genIfStmt will add its own
			old := g.indent
			g.indent = 0
			g.write("if ")
			g.genExpr(e.Condition)
			g.write(" {\n")
			g.indent = old + 1
			for _, stmt := range e.Body.Stmts {
				g.genNode(stmt)
			}
			g.indent = old
			if e.Else != nil {
				g.genElseChain(e.Else)
			} else {
				g.writeln("}")
			}
		case *ast.Block:
			g.writeln("} else {")
			g.indent++
			for _, stmt := range e.Stmts {
				g.genNode(stmt)
			}
			g.indent--
			g.writeln("}")
		}
	} else {
		g.writeln("}")
	}
}

func (g *Generator) genElseChain(node ast.Node) {
	switch e := node.(type) {
	case *ast.IfStmt:
		g.writeIndent()
		g.write("} else if ")
		g.genExpr(e.Condition)
		g.write(" {\n")
		g.indent++
		for _, stmt := range e.Body.Stmts {
			g.genNode(stmt)
		}
		g.indent--
		if e.Else != nil {
			g.genElseChain(e.Else)
		} else {
			g.writeln("}")
		}
	case *ast.Block:
		g.writeln("} else {")
		g.indent++
		for _, stmt := range e.Stmts {
			g.genNode(stmt)
		}
		g.indent--
		g.writeln("}")
	}
}

func (g *Generator) genForStmt(s *ast.ForStmt) {
	g.writeIndent()
	if s.Init == nil && s.Condition == nil && s.Post == nil {
		// Infinite loop
		g.write("for {\n")
	} else if s.Init == nil && s.Post == nil {
		// While-style
		g.write("for ")
		g.genExpr(s.Condition)
		g.write(" {\n")
	} else {
		// C-style
		g.write("for ")
		if s.Init != nil {
			g.genForClause(s.Init)
		}
		g.write("; ")
		if s.Condition != nil {
			g.genExpr(s.Condition)
		}
		g.write("; ")
		if s.Post != nil {
			g.genForClause(s.Post)
		}
		g.write(" {\n")
	}
	g.indent++
	for _, stmt := range s.Body.Stmts {
		g.genNode(stmt)
	}
	g.indent--
	g.writeln("}")
}

func (g *Generator) genForClause(node ast.Node) {
	switch n := node.(type) {
	case *ast.VarStmt:
		g.write(n.Name + " := ")
		g.genExpr(n.Value)
	case *ast.LetStmt:
		g.write(n.Name + " := ")
		g.genExpr(n.Value)
	case *ast.AssignStmt:
		g.genExpr(n.Target)
		switch n.Op {
		case token.Assign:
			g.write(" = ")
		case token.PlusAssign:
			g.write(" += ")
		}
		g.genExpr(n.Value)
	case *ast.MultiAssignStmt:
		for i, target := range n.Targets {
			if i > 0 {
				g.write(", ")
			}
			g.genExpr(target)
		}
		g.write(" = ")
		for i, value := range n.Values {
			if i > 0 {
				g.write(", ")
			}
			g.genExpr(value)
		}
	case *ast.IncDecStmt:
		g.genExpr(n.Operand)
		if n.Op == token.PlusPlus {
			g.write("++")
		} else {
			g.write("--")
		}
	case *ast.ExprStmt:
		g.genExpr(n.Expr)
	}
}

func (g *Generator) genForInStmt(s *ast.ForInStmt) {
	g.writeIndent()
	if s.IterStr {
		// Go's range over string yields (index, rune); wrap value in string().
		g.write("for ")
		if s.Index != "" {
			g.write(s.Index)
		} else {
			g.write("_")
		}
		g.write(", _rune := range ")
		g.genExpr(s.Iterable)
		g.write(" {\n")
		g.indent++
		g.writef("%s := string(_rune)\n", s.Value)
	} else if s.IterHashmap {
		// Go's range over map yields (key, value)
		g.write("for ")
		if s.Index != "" {
			// for k, v in m => for k, v := range m
			g.write(s.Index)
			g.write(", ")
			g.write(s.Value)
		} else {
			// for k in m => for k := range m (keys only)
			g.write(s.Value)
		}
		g.write(" := range ")
		g.genExpr(s.Iterable)
		g.write(" {\n")
		g.indent++
	} else {
		g.write("for ")
		if s.Index != "" {
			g.write(s.Index)
		} else {
			g.write("_")
		}
		g.write(", ")
		g.write(s.Value)
		g.write(" := range ")
		g.genExpr(s.Iterable)
		g.write(" {\n")
		g.indent++
	}
	for _, stmt := range s.Body.Stmts {
		g.genNode(stmt)
	}
	g.indent--
	g.writeln("}")
}

func (g *Generator) genLoopStmt(s *ast.LoopStmt) {
	v := fmt.Sprintf("_loop%d_", g.loopCounter)
	g.loopCounter++
	g.writeIndent()
	g.writef("for %s := 0; %s < ", v, v)
	g.genExpr(s.Count)
	g.writef("; %s++ {\n", v)
	g.indent++
	for _, stmt := range s.Body.Stmts {
		g.genNode(stmt)
	}
	g.indent--
	g.writeln("}")
}

func (g *Generator) genMatchStmt(s *ast.MatchStmt) {
	g.writeIndent()
	g.write("switch ")
	g.genExpr(s.Subject)
	g.write(" {\n")
	for _, arm := range s.Arms {
		g.writeIndent()
		// Check for wildcard _
		if ident, ok := arm.Pattern.(*ast.IdentExpr); ok && ident.Name == "_" {
			g.write("default:\n")
		} else {
			g.write("case ")
			g.genExpr(arm.Pattern)
			g.write(":\n")
		}
		g.indent++
		switch body := arm.Body.(type) {
		case *ast.Block:
			for _, stmt := range body.Stmts {
				g.genNode(stmt)
			}
		default:
			g.genNode(body)
		}
		g.indent--
	}
	g.writeln("}")
}

func (g *Generator) genIncDecStmt(s *ast.IncDecStmt) {
	// For ++/-- on hashmap index with default, emit init-if-missing first
	if idx, ok := s.Operand.(*ast.IndexExpr); ok && idx.HashmapDefaultVal != nil {
		g.writeIndent()
		g.write("if _, ok := ")
		g.genExpr(idx.Object)
		g.write("[")
		if idx.HashmapObjKey {
			g.needsHashmapObjKey = true
			g.write("zenth_hashmap_obj_key(")
			g.genExpr(idx.Index)
			g.write(")")
		} else {
			g.genExpr(idx.Index)
		}
		g.write("]; !ok { ")
		g.genExpr(idx.Object)
		g.write("[")
		if idx.HashmapObjKey {
			g.write("zenth_hashmap_obj_key(")
			g.genExpr(idx.Index)
			g.write(")")
		} else {
			g.genExpr(idx.Index)
		}
		g.write("] = ")
		g.genExpr(idx.HashmapDefaultVal)
		g.write(" }\n")
	}
	g.writeIndent()
	if idx, ok := s.Operand.(*ast.IndexExpr); ok && idx.HashmapDefaultVal != nil {
		g.genRawIndexExpr(idx)
	} else {
		g.genExpr(s.Operand)
	}
	if s.Op == token.PlusPlus {
		g.write("++")
	} else {
		g.write("--")
	}
	g.write("\n")
}

// ---------- Expression generation ----------

func (g *Generator) genExpr(node ast.Node) {
	switch n := node.(type) {
	case *ast.BinaryExpr:
		if n.SliceConcat {
			g.write("append(")
			g.genExpr(n.Left)
			g.write(", ")
			g.genExpr(n.Right)
			g.write("...)")
			break
		}
		if n.PromoteLeft != "" {
			g.write(n.PromoteLeft + "(")
			g.genExpr(n.Left)
			g.write(")")
		} else {
			g.genExpr(n.Left)
		}
		g.write(" " + goOp(n.Op) + " ")
		if n.PromoteRight != "" {
			g.write(n.PromoteRight + "(")
			g.genExpr(n.Right)
			g.write(")")
		} else {
			g.genExpr(n.Right)
		}
	case *ast.UnaryExpr:
		g.write(goOp(n.Op))
		g.genExpr(n.Operand)
	case *ast.CallExpr:
		g.genCallExpr(n)
	case *ast.IndexExpr:
		if n.StrIndex {
			g.needsStrIndex = true
			g.write("zenth_str_index(")
			g.genExpr(n.Object)
			g.write(", ")
			g.genExpr(n.Index)
			g.write(")")
		} else if n.HashmapDefaultVal != nil {
			g.needsMapget = true
			g.write("zenth_mapget(")
			g.genExpr(n.Object)
			g.write(", ")
			if n.HashmapObjKey {
				g.needsHashmapObjKey = true
				g.write("zenth_hashmap_obj_key(")
				g.genExpr(n.Index)
				g.write(")")
			} else {
				g.genExpr(n.Index)
			}
			g.write(", ")
			g.genExpr(n.HashmapDefaultVal)
			g.write(")")
		} else {
			g.genExpr(n.Object)
			g.write("[")
			if n.HashmapObjKey {
				g.needsHashmapObjKey = true
				g.write("zenth_hashmap_obj_key(")
				g.genExpr(n.Index)
				g.write(")")
			} else {
				g.genExpr(n.Index)
			}
			g.write("]")
		}
	case *ast.SliceExpr:
		if n.StrSlice {
			g.needsStrSlice = true
			g.write("zenth_str_slice(")
			g.genExpr(n.Object)
			g.write(", ")
			if n.Low != nil {
				g.genExpr(n.Low)
			} else {
				g.write("0")
			}
			g.write(", ")
			if n.High != nil {
				g.genExpr(n.High)
			} else {
				g.write("-1")
			}
			g.write(")")
		} else {
			g.genExpr(n.Object)
			g.write("[")
			if n.Low != nil {
				g.genExpr(n.Low)
			}
			g.write(":")
			if n.High != nil {
				g.genExpr(n.High)
			}
			g.write("]")
		}
	case *ast.FieldExpr:
		if n.TupleAccess {
			g.write("(")
			g.genExpr(n.Object)
			g.write("[")
			g.write(fmt.Sprintf("%d", n.TupleIndex))
			g.write("]")
			if n.TupleElemGoType != "" && n.TupleElemGoType != "interface{}" {
				g.write(".(")
				g.write(n.TupleElemGoType)
				g.write(")")
			}
			g.write(")")
		} else {
			g.genExpr(n.Object)
			g.write(".")
			g.write(goFieldName(n.Field, n.Object))
		}
	case *ast.IdentExpr:
		g.write(n.Name)
	case *ast.IntLitExpr:
		g.write(fmt.Sprintf("%d", n.Value))
	case *ast.FloatLitExpr:
		s := fmt.Sprintf("%g", n.Value)
		if !strings.Contains(s, ".") && !strings.Contains(s, "e") {
			s += ".0"
		}
		g.write(s)
	case *ast.StringLitExpr:
		g.write(fmt.Sprintf("%q", n.Value))
	case *ast.BoolLitExpr:
		if n.Value {
			g.write("true")
		} else {
			g.write("false")
		}
	case *ast.NilExpr:
		g.write("nil")
	case *ast.ArrayLitExpr:
		g.genArrayLit(n)
	case *ast.TupleLitExpr:
		g.write("[]interface{}{")
		for i, elem := range n.Elements {
			if i > 0 {
				g.write(", ")
			}
			g.genExpr(elem)
		}
		g.write("}")
	case *ast.NamedArgExpr:
		g.genExpr(n.Value)
	case *ast.IfExpr:
		g.genIfExpr(n)
	case *ast.MatchExpr:
		g.genMatchExpr(n)
	case *ast.ClosureExpr:
		g.genClosureExpr(n)
	case *ast.InterpStringExpr:
		g.genInterpString(n)
	default:
		g.write(fmt.Sprintf("/* unhandled expr: %T */", node))
	}
}

func (g *Generator) genCallExpr(c *ast.CallExpr) {
	// Translate built-in functions
	if ident, ok := c.Callee.(*ast.IdentExpr); ok {
		switch ident.Name {
		case "print":
			if len(c.Args) == 2 {
				g.write("if ")
				g.genExpr(c.Args[1])
				g.write(" { fmt.Print(")
				g.genExpr(c.Args[0])
				g.write(") }")
			} else {
				g.write("fmt.Print(")
				g.genArgList(c.Args)
				g.write(")")
			}
			return
		case "println":
			if len(c.Args) == 2 {
				g.write("if ")
				g.genExpr(c.Args[1])
				g.write(" { fmt.Println(")
				g.genExpr(c.Args[0])
				g.write(") }")
			} else {
				g.write("fmt.Println(")
				g.genArgList(c.Args)
				g.write(")")
			}
			return
		case "exit":
			g.imports["os"] = ""
			if len(c.Args) == 0 {
				g.write("os.Exit(0)")
			} else {
				g.write("os.Exit(")
				g.genArgList(c.Args)
				g.write(")")
			}
			return
		case "len":
			g.write("len(")
			g.genArgList(c.Args)
			g.write(")")
			return
		case "str":
			if c.SliceConvFunc == "str" {
				g.needsSliceToStr = true
				g.write("zenth_slice_to_str(")
			} else {
				g.write("fmt.Sprint(")
			}
			g.genArgList(c.Args)
			g.write(")")
			return
		case "range":
			g.needsRange = true
			g.write("zenth_range(")
			g.genArgList(c.Args)
			if len(c.Args) == 2 {
				g.write(", 1")
			}
			g.write(")")
			return
		case "rangei":
			g.needsRangei = true
			g.write("zenth_rangei(")
			g.genArgList(c.Args)
			if len(c.Args) == 2 {
				g.write(", 1")
			}
			g.write(")")
			return
		case "file":
			g.needsFile = true
			g.write("zenth_file(")
			g.genArgList(c.Args)
			g.write(")")
			return
		case "int":
			if c.SliceConvFunc == "int" {
				g.needsSliceToInt = true
				g.write("zenth_slice_to_int(")
			} else {
				g.needsIntConv = true
				g.write("zenth_int(")
			}
			g.genArgList(c.Args)
			g.write(")")
			return
		case "f64":
			if c.SliceConvFunc == "f64" {
				g.needsSliceToF64 = true
				g.write("zenth_slice_to_f64(")
			} else {
				g.needsF64Conv = true
				g.write("zenth_f64(")
			}
			g.genArgList(c.Args)
			g.write(")")
			return
		case "abs":
			switch c.NumericMethod {
			case "abs_int":
				g.needsAbsInt = true
				g.write("zenth_abs_int(")
				g.genArgList(c.Args)
				g.write(")")
			case "abs_f64":
				g.imports["math"] = ""
				g.write("math.Abs(")
				g.genArgList(c.Args)
				g.write(")")
			default:
				g.write("/* invalid abs() */")
			}
			return
		case "min":
			switch c.NumericMethod {
			case "min_int":
				g.needsMinInt = true
				g.write("zenth_min_int(")
				g.genArgList(c.Args)
				g.write(")")
			case "min_f64":
				g.imports["math"] = ""
				g.write("math.Min(")
				g.genArgList(c.Args)
				g.write(")")
			default:
				g.write("/* invalid min() */")
			}
			return
		case "max":
			switch c.NumericMethod {
			case "max_int":
				g.needsMaxInt = true
				g.write("zenth_max_int(")
				g.genArgList(c.Args)
				g.write(")")
			case "max_f64":
				g.imports["math"] = ""
				g.write("math.Max(")
				g.genArgList(c.Args)
				g.write(")")
			default:
				g.write("/* invalid max() */")
			}
			return
		case "clamp":
			switch c.NumericMethod {
			case "clamp_int":
				g.needsClampInt = true
				g.write("zenth_clamp_int(")
				g.genArgList(c.Args)
				g.write(")")
			case "clamp_f64":
				g.imports["math"] = ""
				g.needsClampF64 = true
				g.write("zenth_clamp_f64(")
				g.genArgList(c.Args)
				g.write(")")
			default:
				g.write("/* invalid clamp() */")
			}
			return
		case "round":
			g.imports["math"] = ""
			g.write("math.Round(float64(")
			g.genExpr(c.Args[0])
			g.write("))")
			return
		case "floor":
			g.imports["math"] = ""
			g.write("math.Floor(float64(")
			g.genExpr(c.Args[0])
			g.write("))")
			return
		case "ceil":
			g.imports["math"] = ""
			g.write("math.Ceil(float64(")
			g.genExpr(c.Args[0])
			g.write("))")
			return
		case "pow":
			g.imports["math"] = ""
			g.write("math.Pow(float64(")
			g.genExpr(c.Args[0])
			g.write("), float64(")
			g.genExpr(c.Args[1])
			g.write("))")
			return
		case "sqrt":
			g.imports["math"] = ""
			g.write("math.Sqrt(float64(")
			g.genExpr(c.Args[0])
			g.write("))")
			return
		case "hashmap":
			if c.HashmapCtor {
				g.write("make(map[")
				g.write(c.HashmapKeyGoType)
				g.write("]")
				g.write(c.HashmapValGoType)
				g.write(")")
			} else {
				g.write("/* invalid hashmap() */")
			}
			return
		case "flag":
			g.genFlagCall(c)
			return
		}
	}

	// Handle obj constructor calls
	if ident, ok := c.Callee.(*ast.IdentExpr); ok {
		if decl, ok := g.objs[ident.Name]; ok {
			g.genObjConstructor(decl, c.Args)
			return
		}
	}

	// Handle built-in hashmap methods
	if c.HashmapMethod != "" {
		if field, ok := c.Callee.(*ast.FieldExpr); ok {
			switch c.HashmapMethod {
			case "keys":
				g.needsHashmapKeys = true
				g.write("zenth_hashmap_keys(")
				g.genExpr(field.Object)
				g.write(")")
				return
			case "values":
				g.needsHashmapValues = true
				g.write("zenth_hashmap_values(")
				g.genExpr(field.Object)
				g.write(")")
				return
			}
		}
	}

	// Handle built-in slice methods
	if c.SliceMethod {
		if field, ok := c.Callee.(*ast.FieldExpr); ok {
			switch field.Field {
			case "pop":
				g.needsPop = true
				if len(c.Args) == 0 {
					g.write("zenth_pop(&")
					g.genExpr(field.Object)
					g.write(")")
				} else {
					g.write("zenth_pop_at(&")
					g.genExpr(field.Object)
					g.write(", ")
					g.genExpr(c.Args[0])
					g.write(")")
				}
				return
			case "add":
				g.needsAdd = true
				g.write("zenth_add(&")
				g.genExpr(field.Object)
				g.write(", ")
				g.genExpr(c.Args[0])
				g.write(")")
				return
			case "push":
				g.needsPush = true
				g.write("zenth_push(&")
				g.genExpr(field.Object)
				g.write(", ")
				g.genExpr(c.Args[0])
				g.write(")")
				return
			case "length":
				g.write("len(")
				g.genExpr(field.Object)
				g.write(")")
				return
			case "exists":
				if len(c.Args) == 1 {
					// Slice .exists(elem)
					g.needsContains = true
					g.write("zenth_contains(")
					g.genExpr(field.Object)
					g.write(", ")
					g.genExpr(c.Args[0])
					g.write(")")
				} else {
					// File .exists()
					g.needsFile = true
					g.write("zenth_file_exists(")
					g.genExpr(field.Object)
					g.write(")")
				}
				return
			case "read":
				g.needsFile = true
				g.write("zenth_file_read(")
				g.genExpr(field.Object)
				g.write(")")
				return
			case "lines":
				g.needsFile = true
				g.write("zenth_file_lines(")
				g.genExpr(field.Object)
				g.write(")")
				return
			case "name":
				g.needsFile = true
				g.write("zenth_file_name(")
				g.genExpr(field.Object)
				g.write(")")
				return
			case "ext":
				g.needsFile = true
				g.write("zenth_file_ext(")
				g.genExpr(field.Object)
				g.write(")")
				return
			case "map":
				g.needsMap = true
				g.write("zenth_map(")
				g.genExpr(field.Object)
				g.write(", ")
				g.genExpr(c.Args[0])
				g.write(")")
				return
			case "filter":
				g.needsFilter = true
				g.write("zenth_filter(")
				g.genExpr(field.Object)
				g.write(", ")
				g.genExpr(c.Args[0])
				g.write(")")
				return
			case "to_int":
				if c.SliceConvTarget == "" {
					if len(c.Args) == 1 {
						// Scalar with base: str.to_int(16)
						g.needsIntBase = true
						g.write("zenth_int_base(")
						g.genExpr(field.Object)
						g.write(", ")
						g.genExpr(c.Args[0])
						g.write(")")
					} else {
						// Scalar: str.to_int()
						g.needsIntConv = true
						g.write("zenth_int(")
						g.genExpr(field.Object)
						g.write(")")
					}
				} else {
					// Slice: []str.to_int()
					g.needsSliceToInt = true
					g.write("zenth_slice_to_int(")
					g.genExpr(field.Object)
					g.write(")")
				}
				return
			case "to_f64":
				if c.SliceConvTarget == "" {
					// Scalar: str.to_f64()
					g.needsF64Conv = true
					g.write("zenth_f64(")
				} else {
					// Slice: []str.to_f64()
					g.needsSliceToF64 = true
					g.write("zenth_slice_to_f64(")
				}
				g.genExpr(field.Object)
				g.write(")")
				return
			case "to_str":
				g.needsSliceToStr = true
				g.write("zenth_slice_to_str(")
				g.genExpr(field.Object)
				g.write(")")
				return
			}
		}
	}

	// Handle built-in string methods
	if c.StringMethod != "" {
		if field, ok := c.Callee.(*ast.FieldExpr); ok {
			switch c.StringMethod {
			case "split":
				g.imports["strings"] = ""
				if len(c.Args) == 0 {
					g.write("strings.Fields(")
					g.genExpr(field.Object)
					g.write(")")
				} else {
					g.write("strings.Split(")
					g.genExpr(field.Object)
					g.write(", ")
					g.genExpr(c.Args[0])
					g.write(")")
				}
				return
			case "split_once":
				g.imports["strings"] = ""
				g.needsSplitOnce = true
				g.write("zenth_split_once(")
				g.genExpr(field.Object)
				g.write(", ")
				g.genExpr(c.Args[0])
				g.write(")")
				return
			case "upper":
				g.imports["strings"] = ""
				g.write("strings.ToUpper(")
				g.genExpr(field.Object)
				g.write(")")
				return
			case "lower":
				g.imports["strings"] = ""
				g.write("strings.ToLower(")
				g.genExpr(field.Object)
				g.write(")")
				return
			case "starts_with":
				g.imports["strings"] = ""
				g.write("strings.HasPrefix(")
				g.genExpr(field.Object)
				g.write(", ")
				g.genExpr(c.Args[0])
				g.write(")")
				return
			case "ends_with":
				g.imports["strings"] = ""
				g.write("strings.HasSuffix(")
				g.genExpr(field.Object)
				g.write(", ")
				g.genExpr(c.Args[0])
				g.write(")")
				return
			case "strip":
				g.imports["strings"] = ""
				if len(c.Args) == 0 {
					g.write("strings.TrimSpace(")
					g.genExpr(field.Object)
					g.write(")")
				} else {
					g.write("strings.Trim(")
					g.genExpr(field.Object)
					g.write(", ")
					g.genExpr(c.Args[0])
					g.write(")")
				}
				return
			case "find":
				g.imports["strings"] = ""
				g.write("strings.Index(")
				g.genExpr(field.Object)
				g.write(", ")
				g.genExpr(c.Args[0])
				g.write(")")
				return
			case "count":
				g.imports["strings"] = ""
				g.write("strings.Count(")
				g.genExpr(field.Object)
				g.write(", ")
				g.genExpr(c.Args[0])
				g.write(")")
				return
			case "replace":
				g.imports["strings"] = ""
				if len(c.Args) == 2 {
					g.write("strings.ReplaceAll(")
					g.genExpr(field.Object)
					g.write(", ")
					g.genExpr(c.Args[0])
					g.write(", ")
					g.genExpr(c.Args[1])
					g.write(")")
				} else {
					g.write("strings.Replace(")
					g.genExpr(field.Object)
					g.write(", ")
					g.genExpr(c.Args[0])
					g.write(", ")
					g.genExpr(c.Args[1])
					g.write(", ")
					g.genExpr(c.Args[2])
					g.write(")")
				}
				return
			case "length":
				g.write("len(")
				g.genExpr(field.Object)
				g.write(")")
				return
			case "contains":
				g.imports["strings"] = ""
				g.write("strings.Contains(")
				g.genExpr(field.Object)
				g.write(", ")
				g.genExpr(c.Args[0])
				g.write(")")
				return
			case "to_base":
				g.needsToBase = true
				g.write("zenth_to_base(int64(")
				g.genExpr(field.Object)
				g.write("), ")
				g.genExpr(c.Args[0])
				g.write(")")
				return
			}
		}
	}

	// Translate module function calls (fmt.println -> fmt.Println)
	if field, ok := c.Callee.(*ast.FieldExpr); ok {
		if ident, ok := field.Object.(*ast.IdentExpr); ok {
			mapped := mapModuleCall(ident.Name, field.Field)
			if mapped != "" {
				g.write(mapped)
				g.write("(")
				g.genArgList(c.Args)
				g.write(")")
				return
			}
		}
	}

	args := c.Args
	if decl := g.resolveCallDecl(c); decl != nil {
		args = g.arrangeCallArgs(c.Args, decl)
	}

	g.genExpr(c.Callee)
	g.write("(")
	g.genArgList(args)
	g.write(")")
}

func (g *Generator) genFlagCall(c *ast.CallExpr) {
	// Capture the default value as a Go expression
	var defaultBuf strings.Builder
	oldBuf := g.buf
	g.buf = defaultBuf
	// Find the default named arg
	for _, arg := range c.Args {
		if named, ok := arg.(*ast.NamedArgExpr); ok && named.Name == "default" {
			g.genExpr(named.Value)
			break
		}
	}
	defaultExpr := g.buf.String()
	g.buf = oldBuf

	g.needsFlag = true
	g.flagDecls = append(g.flagDecls, flagDecl{
		name:       c.FlagName,
		goType:     c.FlagGoType,
		defaultVal: defaultExpr,
	})

	// Emit the dereference expression inline
	g.write("*_zflag_" + c.FlagName)
}

func (g *Generator) resolveCallDecl(c *ast.CallExpr) *ast.FnDecl {
	if c.ResolvedFunc != "" {
		if decl, ok := g.funcs[c.ResolvedFunc]; ok {
			return decl
		}
	}
	if ident, ok := c.Callee.(*ast.IdentExpr); ok {
		if decl, ok := g.funcs[ident.Name]; ok {
			return decl
		}
	}
	return nil
}

func (g *Generator) arrangeCallArgs(args []ast.Node, decl *ast.FnDecl) []ast.Node {
	ordered := make([]ast.Node, len(decl.Params))
	nextPositional := 0
	for _, arg := range args {
		if named, ok := arg.(*ast.NamedArgExpr); ok {
			for i, p := range decl.Params {
				if p.Name == named.Name {
					ordered[i] = named.Value
					break
				}
			}
			continue
		}
		if nextPositional < len(ordered) {
			ordered[nextPositional] = arg
		}
		nextPositional++
	}
	// Fill missing optional parameters with defaults.
	for i, p := range decl.Params {
		if ordered[i] == nil {
			ordered[i] = p.Default
		}
	}
	// Keep only the leading arguments that are present.
	result := make([]ast.Node, 0, len(ordered))
	for _, arg := range ordered {
		if arg == nil {
			break
		}
		result = append(result, arg)
	}
	return result
}

func (g *Generator) genArgList(args []ast.Node) {
	for i, arg := range args {
		if i > 0 {
			g.write(", ")
		}
		g.genExpr(arg)
	}
}

func (g *Generator) genArrayLit(a *ast.ArrayLitExpr) {
	if a.GoType != "" {
		g.write(a.GoType + "{")
	} else {
		g.write("[]interface{}{")
	}
	for i, elem := range a.Elements {
		if i > 0 {
			g.write(", ")
		}
		g.genExpr(elem)
	}
	g.write("}")
}

func (g *Generator) genObjConstructor(decl *ast.ObjDecl, args []ast.Node) {
	// Build map of provided named args
	provided := make(map[string]ast.Node)
	for _, arg := range args {
		if named, ok := arg.(*ast.NamedArgExpr); ok {
			provided[named.Name] = named.Value
		}
	}

	g.write("&" + decl.Name + "{")
	first := true
	for _, f := range decl.Fields {
		val, ok := provided[f.Name]
		if !ok {
			val = f.Default
		}
		if val == nil {
			continue
		}
		if !first {
			g.write(", ")
		}
		first = false
		g.write(exportName(f.Name) + ": ")
		g.genExpr(val)
	}
	g.write("}")
}

func (g *Generator) genClosureExpr(c *ast.ClosureExpr) {
	g.write("func(")
	g.write(c.GoParams)
	g.write(")")
	if c.GoReturn != "" && c.GoReturn != "interface{}" {
		g.write(" " + c.GoReturn)
	}
	if _, isBlock := c.Body.(*ast.Block); isBlock {
		g.write(" {\n")
		g.indent++
		block := c.Body.(*ast.Block)
		for _, stmt := range block.Stmts {
			g.genNode(stmt)
		}
		g.indent--
		g.writeIndent()
		g.write("}")
	} else {
		g.write(" { return ")
		g.genExpr(c.Body)
		g.write(" }")
	}
}

func (g *Generator) genInterpString(s *ast.InterpStringExpr) {
	// Build fmt.Sprintf("...%v...", arg1, arg2, ...)
	var format strings.Builder
	var args []ast.Node
	for _, part := range s.Parts {
		if part.IsExpr {
			format.WriteString("%v")
			args = append(args, part.Expr)
		} else {
			// Escape % as %% and quote special chars for Go string literal
			for _, ch := range part.Lit {
				if ch == '%' {
					format.WriteString("%%")
				} else {
					format.WriteRune(ch)
				}
			}
		}
	}
	g.write("fmt.Sprintf(")
	g.write(fmt.Sprintf("%q", format.String()))
	for _, arg := range args {
		g.write(", ")
		g.genExpr(arg)
	}
	g.write(")")
}

func (g *Generator) genIfExpr(e *ast.IfExpr) {
	goType := e.GoType
	if goType == "" {
		goType = "interface{}"
	}
	g.write("func() " + goType + " { ")
	g.genIfExprBody(e)
	g.write(" }()")
}

func (g *Generator) genIfExprBody(e *ast.IfExpr) {
	g.write("if ")
	g.genExpr(e.Condition)
	g.write(" { return ")
	g.genExpr(e.Then)
	g.write(" }")
	// else branch
	if inner, ok := e.Else.(*ast.IfExpr); ok {
		g.write(" else ")
		g.genIfExprBody(inner)
	} else {
		g.write(" else { return ")
		g.genExpr(e.Else)
		g.write(" }")
	}
}

func (g *Generator) genMatchExpr(e *ast.MatchExpr) {
	goType := e.GoType
	if goType == "" {
		goType = "interface{}"
	}
	g.write("func() " + goType + " {\n")
	g.write("switch ")
	g.genExpr(e.Subject)
	g.write(" {\n")
	for _, arm := range e.Arms {
		if ident, ok := arm.Pattern.(*ast.IdentExpr); ok && ident.Name == "_" {
			g.write("default:\n")
		} else {
			g.write("case ")
			g.genExpr(arm.Pattern)
			g.write(":\n")
		}
		g.write("return ")
		g.genExpr(arm.Value)
		g.write("\n")
	}
	g.write("}\n")
	g.write("return *new(" + goType + ")\n")
	g.write("}()")
}

// ---------- Helpers ----------

func goOp(op token.Type) string {
	switch op {
	case token.Plus:
		return "+"
	case token.Minus:
		return "-"
	case token.Star:
		return "*"
	case token.Slash:
		return "/"
	case token.Percent:
		return "%"
	case token.Eq:
		return "=="
	case token.Neq:
		return "!="
	case token.Lt:
		return "<"
	case token.Gt:
		return ">"
	case token.Lte:
		return "<="
	case token.Gte:
		return ">="
	case token.And:
		return "&&"
	case token.Or:
		return "||"
	case token.Not:
		return "!"
	default:
		return "?"
	}
}

func (g *Generator) genType(t *ast.TypeExpr) string {
	return genTypeExprResolved(t, g.typeAliases)
}

func genTypeExprResolved(t *ast.TypeExpr, aliases map[string]*ast.TypeExpr) string {
	if t == nil {
		return ""
	}
	// Resolve alias: plain name with no params/flags
	if !t.IsSlice && !t.IsHashmap && !t.IsTuple && !t.IsArray && len(t.Params) == 0 {
		if alias, ok := aliases[t.Name]; ok {
			return genTypeExprResolved(alias, aliases)
		}
	}
	if t.IsHashmap && len(t.Params) == 2 {
		return "map[" + genTypeExprResolved(t.Params[0], aliases) + "]" + genTypeExprResolved(t.Params[1], aliases)
	}
	if t.IsTuple {
		return "[]interface{}"
	}
	if t.IsSlice && len(t.Params) > 0 {
		return "[]" + genTypeExprResolved(t.Params[0], aliases)
	}
	return mapTypeName(t.Name)
}

func mapTypeName(name string) string {
	switch name {
	case "int":
		return "int"
	case "i8":
		return "int8"
	case "i16":
		return "int16"
	case "i32":
		return "int32"
	case "i64":
		return "int64"
	case "u8":
		return "uint8"
	case "u16":
		return "uint16"
	case "u32":
		return "uint32"
	case "u64":
		return "uint64"
	case "f32":
		return "float32"
	case "f64":
		return "float64"
	case "bool":
		return "bool"
	case "str":
		return "string"
	case "byte":
		return "byte"
	default:
		// User-defined obj types are always pointers
		return "*" + name
	}
}

func exportName(name string) string {
	if len(name) == 0 {
		return name
	}
	return strings.ToUpper(name[:1]) + name[1:]
}

func goFieldName(field string, _ ast.Node) string {
	// All field/method access on objects gets capitalized for Go export
	return exportName(field)
}

func mapModuleCall(module, fn string) string {
	switch module {
	case "fmt":
		switch fn {
		case "println":
			return "fmt.Println"
		case "print":
			return "fmt.Print"
		case "printf":
			return "fmt.Printf"
		case "sprintf", "format":
			return "fmt.Sprintf"
		}
	case "math":
		switch fn {
		case "sqrt":
			return "math.Sqrt"
		case "abs":
			return "math.Abs"
		case "pow":
			return "math.Pow"
		case "min":
			return "math.Min"
		case "max":
			return "math.Max"
		case "pi":
			return "math.Pi"
		}
	case "os":
		switch fn {
		case "exit":
			return "os.Exit"
		case "args":
			return "os.Args"
		}
	case "strings", "str":
		switch fn {
		case "split":
			return "strings.Split"
		case "join":
			return "strings.Join"
		case "contains":
			return "strings.Contains"
		case "replace":
			return "strings.ReplaceAll"
		case "trim":
			return "strings.TrimSpace"
		}
	}
	return ""
}
