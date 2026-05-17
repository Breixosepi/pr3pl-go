/*Eugenio Giusepi Montilla Russo*/
/*29958321*/

package main

import (
	"fmt"
	"os"
	"pr3pl/ast"
	"pr3pl/evaluator"
	"pr3pl/lexer"
	"pr3pl/object"
	"pr3pl/parser"
	"pr3pl/repl"
	"pr3pl/semantic"
	"pr3pl/transpiler"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
)

func main() {

	args := os.Args

	if len(args) > 1 {
		filePath := args[1]
		content, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Error reading file: %s\n", err)
			return
		}
		runSource(string(content))
	} else {
		fmt.Printf("PR3PL Interpreter - Interactive mode\n")
		repl.Start(os.Stdin, os.Stdout)
	}
}

func runSource(input string) {

	staticEnv := semantic.NewEnvironment()
	dynamicEnv := object.NewEnvironment()

	preludeAST := loadPrelude(staticEnv, dynamicEnv)

	l := lexer.New("file", input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		fmt.Printf("%s-------------------------------------------------------------------------------%s\n", ColorRed, ColorReset)
		fmt.Printf("%sSyntax errors:%s\n", ColorRed, ColorReset)
		for _, msg := range p.Errors() {
			fmt.Printf("%s\t%s%s\n", ColorRed, msg, ColorReset)
		}
		return
	}

	var originalCode string
	if preludeAST != nil {
		originalCode += transpiler.ToOriginalPR3PL(preludeAST) + "\n"
	}
	originalCode += transpiler.ToOriginalPR3PL(program)

	errFile := os.WriteFile("transpiled.txt", []byte(originalCode), 0644)
	if errFile != nil {
		fmt.Printf("%sError generating transpilation file: %v%s\n", ColorRed, errFile, ColorReset)
	} else {
		fmt.Printf("\n%ttranspilation saved in'transpiled.txt' %s\n", ColorGreen, ColorReset)
	}

	typeResult, err := semantic.TypeCheck(program, staticEnv)
	if err != nil {
		fmt.Printf("%s-------------------------------------------------------------------------------%s\n", ColorRed, ColorReset)
		fmt.Printf("%sType error: %s%s\n", ColorRed, err, ColorReset)
		return
	}

	result := evaluator.Eval(program, dynamicEnv)
	if result != nil {
		fmt.Printf("%sType: %s%s\n", ColorYellow, typeResult.Signature(), ColorReset)
		fmt.Printf("%sResult: %s%s\n\n", ColorBlue, result.Inspect(), ColorReset)
	}
}

func loadPrelude(staticEnv *semantic.Environment, dynamicEnv *object.Environment) *ast.Program {

	preludePath := "stdlib/prelude.pr3pl"
	content, err := os.ReadFile(preludePath)
	if err != nil {
		return nil
	}

	l := lexer.New("stdlib", string(content))
	p := parser.New(l)
	preludeProg := p.ParseProgram()

	if len(p.Errors()) == 0 {
		_, err := semantic.TypeCheck(preludeProg, staticEnv)
		if err == nil {
			evaluator.Eval(preludeProg, dynamicEnv)
		}
		return preludeProg
	}

	return nil
}
