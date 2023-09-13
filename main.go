package main

import (
	"fmt"
	"os"

	"github.com/Jonak-Adipta-Kalita/JAK-Programming-Language/evaluator"
	"github.com/Jonak-Adipta-Kalita/JAK-Programming-Language/file"
	"github.com/Jonak-Adipta-Kalita/JAK-Programming-Language/lexer"
	"github.com/Jonak-Adipta-Kalita/JAK-Programming-Language/object"
	"github.com/Jonak-Adipta-Kalita/JAK-Programming-Language/parser"
	"github.com/Jonak-Adipta-Kalita/JAK-Programming-Language/repl"
)

func main() {
	if len(os.Args) != 2 {
		repl.Start(os.Stdin, os.Stdout)
	} else {
		filePath := os.Args[1]
		file.SetMainFileName(filePath)
		file.SetFileName(filePath)
		contents, err := os.ReadFile(filePath)

		if err != nil {
			fmt.Printf("Failure to read file '%s'. Err: %s", string(contents), err)
			return
		}

		env := object.NewEnvironment()
		macroEnv := object.NewEnvironment()

		l := lexer.New(string(contents))
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			evaluator.PrintParserErrors(os.Stdout, p.Errors())
			return
		}

		evaluator.DefineMacros(program, macroEnv)
		expanded := evaluator.ExpandMacros(program, macroEnv)

		evaluator.Eval(expanded, env)
	}
}
