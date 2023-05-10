package object

type Error struct {
	Message  string
	File     string
	Line     int
	StartPos int
	EndPos   int
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return e.Message }
func (e *Error) InvokeMethod(method string, args ...Object) Object {
	return nil
}
