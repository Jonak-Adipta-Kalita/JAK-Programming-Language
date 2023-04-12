package evaluator

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Jonak-Adipta-Kalita/JAK-Programming-Language/object"
)

var builtins = map[string]*object.Builtin{
	"first": {
		Fn: func(file string, line int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", file, line,
					len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `first` must be ARRAY, got %s", file, line,
					args[0].Type())
			}
			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}
			return NULL
		},
	},
	"last": {
		Fn: func(file string, line int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", file, line,
					len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `last` must be ARRAY, got %s", file, line,
					args[0].Type())
			}
			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				return arr.Elements[length-1]
			}
			return NULL
		},
	},
	"append": {
		Fn: func(file string, line int, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", file, line,
					len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `push` must be ARRAY, got %s", file, line,
					args[0].Type())
			}
			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			newElements := make([]object.Object, length+1)
			copy(newElements, arr.Elements)
			newElements[length] = args[1]
			return &object.Array{Elements: newElements}
		},
	},
	"detach": {
		Fn: func(file string, line int, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", file, line,
					len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `detach` must be ARRAY, got %s", file, line,
					args[0].Type())
			}
			arr := args[0].(*object.Array)
			index := args[1].(*object.Integer).Value

			length := len(arr.Elements)
			if int(index) >= length {
				return newError("index out of range", file, line)
			}
			newElements := make([]object.Object, length-1)
			copy(newElements[:index], arr.Elements[:index])
			copy(newElements[index:], arr.Elements[index+1:])
			return &object.Array{Elements: newElements}
		},
	},

	"keys": {
		Fn: func(file string, line int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", file, line,
					len(args))
			}
			if args[0].Type() != object.HASH_OBJ {
				return newError("argument to `keys` must be HASH, got %s", file, line,
					args[0].Type())
			}
			hash := args[0].(*object.Hash)
			pairs := []object.Object{}

			for _, pair := range hash.Pairs {
				pairs = append(pairs, pair.Key)
			}

			return &object.Array{Elements: pairs}
		},
	},

	"len": {
		Fn: func(file string, line int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", file, line,
					len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to `len` not supported, got %s", file, line,
					args[0].Type())
			}
		},
	},
	"print": {
		Fn: func(file string, line int, args ...object.Object) object.Object {
			for i, arg := range args {
				if i == len(args)-1 {
					fmt.Println(arg.Inspect())
				} else {
					fmt.Print(arg.Inspect())
				}
			}
			return NULL
		},
	},
	"println": {
		Fn: func(file string, line int, args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NULL
		},
	},
	"input": {
		Fn: func(file string, line int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", file, line, len(args))
			}
			var input string
			fmt.Print(args[0].Inspect())
			fmt.Scanln(&input)
			return &object.String{Value: input}
		},
	},
	"range": {
		Fn: func(file string, line int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", file, line, len(args))
			}
			if args[0].Type() != object.INTEGER_OBJ {
				return newError("argument to `range` must be INTEGER, got %s", file, line, args[0].Type())
			}
			integer := args[0].(*object.Integer)
			var elements []object.Object
			for i := 0; i < int(integer.Value); i++ {
				elements = append(elements, &object.Integer{Value: int64(i)})
			}
			return &object.Array{Elements: elements}
		},
	},
	"typeof": {
		Fn: func(file string, line int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", file, line, len(args))
			}
			return &object.String{Value: string(args[0].Type())}
		},
	},
	"exit": {
		Fn: func(file string, line int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", file, line, len(args))
			}
			if args[0].Type() != object.INTEGER_OBJ {
				return newError("argument to `exit` must be INTEGER, got %s", file, line, args[0].Type())
			}
			integer := args[0].(*object.Integer)
			os.Exit(int(integer.Value))
			return NULL
		},
	},
	"int": {
		Fn: func(file string, line int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", file, line, len(args))
			}
			switch arg := args[0].(type) {
			case *object.Integer:
				return arg
			case *object.String:
				integer, err := strconv.ParseInt(arg.Value, 10, 64)
				if err != nil {
					return newError("could not convert string to integer", file, line)
				}
				return &object.Integer{Value: integer}
			default:
				return newError("argument to `int` not supported, got %s", file, line, args[0].Type())
			}
		},
	},
	"str": {
		Fn: func(file string, line int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", file, line, len(args))
			}
			switch arg := args[0].(type) {
			case *object.Integer:
				return &object.String{Value: strconv.FormatInt(arg.Value, 10)}
			case *object.String:
				return arg
			default:
				return newError("argument to `str` not supported, got %s", file, line, args[0].Type())
			}
		},
	},
	"bool": {
		Fn: func(file string, line int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", file, line, len(args))
			}
			switch arg := args[0].(type) {
			case *object.Boolean:
				return arg
			case *object.String:
				if arg.Value == "true" {
					return TRUE
				}
				if arg.Value == "false" {
					return FALSE
				}
				return newError("could not convert string to boolean", file, line)
			default:
				return newError("argument to `bool` not supported, got %s", file, line, args[0].Type())
			}
		},
	},
}
