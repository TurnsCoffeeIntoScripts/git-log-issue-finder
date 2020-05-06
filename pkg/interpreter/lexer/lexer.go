// Package lexer contains the definition of the lexer element of the glif interpreter
package lexer

import (
	"github.com/TurnsCoffeeIntoScripts/git-log-issue-finder/pkg/interpreter/gitoken"
)

// Lexer is the representation of the glif lexer
// It contains the following:
//	- An input string (the program)
//	- Positions indicators
//	- A byte representing the current character under examination
type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

// New creates a new lexer from an input string (program/code)
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func newToken(tokenType gitoken.TokenType, ch byte) gitoken.Token {
	return gitoken.Token{Type: tokenType, Literal: string(ch)}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// NextToken reads the next gitoken.Token by parsing the next one (or two) bytes
func (l *Lexer) NextToken() gitoken.Token {
	var tkn gitoken.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tkn = gitoken.Token{Type: gitoken.EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tkn = newToken(gitoken.ASSIGN, l.ch)
		}
	case '+':
		tkn = newToken(gitoken.PLUS, l.ch)
	case '-':
		if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			tkn = gitoken.Token{Type: gitoken.TO, Literal: string(ch) + string(l.ch)}
		} else {
			tkn = newToken(gitoken.MINUS, l.ch)
		}
	case ';':
		tkn = newToken(gitoken.SEMICOLON, l.ch)
	case ',':
		tkn = newToken(gitoken.COMMA, l.ch)
	case '(':
		tkn = newToken(gitoken.LPAREN, l.ch)
	case ')':
		tkn = newToken(gitoken.RPAREN, l.ch)
	case '{':
		tkn = newToken(gitoken.LBRACE, l.ch)
	case '}':
		tkn = newToken(gitoken.RBRACE, l.ch)
	case '[':
		tkn = newToken(gitoken.LBRAKET, l.ch)
	case ']':
		tkn = newToken(gitoken.RBRAKET, l.ch)
	case ':':
		tkn = newToken(gitoken.COLON, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tkn = gitoken.Token{Type: gitoken.NOTEQ, Literal: string(ch) + string(l.ch)}
		} else {
			tkn = newToken(gitoken.BANG, l.ch)
		}
	case '*':
		tkn = newToken(gitoken.ASTERISK, l.ch)
	case '/':
		tkn = newToken(gitoken.SLASH, l.ch)
	case '<':
		tkn = newToken(gitoken.LT, l.ch)
	case '>':
		tkn = newToken(gitoken.GT, l.ch)

	case '"':
		tkn.Type = gitoken.STRING
		tkn.Literal = l.readString()
	case 0:
		tkn.Literal = ""
		tkn.Type = gitoken.EOF
	default:
		if isLetter(l.ch) {
			tkn.Literal = l.readIdentifier()
			tkn.Type = gitoken.LookupIdent(tkn.Literal)
			return tkn
		} else if isDigit(l.ch) {
			tkn.Type = gitoken.INT
			tkn.Literal = l.readNumber()
			return tkn
		} else {
			tkn = newToken(gitoken.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tkn
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}

	return l.input[position:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}

	return l.input[l.readPosition]
}
