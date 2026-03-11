package checker

import (
	"fmt"
	"strings"
)

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

func (b *BuiltinType) String() string { return b.Name }
func (b *BuiltinType) Equals(other ZType) bool {
	if o, ok := other.(*BuiltinType); ok {
		return b.Name == o.Name
	}
	return false
}
func (b *BuiltinType) typeMarker() {}

var (
	TypeInt   = &BuiltinType{"int"}
	TypeI8    = &BuiltinType{"i8"}
	TypeI16   = &BuiltinType{"i16"}
	TypeI32   = &BuiltinType{"i32"}
	TypeI64   = &BuiltinType{"i64"}
	TypeU8    = &BuiltinType{"u8"}
	TypeU16   = &BuiltinType{"u16"}
	TypeU32   = &BuiltinType{"u32"}
	TypeU64   = &BuiltinType{"u64"}
	TypeF32   = &BuiltinType{"f32"}
	TypeF64   = &BuiltinType{"f64"}
	TypeBool  = &BuiltinType{"bool"}
	TypeStr   = &BuiltinType{"str"}
	TypeByte  = &BuiltinType{"byte"}
	TypeVoid  = &BuiltinType{"void"}
	TypeError = &BuiltinType{"error"}
	TypeNil   = &BuiltinType{"nil"}
	TypeFile  = &BuiltinType{"file"}
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

// SliceType represents an array(T) type.
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

// HashmapType represents a hashmap[K]V type.
type HashmapType struct {
	Key   ZType
	Value ZType
}

func (m *HashmapType) String() string { return fmt.Sprintf("hashmap[%s]%s", m.Key, m.Value) }
func (m *HashmapType) Equals(other ZType) bool {
	if o, ok := other.(*HashmapType); ok {
		return m.Key.Equals(o.Key) && m.Value.Equals(o.Value)
	}
	return false
}
func (m *HashmapType) typeMarker() {}

// SetType represents a set(T) type.
type SetType struct {
	Elem ZType
}

func (s *SetType) String() string { return fmt.Sprintf("set(%s)", s.Elem) }
func (s *SetType) Equals(other ZType) bool {
	if o, ok := other.(*SetType); ok {
		return s.Elem.Equals(o.Elem)
	}
	return false
}
func (s *SetType) typeMarker() {}

// TupleType represents a heterogenous tuple type.
type TupleType struct {
	Elems []ZType
	Names []string // nil for positional tuples; same length as Elems when named
}

func (t *TupleType) String() string {
	if len(t.Elems) == 0 {
		return "tuple[]"
	}
	parts := make([]string, len(t.Elems))
	for i, elem := range t.Elems {
		if t.Names != nil {
			parts[i] = t.Names[i] + ": " + elem.String()
		} else {
			parts[i] = elem.String()
		}
	}
	return "tuple[" + strings.Join(parts, ", ") + "]"
}

func (t *TupleType) Equals(other ZType) bool {
	o, ok := other.(*TupleType)
	if !ok || len(t.Elems) != len(o.Elems) {
		return false
	}
	// Named vs positional are distinct types
	if (t.Names == nil) != (o.Names == nil) {
		return false
	}
	if t.Names != nil {
		for i := range t.Names {
			if t.Names[i] != o.Names[i] {
				return false
			}
		}
	}
	for i := range t.Elems {
		if !t.Elems[i].Equals(o.Elems[i]) {
			return false
		}
	}
	return true
}

func (t *TupleType) typeMarker() {}

// EnumType represents a user-defined enum type.
type EnumType struct {
	Name     string
	Variants map[string]int64 // variant name -> integer value
}

func (e *EnumType) String() string { return e.Name }
func (e *EnumType) Equals(other ZType) bool {
	if o, ok := other.(*EnumType); ok {
		return e.Name == o.Name
	}
	return false
}
func (e *EnumType) typeMarker() {}

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
