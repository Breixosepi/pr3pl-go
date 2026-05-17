package repl

import (
	"bufio"
	"fmt"
	"io"
	"pr3pl/evaluator"
	"pr3pl/lexer"
	"pr3pl/object"
	"pr3pl/parser"
	"pr3pl/semantic"
)

const PROMPT = "Input: "

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
)

func Start(in io.Reader, out io.Writer) {

	scanner := bufio.NewScanner(in)
	staticEnv := semantic.NewEnvironment()
	dynamicEnv := object.NewEnvironment()

	fmt.Fprintf(out, "%sCommand Line Interface for PR3PL (Go Version)%s\n", ColorBlue, ColorReset)
	fmt.Fprintf(out, "If you want to quit, type exit and press enter or press ctrl+C\n\n")

	for {

		fmt.Fprintf(out, "%s%s%s", ColorGreen, PROMPT, ColorReset)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if line == "exit" {
			return
		}

		l := lexer.New("repl", line)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		typeResult, err := semantic.TypeCheck(program, staticEnv)
		if err != nil {
			fmt.Fprintf(out, "%s-------------------------------------------------------------------------------%s\n", ColorRed, ColorReset)
			fmt.Fprintf(out, "%sType error: %s%s\n\n", ColorRed, err.Error(), ColorReset)
			continue
		}

		result := evaluator.Eval(program, dynamicEnv)
		if result != nil {
			fmt.Fprintf(out, "%sType: %s%s\n", ColorYellow, typeResult.Signature(), ColorReset)
			fmt.Fprintf(out, "%sOutput: %s%s\n\n", ColorBlue, result.Inspect(), ColorReset)
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	fmt.Fprintf(out, "%s-------------------------------------------------------------------------------%s\n", ColorRed, ColorReset)
	fmt.Fprintf(out, "%sSyntax errors:%s\n", ColorRed, ColorReset)
	for _, msg := range errors {
		fmt.Fprintf(out, "%s\t%s%s\n", ColorRed, msg, ColorReset)
	}
	fmt.Fprintf(out, "\n")
}
