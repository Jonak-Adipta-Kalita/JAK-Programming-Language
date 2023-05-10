package object

import (
	"fmt"

	"github.com/Jonak-Adipta-Kalita/JAK-Programming-Language/position"
)

type Error struct {
	Message   string
	ErrorName string
	StartPos  *position.Position
	EndPos    *position.Position
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string {
	message := fmt.Sprintf("Error: `%s`", e.Message)
	message += fmt.Sprintf("\n\tat %s: %d", e.StartPos.FileName, e.StartPos.Line+1)

	return message
}
func (e *Error) InvokeMethod(method string, args ...Object) Object {
	return nil
}
