package repl

import (
	"bufio"
	"fmt"
	"io"
	"pr3pl/lexer"
	"pr3pl/parser"
	"strings"
)

const PROMPT = ">> "
const MULTILINE_PROMPT = ".. "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

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

		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "¡Oops! Tienes errores de sintaxis en PR3PL:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
