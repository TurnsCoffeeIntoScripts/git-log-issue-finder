package object

import (
	"bytes"
	"github.com/TurnsCoffeeIntoScripts/git-log-issue-finder/pkg/interpreter/ast"
	"strings"
)

// BuiltinFunction represent a function builtin to the interpreter (meaning not user-defined)
type BuiltinFunction func(args ...Object) Object

// Function is an executable block
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

// Type returns FunctionObj (FUNCTION)
func (f *Function) Type() Type {
	return FunctionObj
}

// Inspect the body and parameters of the Function type
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}
