package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/Jonak-Adipta-Kalita/JAK-Programming-Language/ast"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

const (
	INTEGER_OBJ = "INTEGER"
	FLOAT_OBJ   = "FLOAT"
	STRING_OBJ  = "STRING"
	BOOLEAN_OBJ = "BOOLEAN"
	ARRAY_OBJ   = "ARRAY"
	HASH_OBJ    = "HASH"

	NULL_OBJ  = "NULL"
	ERROR_OBJ = "ERROR"

	RETURN_VALUE_OBJ = "RETURN_VALUE"
	FUNCTION_OBJ     = "FUNCTION"

	BUILTIN_OBJ = "BUILTIN"
	QUOTE_OBJ   = "QUOTE"
	MACRO_OBJ   = "MACRO"
	FILE_OBJ    = "FILE"
)

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("func")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")
	return out.String()
}

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

type String struct {
	Value  string
	offset int
}

func (s *String) Reset() {
	s.offset = 0
}
func (s *String) Next() (Object, Object, bool) {

	if s.offset < utf8.RuneCountInString(s.Value) {
		s.offset++
		chars := []rune(s.Value)
		val := &String{Value: string(chars[s.offset-1])}

		return val, &Integer{Value: int64(s.offset - 1)}, true
	}

	return nil, &Integer{Value: 0}, false
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }

type BuiltinFunction func(file string, line int, args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }

type Array struct {
	Elements []Object
	offset   int
}

func (ao *Array) Type() ObjectType { return ARRAY_OBJ }
func (ao *Array) Inspect() string {
	var out bytes.Buffer
	elements := []string{}
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}
func (ao *Array) Reset() { ao.offset = 0 }
func (ao *Array) Next() (Object, Object, bool) {
	if ao.offset < len(ao.Elements) {
		ao.offset++

		element := ao.Elements[ao.offset-1]
		return element, &Integer{Value: int64(ao.offset - 1)}, true
	}

	return nil, &Integer{Value: 0}, false
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

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
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

type HashPair struct {
	Key   Object
	Value Object
}
type Hash struct {
	Pairs  map[HashKey]HashPair
	offset int
}

func (h *Hash) Reset() {
	h.offset = 0
}
func (h *Hash) Next() (Object, Object, bool) {
	if h.offset < len(h.Pairs) {
		idx := 0

		for _, pair := range h.Pairs {
			if h.offset == idx {
				h.offset++
				return pair.Key, pair.Value, true
			}
			idx++
		}
	}

	return nil, &Integer{Value: 0}, false
}

type Hashable interface {
	HashKey() HashKey
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	var out bytes.Buffer
	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s",
			pair.Key.Inspect(), pair.Value.Inspect()))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

type Iterable interface {
	Reset()
	Next() (Object, Object, bool)
}

type Quote struct {
	Node ast.Node
}

func (q *Quote) Type() ObjectType { return QUOTE_OBJ }
func (q *Quote) Inspect() string {
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

type Float struct {
	Value float64
}

func (f *Float) Inspect() string  { return strconv.FormatFloat(f.Value, 'f', -1, 64) }
func (f *Float) Type() ObjectType { return FLOAT_OBJ }

type File struct {
	File *os.File
}

func (f *File) Inspect() string  { return f.File.Name() }
func (f *File) Type() ObjectType { return FILE_OBJ }
