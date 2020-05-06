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

	checkTest(t, program, "let myVar = anotherVar;")
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

	checkTest(t, program, "let x = 42;")
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

	checkTest(t, program, "set repopath \".\";")
}

func TestReturnStatement1(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&ReturnStatement{
				Token: gitoken.Token{Type: gitoken.RETURN, Literal: "return"},
				ReturnValue: &IntegerLiteral{
					Token: gitoken.Token{Type: gitoken.IDENT, Literal: "23"},
				},
			},
		},
	}

	checkTest(t, program, "return 23;")
}

func TestReturnStatement2(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&ReturnStatement{
				Token: gitoken.Token{Type: gitoken.RETURN, Literal: "return"},
				ReturnValue: &Identifier{
					Token: gitoken.Token{Type: gitoken.IDENT, Literal: "x"},
					Value: "x",
				},
			},
		},
	}

	checkTest(t, program, "return x;")
}

func TestPrefixExpression1(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&ExpressionStatement{
				Token:      gitoken.Token{},
				Expression: &PrefixExpression{
					Token:    gitoken.Token{Type: gitoken.BANG, Literal: "!"},
					Operator: gitoken.BANG,
					Right:    &Boolean{
						Token: gitoken.Token{Type: gitoken.IDENT, Literal: "true"},
						Value: true,
					},
				},
			},
		},
	}

	checkTest(t, program, "(!true)")
}

func TestPrefixExpression2(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&ExpressionStatement{
				Token:      gitoken.Token{},
				Expression: &PrefixExpression{
					Token:    gitoken.Token{Type: gitoken.MINUS, Literal: "-"},
					Operator: gitoken.MINUS,
					Right:    &IntegerLiteral{
						Token: gitoken.Token{Type: gitoken.IDENT, Literal: "73"},
						Value: 73,
					},
				},
			},
		},
	}

	checkTest(t, program, "(-73)")
}

func checkTest(t *testing.T, program *Program, expected string) {
	if program.String() != expected {
		t.Errorf("program.String() wrong. \ngot=%q\nexpected=%q", program.String(), expected)
	}
}