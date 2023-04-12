package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/Jonak-Adipta-Kalita/JAK-Programming-Language/evaluator"
	"github.com/Jonak-Adipta-Kalita/JAK-Programming-Language/file"
	"github.com/Jonak-Adipta-Kalita/JAK-Programming-Language/lexer"
	"github.com/Jonak-Adipta-Kalita/JAK-Programming-Language/object"
	"github.com/Jonak-Adipta-Kalita/JAK-Programming-Language/parser"
)

const PROMPT = ">>> "

func Start(in io.Reader, out io.Writer) {
	file.SetFileName("STDIN")
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			evaluator.PrintParserErrors(out, p.Errors())
			continue
		}

		evaluator.Eval(program, env)
	}
}
