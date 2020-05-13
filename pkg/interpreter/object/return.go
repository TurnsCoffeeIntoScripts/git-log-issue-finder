package object

// ReturnValue is a wrapper for whatever value is passed by a return call
type ReturnValue struct {
	Value Object
}

// Type returns ReturnValueObj (RETURN_VALUE)
func (rv *ReturnValue) Type() Type {
	return ReturnValueObj
}

// Inspect the inner value of the Object (any object implementing the interface Object)
func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}
