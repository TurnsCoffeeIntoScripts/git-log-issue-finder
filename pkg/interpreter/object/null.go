package object

// Null is the construct used to represent null values in the AST
type Null struct {
}

// Type returns NullObj (NULL)
func (n *Null) Type() Type {
	return NullObj
}

// Inspect simply returns "null"
func (n *Null) Inspect() string {
	return "null"
}
