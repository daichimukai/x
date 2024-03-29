package eval_test

import (
	"testing"

	"github.com/daichimukai/x/syakyo/monkey/eval"
	"github.com/daichimukai/x/syakyo/monkey/lexer"
	"github.com/daichimukai/x/syakyo/monkey/object"
	"github.com/daichimukai/x/syakyo/monkey/parser"
	"github.com/stretchr/testify/require"
)

func TestEvalBooleanExpression(t *testing.T) {
	testcases := []struct {
		input  string
		expect bool
	}{
		{
			input:  `true`,
			expect: true,
		},
		{
			input:  `false`,
			expect: false,
		},
		{
			input:  `0 == 0`,
			expect: true,
		},
		{
			input:  `0 == 1`,
			expect: false,
		},
		{
			input:  `1 == 0`,
			expect: false,
		},
		{
			input:  `0 != 0`,
			expect: false,
		},
		{
			input:  `0 != 1`,
			expect: true,
		},
		{
			input:  `1 != 0`,
			expect: true,
		},
		{
			input:  `0 > 0`,
			expect: false,
		},
		{
			input:  `0 > 1`,
			expect: false,
		},
		{
			input:  `1 > 0`,
			expect: true,
		},
		{
			input:  `0 < 0`,
			expect: false,
		},
		{
			input:  `0 < 1`,
			expect: true,
		},
		{
			input:  `1 < 0`,
			expect: false,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.input, func(t *testing.T) {
			evaluated := testEval(t, tt.input)
			testBooleanObject(t, tt.expect, evaluated)
		})
	}
}

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
		{
			input:  `-5;`,
			expect: -5,
		},
		{
			input:  `-10;`,
			expect: -10,
		},
		{
			input:  `1 + 2;`,
			expect: 3,
		},
		{
			input:  `1 - 2;`,
			expect: -1,
		},
		{
			input:  `1 * 2;`,
			expect: 2,
		},
		{
			input:  `1 / 2;`,
			expect: 0,
		},
		{
			input:  `5 + 5 + 5 + 5 - 10;`,
			expect: 10,
		},
		{
			input:  `-50 + 100 + -50;`,
			expect: 0,
		},
		{
			input:  `2 * (5 + 10);`,
			expect: 30,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.input, func(t *testing.T) {
			evaluated := testEval(t, tt.input)
			testIntegerObject(t, tt.expect, evaluated)
		})
	}
}

func TestEvalStringLiteral(t *testing.T) {
	input := `"Hello world!"`

	evaluated := testEval(t, input)
	str, ok := evaluated.(*object.String)
	require.True(t, ok)
	require.Equal(t, "Hello world!", str.Value)
}

func TestEvalArrayLiteral(t *testing.T) {
	input := `[1, 2 * 2, 3 + 3]`

	evaluated := testEval(t, input)
	arr, ok := evaluated.(*object.Array)
	require.True(t, ok)

	require.Len(t, arr.Elements, 3)
	testIntegerObject(t, 1, arr.Elements[0])
	testIntegerObject(t, 4, arr.Elements[1])
	testIntegerObject(t, 6, arr.Elements[2])
}

func TestEvalArrayIndexExpressions(t *testing.T) {
	testcases := []struct {
		input  string
		expect any
	}{
		{
			input:  `[1,2,3][0]`,
			expect: 1,
		},
		{
			input:  `[1,2,3][1]`,
			expect: 2,
		},
		{
			input:  `[1,2,3][2]`,
			expect: 3,
		},
		{
			input:  `let i = 0; [1][i]`,
			expect: 1,
		},
		{
			input:  `[1, 2, 3][1 + 1]`,
			expect: 3,
		},
		{
			input:  `let myArray = [1, 2, 3]; myArray[2];`,
			expect: 3,
		},
		{
			input:  `let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];`,
			expect: 6,
		},
		{
			input:  `let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i];`,
			expect: 2,
		},
		{
			input:  `[1,2,3][3]`,
			expect: nil,
		},
		{
			input:  `[1,2,3][-1]`,
			expect: nil,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.input, func(t *testing.T) {
			evaluated := testEval(t, tt.input)
			integer, ok := tt.expect.(int)
			if ok {
				testIntegerObject(t, int64(integer), evaluated)
			} else {
				testNullObject(t, evaluated)
			}
		})
	}
}

func TestEvalStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`

	evaluated := testEval(t, input)
	str, ok := evaluated.(*object.String)
	require.True(t, ok)
	require.Equal(t, "Hello World!", str.Value)
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

func TestIfExpression(t *testing.T) {
	testcases := []struct {
		input  string
		expect interface{}
	}{
		{
			input:  `if (true) { 10 }`,
			expect: 10,
		},
		{
			input:  `if (false) { 10 }`,
			expect: nil,
		},
		{
			input:  `if (true) { 10 } else { 20 }`,
			expect: 10,
		},
		{
			input:  `if (false) { 10 } else { 20 }`,
			expect: 20,
		},
		{
			input:  `if (1) { 10 }`,
			expect: 10,
		},
		{
			input:  `if (0) { 10 }`, // ¯\_(ツ)_/¯
			expect: 10,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.input, func(t *testing.T) {
			evaluated := testEval(t, tt.input)
			integer, ok := tt.expect.(int)
			if ok {
				testIntegerObject(t, int64(integer), evaluated)
			} else {
				testNullObject(t, evaluated)
			}
		})
	}
}

func TestReturnStatements(t *testing.T) {
	testcases := []struct {
		input  string
		expect int64
	}{
		{
			input:  `return 10;`,
			expect: 10,
		},
		{
			input:  `return 10; 9;`,
			expect: 10,
		},
		{
			input:  `9; return 2 * 5; 10;`,
			expect: 10,
		},
		{
			input: `
			if (10 > 1) {
				if (10 > 1) {
					return 10;
				}
				return 1;
			}`,
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

func TestLetStatements(t *testing.T) {
	testcases := []struct {
		input  string
		expect int64
	}{
		{
			input:  `let a = 5; a;`,
			expect: 5,
		},
		{
			input:  `let a = 5 * 5; a;`,
			expect: 25,
		},
		{
			input:  `let a = 5; let b = a; b;`,
			expect: 5,
		},
		{
			input: `
			let a = 5;
			let b = a;
			let c = a + b + 5;
			c;`,
			expect: 15,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.input, func(t *testing.T) {
			testIntegerObject(t, tt.expect, testEval(t, tt.input))
		})
	}
}

func TestFunctionObject(t *testing.T) {
	input := `fn(x) { x + 2; };`

	evaluated := testEval(t, input)
	fn, ok := evaluated.(*object.Function)
	require.True(t, ok)

	require.Equal(t, 1, len(fn.Parameters))
	require.Equal(t, "x", fn.Parameters[0].Value)
	require.Equal(t, "(x + 2)", fn.Body.String())
}

func TestFunctionApplication(t *testing.T) {
	testcases := []struct {
		input  string
		expect int64
	}{
		{
			input:  `let id = fn(x) { x; }; id(5);`,
			expect: 5,
		},
		{
			input:  `let id = fn(x) { return x; }; id(5);`,
			expect: 5,
		},
		{
			input:  `let double = fn(x) { x * 2; }; double(5);`,
			expect: 10,
		},
		{
			input:  `let add = fn(x, y) { x + y; }; add(5, 5);`,
			expect: 10,
		},
		{
			input: `
			let add = fn(x, y) { x + y; };
			add(5 + 5, add(5, 5));`,
			expect: 20,
		},
		{
			input:  `fn(x) { x; }(5);`,
			expect: 5,
		},
		{
			input: `
			let newAdder = fn(x) { fn(y) { x + y; } };
			let addTwo = newAdder(2);
			addTwo(2);`,
			expect: 4,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.input, func(t *testing.T) {
			testIntegerObject(t, tt.expect, testEval(t, tt.input))
		})
	}

}

func TestBuiltinFunctions(t *testing.T) {
	testcases := []struct {
		input  string
		expect int64
	}{
		{
			input:  `len("")`,
			expect: 0,
		},
		{
			input:  `len("four")`,
			expect: 4,
		},
		{
			input:  `len("Hello World!")`,
			expect: 12,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.input, func(t *testing.T) {
			evaluated := testEval(t, tt.input)
			i, ok := evaluated.(*object.Integer)
			require.True(t, ok)
			require.Equal(t, tt.expect, i.Value)
		})
	}
}

func TestErrorHandling(t *testing.T) {
	testcases := []struct {
		input  string
		expect string
	}{
		{
			input:  `5 + true`,
			expect: "type mismatch: INTEGER + BOOLEAN",
		},
		{
			input:  `5 + true; 5;`,
			expect: "type mismatch: INTEGER + BOOLEAN",
		},
		{
			input:  `-true`,
			expect: "unknown operator: -BOOLEAN",
		},
		{
			input:  `true + false`,
			expect: "unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			input:  `5; true + false; 5;`,
			expect: "unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			input:  `if (10 > 1) { true + false; }`,
			expect: "unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			input:  `foobar`,
			expect: "identifier not found: foobar",
		},
		{
			input:  `"foo" - "bar"`,
			expect: "unknown operator: STRING - STRING",
		},
		{
			input:  `len()`,
			expect: "wrong number of arguments: got=0, want=1",
		},
		{
			input:  `len(1)`,
			expect: "argument to `len` not supported: got INTEGER",
		},
		{
			input:  `len("one", "two")`,
			expect: "wrong number of arguments: got=2, want=1",
		},
	}

	for _, tt := range testcases {
		t.Run(tt.input, func(t *testing.T) {
			evaluated := testEval(t, tt.input)

			errObj, ok := evaluated.(*object.Error)
			require.True(t, ok)
			require.Equal(t, tt.expect, errObj.Message)
		})
	}
}

func testEval(t *testing.T, input string) object.Object {
	t.Helper()
	l := lexer.New(input)
	p := parser.New(l)
	env := eval.NewEnvironment()
	program, err := p.ParseProgram()
	require.NoError(t, err)
	return env.Eval(program)
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

func testNullObject(t *testing.T, obj object.Object) {
	t.Helper()
	require.Equal(t, object.Null, obj)
}
