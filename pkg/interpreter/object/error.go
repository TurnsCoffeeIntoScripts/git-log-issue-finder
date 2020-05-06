package object

// Error is an object that encloses an error message and allows error management within the interpreter
type Error struct {
	Message string
}

// Type returns ErrorObj (ERROR)
func (e *Error) Type() Type {
	return ErrorObj
}

// Inspect returns the error message enclosed in the error object
func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}
