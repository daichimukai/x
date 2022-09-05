// Package eval implements an evaluator of the monkey language.
package eval

import (
	"github.com/daichimukai/x/syakyo/monkey/ast"
	"github.com/daichimukai/x/syakyo/monkey/object"
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return object.BooleanFromNative(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)
	default:
		return nil
	}
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement)
	}

	return result
}

func evalPrefixExpression(op string, right object.Object) object.Object {
	switch op {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusOperatorExpression(right)
	default:
		return object.Null
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	if right == object.False || right == object.Null {
		return object.True
	}
	return object.False
}

func evalMinusOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.IntegerObjectType {
		return object.Null
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(op string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.IntegerObjectType && right.Type() == object.IntegerObjectType:
		return evalIntegerInfixExpression(op, left, right)
	default:
		return object.Null
	}
}

func evalIntegerInfixExpression(op string, left, right object.Object) object.Object {
	lvalue := left.(*object.Integer).Value
	rvalue := right.(*object.Integer).Value
	var value int64
	switch op {
	case "+":
		value = lvalue + rvalue
	case "-":
		value = lvalue - rvalue
	case "*":
		value = lvalue * rvalue
	case "/":
		value = lvalue / rvalue
	case "==":
		return object.BooleanFromNative(lvalue == rvalue)
	case "!=":
		return object.BooleanFromNative(lvalue != rvalue)
	case "<":
		return object.BooleanFromNative(lvalue < rvalue)
	case ">":
		return object.BooleanFromNative(lvalue > rvalue)
	default:
		return object.Null
	}
	return &object.Integer{Value: value}
}
