package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/Jonak-Adipta-Kalita/JAK-Programming-Language/ast"
	"github.com/Jonak-Adipta-Kalita/JAK-Programming-Language/token"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
	InvokeMethod(method string, args ...Object) Object
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
func (b *Boolean) InvokeMethod(method string, args ...Object) Object {
	return nil
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) InvokeMethod(method string, args ...Object) Object {
	return nil
}

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }
func (n *Null) InvokeMethod(method string, args ...Object) Object {
	return nil
}

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }
func (rv *ReturnValue) InvokeMethod(method string, args ...Object) Object {
	return nil
}

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
func (f *Function) InvokeMethod(method string, args ...Object) Object {
	return nil
}

type String struct {
	Value  string
	offset int
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
func (s *String) Reset()           { s.offset = 0 }
func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }
func (s *String) InvokeMethod(method string, args ...Object) Object {
	switch method {
	case "count":
		if len(args) < 1 {
			return &Error{Message: "Missing argument to count()!"}
		}
		arg := args[0].Inspect()
		return &Integer{Value: int64(strings.Count(s.Value, arg))}
	case "find":
		if len(args) < 1 {
			return &Error{Message: "Missing argument to find()!"}
		}
		arg := args[0].Inspect()
		return &Integer{Value: int64(strings.Index(s.Value, arg))}
	case "replace":
		if len(args) < 2 {
			return &Error{Message: "Missing arguments to replace()!"}
		}
		oldS := args[0].Inspect()
		newS := args[1].Inspect()
		return &String{Value: strings.Replace(s.Value, oldS, newS, -1)}
	case "reverse":
		out := make([]rune, utf8.RuneCountInString(s.Value))
		i := len(out)
		for _, c := range s.Value {
			i--
			out[i] = c
		}
		return &String{Value: string(out)}
	case "split":
		sep := " "
		if len(args) >= 1 {
			sep = args[0].(*String).Value
		}
		fields := strings.Split(s.Value, sep)
		l := len(fields)
		result := make([]Object, l)
		for i, txt := range fields {
			result[i] = &String{Value: txt}
		}
		return &Array{Elements: result}
	case "trim":
		return &String{Value: strings.TrimSpace(s.Value)}
	case "toLower":
		return &String{Value: strings.ToLower(s.Value)}
	case "toUpper":
		return &String{Value: strings.ToUpper(s.Value)}
	case "toTitle":
		return &String{Value: strings.ToTitle(s.Value)}
	default:
		return nil
	}
}

type Builtin struct {
	Fn func(token token.Token, args ...Object) Object
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }
func (b *Builtin) InvokeMethod(method string, args ...Object) Object {
	return nil
}

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
func (ao *Array) InvokeMethod(method string, args ...Object) Object {
	switch method {
	case "find":
		if len(args) < 1 {
			return &Error{Message: "Missing argument to find()!"}
		}

		arg := args[0].Inspect()
		result := -1
		for idx, entry := range ao.Elements {
			if entry.Inspect() == arg {
				result = idx
				break
			}
		}
		return &Integer{Value: int64(result)}
	case "append":
		if len(args) < 1 {
			return &Error{Message: "Missing argument to append()!"}
		}

		ao.Elements = append(ao.Elements, args[0])
		return &Null{}
	case "detach":
		if len(args) < 2 {
			return &Error{Message: "Missing argument to append()!"}
		}
		if args[0].Type() != INTEGER_OBJ {
			return &Error{Message: "First argument to detach() must be an integer!"}
		}

		idx := args[0].(*Integer).Value
		if idx < 0 || idx >= int64(len(ao.Elements)) {
			return &Error{Message: "Index out of range!"}
		}

		ao.Elements = append(ao.Elements[:idx], ao.Elements[idx+1:]...)
		return &Null{}
	default:
		return nil
	}
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
func (h *Hash) InvokeMethod(method string, args ...Object) Object {
	switch method {
	case "keys":
		ents := len(h.Pairs)
		array := make([]Object, ents)

		i := 0
		for _, ent := range h.Pairs {
			array[i] = ent.Key
			i++
		}

		return &Array{Elements: array}
	case "values":
		pairs := []Object{}

		for _, pair := range h.Pairs {
			pairs = append(pairs, pair.Value)
		}

		return &Array{Elements: pairs}
	default:
		return nil
	}
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
func (q *Quote) InvokeMethod(method string, args ...Object) Object {
	return nil
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
func (q *Macro) InvokeMethod(method string, args ...Object) Object {
	return nil
}

type Float struct {
	Value float64
}

func (f *Float) Inspect() string  { return strconv.FormatFloat(f.Value, 'f', -1, 64) }
func (f *Float) Type() ObjectType { return FLOAT_OBJ }
func (f *Float) InvokeMethod(method string, args ...Object) Object {
	return nil
}

type File struct {
	File *os.File
}

func (f *File) Inspect() string  { return f.File.Name() }
func (f *File) Type() ObjectType { return FILE_OBJ }
func (f *File) InvokeMethod(method string, args ...Object) Object {
	switch method {
	case "open":
		if len(args) < 2 {
			return &Error{Message: "Missing argument to open()!"}
		}
		if args[0].Type() != STRING_OBJ {
			return &Error{Message: "First argument to open() must be a string!"}
		}
		if args[1].Type() != STRING_OBJ {
			return &Error{Message: "Second argument to open() must be a string!"}
		}

		filename := args[0].(*String).Value
		mode := args[1].(*String).Value
		var flag int
		switch mode {
		case "r":
			flag = os.O_RDONLY
		case "w":
			flag = os.O_WRONLY | os.O_CREATE
		case "a":
			flag = os.O_WRONLY | os.O_APPEND | os.O_CREATE
		case "rw":
			flag = os.O_RDWR | os.O_CREATE
		case "ra":
			flag = os.O_RDWR | os.O_APPEND | os.O_CREATE
		default:
			return &Error{Message: fmt.Sprintf("Invalid mode for open()! Mode given %s", mode)}
		}
		openedFile, err := os.OpenFile(filename, flag, 0644)
		if err != nil {
			return &Error{Message: fmt.Sprintf("Could not open file: %s", filename)}
		}
		return &File{File: openedFile}
	case "close":
		err := f.File.Close()

		if err != nil {
			return &Error{Message: fmt.Sprintf("Could not close file: %s", f.File.Name())}
		}

		return &Null{}
	case "read":
		content, err := ioutil.ReadAll(f.File)

		if err != nil {
			return &Error{Message: fmt.Sprintf("Could not read file: %s", f.File.Name())}
		}

		return &String{Value: string(content)}
	case "write":
		if len(args) < 1 {
			return &Error{Message: "Missing argument to write()!"}
		}
		if args[0].Type() != STRING_OBJ {
			return &Error{Message: "First argument to write() must be a string!"}
		}

		content := args[0].(*String).Value
		_, err := f.File.WriteString(content)

		if err != nil {
			return &Error{Message: fmt.Sprintf("Could not write to file: %s", f.File.Name())}
		}

		return &Null{}
	case "readlines":
		content, err := ioutil.ReadAll(f.File)
		if err != nil {
			return &Error{Message: fmt.Sprintf("Could not read file: %s", f.File.Name())}
		}

		lines := strings.Split(string(content), "\n")
		array := make([]Object, len(lines))
		for i, line := range lines {
			array[i] = &String{Value: line}
		}

		return &Array{Elements: array}
	default:
		return nil
	}
}
