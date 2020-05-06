package object

import "fmt"

// Integer is a wrapper on a native int64 type
type Integer struct {
	Value int64
}

// Type returns IntegerObj (INTEGER)
func (i *Integer) Type() Type {
	return IntegerObj
}

// Inspect the native int64 tyoe
func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}
