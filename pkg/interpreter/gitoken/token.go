// Package gitoken (Glif Interpreter TOKEN) defines the list of recognized tokens
package gitoken

// Definition if the various constant for the tokens to be returned by the lexer
const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers and litterals
	IDENT  = "IDENT" //
	INT    = "INT"   // 0123456789
	STRING = "STRING"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT = "<"
	GT = ">"

	EQ    = "=="
	NOTEQ = "!="

	TO = "->"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	LBRAKET   = "["
	RBRAKET   = "]"
	COLON     = ":"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	SET      = "SET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	REPOPATH = "REPOPATH"
)

// TokenType is a simple string to store the type of the token object
type TokenType string

// Token is a series of character forming a word (literal)
type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"fn":       FUNCTION,
	"let":      LET,
	"set":      SET,
	"true":     TRUE,
	"false":    FALSE,
	"if":       IF,
	"else":     ELSE,
	"return":   RETURN,
	"repopath": REPOPATH,
}

// LookupIdent returns the TokenType of the keyword if the ident string is a keyword.
// Otherwise it returns the contant IDEN
func LookupIdent(ident string) TokenType {
	if tkn, ok := keywords[ident]; ok {
		return tkn
	}

	return IDENT
}
