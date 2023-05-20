package object

import (
	"fmt"
	"strings"

	"github.com/Jonak-Adipta-Kalita/JAK-Programming-Language/token"
)

type Error struct {
	Message   string
	ErrorName string
	FileName  string
	Token     token.Token
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string {
	message := fmt.Sprintf("Error: `%s`", e.Message)
	if e.FileName != "<stdin>" {
		message += fmt.Sprintf("\n\tat %s: %d", e.FileName, e.Token.Line+1)
	}
	message += fmt.Sprintf("\n\t%s", stringWithArrows(e.Token))

	return message
}
func (e *Error) InvokeMethod(method string, args ...Object) Object {
	return nil
}

func stringWithArrows(charToken token.Token) string {
	lines := strings.Split(charToken.Literal, "\n")
	line := lines[charToken.Line]

	return fmt.Sprintf("%s\n\t%s", line, strings.Repeat("^", charToken.PosStart))
}
