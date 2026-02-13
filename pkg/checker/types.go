package checker

import "fmt"

// ZType represents a Zenth type.
type ZType interface {
	String() string
	Equals(other ZType) bool
	typeMarker()
}

// Built-in types
type BuiltinType struct {
	Name string
}

func (b *BuiltinType) String() string         { return b.Name }
func (b *BuiltinType) Equals(other ZType) bool {
	if o, ok := other.(*BuiltinType); ok {
		return b.Name == o.Name
	}
	return false
}
func (b *BuiltinType) typeMarker() {}

var (
	TypeInt    = &BuiltinType{"int"}
	TypeI8     = &BuiltinType{"i8"}
	TypeI16    = &BuiltinType{"i16"}
	TypeI32    = &BuiltinType{"i32"}
	TypeI64    = &BuiltinType{"i64"}
	TypeU8     = &BuiltinType{"u8"}
	TypeU16    = &BuiltinType{"u16"}
	TypeU32    = &BuiltinType{"u32"}
	TypeU64    = &BuiltinType{"u64"}
	TypeF32    = &BuiltinType{"f32"}
	TypeF64    = &BuiltinType{"f64"}
	TypeBool   = &BuiltinType{"bool"}
	TypeStr    = &BuiltinType{"str"}
	TypeByte   = &BuiltinType{"byte"}
	TypeVoid   = &BuiltinType{"void"}
	TypeError  = &BuiltinType{"error"}
	TypeNil    = &BuiltinType{"nil"}
)

// ObjType represents a user-defined obj type.
type ObjType struct {
	Name   string
	Fields map[string]ZType
}

func (s *ObjType) String() string { return s.Name }
func (s *ObjType) Equals(other ZType) bool {
	if o, ok := other.(*ObjType); ok {
		return s.Name == o.Name
	}
	return false
}
func (s *ObjType) typeMarker() {}

// SliceType represents a []T type.
type SliceType struct {
	Elem ZType
}

func (s *SliceType) String() string { return fmt.Sprintf("[]%s", s.Elem) }
func (s *SliceType) Equals(other ZType) bool {
	if o, ok := other.(*SliceType); ok {
		return s.Elem.Equals(o.Elem)
	}
	return false
}
func (s *SliceType) typeMarker() {}

// FuncType represents a function type (for passing functions as values).
type FuncType struct {
	Params  []ZType
	Returns ZType
}

func (f *FuncType) String() string {
	return fmt.Sprintf("fn(...) -> %s", f.Returns)
}
func (f *FuncType) Equals(other ZType) bool {
	o, ok := other.(*FuncType)
	if !ok {
		return false
	}
	if len(f.Params) != len(o.Params) {
		return false
	}
	for i := range f.Params {
		if !f.Params[i].Equals(o.Params[i]) {
			return false
		}
	}
	return f.Returns.Equals(o.Returns)
}
func (f *FuncType) typeMarker() {}

// LookupBuiltinType maps type names to ZType.
func LookupBuiltinType(name string) ZType {
	switch name {
	case "int":
		return TypeInt
	case "i8":
		return TypeI8
	case "i16":
		return TypeI16
	case "i32":
		return TypeI32
	case "i64":
		return TypeI64
	case "u8":
		return TypeU8
	case "u16":
		return TypeU16
	case "u32":
		return TypeU32
	case "u64":
		return TypeU64
	case "f32":
		return TypeF32
	case "f64":
		return TypeF64
	case "bool":
		return TypeBool
	case "str":
		return TypeStr
	case "byte":
		return TypeByte
	case "error":
		return TypeError
	default:
		return nil
	}
}

// IsNumeric returns true if the type is a numeric type.
func IsNumeric(t ZType) bool {
	b, ok := t.(*BuiltinType)
	if !ok {
		return false
	}
	switch b.Name {
	case "int", "i8", "i16", "i32", "i64", "u8", "u16", "u32", "u64", "f32", "f64":
		return true
	}
	return false
}

// IsInteger returns true if the type is an integer type.
func IsInteger(t ZType) bool {
	b, ok := t.(*BuiltinType)
	if !ok {
		return false
	}
	switch b.Name {
	case "int", "i8", "i16", "i32", "i64", "u8", "u16", "u32", "u64":
		return true
	}
	return false
}

// IsFloat returns true if the type is a floating-point type.
func IsFloat(t ZType) bool {
	b, ok := t.(*BuiltinType)
	if !ok {
		return false
	}
	return b.Name == "f32" || b.Name == "f64"
}

// PromoteNumeric returns the promoted type when mixing numeric types.
// Returns nil if promotion is not possible.
func PromoteNumeric(a, b ZType) ZType {
	if !IsNumeric(a) || !IsNumeric(b) {
		return nil
	}
	if a.Equals(b) {
		return a
	}
	aFloat, bFloat := IsFloat(a), IsFloat(b)
	// Float + integer → float type wins
	if aFloat && !bFloat {
		return a
	}
	if !aFloat && bFloat {
		return b
	}
	// Both floats, different sizes → f64
	if aFloat && bFloat {
		return TypeF64
	}
	// Both integers, different sizes → not auto-promoted
	return nil
}
