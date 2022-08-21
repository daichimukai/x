package parser_test

import (
	"testing"

	"github.com/daichimukai/x/syakyo/monkey/ast"
	"github.com/daichimukai/x/syakyo/monkey/lexer"
	"github.com/daichimukai/x/syakyo/monkey/parser"
	"github.com/stretchr/testify/require"
)

func TestLetStatement(t *testing.T) {
	testcases := map[string]struct {
		input         string
		expectedIdent string
	}{
		"right hand side is int": {
			input:         `let x = 5;`,
			expectedIdent: "x",
		},
	}

	for name, tt := range testcases {
		t.Run(name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program, err := p.ParseProgram()
			require.NoError(t, err)
			require.Equal(t, 1, len(program.Statements))

			s := program.Statements[0]
			require.Equal(t, "let", s.TokenLiteral())
			letStmt, ok := s.(*ast.LetStatement)
			require.True(t, ok)

			require.Equal(t, tt.expectedIdent, letStmt.Name.Value)
			require.Equal(t, tt.expectedIdent, letStmt.Name.TokenLiteral())
		})
	}
}
