package eval_test

import (
	"testing"

	"github.com/daichimukai/x/syakyo/monkey/eval"
	"github.com/daichimukai/x/syakyo/monkey/lexer"
	"github.com/daichimukai/x/syakyo/monkey/object"
	"github.com/daichimukai/x/syakyo/monkey/parser"
	"github.com/stretchr/testify/require"
)

func TestEvalIntegerExpression(t *testing.T) {
	testcases := []struct {
		input  string
		expect int64
	}{
		{
			input:  `5;`,
			expect: 5,
		},
		{
			input:  `10;`,
			expect: 10,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.input, func(t *testing.T) {
			evaluated := testEval(t, tt.input)
			testIntegerObject(t, tt.expect, evaluated)
		})
	}
}

func TestEvalBangOperator(t *testing.T) {
	testcases := []struct {
		input  string
		expect bool
	}{
		{
			input:  `!true`,
			expect: false,
		},
		{
			input:  `!false`,
			expect: true,
		},
		{
			input:  `!5`,
			expect: false,
		},
		{
			input:  `!!true`,
			expect: true,
		},
		{
			input:  `!!false`,
			expect: false,
		},
		{
			input:  `!!5`,
			expect: true,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.input, func(t *testing.T) {
			evaluated := testEval(t, tt.input)
			testBooleanObject(t, tt.expect, evaluated)
		})
	}
}

func testEval(t *testing.T, input string) object.Object {
	t.Helper()
	l := lexer.New(input)
	p := parser.New(l)
	program, err := p.ParseProgram()
	require.NoError(t, err)
	return eval.Eval(program)
}

func testIntegerObject(t *testing.T, expect int64, obj object.Object) {
	t.Helper()
	result, ok := obj.(*object.Integer)
	require.True(t, ok)
	require.Equal(t, expect, result.Value)
}

func testBooleanObject(t *testing.T, expect bool, obj object.Object) {
	t.Helper()
	require.Equal(t, object.BooleanFromNative(expect), obj)
}
