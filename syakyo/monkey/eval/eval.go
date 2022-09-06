// Package eval implements an evaluator of the monkey language.
package eval

import (
	"github.com/daichimukai/x/syakyo/monkey/ast"
	"github.com/daichimukai/x/syakyo/monkey/object"
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return object.BooleanFromNative(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		if isError(left) {
			return left
		}
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.BlockStatement:
		return evalBlockStatement(node.Statements)
	case *ast.IfExpression:
		return evalIfExpression(node)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	default:
		return nil
	}
}

func evalProgram(program *ast.Program) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement)
		if result == nil {
			continue
		}

		rt := result.Type()
		if rt == object.ReturnValueObjectType || rt == object.ErrorObjectType {
			return result
		}
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
		return object.NewError("unknown operator: %s%s", op, right.Type().String())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	if isTruthy(right) {
		return object.False
	} else {
		return object.True
	}
}

func evalMinusOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.IntegerObjectType {
		return object.NewError("unknown operator: -%s", right.Type().String())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(op string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.IntegerObjectType && right.Type() == object.IntegerObjectType:
		return evalIntegerInfixExpression(op, left, right)
	case left.Type() != right.Type():
		return object.NewError(
			"type mismatch: %s %s %s",
			left.Type().String(), op, right.Type().String(),
		)
	default:
		return object.NewError(
			"unknown operator: %s %s %s",
			left.Type().String(), op, right.Type().String(),
		)
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
		return object.NewError(
			"unknown operator: %s %s %s",
			left.Type().String(), op, right.Type().String(),
		)
	}
	return &object.Integer{Value: value}
}

func evalIfExpression(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	} else {
		return object.Null
	}
}

func isTruthy(obj object.Object) bool {
	return !(obj == object.False || obj == object.Null)
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ErrorObjectType
	}
	return false
}
