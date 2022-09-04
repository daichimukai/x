package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/daichimukai/x/syakyo/monkey/eval"
	"github.com/daichimukai/x/syakyo/monkey/lexer"
	"github.com/daichimukai/x/syakyo/monkey/parser"
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
		p := parser.New(l)

		program, err := p.ParseProgram()
		if err != nil {
			io.WriteString(out, fmt.Sprintf("parse error: %+s\n", err))
			continue
		}

		if evaluated := eval.Eval(program); evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}
