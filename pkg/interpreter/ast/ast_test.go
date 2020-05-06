package ast

import (
	"github.com/TurnsCoffeeIntoScripts/git-log-issue-finder/pkg/interpreter/gitoken"
	"testing"
)

func TestString(t *testing.T) {
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
