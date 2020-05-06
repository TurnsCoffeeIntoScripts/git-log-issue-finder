// Package parser contains the definition of the parser element of the glif interpreter.
//
// If also declare the following:
//	- Precedence constants
//	- map of token/precedence
// 	- prefix and infix functions
//	- parsing method for recognize tokens
package parser

import (
	"fmt"
	"github.com/TurnsCoffeeIntoScripts/git-log-issue-finder/pkg/interpreter/ast"
	"github.com/TurnsCoffeeIntoScripts/git-log-issue-finder/pkg/interpreter/gitoken"
	"github.com/TurnsCoffeeIntoScripts/git-log-issue-finder/pkg/interpreter/lexer"
	"strconv"
)

// Contant used in the precedence parsing for operator
const (
	_ int = iota
	LOWEST
	TO          // ->
	EQUALS      // ==
	LESSGREATER // < or >
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INDEX       // array[index]
)

// Map associating tokens given by the lexer to a specific precedence
var precedences = map[gitoken.TokenType]int{
	gitoken.TO:       TO,
	gitoken.EQ:       EQUALS,
	gitoken.NOTEQ:    EQUALS,
	gitoken.LT:       LESSGREATER,
	gitoken.GT:       LESSGREATER,
	gitoken.PLUS:     SUM,
	gitoken.MINUS:    SUM,
	gitoken.SLASH:    PRODUCT,
	gitoken.ASTERISK: PRODUCT,
	gitoken.LPAREN:   CALL,
	gitoken.LBRAKET:  INDEX,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// Parser is the representation of the glif parser
// It contains the following:
//	- A slice of errors (parsing error)
//	- The current token AND the next (peek) token
//	- Maps containing the prefix AND infix parse functions
type Parser struct {
	l         *lexer.Lexer
	errors    []string
	showTrace bool

	currentToken gitoken.Token
	peekToken    gitoken.Token

	prefixParseFns map[gitoken.TokenType]prefixParseFn
	infixParseFns  map[gitoken.TokenType]infixParseFn
}

// NewWithOptions creates a new parser from a lexer and additionnal bool options
// The creation of the parser is done via the New function and not directly here,
// that way there is only one place where the parser is instantiated.
func NewWithOptions(l *lexer.Lexer, showTrace bool) *Parser {
	p := New(l)
	p.showTrace = showTrace

	return p
}

// New creates a new parser from a lexer and register the prefix/infix parsing functions
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:         l,
		errors:    []string{},
		showTrace: false,
	}
	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[gitoken.TokenType]prefixParseFn)
	p.registerPrefix(gitoken.IDENT, p.parseIdentifier)
	p.registerPrefix(gitoken.INT, p.parseIntegerLiteral)
	p.registerPrefix(gitoken.BANG, p.parsePrefixExpression)
	p.registerPrefix(gitoken.MINUS, p.parsePrefixExpression)
	p.registerPrefix(gitoken.TRUE, p.parseBoolean)
	p.registerPrefix(gitoken.FALSE, p.parseBoolean)
	p.registerPrefix(gitoken.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(gitoken.IF, p.parseIfExpression)
	p.registerPrefix(gitoken.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(gitoken.STRING, p.parseStringLiteral)
	p.registerPrefix(gitoken.LBRAKET, p.parseArrayLiteral)
	p.registerPrefix(gitoken.LBRACE, p.parseHashLiteral)

	p.infixParseFns = make(map[gitoken.TokenType]infixParseFn)
	p.registerInfix(gitoken.PLUS, p.parseInfixExpression)
	p.registerInfix(gitoken.MINUS, p.parseInfixExpression)
	p.registerInfix(gitoken.SLASH, p.parseInfixExpression)
	p.registerInfix(gitoken.ASTERISK, p.parseInfixExpression)
	p.registerInfix(gitoken.EQ, p.parseInfixExpression)
	p.registerInfix(gitoken.NOTEQ, p.parseInfixExpression)
	p.registerInfix(gitoken.LT, p.parseInfixExpression)
	p.registerInfix(gitoken.GT, p.parseInfixExpression)
	p.registerInfix(gitoken.LPAREN, p.parseCallExpression)
	p.registerInfix(gitoken.LBRAKET, p.parseIndexExpression)
	p.registerInfix(gitoken.TO, p.parseInfixExpression)

	return p
}

// Errors returns the errors found (if any) while parsing
func (p *Parser) Errors() []string {
	return p.errors
}

// ParseProgram is the entry point of the parser. It returns the *ast.Program node which represents the entire script.
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.currentTokenIs(gitoken.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

// Advance the current token the next one (peek) and set the next one (peek) to the lexer's next token.
func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) registerPrefix(tokenType gitoken.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType gitoken.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) noPrefixParseFnError(t gitoken.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

// Parse an ast.Statement (LET/SET/RETURN) or an expression statement by default.
func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case gitoken.LET:
		return p.parseLetStatement()
	case gitoken.SET:
		return p.parseSetStatement()
	case gitoken.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// Parse an expression, either prefix or infix.
func (p *Parser) parseExpression(precedence int) ast.Expression {
	if p.showTrace {
		defer untrace(trace("parseExpression " + p.currentToken.Literal))
	}
	prefix := p.prefixParseFns[p.currentToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.currentToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(gitoken.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	if p.showTrace {
		defer untrace(trace("parseLetStatement " + p.currentToken.Literal))
	}
	stmt := &ast.LetStatement{Token: p.currentToken}
	if !p.expectPeek(gitoken.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
	if !p.expectPeek(gitoken.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)
	if p.peekTokenIs(gitoken.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseSetStatement() *ast.SetStatement {
	if p.showTrace {
		defer untrace(trace("parseSetStatement " + p.currentToken.Literal))
	}
	stmt := &ast.SetStatement{Token: p.currentToken}
	if !p.expectPeek(gitoken.REPOPATH) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
	p.nextToken()

	stmt.Value = p.parseStringLiteral()
	if p.peekTokenIs(gitoken.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	if p.showTrace {
		defer untrace(trace("parseReturnStatement " + p.currentToken.Literal))
	}
	stmt := &ast.ReturnStatement{Token: p.currentToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)
	if p.peekTokenIs(gitoken.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	if p.showTrace {
		defer untrace(trace("parseExpressionStatement " + p.currentToken.Literal))
	}
	stmt := &ast.ExpressionStatement{Token: p.currentToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(gitoken.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	if p.showTrace {
		defer untrace(trace("parseBlockStatement " + p.currentToken.Literal))
	}
	block := &ast.BlockStatement{Token: p.currentToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.currentTokenIs(gitoken.RBRACE) && !p.currentTokenIs(gitoken.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseIdentifier() ast.Expression {
	if p.showTrace {
		defer untrace(trace("parseIdentifier " + p.currentToken.Literal))
	}
	return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	if p.showTrace {
		defer untrace(trace("parseIntegerLiteral " + p.currentToken.Literal))
	}
	il := &ast.IntegerLiteral{Token: p.currentToken}

	value, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.currentToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	il.Value = value
	return il
}

func (p *Parser) parseStringLiteral() ast.Expression {
	if p.showTrace {
		defer untrace(trace("parseStringLiteral " + p.currentToken.Literal))
	}
	return &ast.StringLiteral{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	if p.showTrace {
		defer untrace(trace("parseArrayLiteral " + p.currentToken.Literal))
	}
	array := &ast.ArrayLiteral{Token: p.currentToken}
	array.Elements = p.parseExpressionList(gitoken.RBRAKET)
	return array
}

func (p *Parser) parseHashLiteral() ast.Expression {
	if p.showTrace {
		defer untrace(trace("parseHashLiteral " + p.currentToken.Literal))
	}
	hash := &ast.HashLiteral{Token: p.currentToken}
	hash.Pairs = make(map[ast.Expression]ast.Expression)

	for !p.peekTokenIs(gitoken.RBRACE) {
		p.nextToken()
		key := p.parseExpression(LOWEST)

		if !p.expectPeek(gitoken.COLON) {
			return nil
		}

		p.nextToken()
		value := p.parseExpression(LOWEST)

		hash.Pairs[key] = value

		if !p.peekTokenIs(gitoken.RBRACE) && !p.expectPeek(gitoken.COMMA) {
			return nil
		}
	}

	if !p.expectPeek(gitoken.RBRACE) {
		return nil
	}

	return hash
}

func (p *Parser) parseBoolean() ast.Expression {
	if p.showTrace {
		defer untrace(trace("parseBoolean " + p.currentToken.Literal))
	}
	return &ast.Boolean{Token: p.currentToken, Value: p.currentTokenIs(gitoken.TRUE)}
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	if p.showTrace {
		defer untrace(trace("parseFunctionLiteral " + p.currentToken.Literal))
	}
	fl := &ast.FunctionLiteral{Token: p.currentToken}

	if !p.expectPeek(gitoken.LPAREN) {
		return nil
	}

	fl.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(gitoken.LBRACE) {
		return nil
	}

	fl.Body = p.parseBlockStatement()

	return fl
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	if p.showTrace {
		defer untrace(trace("parseFunctionParameters " + p.currentToken.Literal))
	}
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(gitoken.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	id := &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
	identifiers = append(identifiers, id)

	for p.peekTokenIs(gitoken.COMMA) {
		p.nextToken()
		p.nextToken()
		id := &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
		identifiers = append(identifiers, id)
	}

	if !p.expectPeek(gitoken.RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	if p.showTrace {
		defer untrace(trace("parsePrefixExpression " + p.currentToken.Literal))
	}
	expression := &ast.PrefixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)
	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	if p.showTrace {
		defer untrace(trace("parseInfixExpression " + p.currentToken.Literal))
	}
	expression := &ast.InfixExpression{
		Token:    p.currentToken,
		Left:     left,
		Operator: p.currentToken.Literal,
	}

	precedence := p.currentPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	if p.showTrace {
		defer untrace(trace("parseGroupedExpression " + p.currentToken.Literal))
	}
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(gitoken.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	if p.showTrace {
		defer untrace(trace("parseIfExpression " + p.currentToken.Literal))
	}
	expression := &ast.IfExpression{Token: p.currentToken}

	if !p.expectPeek(gitoken.LPAREN) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(gitoken.RPAREN) {
		return nil
	}

	if !p.expectPeek(gitoken.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(gitoken.ELSE) {
		p.nextToken()

		if !p.expectPeek(gitoken.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseCallExpression(fn ast.Expression) ast.Expression {
	if p.showTrace {
		defer untrace(trace("parseCallExpression " + p.currentToken.Literal))
	}
	exp := &ast.CallExpression{Token: p.currentToken, Function: fn}
	exp.Arguments = p.parseExpressionList(gitoken.RPAREN)
	return exp
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	if p.showTrace {
		defer untrace(trace("parseIndexExpression " + p.currentToken.Literal))
	}
	exp := &ast.IndexExpression{Token: p.currentToken, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(gitoken.RBRAKET) {
		return nil
	}

	return exp
}

func (p *Parser) parseExpressionList(end gitoken.TokenType) []ast.Expression {
	if p.showTrace {
		defer untrace(trace("parseExpressionList " + p.currentToken.Literal))
	}
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(gitoken.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

// Return the precedence of the current token
func (p *Parser) currentPrecedence() int {
	if p, ok := precedences[p.currentToken.Type]; ok {
		return p
	}

	return LOWEST
}

// Return the precedence of the next (peek) token
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

// Check whether the current token is of a certain type
func (p *Parser) currentTokenIs(t gitoken.TokenType) bool {
	return p.currentToken.Type == t
}

// Check whether the next (peek) token is of a certain type
func (p *Parser) peekTokenIs(t gitoken.TokenType) bool {
	return p.peekToken.Type == t
}

// Check whether the next (peek) token is of a certain type and advances the current token if true. Otherwise returns
// an error
func (p *Parser) expectPeek(t gitoken.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}

	p.peekError(t)
	return false
}

// Add a parser error because next (peek) token is not what was expected by the parser
func (p *Parser) peekError(t gitoken.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
