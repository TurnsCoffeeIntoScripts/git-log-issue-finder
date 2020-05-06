// Package ast defines the various building blocks of the interpreter's abstract syntax tree (ast)
package ast

import (
	"bytes"
	"github.com/TurnsCoffeeIntoScripts/git-log-issue-finder/pkg/interpreter/gitoken"
	"strings"
)

// Node is an interface that needs to be implemented by every element that the AST will contain
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement is an interface enclosing a node and to be implemented by the language's statements
type Statement interface {
	Node
	statementNode()
}

// Expression is an interface enclosing a node and to be implemented by the language's expression
type Expression interface {
	Node
	expressionNode()
}

// Program is a collection (slice) of statements that represents the input program
type Program struct {
	Statements []Statement
}

// Identifier is an identifier node
type Identifier struct {
	Token gitoken.Token
	Value string
}

// LetStatement is an ast node representing a statement of the form: 'let <IDENT> = <EXPR>
type LetStatement struct {
	Token gitoken.Token
	Name  *Identifier
	Value Expression
}

// SetStatement is an ast node representing a statement of the form: 'set <KEYWORD> <LITTERAL>
type SetStatement struct {
	Token gitoken.Token
	Name  *Identifier
	Value Expression
}

// ReturnStatement is an ast node representing a statement of the form: 'return <EXPR>'
type ReturnStatement struct {
	Token       gitoken.Token
	ReturnValue Expression
}

// ExpressionStatement is an ast node representing a statement of the form: '<EXPR>'
type ExpressionStatement struct {
	Token      gitoken.Token
	Expression Expression
}

// BlockStatement is an ast node representing a collection of statement and an initiating token
type BlockStatement struct {
	Token      gitoken.Token
	Statements []Statement
}

// FunctionLiteral is an ast node representing a function of the form: 'fn(<PARAMS>) ast.BlockStatement'
type FunctionLiteral struct {
	Token      gitoken.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

// ArrayLiteral is an ast node representing an array of the form: '[<ELEM1>, <ELEME2>, ..., <ELEMn>]'
type ArrayLiteral struct {
	Token    gitoken.Token
	Elements []Expression
}

// HashLiteral is an ast node representing a hash of the form: '{<KEY1>: <VAL1>, <KEY2>: <VAL2>, ..., <KEYn>: <VALn>}'
type HashLiteral struct {
	Token gitoken.Token
	Pairs map[Expression]Expression
}

// IntegerLiteral is an ast node representing an integer value
type IntegerLiteral struct {
	Token gitoken.Token
	Value int64
}

// StringLiteral is an ast node representing a string value
type StringLiteral struct {
	Token gitoken.Token
	Value string
}

// Boolean is an ast node representing a boolean value
type Boolean struct {
	Token gitoken.Token
	Value bool
}

// PrefixExpression is an ast node representing an expression of the form: '<OPERATOR><EXPR>'
type PrefixExpression struct {
	Token    gitoken.Token
	Operator string
	Right    Expression
}

// InfixExpression is an ast node representing an expression of the form: '<EXPR> <OPERATOR> <EXPR>'
type InfixExpression struct {
	Token    gitoken.Token
	Left     Expression
	Operator string
	Right    Expression
}

// IfExpression is an ast node representing an expression of the form: 'if(<EXPR>) ast.BlockStatement else ast.BlockStatement '
type IfExpression struct {
	Token       gitoken.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

// CallExpression is an ast node representing an expression of the form: 'FNIDENT(PARAMS)'
type CallExpression struct {
	Token     gitoken.Token
	Function  Expression
	Arguments []Expression
}

// IndexExpression is an ast node representing an expression of the form: 'ARRAY[<EXPR>]'
type IndexExpression struct {
	Token gitoken.Token
	Left  Expression
	Index Expression
}

// TokenLiteral returns the literal string of the token
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}

	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

func (i *Identifier) expressionNode() {

}

// TokenLiteral returns the literal string of the token
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string {
	return i.Value
}

func (ls *LetStatement) statementNode() {

}

func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")
	return out.String()
}

// TokenLiteral returns the literal string of the token
func (ls *SetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

func (ls *SetStatement) statementNode() {

}

func (ls *SetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" ")

	if ls.Value != nil {
		out.WriteString("\"")
		out.WriteString(ls.Value.String())
		out.WriteString("\"")
	}

	out.WriteString(";")
	return out.String()
}

// TokenLiteral returns the literal string of the token
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

func (rs *ReturnStatement) statementNode() {

}

// TokenLiteral returns the literal string of the token
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

func (es *ExpressionStatement) statementNode() {

}

// TokenLiteral returns the literal string of the token
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}

func (bs *BlockStatement) statementNode() {

}

// TokenLiteral returns the literal string of the token
func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

func (fl *FunctionLiteral) expressionNode() {

}

// TokenLiteral returns the literal string of the token
func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	var params []string
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(fl.Body.String())

	return out.String()
}

func (al *ArrayLiteral) expressionNode() {

}

// TokenLiteral returns the literal string of the token
func (al *ArrayLiteral) TokenLiteral() string {
	return al.Token.Literal
}

func (al *ArrayLiteral) String() string {
	var out bytes.Buffer
	var elements []string
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

func (hl *HashLiteral) expressionNode() {

}

// TokenLiteral returns the literal string of the token
func (hl *HashLiteral) TokenLiteral() string {
	return hl.Token.Literal
}

func (hl *HashLiteral) String() string {
	var out bytes.Buffer

	var pairs []string
	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

func (il *IntegerLiteral) expressionNode() {

}

// TokenLiteral returns the literal string of the token
func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}

func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

func (sl *StringLiteral) expressionNode() {

}

// TokenLiteral returns the literal string of the token
func (sl *StringLiteral) TokenLiteral() string {
	return sl.Token.Literal
}

func (sl *StringLiteral) String() string {
	return sl.Token.Literal
}

func (b *Boolean) expressionNode() {

}

// TokenLiteral returns the literal string of the token
func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

func (b *Boolean) String() string {
	return b.Token.Literal
}

func (pe *PrefixExpression) expressionNode() {

}

// TokenLiteral returns the literal string of the token
func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

func (ie *InfixExpression) expressionNode() {

}

// TokenLiteral returns the literal string of the token
func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

func (ie *IfExpression) expressionNode() {

}

// TokenLiteral returns the literal string of the token
func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

func (ce *CallExpression) expressionNode() {

}

// TokenLiteral returns the literal string of the token
func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}

func (ce *CallExpression) String() string {
	var out bytes.Buffer

	var args []string
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

func (ie *IndexExpression) expressionNode() {

}

// TokenLiteral returns the literal string of the token
func (ie *IndexExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("]")
	out.WriteString(")")

	return out.String()

}
