package object

import "fmt"

// Boolean is a wrapper on a native bool type
type Boolean struct {
	Value bool
}

// Type returns BooleanObj (BOOLEAN)
func (b *Boolean) Type() Type {
	return BooleanObj
}

// Inspect the native bool value
func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}
