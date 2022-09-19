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
	testcases := []struct {
		input         string
		expectedIdent string
		expectedValue interface{}
	}{
		{
			input:         `let x = 5;`,
			expectedIdent: "x",
			expectedValue: 5,
		},
		{
			input:         `let y = true;`,
			expectedIdent: "y",
			expectedValue: true,
		},
		{
			input:         `let foobar = y;`,
			expectedIdent: "foobar",
			expectedValue: "y",
		},
	}

	for _, tt := range testcases {
		t.Run(tt.input, func(t *testing.T) {
			program := parseProgram(t, tt.input)
			require.Equal(t, 1, len(program.Statements))

			s := program.Statements[0]
			require.Equal(t, "let", s.TokenLiteral())
			letStmt, ok := s.(*ast.LetStatement)
			require.True(t, ok)

			testLiteralExpression(t, tt.expectedIdent, letStmt.Name)
			testLiteralExpression(t, tt.expectedValue, letStmt.Value)
		})
	}
}

func TestReturnStatement(t *testing.T) {
	testcases := []struct {
		input  string
		expect interface{}
	}{
		{
			input:  `return 5;`,
			expect: 5,
		},
		{
			input:  `return true;`,
			expect: true,
		},
		{
			input:  `return x;`,
			expect: "x",
		},
	}

	for _, tt := range testcases {
		t.Run(tt.input, func(t *testing.T) {
			program := parseProgram(t, tt.input)
			require.Equal(t, 1, len(program.Statements))
			s := program.Statements[0]
			require.Equal(t, "return", s.TokenLiteral())

			retStmt, ok := s.(*ast.ReturnStatement)
			require.True(t, ok)
			testLiteralExpression(t, tt.expect, retStmt.ReturnValue)
		})
	}
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

	testIdentifier(t, "foobar", stmt.Expression)
}

func testIdentifier(t *testing.T, value string, exp ast.Expression) {
	ident, ok := exp.(*ast.Identifier)
	require.True(t, ok)
	require.Equal(t, value, ident.Value)
	require.Equal(t, value, ident.TokenLiteral())
}

func TestBooleanExpression(t *testing.T) {
	testcases := map[string]struct {
		input  string
		expect bool
	}{
		"true": {
			input:  `true;`,
			expect: true,
		},
		"false": {
			input:  `false;`,
			expect: false,
		},
	}

	for scenario, tt := range testcases {
		t.Run(scenario, func(t *testing.T) {
			program := parseProgram(t, tt.input)
			require.Equal(t, 1, len(program.Statements))
			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			require.True(t, ok)
			testLiteralExpression(t, tt.expect, stmt.Expression)
		})
	}
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

	testLiteralExpression(t, 5, stmt.Expression)
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world"`

	l := lexer.New(input)
	p := parser.New(l)
	program, err := p.ParseProgram()
	require.NoError(t, err)

	require.Equal(t, 1, len(program.Statements))
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	require.True(t, ok)

	require.Equal(t, "hello world", literal.Value)
}

func TestArrayLiteralExpression(t *testing.T) {
	input := `[1, 2 * 2, 3 + 3]`

	l := lexer.New(input)
	p := parser.New(l)
	program, err := p.ParseProgram()
	require.NoError(t, err)

	require.Equal(t, 1, len(program.Statements))
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	require.True(t, ok)
	require.Equal(t, 3, len(array.Elements))

	testLiteralExpression(t, 1, array.Elements[0])
	testInfixExpression(t, 2, "*", 2, array.Elements[1])
	testInfixExpression(t, 3, "+", 3, array.Elements[2])
}

func TestParsingIndexExpression(t *testing.T) {
	input := `myArray[1 + 1];`

	l := lexer.New(input)
	p := parser.New(l)
	program, err := p.ParseProgram()
	require.NoError(t, err)

	require.Len(t, program.Statements, 1)
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	indexExpr, ok := stmt.Expression.(*ast.IndexExpression)
	require.True(t, ok)
	testIdentifier(t, "myArray", indexExpr.Left)
	testInfixExpression(t, 1, "+", 1, indexExpr.Index)
}

func testBoolean(t *testing.T, expected bool, exp ast.Expression) {
	t.Helper()

	b, ok := exp.(*ast.Boolean)
	require.True(t, ok)
	require.Equal(t, expected, b.Value)
	require.Equal(t, fmt.Sprintf("%t", expected), b.TokenLiteral())
}

