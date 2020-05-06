// Package repl is for the read-eval-print-loop
package repl

import (
	"bufio"
	"fmt"
	"github.com/TurnsCoffeeIntoScripts/git-log-issue-finder/pkg/interpreter/evaluator"
	"github.com/TurnsCoffeeIntoScripts/git-log-issue-finder/pkg/interpreter/lexer"
	"github.com/TurnsCoffeeIntoScripts/git-log-issue-finder/pkg/interpreter/object"
	"github.com/TurnsCoffeeIntoScripts/git-log-issue-finder/pkg/interpreter/parser"
	"github.com/TurnsCoffeeIntoScripts/git-log-issue-finder/pkg/version"
	"io"
)

// Prompt are the characters that are displayed for the REPL
const (
	InitialPrompt = "-- Type 'exit' to close the repl.\n----------\n"
	Prompt        = ">> "
	trace         = true
)

// CatBug is the "picture" displayed when an error has been detected by either the parse or the evaluator
const CatBug = "" +
	"CAN I HAS NO BUGZ PLEASE?!\n" +
	"       _                        \n" +
	"       \\`*-.                    \n" +
	"        )  _`-.                 \n" +
	"       .  : `. .                \n" +
	"       : _   '  \\               \n" +
	"       ; *` _.   `*-._          \n" +
	"       `-.-'          `-.       \n" +
	"         ;       `       `.     \n" +
	"         :.       .        \\    \n" +
	"         . \\  .   :   .-'   .   \n" +
	"         '  `+.;  ;  '      :   \n" +
	"         :  '  |    ;       ;-. \n" +
	"         ; '   : :`-:     _.`* ;\n" +
	"[bug] .*' /  .*' ; .*`- +'  `*' \n" +
	"      `*-*   `*-*  `*-*'\n"

// Start begin the repl loop
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironmentWithParams("*")

	fmt.Println(version.Get())
	fmt.Printf(InitialPrompt)
	for {
		fmt.Printf(Prompt)
		if scanned := scanner.Scan(); !scanned {
			return
		}

		line := scanner.Text()

		if line == "exit" {
			break
		}

		l := lexer.New(line)
		p := parser.NewWithOptions(l, trace)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			PrintParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

// PrintParserErrors simply prints the error returned by the parser and also prints the "picture" CatBug
func PrintParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, CatBug)
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
