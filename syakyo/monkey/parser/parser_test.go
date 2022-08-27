package parser_test

import (
	"fmt"
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

func TestReturnStatement(t *testing.T) {
	input := `
return 5;
`

	l := lexer.New(input)
	p := parser.New(l)

	program, err := p.ParseProgram()
	require.NoError(t, err)
	require.Equal(t, 1, len(program.Statements))

	s := program.Statements[0]
	require.Equal(t, "return", s.TokenLiteral())
	_, ok := s.(*ast.ReturnStatement)
	require.True(t, ok)
}

func TestIdentifier(t *testing.T) {
	input := `foobar;`

	l := lexer.New(input)
	p := parser.New(l)
	program, err := p.ParseProgram()
	require.NoError(t, err)

	require.Equal(t, 1, len(program.Statements))
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	ident, ok := stmt.Expression.(*ast.Identifier)
	require.True(t, ok)
	require.Equal(t, "foobar", ident.Value)
	require.Equal(t, "foobar", ident.TokenLiteral())
}

func TestIntegralLiteralExpression(t *testing.T) {
	input := `5;`

	l := lexer.New(input)
	p := parser.New(l)
	program, err := p.ParseProgram()
	require.NoError(t, err)

	require.Equal(t, 1, len(program.Statements))
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	require.True(t, ok)
	require.Equal(t, int64(5), literal.Value)
	require.Equal(t, "5", literal.TokenLiteral())
}

func TestParsingPrefixExpression(t *testing.T) {
	testcases := map[string]struct {
		input        string
		operator     string
		integerValue int64
	}{
		"parse !5;": {
			input:        "!5;",
			operator:     "!",
			integerValue: 5,
		},
		"parse -15;": {
			input:        "-15;",
			operator:     "-",
			integerValue: 15,
		},
	}

	for scenario, tt := range testcases {
		t.Run(scenario, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program, err := p.ParseProgram()
			require.NoError(t, err)

			require.Equal(t, 1, len(program.Statements))
			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			require.True(t, ok)

			expr, ok := stmt.Expression.(*ast.PrefixExpression)
			require.True(t, ok)
			require.Equal(t, tt.operator, expr.Operator)
			testIntegerLiteral(t, tt.integerValue, expr.Right)
		})
	}
}

func testIntegerLiteral(t *testing.T, expect int64, il ast.Expression) {
	integ, ok := il.(*ast.IntegerLiteral)
	require.True(t, ok)
	require.Equal(t, expect, integ.Value)
	require.Equal(t, fmt.Sprintf("%d", expect), integ.TokenLiteral())
}

func TestParsingInfixExpressions(t *testing.T) {
	testcases := map[string]struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		"plus":  {"5 + 5;", 5, "+", 5},
		"minus": {"5 - 5;", 5, "-", 5},
		"mult":  {"5 * 5;", 5, "*", 5},
		"div":   {"5 / 5;", 5, "/", 5},
		"gt":    {"5 > 5;", 5, ">", 5},
		"lt":    {"5 < 5;", 5, "<", 5},
		"eq":    {"5 == 5;", 5, "==", 5},
		"neq":   {"5 != 5;", 5, "!=", 5},
	}

	for scenario, tt := range testcases {
		t.Run(scenario, func(t *testing.T) {
			program := parseProgram(t, tt.input)
			require.Equal(t, 1, len(program.Statements))

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			require.True(t, ok)

			exp, ok := stmt.Expression.(*ast.InfixExpression)
			require.True(t, ok)

			testIntegerLiteral(t, tt.leftValue, exp.Left)
			testIntegerLiteral(t, tt.rightValue, exp.Right)
			require.Equal(t, tt.operator, exp.Operator)
		})
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	testcases := []struct {
		input  string
		expect string
	}{
		{
			input:  "-a * b",
			expect: "((-a) * b)",
		},
		{
			input:  "!-a",
			expect: "(!(-a))",
		},
		{
			input:  "a + b + c",
			expect: "((a + b) + c)",
		},
		{
			input:  "a + b - c",
			expect: "((a + b) - c)",
		},
		{
			input:  "a * b * c",
			expect: "((a * b) * c)",
		},
		{
			input:  "a * b / c",
			expect: "((a * b) / c)",
		},
		{
			input:  "a + b / c",
			expect: "(a + (b / c))",
		},
		{
			input:  "a + b * c + d / e - f",
			expect: "(((a + (b * c)) + (d / e)) - f)",
		},
		{
			input:  "3 + 4; -5 * 5",
			expect: "(3 + 4)((-5) * 5)",
		},
		{
			input:  "5 > 4 == 3 < 4",
			expect: "((5 > 4) == (3 < 4))",
		},
		{
			input:  "5 < 4 != 3 > 4",
			expect: "((5 < 4) != (3 > 4))",
		},
		{
			input:  "3 + 4 * 5 == 3 * 1 + 4 * 5",
			expect: "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
	}

	for _, tt := range testcases {
		t.Run(tt.input, func(t *testing.T) {
			program := parseProgram(t, tt.input)
			require.Equal(t, tt.expect, program.String())
		})
	}
}

func parseProgram(t *testing.T, input string) ast.Program {
	t.Helper()

	l := lexer.New(input)
	p := parser.New(l)
	program, err := p.ParseProgram()
	require.NoError(t, err)
	return *program
}
