/*Eugenio Giusepi Montilla Russo*/
/*29958321*/

package main

import (
	"fmt"
	"os"
	"pr3pl/evaluator"
	"pr3pl/lexer"
	"pr3pl/object"
	"pr3pl/parser"
	"pr3pl/repl"
	"pr3pl/semantic"
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
	staticEnv := semantic.NewEnvironment()
	dynamicEnv := object.NewEnvironment()

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
