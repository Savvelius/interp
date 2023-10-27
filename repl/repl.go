package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/Savvelius/go-interp/lexer"
	"github.com/Savvelius/go-interp/token"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		// try scanning a line
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)

		// lex a line
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("Type: %v Value: %v\n", tok.Type, tok.Literal)
		}
	}
}
