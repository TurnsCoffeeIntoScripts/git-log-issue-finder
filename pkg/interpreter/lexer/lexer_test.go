package lexer

import (
	"github.com/TurnsCoffeeIntoScripts/git-log-issue-finder/pkg/interpreter/gitoken"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `let five = 5;
let ten = 10;

let add = fn(x, y) {
	x + y;
};

let result = add(five, ten);
!-/*5;
5 < 10 > 5;

if (5 < 10) {
	return true;
} else {
	return false;
}

10 == 10;
10 != 9;
"foobar"
"foo bar"
[1, 2];
{"foo": "bar"}
set repopath "abc/def"
return b->c
`

	tests := []struct {
		expectedType     gitoken.TokenType
		expectedLitteral string
	}{
		{gitoken.LET, "let"},
		{gitoken.IDENT, "five"},
		{gitoken.ASSIGN, "="},
		{gitoken.INT, "5"},
		{gitoken.SEMICOLON, ";"},
		{gitoken.LET, "let"},
		{gitoken.IDENT, "ten"},
		{gitoken.ASSIGN, "="},
		{gitoken.INT, "10"},
		{gitoken.SEMICOLON, ";"},
		{gitoken.LET, "let"},
		{gitoken.IDENT, "add"},
		{gitoken.ASSIGN, "="},
		{gitoken.FUNCTION, "fn"},
		{gitoken.LPAREN, "("},
		{gitoken.IDENT, "x"},
		{gitoken.COMMA, ","},
		{gitoken.IDENT, "y"},
		{gitoken.RPAREN, ")"},
		{gitoken.LBRACE, "{"},
		{gitoken.IDENT, "x"},
		{gitoken.PLUS, "+"},
		{gitoken.IDENT, "y"},
		{gitoken.SEMICOLON, ";"},
		{gitoken.RBRACE, "}"},
		{gitoken.SEMICOLON, ";"},
		{gitoken.LET, "let"},
		{gitoken.IDENT, "result"},
		{gitoken.ASSIGN, "="},
		{gitoken.IDENT, "add"},
		{gitoken.LPAREN, "("},
		{gitoken.IDENT, "five"},
		{gitoken.COMMA, ","},
		{gitoken.IDENT, "ten"},
		{gitoken.RPAREN, ")"},
		{gitoken.SEMICOLON, ";"},
		{gitoken.BANG, "!"},
		{gitoken.MINUS, "-"},
		{gitoken.SLASH, "/"},
		{gitoken.ASTERISK, "*"},
		{gitoken.INT, "5"},
		{gitoken.SEMICOLON, ";"},
		{gitoken.INT, "5"},
		{gitoken.LT, "<"},
		{gitoken.INT, "10"},
		{gitoken.GT, ">"},
		{gitoken.INT, "5"},
		{gitoken.SEMICOLON, ";"},
		{gitoken.IF, "if"},
		{gitoken.LPAREN, "("},
		{gitoken.INT, "5"},
		{gitoken.LT, "<"},
		{gitoken.INT, "10"},
		{gitoken.RPAREN, ")"},
		{gitoken.LBRACE, "{"},
		{gitoken.RETURN, "return"},
		{gitoken.TRUE, "true"},
		{gitoken.SEMICOLON, ";"},
		{gitoken.RBRACE, "}"},
		{gitoken.ELSE, "else"},
		{gitoken.LBRACE, "{"},
		{gitoken.RETURN, "return"},
		{gitoken.FALSE, "false"},
		{gitoken.SEMICOLON, ";"},
		{gitoken.RBRACE, "}"},
		{gitoken.INT, "10"},
		{gitoken.EQ, "=="},
		{gitoken.INT, "10"},
		{gitoken.SEMICOLON, ";"},
		{gitoken.INT, "10"},
		{gitoken.NOTEQ, "!="},
		{gitoken.INT, "9"},
		{gitoken.SEMICOLON, ";"},
		{gitoken.STRING, "foobar"},
		{gitoken.STRING, "foo bar"},
		{gitoken.LBRAKET, "["},
		{gitoken.INT, "1"},
		{gitoken.COMMA, ","},
		{gitoken.INT, "2"},
		{gitoken.RBRAKET, "]"},
		{gitoken.SEMICOLON, ";"},

		{gitoken.LBRACE, "{"},
		{gitoken.STRING, "foo"},
		{gitoken.COLON, ":"},
		{gitoken.STRING, "bar"},
		{gitoken.RBRACE, "}"},

		{gitoken.SET, "set"},
		{gitoken.REPOPATH, "repopath"},
		{gitoken.STRING, "abc/def"},

		{gitoken.RETURN, "return"},
		{gitoken.IDENT, "b"},
		{gitoken.TO, "->"},
		{gitoken.IDENT, "c"},

		{gitoken.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tkn := l.NextToken()

		if tkn.Type != tt.expectedType {
			t.Fatalf("tests[%d] - gitokentype wrong. exptected=%q, got=%q",
				i, tt.expectedType, tkn.Type)
		}

		if tkn.Literal != tt.expectedLitteral {
			t.Fatalf("tests[%d] - litearl wrong. expected=%q, got=%q",
				i, tt.expectedLitteral, tkn.Literal)
		}
	}
}
