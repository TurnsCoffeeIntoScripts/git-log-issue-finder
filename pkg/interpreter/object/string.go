package object

// String is a wraper on a native string type
type String struct {
	Value string
}

// Type returns StringObj (STRING)
func (s *String) Type() Type {
	return StringObj
}

// Inspect the native string value
func (s *String) Inspect() string {
	return s.Value
}
