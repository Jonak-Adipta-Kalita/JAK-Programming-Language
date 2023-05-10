package evaluator

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/Jonak-Adipta-Kalita/JAK-Programming-Language/object"
	"github.com/Jonak-Adipta-Kalita/JAK-Programming-Language/token"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(token token.Token, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", token,
					len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			case *object.String:
				return &object.Integer{Value: int64(utf8.RuneCountInString(arg.Value))}
			default:
				return newError("argument to `len` not supported, got %s", token,
					args[0].Type())
			}
		},
	},
	"print": {
		Fn: func(token token.Token, args ...object.Object) object.Object {
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
		Fn: func(token token.Token, args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NULL
		},
	},
	"input": {
		Fn: func(token token.Token, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", token, len(args))
			}
			var input string
			fmt.Print(args[0].Inspect())
			fmt.Scanln(&input)
			return &object.String{Value: input}
		},
	},
	"format": {
		Fn: func(token token.Token, args ...object.Object) object.Object {
			if len(args) < 1 {
				return newError("wrong number of arguments. got=%d, want=1(at least)", token, len(args))
			}
			if args[0].Type() != object.STRING_OBJ {
				return newError("argument to `format` must be STRING, got %s", token, args[0].Type())
			}

			formatString := args[0].(*object.String).Value
			formattedString := formatString

			for i, arg := range args[1:] {
				replacement := "{" + strconv.Itoa(i) + "}"
				formattedString = strings.ReplaceAll(formattedString, replacement, arg.Inspect())
			}

			return &object.String{Value: formattedString}
		},
	},
	"range": {
		Fn: func(token token.Token, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", token, len(args))
			}
			if args[0].Type() != object.INTEGER_OBJ {
				return newError("argument to `range` must be INTEGER, got %s", token, args[0].Type())
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
		Fn: func(token token.Token, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", token, len(args))
			}
			return &object.String{Value: string(args[0].Type())}
		},
	},
	"exit": {
		Fn: func(token token.Token, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", token, len(args))
			}
			if args[0].Type() != object.INTEGER_OBJ {
				return newError("argument to `exit` must be INTEGER, got %s", token, args[0].Type())
			}
			integer := args[0].(*object.Integer)
			os.Exit(int(integer.Value))
			return NULL
		},
	},
	"int": {
		Fn: func(token token.Token, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", token, len(args))
			}
			switch arg := args[0].(type) {
			case *object.Integer:
				return arg
			case *object.String:
				integer, err := strconv.ParseInt(arg.Value, 10, 64)
				if err != nil {
					return newError("could not convert string to integer", token)
				}
				return &object.Integer{Value: integer}
			case *object.Float:
				return &object.Integer{Value: int64(arg.Value)}
			default:
				return newError("argument to `int` not supported, got %s", token, args[0].Type())
			}
		},
	},
	"float": {
		Fn: func(token token.Token, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", token, len(args))
			}
			switch arg := args[0].(type) {
			case *object.Integer:
				return &object.Float{Value: float64(arg.Value)}
			case *object.Float:
				return arg
			case *object.String:
				float, err := strconv.ParseFloat(arg.Value, 64)
				if err != nil {
					return newError("could not convert string to float", token)
				}
				return &object.Float{Value: float}
			default:
				return newError("argument to `float` not supported, got %s", token, args[0].Type())
			}
		},
	},
	"str": {
		Fn: func(token token.Token, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", token, len(args))
			}
			switch arg := args[0].(type) {
			case *object.Integer:
				return &object.String{Value: strconv.FormatInt(arg.Value, 10)}
			case *object.Float:
				return &object.String{Value: strconv.FormatFloat(arg.Value, 'f', -1, 64)}
			case *object.String:
				return arg
			default:
				return newError("argument to `str` not supported, got %s", token, args[0].Type())
			}
		},
	},
	"bool": {
		Fn: func(token token.Token, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", token, len(args))
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
				return newError("could not convert string to boolean", token)
			default:
				return newError("argument to `bool` not supported, got %s", token, args[0].Type())
			}
		},
	},
	"mkdir": {
		Fn: func(token token.Token, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", token, len(args))
			}
			if args[0].Type() != object.STRING_OBJ {
				return newError("argument to `mkdir` must be STRING, got %s", token, args[0].Type())
			}
			dirname := args[0].(*object.String).Value
			err := os.Mkdir(dirname, 0755)
			if err != nil {
				return newError("could not create directory %s", token, dirname)
			}
			return NULL
		},
	},
	"rmdir": {
		Fn: func(token token.Token, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", token, len(args))
			}
			if args[0].Type() != object.STRING_OBJ {
				return newError("argument to `rmdir` must be STRING, got %s", token, args[0].Type())
			}
			dirname := args[0].(*object.String).Value
			err := os.Remove(dirname)
			if err != nil {
				return newError("could not remove directory %s", token, dirname)
			}
			return NULL
		},
	},
	"mkfile": {
		Fn: func(token token.Token, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", token, len(args))
			}
			if args[0].Type() != object.STRING_OBJ {
				return newError("argument to `mkfile` must be STRING, got %s", token, args[0].Type())
			}
			filename := args[0].(*object.String).Value
			fileObj, err := os.Create(filename)
			fileObj.Close()
			if err != nil {
				return newError("could not create file %s", token, filename)
			}
			return NULL
		},
	},
	"rmfile": {
		Fn: func(token token.Token, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", token, len(args))
			}
			if args[0].Type() != object.STRING_OBJ {
				return newError("argument to `rmfile` must be STRING, got %s", token, args[0].Type())
			}
			filename := args[0].(*object.String).Value
			err := os.Remove(filename)
			if err != nil {
				return newError("could not remove file %s", token, filename)
			}
			return NULL
		},
	},
}
