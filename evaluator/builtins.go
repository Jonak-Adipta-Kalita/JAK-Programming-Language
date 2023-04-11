package evaluator

import (
	"fmt"

	"github.com/Jonak-Adipta-Kalita/JAK-Programming-Language/object"
)

var builtins = map[string]*object.Builtin{
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
	"push": {
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
}
