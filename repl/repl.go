package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/Savvelius/go-interp/evaluator"
	"github.com/Savvelius/go-interp/lexer"
	"github.com/Savvelius/go-interp/object"
	"github.com/Savvelius/go-interp/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParseErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParseErrors(writer io.Writer, errors []string) {
	io.WriteString(writer, "Error:\n")
	for _, err := range errors {
		io.WriteString(writer, "\t"+err+"\n")
	}
}
