package object

import (
	"bytes"
	"strings"
)

// Array represent a standard collection of object
type Array struct {
	Elements []Object
}

// Type returns ArrayObj (ARRAY)
func (a *Array) Type() Type {
	return ArrayObj
}

// Inspect the content the Array type
func (a *Array) Inspect() string {
	var out bytes.Buffer

	var elements []string
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}
