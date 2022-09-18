// Package eval implements an evaluator of the monkey language.
package eval

import (
	"github.com/daichimukai/x/syakyo/monkey/ast"
	"github.com/daichimukai/x/syakyo/monkey/object"
)

type Environment struct {
	store map[string]object.Object
	outer *Environment
}

func NewEnvironment() *Environment {
	return &Environment{
		store: make(map[string]object.Object),
	}
}

func (e *Environment) NewEnclosedEnvironment() object.Environment {
	env := NewEnvironment()
	env.outer = e
	return env
}

func (e *Environment) Get(name string) (object.Object, bool) {
	val, ok := e.store[name]
	if !ok && e.outer != nil {
		return e.outer.Get(name)
	}
	return val, ok
}

func (e *Environment) Set(name string, val object.Object) object.Object {
	e.store[name] = val
	return val
}

func (e *Environment) Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return e.evalProgram(node)
	case *ast.ExpressionStatement:
		return e.Eval(node.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.Boolean:
		return object.BooleanFromNative(node.Value)
	case *ast.PrefixExpression:
		right := e.Eval(node.Right)
		if isError(right) {
			return right
		}
		return e.evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := e.Eval(node.Left)
		if isError(left) {
			return left
		}
		right := e.Eval(node.Right)
		if isError(right) {
			return right
		}
		return e.evalInfixExpression(node.Operator, left, right)
	case *ast.BlockStatement:
		return e.evalBlockStatement(node.Statements)
	case *ast.IfExpression:
		return e.evalIfExpression(node)
	case *ast.ReturnStatement:
		val := e.Eval(node.ReturnValue)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.LetStatement:
		val := e.Eval(node.Value)
		if isError(val) {
			return val
		}
		e.Set(node.Name.Value, val)
		return nil
	case *ast.Identifier:
		return e.evalIdentifier(node)
	case *ast.FunctionLiteral:
		return &object.Function{
			Parameters: node.Parameters,
			Body:       node.Body,
			Env:        e,
		}
	case *ast.CallExpression:
		function := e.Eval(node.Function)
		if isError(function) {
			return function
		}
		args := e.evalExpressions(node.Arguments)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return object.ApplyFunction(function, args)
	default:
		return nil
	}
}

func (e *Environment) evalProgram(program *ast.Program) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = e.Eval(statement)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func (e *Environment) evalBlockStatement(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = e.Eval(statement)
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

func (e *Environment) evalPrefixExpression(op string, right object.Object) object.Object {
	switch op {
	case "!":
		return e.evalBangOperatorExpression(right)
	case "-":
		return e.evalMinusOperatorExpression(right)
	default:
		return object.NewError("unknown operator: %s%s", op, right.Type().String())
	}
}

func (e *Environment) evalBangOperatorExpression(right object.Object) object.Object {
	if isTruthy(right) {
		return object.False
	} else {
		return object.True
	}
}

func (e *Environment) evalMinusOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.IntegerObjectType {
		return object.NewError("unknown operator: -%s", right.Type().String())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func (e *Environment) evalInfixExpression(op string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.IntegerObjectType && right.Type() == object.IntegerObjectType:
		return e.evalIntegerInfixExpression(op, left, right)
	case left.Type() == object.StringObjectType && right.Type() == object.StringObjectType:
		return e.evalStringInfixExpression(op, left, right)
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

func (e *Environment) evalIntegerInfixExpression(op string, left, right object.Object) object.Object {
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

func (e *Environment) evalStringInfixExpression(op string, left, right object.Object) object.Object {
	if op != "+" {
		return object.NewError("unknown operator: %s %s %s", left.Type(), op, right.Type())
	}

	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	return &object.String{
		Value: leftVal + rightVal,
	}
}

func (e *Environment) evalIfExpression(ie *ast.IfExpression) object.Object {
	condition := e.Eval(ie.Condition)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return e.Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return e.Eval(ie.Alternative)
	} else {
		return object.Null
	}
}

func (e *Environment) evalIdentifier(node *ast.Identifier) object.Object {
	val, ok := e.Get(node.Value)
	if !ok {
		return object.NewError("identifier not found: %s", node.Value)
	}
	return val
}

func (e *Environment) evalExpressions(exprs []ast.Expression) []object.Object {
	var result []object.Object

	for _, expr := range exprs {
		evaluated := e.Eval(expr)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
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
