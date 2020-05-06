package main

import (
	"errors"
	"fmt"
	"github.com/TurnsCoffeeIntoScripts/git-log-issue-finder/pkg/configuration"
	"github.com/TurnsCoffeeIntoScripts/git-log-issue-finder/pkg/helpers"
	"github.com/TurnsCoffeeIntoScripts/git-log-issue-finder/pkg/interpreter/evaluator"
	"github.com/TurnsCoffeeIntoScripts/git-log-issue-finder/pkg/interpreter/lexer"
	iobject "github.com/TurnsCoffeeIntoScripts/git-log-issue-finder/pkg/interpreter/object"
	"github.com/TurnsCoffeeIntoScripts/git-log-issue-finder/pkg/interpreter/parser"
	"github.com/TurnsCoffeeIntoScripts/git-log-issue-finder/pkg/interpreter/repl"
	"github.com/TurnsCoffeeIntoScripts/git-log-issue-finder/pkg/script"
	"log"
	"os"
)

func main() {
	glifParam := configuration.GlifParameters{}
	forceRepl := false
	if len(os.Args) == 1 {
		forceRepl = true
	}

	if ok := glifParam.Parse(forceRepl); !ok {
		log.Fatal("failed to properly parse input flags")
	}

	if err := start(glifParam); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func start(glifParam configuration.GlifParameters) error {
	if helpers.IsBoolPtrTrue(glifParam.Flags.REPL) {
		repl.Start(os.Stdin, os.Stdout)
		return nil
	}

	var input *string

	if helpers.IsBoolPtrTrue(glifParam.Scripts.UseDiffLatestSemverWithLatestBuilds) {
		input = &script.DiffLatestSemverWithLatestBuilds
	} else if helpers.IsBoolPtrTrue(glifParam.Scripts.UseDiffLatestSemverWithLatestRCs) {
		input = &script.DiffLatestSemverWithLatestRCs
	} else if helpers.IsBoolPtrTrue(glifParam.Scripts.UseDiffLatestSemver) {
		input = &script.DiffLatestSemver
	} else if helpers.IsBoolPtrTrue(glifParam.Scripts.UseUserSpecifiedScript) {
		input = &glifParam.UserSpecifiedScript
	} else {
		return errors.New("unknown scripting mode")
	}

	l := lexer.New(*input)
	p := parser.NewWithOptions(l, false)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		repl.PrintParserErrors(os.Stdout, p.Errors())
		return fmt.Errorf("error parsing script")
	}

	env := iobject.NewEnvironmentWithParams(*glifParam.Tickets)
	evaluated := evaluator.Eval(program, env)

	switch evaluated.(type) {
	case *iobject.Error:
		return fmt.Errorf(fmt.Sprintf("%s", evaluated.Inspect()))
	default:
		return nil
	}
}
