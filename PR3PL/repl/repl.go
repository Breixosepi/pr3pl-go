package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"pr3pl/evaluator"
	"pr3pl/lexer"
	"pr3pl/object"
	"pr3pl/parser"
	"pr3pl/semantic"
)

const PROMPT = ">> "
const MULTILINE_PROMPT = ".. "

func Start(in io.Reader, out io.Writer) {

	scanner := bufio.NewScanner(in)
	staticEnv := semantic.NewEnvironment()
	dynamicEnv := object.NewEnvironment()

	for {
		fmt.Fprintf(out, PROMPT)

		var inputBuffer strings.Builder

		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				break
			}
			inputBuffer.WriteString(line)
			inputBuffer.WriteString("\n")
			fmt.Fprintf(out, MULTILINE_PROMPT)
		}

		if err := scanner.Err(); err != nil {
			return
		}

		if inputBuffer.Len() == 0 {
			continue
		}

		l := lexer.New("repl", inputBuffer.String())
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		typeResult, err := semantic.TypeCheck(program, staticEnv)
		if err != nil {
			io.WriteString(out, "Error de Tipos (Semántico):\n\t"+err.Error()+"\n")
			continue
		}

		func() {
			defer func() {
				if r := recover(); r != nil {
					io.WriteString(out, fmt.Sprintf("Error de Ejecución (Runtime):\n\t%v\n", r))
				}
			}()

			evalResult := evaluator.Eval(program, dynamicEnv)
			if evalResult != nil {
				io.WriteString(out, fmt.Sprintf("Tipo: %s\nResultado: %s\n", typeResult.Signature(), evalResult.Inspect()))
			}
		}()
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Se encontraron errores de sintaxis:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
