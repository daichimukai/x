package lexer_test

import (
	"testing"

	"github.com/daichimukai/x/syakyo/monkey/lexer"
	"github.com/daichimukai/x/syakyo/monkey/token"
	"github.com/stretchr/testify/require"
)

func TestNextToken_SingleToken(t *testing.T) {
	testCases := map[string]struct {
		input           string
		expectedType    token.TokenType
		expectedLiteral string
	}{
		"eof":          {"", token.TypeEof, ""},
		"ident":        {"foo", token.TypeIdent, "foo"},
		"int":          {"0", token.TypeInt, "0"},
		"string":       {`"foo"`, token.TypeString, `foo`},
		"assign":       {"=", token.TypeAssign, "="},
		"plus":         {"+", token.TypePlus, "+"},
		"minus":        {"-", token.TypeMinus, "-"},
		"bang":         {"!", token.TypeBang, "!"},
		"asterisk":     {"*", token.TypeAsterisk, "*"},
		"slash":        {"/", token.TypeSlash, "/"},
		"less than":    {"<", token.TypeLt, "<"},
		"greater than": {">", token.TypeGt, ">"},
		"equal":        {"==", token.TypeEq, "=="},
		"not equal":    {"!=", token.TypeNotEq, "!="},
		"comma":        {",", token.TypeComma, ","},
		"semicolon":    {";", token.TypeSemicolon, ";"},
		"left paren":   {"(", token.TypeLeftParen, "("},
		"right paren":  {")", token.TypeRightParen, ")"},
		"left brace":   {"{", token.TypeLeftBrace, "{"},
		"right brace":  {"}", token.TypeRightBrace, "}"},
		"left braket":  {"[", token.TypeLeftBraket, "["},
		"right braket": {"]", token.TypeRightBraket, "]"},
		"function":     {"fn", token.TypeFunction, "fn"},
		"let":          {"let", token.TypeLet, "let"},
		"true":         {"true", token.TypeTrue, "true"},
		"false":        {"false", token.TypeFalse, "false"},
		"if":           {"if", token.TypeIf, "if"},
		"else":         {"else", token.TypeElse, "else"},
		"return":       {"return", token.TypeReturn, "return"},
	}

	for name, tt := range testCases {
		l := lexer.New(tt.input)
		t.Run(name, func(t *testing.T) {
			tok := l.NextToken()
			require.Equal(t, tt.expectedType, tok.Type)
			require.Equal(t, tt.expectedLiteral, tok.Literal)

			tok = l.NextToken()
			require.Equal(t, token.TypeEof, tok.Type)
			require.Equal(t, "", tok.Literal)
		})
	}
}

func TestNextToken_Whitespace(t *testing.T) {
	testCases := map[string]string{
		"space": ` `,
		"tab": `	`,
		"newline": `
`,
	}

	for name, input := range testCases {
		l := lexer.New(input)
		t.Run(name, func(t *testing.T) {
			tok := l.NextToken()
			require.Equal(t, token.TypeEof, tok.Type)
			require.Equal(t, "", tok.Literal)
		})
	}
}

func TestNextToken_SimpleProgram(t *testing.T) {
	input := `let five = 5;
let ten = 10;
let add = fn(x, y) {
  x + y;
};

let result = add(five, ten);
let str = "foo";
`

	expectedTokens := []struct {
		expectedTokenType token.TokenType
		expectedLiteral   string
	}{
		{token.TypeLet, "let"},
		{token.TypeIdent, "five"},
		{token.TypeAssign, "="},
		{token.TypeInt, "5"},
		{token.TypeSemicolon, ";"},
		{token.TypeLet, "let"},
		{token.TypeIdent, "ten"},
		{token.TypeAssign, "="},
		{token.TypeInt, "10"},
		{token.TypeSemicolon, ";"},
		{token.TypeLet, "let"},
		{token.TypeIdent, "add"},
		{token.TypeAssign, "="},
		{token.TypeFunction, "fn"},
		{token.TypeLeftParen, "("},
		{token.TypeIdent, "x"},
		{token.TypeComma, ","},
		{token.TypeIdent, "y"},
		{token.TypeRightParen, ")"},
		{token.TypeLeftBrace, "{"},
		{token.TypeIdent, "x"},
		{token.TypePlus, "+"},
		{token.TypeIdent, "y"},
		{token.TypeSemicolon, ";"},
		{token.TypeRightBrace, "}"},
		{token.TypeSemicolon, ";"},
		{token.TypeLet, "let"},
		{token.TypeIdent, "result"},
		{token.TypeAssign, "="},
		{token.TypeIdent, "add"},
		{token.TypeLeftParen, "("},
		{token.TypeIdent, "five"},
		{token.TypeComma, ","},
		{token.TypeIdent, "ten"},
		{token.TypeRightParen, ")"},
		{token.TypeSemicolon, ";"},
		{token.TypeLet, "let"},
		{token.TypeIdent, "str"},
		{token.TypeAssign, "="},
		{token.TypeString, "foo"},
		{token.TypeSemicolon, ";"},
		{token.TypeEof, ""},
	}

	l := lexer.New(input)
	for i, tt := range expectedTokens {
		tok := l.NextToken()
		require.Equalf(t, tt.expectedLiteral, tok.Literal, "literal differs for the token at %d", i)
		require.Equalf(t, tt.expectedTokenType, tok.Type, "type differs for the token at %d", i)
	}
}
