package object

import "fmt"

type Error struct {
	Message   string
	File      string
	ErrorName string
	Line      int
	StartPos  int
	EndPos    int
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string {
	message := fmt.Sprintf("Error: `%s`", e.Message)
	message += fmt.Sprintf("\n\tat %s: %d", e.File, e.Line+1)

	return message
}
func (e *Error) InvokeMethod(method string, args ...Object) Object {
	return nil
}
