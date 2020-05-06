package object

// Builtin represent a builtin construct of the interpreter
type Builtin struct {
	Fn         BuiltinFunction
	RequireEnv bool
	EnvName    string
}

// Type returns BuiltinObj (BUILTIN)
func (b *Builtin) Type() Type {
	return BuiltinObj
}

// Inspect simply returns the string "builtin function"
func (b *Builtin) Inspect() string {
	return "builtin function"
}
