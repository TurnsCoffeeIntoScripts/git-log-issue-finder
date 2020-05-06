package ast

import (
	"github.com/TurnsCoffeeIntoScripts/git-log-issue-finder/pkg/interpreter/gitoken"
	"testing"
)

func TestLetStatement1(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: gitoken.Token{Type: gitoken.LET, Literal: "let"},
				Name: &Identifier{
					Token: gitoken.Token{Type: gitoken.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: gitoken.Token{Type: gitoken.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != "let myVar = anotherVar;" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}

func TestLetStatement2(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: gitoken.Token{Type: gitoken.LET, Literal: "let"},
				Name: &Identifier{
					Token: gitoken.Token{Type: gitoken.IDENT, Literal: "x"},
					Value: "x",
				},
				Value: &IntegerLiteral{
					Token: gitoken.Token{Type: gitoken.IDENT, Literal: "42"},
					Value: 42,
				},
			},
		},
	}

	if program.String() != "let x = 42;" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}

func TestSetStatement1(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&SetStatement{
				Token: gitoken.Token{Type: gitoken.SET, Literal: "set"},
				Name: &Identifier{
					Token: gitoken.Token{Type: gitoken.REPOPATH, Literal: "repopath"},
					Value: "repopath",
				},
				Value: &StringLiteral{
					Token: gitoken.Token{Type: gitoken.IDENT, Literal: "."},
					Value: ".",
				},
			},
		},
	}

	if program.String() != "set repopath \".\";" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}
