package object

import (
	"fmt"

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

	return message
}
func (e *Error) InvokeMethod(method string, args ...Object) Object {
	return nil
}
