package ast_test

import (
	"testing"

	"github.com/daichimukai/x/syakyo/monkey/ast"
	"github.com/daichimukai/x/syakyo/monkey/token"
	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.LetStatement{
				Token: token.Token{
					Type:    token.TypeLet,
					Literal: "let",
				},
				Name: &ast.Identifier{
					Token: token.Token{
						Type:    token.TypeIdent,
						Literal: "myVar",
					},
					Value: "myVar",
				},
				Value: &ast.Identifier{
					Token: token.Token{
						Type:    token.TypeIdent,
						Literal: "anotherVar",
					},
					Value: "anotherVar",
				},
			},
		},
	}

	require.Equal(t, "let myVar = anotherVar;", program.String())
}
