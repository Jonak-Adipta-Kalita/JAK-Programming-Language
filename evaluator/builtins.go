package evaluator

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/Jonak-Adipta-Kalita/JAK-Programming-Language/object"
)

var builtins = map[string]*object.Builtin{
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
	"values": {
		Fn: func(file string, line int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", file, line,
					len(args))
			}
			if args[0].Type() != object.HASH_OBJ {
				return newError("argument to `values` must be HASH, got %s", file, line,
					args[0].Type())
			}
			hash := args[0].(*object.Hash)
			pairs := []object.Object{}

			for _, pair := range hash.Pairs {
				pairs = append(pairs, pair.Value)
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
				return &object.Integer{Value: int64(utf8.RuneCountInString(arg.Value))}
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
	"format": {
		Fn: func(file string, line int, args ...object.Object) object.Object {
			if len(args) < 1 {
				return newError("wrong number of arguments. got=%d, want=1(at least)", file, line, len(args))
			}
			if args[0].Type() != object.STRING_OBJ {
				return newError("argument to `format` must be STRING, got %s", file, line, args[0].Type())
			}

			formatString := args[0].(*object.String).Value
			formattedString := formatString

			// Replace f-string-style expressions with corresponding argument values
			for i, arg := range args[1:] {
				replacement := "{" + strconv.Itoa(i) + "}"
				formattedString = strings.ReplaceAll(formattedString, replacement, arg.Inspect())
			}

			return &object.String{Value: formattedString}
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
			case *object.Float:
				return &object.Integer{Value: int64(arg.Value)}
			default:
				return newError("argument to `int` not supported, got %s", file, line, args[0].Type())
			}
		},
	},
	"float": {
		Fn: func(file string, line int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", file, line, len(args))
			}
			switch arg := args[0].(type) {
			case *object.Integer:
				return &object.Float{Value: float64(arg.Value)}
			case *object.Float:
				return arg
			case *object.String:
				float, err := strconv.ParseFloat(arg.Value, 64)
				if err != nil {
					return newError("could not convert string to float", file, line)
				}
				return &object.Float{Value: float}
			default:
				return newError("argument to `float` not supported, got %s", file, line, args[0].Type())
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
			case *object.Float:
				return &object.String{Value: strconv.FormatFloat(arg.Value, 'f', -1, 64)}
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

	"open": {
		Fn: func(file string, line int, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", file, line, len(args))
			}
			if args[0].Type() != object.STRING_OBJ {
				return newError("argument to `open` must be STRING, got %s", file, line, args[0].Type())
			}
			if args[1].Type() != object.STRING_OBJ {
				return newError("argument to `open` must be STRING, got %s", file, line, args[1].Type())
			}
			filename := args[0].(*object.String).Value
			mode := args[1].(*object.String).Value
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
				return newError("invalid mode %s", file, line, mode)
			}
			openedFile, err := os.OpenFile(filename, flag, 0644)
			if err != nil {
				return newError("could not open file %s", file, line, filename)
			}
			return &object.File{File: openedFile}
		},
	},
	"close": {
		Fn: func(file string, line int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", file, line, len(args))
			}
			if args[0].Type() != object.FILE_OBJ {
				return newError("argument to `close` must be FILE, got %s", file, line, args[0].Type())
			}
			fileObj := args[0].(*object.File)
			err := fileObj.File.Close()
			if err != nil {
				return newError("could not close file %s", file, line, fileObj.File.Name())
			}
			return NULL
		},
	},
	"read": {
		Fn: func(file string, line int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", file, line, len(args))
			}
			if args[0].Type() != object.FILE_OBJ {
				return newError("argument to `read` must be FILE, got %s", file, line, args[0].Type())
			}
			fileObj := args[0].(*object.File)
			content, err := ioutil.ReadAll(fileObj.File)
			if err != nil {
				return newError("could not read file %s", file, line, fileObj.File.Name())
			}
			return &object.String{Value: string(content)}
		},
	},
	"write": {
		Fn: func(file string, line int, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", file, line, len(args))
			}
			if args[0].Type() != object.FILE_OBJ {
				return newError("argument to `write` must be FILE, got %s", file, line, args[0].Type())
			}
			if args[1].Type() != object.STRING_OBJ {
				return newError("argument to `write` must be STRING, got %s", file, line, args[1].Type())
			}
			fileObj := args[0].(*object.File)
			_, err := fileObj.File.WriteString(args[1].(*object.String).Value)
			if err != nil {
				return newError("could not write to file %s", file, line, fileObj.File.Name())
			}
			return NULL
		},
	},
	"readlines": {
		Fn: func(file string, line int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", file, line, len(args))
			}
			if args[0].Type() != object.FILE_OBJ {
				return newError("argument to `readlines` must be FILE, got %s", file, line, args[0].Type())
			}
			fileObj := args[0].(*object.File)
			content, err := ioutil.ReadAll(fileObj.File)
			if err != nil {
				return newError("could not read file %s", file, line, fileObj.File.Name())
			}
			var elements []object.Object
			for _, line := range strings.Split(string(content), "\n") {
				elements = append(elements, &object.String{Value: line})
			}

			return &object.Array{Elements: elements}
		},
	},

	"mkdir": {
		Fn: func(file string, line int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", file, line, len(args))
			}
			if args[0].Type() != object.STRING_OBJ {
				return newError("argument to `mkdir` must be STRING, got %s", file, line, args[0].Type())
			}
			dirname := args[0].(*object.String).Value
			err := os.Mkdir(dirname, 0755)
			if err != nil {
				return newError("could not create directory %s", file, line, dirname)
			}
			return NULL
		},
	},
	"rmdir": {
		Fn: func(file string, line int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", file, line, len(args))
			}
			if args[0].Type() != object.STRING_OBJ {
				return newError("argument to `rmdir` must be STRING, got %s", file, line, args[0].Type())
			}
			dirname := args[0].(*object.String).Value
			err := os.Remove(dirname)
			if err != nil {
				return newError("could not remove directory %s", file, line, dirname)
			}
			return NULL
		},
	},
	"mkfile": {
		Fn: func(file string, line int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", file, line, len(args))
			}
			if args[0].Type() != object.STRING_OBJ {
				return newError("argument to `mkfile` must be STRING, got %s", file, line, args[0].Type())
			}
			filename := args[0].(*object.String).Value
			fileObj, err := os.Create(filename)
			fileObj.Close()
			if err != nil {
				return newError("could not create file %s", file, line, filename)
			}
			return NULL
		},
	},
	"rmfile": {
		Fn: func(file string, line int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", file, line, len(args))
			}
			if args[0].Type() != object.STRING_OBJ {
				return newError("argument to `rmfile` must be STRING, got %s", file, line, args[0].Type())
			}
			filename := args[0].(*object.String).Value
			err := os.Remove(filename)
			if err != nil {
				return newError("could not remove file %s", file, line, filename)
			}
			return NULL
		},
	},
}
