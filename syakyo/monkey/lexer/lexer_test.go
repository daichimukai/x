package lexer_test

import (
	"testing"

	"github.com/daichimukai/x/syakyo/monkey/lexer"
	"github.com/daichimukai/x/syakyo/monkey/token"
	"github.com/stretchr/testify/require"
)

func TestNextToken(t *testing.T) {
	input := `=+(){},;`

	testCases := []struct {
		name            string
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{"assign", token.TypeAssign, "="},
		{"plus", token.TypePlus, "+"},
		{"left paren", token.TypeLeftParen, "("},
		{"right paren", token.TypeRightParen, ")"},
		{"left brace", token.TypeLeftBrace, "{"},
		{"right brace", token.TypeRightBrace, "}"},
		{"comma", token.TypeComma, ","},
		{"semicolon", token.TypeSemicolon, ";"},
		{"eof", token.TypeEof, ""},
	}

	l := lexer.New(input)
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tok := l.NextToken()
			require.Equal(t, tt.expectedType, tok.Type)
			require.Equal(t, tt.expectedLiteral, tok.Literal)
		})
	}
}
