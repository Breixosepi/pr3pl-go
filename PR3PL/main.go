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

func main() {
	args := os.Args

	if len(args) > 1 {
		filePath := args[1]
		content, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Error al leer el archivo: %s\n", err)
			return
		}
		runSource(string(content))
	} else {
		fmt.Printf("PR3PL Interpreter - Modo Interactivo\n")
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
		fmt.Printf("Errores de Sintaxis:\n")
		for _, msg := range p.Errors() {
			fmt.Printf("\t%s\n", msg)
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
		fmt.Printf("Error al generar archivo de transpilación: %v\n", errFile)
	} else {
		fmt.Println("Transpilación guardada en 'transpiled.txt' ")
	}

	typeResult, err := semantic.TypeCheck(program, staticEnv)
	if err != nil {
		fmt.Printf("Error Semántico: %s\n", err)
		return
	}

	result := evaluator.Eval(program, dynamicEnv)
	if result != nil {
		fmt.Printf("Tipo Final: %s\n", typeResult.Signature())
		fmt.Printf("Resultado: %s\n", result.Inspect())
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
