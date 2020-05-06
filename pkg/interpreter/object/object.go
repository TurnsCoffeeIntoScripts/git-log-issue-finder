// Package object contains the various implementations of the object.Object interface
// List of the supported object
//	- Array
//	- Boolean
//	- Builtin (function)
//	- Environment (for variable definition and such)
//	- Error (for parser handling)
//	- Function (user defined, not builtins)
//	- Hash
//	- Integer
//	- Null
//	- Repo (a go-git git repository)
// 	- Return
//	- String
//	- Tag (a go-git tag object)
package object

// Definition of constants for "Type"
const (
	IntegerObj     = "INTEGER"
	BooleanObj     = "BOOLEAN"
	NullObj        = "NULL"
	ReturnValueObj = "RETURN_VALUE"
	ErrorObj       = "ERROR"
	FunctionObj    = "FUNCTION"
	StringObj      = "STRING"
	BuiltinObj     = "BUILTIN"
	ArrayObj       = "ARRAY"
	HashObj        = "HASH"
	RepoObj        = "REPO"
	TagObj         = "TAG"
)

// Type refers to the constant which defines an internal type
type Type string

// Object is the interface to be implement for any and all object that will be stored in the AST
type Object interface {
	Type() Type
	Inspect() string
}
