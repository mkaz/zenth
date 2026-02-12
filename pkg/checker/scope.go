package checker

// Symbol represents a declared name in a scope.
type Symbol struct {
	Name    string
	Type    ZType
	Mutable bool // true for var, false for let/const
	IsConst bool
}

// Scope represents a lexical scope.
type Scope struct {
	parent  *Scope
	symbols map[string]*Symbol
}

// NewScope creates a new scope with the given parent.
func NewScope(parent *Scope) *Scope {
	return &Scope{parent: parent, symbols: make(map[string]*Symbol)}
}

// Define adds a symbol to this scope.
func (s *Scope) Define(sym *Symbol) bool {
	if _, exists := s.symbols[sym.Name]; exists {
		return false // already defined in this scope
	}
	s.symbols[sym.Name] = sym
	return true
}

// Lookup looks up a name, walking up parent scopes.
func (s *Scope) Lookup(name string) *Symbol {
	if sym, ok := s.symbols[name]; ok {
		return sym
	}
	if s.parent != nil {
		return s.parent.Lookup(name)
	}
	return nil
}
