package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/daichimukai/x/syakyo/monkey/lexer"
	"github.com/daichimukai/x/syakyo/monkey/token"
)

const prompt = "monkey> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(prompt)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)

		for tok := l.NextToken(); tok.Type != token.TypeEof; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