func testLiteralExpression(t *testing.T, expected interface{}, exp ast.Expression) {
	t.Helper()

	switch v := expected.(type) {
	case int:
		testIntegerLiteral(t, int64(v), exp)
	case int64:
		testIntegerLiteral(t, v, exp)
	case string:
		testIdentifier(t, v, exp)
	case bool:
		testBoolean(t, v, exp)
	default:
		t.Fatalf("bug: unknown type: %q", expected)
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	program := parseProgram(t, input)
	require.Equal(t, 1, len(program.Statements))

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	exp, ok := stmt.Expression.(*ast.IfExpression)
	require.True(t, ok)

	require.Equal(t, 1, len(exp.Consequence.Statements))
	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	testInfixExpression(t, "x", "<", "y", exp.Condition)
	testIdentifier(t, "x", consequence.Expression)
	require.Nil(t, exp.Alternative)
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	program := parseProgram(t, input)
	require.Equal(t, 1, len(program.Statements))

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	exp, ok := stmt.Expression.(*ast.IfExpression)
	require.True(t, ok)

	require.Equal(t, 1, len(exp.Consequence.Statements))
	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	require.Equal(t, 1, len(exp.Alternative.Statements))
	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	testInfixExpression(t, "x", "<", "y", exp.Condition)
	testIdentifier(t, "x", consequence.Expression)
	testIdentifier(t, "y", alternative.Expression)
}

func TestFunctionLiteral(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	program := parseProgram(t, input)
	require.Equal(t, 1, len(program.Statements))

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	require.True(t, ok)

	require.Equal(t, 2, len(function.Parameters))
	testLiteralExpression(t, "x", function.Parameters[0])
	testLiteralExpression(t, "y", function.Parameters[1])

	require.Equal(t, 1, len(function.Body.Statements))
	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)
	testInfixExpression(t, "x", "+", "y", bodyStmt.Expression)
}

func TestFunctionParameters(t *testing.T) {
	testcases := []struct {
		input    string
		expected []string
	}{
		{
			input:    `fn() {};`,
			expected: nil,
		},
		{
			input:    `fn(x) {};`,
			expected: []string{"x"},
		},
		{
			input:    `fn(x,y) {};`,
			expected: []string{"x", "y"},
		},
		{
			input:    `fn(x,y,z) {};`,
			expected: []string{"x", "y", "z"},
		},
	}

	for _, tt := range testcases {
		program := parseProgram(t, tt.input)
		require.Equal(t, 1, len(program.Statements))
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		require.True(t, ok)
		function, ok := stmt.Expression.(*ast.FunctionLiteral)
		require.True(t, ok)

		require.Equal(t, len(tt.expected), len(function.Parameters))
		for i := 0; i < len(tt.expected); i++ {
			testLiteralExpression(t, tt.expected[i], function.Parameters[i])
		}
	}
}

func TestCallExpression(t *testing.T) {
	input := `add(1, 2 * 3, 4 + 5);`
	program := parseProgram(t, input)
	require.Equal(t, 1, len(program.Statements))
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	exp, ok := stmt.Expression.(*ast.CallExpression)
	require.True(t, ok)

	testLiteralExpression(t, "add", exp.Function)
	require.Equal(t, 3, len(exp.Arguments))
	testLiteralExpression(t, 1, exp.Arguments[0])
	testInfixExpression(t, 2, "*", 3, exp.Arguments[1])
	testInfixExpression(t, 4, "+", 5, exp.Arguments[2])
}

func TestParsingPrefixExpression(t *testing.T) {
	testcases := []struct {
		input    string
		operator string
		expected interface{}
	}{
		{
			input:    "!5;",
			operator: "!",
			expected: 5,
		},
		{
			input:    "-15;",
			operator: "-",
			expected: 15,
		},
		{
			input:    "!true;",
			operator: "!",
			expected: true,
		},
		{
			input:    "!false;",
			operator: "!",
			expected: false,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.input, func(t *testing.T) {
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
			testLiteralExpression(t, tt.expected, expr.Right)
		})
	}
}

func testIntegerLiteral(t *testing.T, expect int64, il ast.Expression) {
	t.Helper()

	integ, ok := il.(*ast.IntegerLiteral)
	require.True(t, ok)
	require.Equal(t, expect, integ.Value)
	require.Equal(t, fmt.Sprintf("%d", expect), integ.TokenLiteral())
}

func TestParsingInfixExpressions(t *testing.T) {
	testcases := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"false != false", false, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range testcases {
		t.Run(tt.input, func(t *testing.T) {
			program := parseProgram(t, tt.input)
			require.Equal(t, 1, len(program.Statements))

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			require.True(t, ok)

			testInfixExpression(t, tt.leftValue, tt.operator, tt.rightValue, stmt.Expression)
		})
	}
}

func testInfixExpression(
	t *testing.T,
	left interface{}, operator string, right interface{},
	exp ast.Expression,
) {
	infix, ok := exp.(*ast.InfixExpression)
	require.True(t, ok)

	testLiteralExpression(t, left, infix.Left)
	testLiteralExpression(t, right, infix.Right)
	require.Equal(t, operator, infix.Operator)
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	testcases := []struct {
		input  string
		expect string
	}{
		{
			input:  "true",
			expect: "true",
		},
		{
			input:  "false",
			expect: "false",
		},
		{
			input:  "3 > 5 == false",
			expect: "((3 > 5) == false)",
		},
		{
			input:  "3 < 5 == true",
			expect: "((3 < 5) == true)",
		},
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
		{
			input:  "1 + (2 + 3) + 4",
			expect: "((1 + (2 + 3)) + 4)",
		},
		{
			input:  "(5 + 5) * 2",
			expect: "((5 + 5) * 2)",
		},
		{
			input:  "-(5 + 5)",
			expect: "(-(5 + 5))",
		},
		{
			input:  "!(true == true)",
			expect: "(!(true == true))",
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
