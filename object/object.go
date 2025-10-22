package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"s8/ast"
	"s8/code"
)

type ObjectType string

// For the Type() method
const (
	INTERGER_OBJ     = "INTEGER"
	FLOAT_OBJ        = "FLOAT"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	STRING_OBJ       = "STRING"
	BUILTIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY"
	HASH_OBJ         = "HASH"
	IDENT_OBJ        = "IDENTIFIER"
	QUOTE_OBJ        = "QUOTE"
	MACRO_OBJ        = "MACRO"
	BREAK_OBJ        = "BREAK"
	CONTINUE_OBJ     = "CONTINUE"
	// Hold the instructions of a compiled function
	// then we pass it as a constant
	COMPILED_FUNCTION_OBJ = "COMPILED_FUNCTION"
	CLOSURE_OBJ           = "CLOSURE"
)

// Interface instead of struct
// As each value needs a different internal representation
type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

func (i *Integer) Type() ObjectType { return INTERGER_OBJ }

type Float struct {
	Value float64
}

func (f *Float) Inspect() string { return fmt.Sprintf("%f", f.Value) }

func (f *Float) Type() ObjectType { return FLOAT_OBJ }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }

func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }

func (n *Null) Inspect() string { return "null" }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }

func (rv *ReturnValue) Inspect() string { return rv.Value.Inspect() }

type Break struct{}

func (b *Break) Type() ObjectType { return BREAK_OBJ }

func (b *Break) Inspect() string { return "break" }

type Continue struct{}

func (c *Continue) Type() ObjectType { return CONTINUE_OBJ }

func (c *Continue) Inspect() string { return "continue" }

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }

func (e *Error) Inspect() string { return "ERROR: " + e.Message }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment // A function's very own environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }

func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}

	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("funk")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))

	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }

func (s *String) Inspect() string { return s.Value }

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }

func (b *Builtin) Inspect() string { return "builtin function" }

type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType { return ARRAY_OBJ }

func (a *Array) Inspect() string {
	var out bytes.Buffer

	elems := []string{}

	for _, e := range a.Elements {
		elems = append(elems, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elems, ", "))
	out.WriteString("]")

	return out.String()
}

// Help keys of the same type sharing the same hash pointing to the same memory location
// Example: {"FB": 1, "Ea": 2} and hash("FB") == hash("Ea")
// as they produce identical hash value
type HashKey struct {
	Type  ObjectType // Scope the HashKey to different object types i.e., string, integer and boolean
	Value uint64     // Original key of a hash - Needed to store hashed values of strings which need a large range
}

// Return numeric value of booleans
func (b *Boolean) HashKey() HashKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: b.Type(), Value: value}
}

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()

	// Hash the string value
	// There is a small chance that different values result n the same hash i.e. hash collision
	// We work around that with separate chaining, open addressing, etc.
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

// Responsible for generating the HashKey
type HashPair struct {
	Key   Object // Original key object
	Value Object // Value associated with the key
}

// Look like this:
//
//	hash := &Hash{
//	    Pairs: map[HashKey]HashPair{
//	        HashKey{Type: STRING_OBJ, Value: 234892348}: HashPair{
//	            Key:   StringObject("foo"),
//	            Value: StringObject("hello"),
//	        },
//	        HashKey{Type: STRING_OBJ, Value: 987654321}: HashPair{
//	            Key:   StringObject("bar"),
//	            Value: StringObject("world"),
//	        },
//	    }
//	}
//
// Why don't we change the name to HashMap here?
// Because HashMap refers to the internal Go map,
// while Hash represents the overall object for the literal (more appropriate)
type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }

func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}

	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

// Since the hash key could be of any type,
// check if the given object could be used as a hash key
type Hashable interface {
	HashKey() HashKey
}

type Quote struct {
	// When we evaluate a call to quote
	// We can prevent the argument (as a call) from being evaluated immediately
	Node ast.Node
}

func (q *Quote) Type() ObjectType { return QUOTE_OBJ }

func (q *Quote) Inspect() string {
	// Another abstraction layer in String() here
	return "QUOTE(" + q.Node.String() + ")"
}

type Macro struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (m *Macro) Type() ObjectType { return MACRO_OBJ }

func (m *Macro) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range m.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("macro")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(m.Body.String())
	out.WriteString("\n}")

	return out.String()
}

// Hold bytecode instead of AST nodes
type CompiledFunction struct {
	Instructions code.Instructions
	// The number of local bindings the function is going to create
	NumLocals     int
	NumParameters int
}

func (cf *CompiledFunction) Type() ObjectType { return COMPILED_FUNCTION_OBJ }

func (cf *CompiledFunction) Inspect() string {
	return fmt.Sprintf("CompiledFunction[%p]", cf)
}

type Closure struct {
	Fn *CompiledFunction
	// Free variables.
	// Equivalent to Env field in *object.Function.
	Free []Object
}

func (c *Closure) Type() ObjectType { return CLOSURE_OBJ }
func (c *Closure) Inspect() string  { return fmt.Sprintf("Closure[%p]", c) }
