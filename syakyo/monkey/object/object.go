// Package object implements the object system of the monkey language.
package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/daichimukai/x/syakyo/monkey/ast"
)

//go:generate stringer -type ObjectType -linecomment
type ObjectType int

const (
	IntegerObjectType     ObjectType = iota // INTEGER
	StringObjectType                        // STRING
	BooleanObjectType                       // BOOLEAN
	NullObjectType                          // NULL
	ReturnValueObjectType                   // RETURN_VALUE
	ErrorObjectType                         // ERROR
	FunctionObjectType                      // FUNCTION
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return IntegerObjectType }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return StringObjectType }
func (s *String) Inspect() string  { return s.Value }

var (
	True  = &boolean{Value: true}
	False = &boolean{Value: false}
)

func BooleanFromNative(b bool) *boolean {
	if b {
		return True
	} else {
		return False
	}
}

type boolean struct {
	Value bool
}

func (b *boolean) Type() ObjectType { return BooleanObjectType }
func (b *boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

var Null = &null{}

type null struct{}

func (n *null) Type() ObjectType { return NullObjectType }
func (n *null) Inspect() string  { return "null" }

// ReturnValue wraps an object that is returned from a block.
type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return ReturnValueObjectType }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

// Error is an object that means some error happend.
type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ErrorObjectType }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

func NewError(format string, a ...interface{}) *Error {
	return &Error{
		Message: fmt.Sprintf(format, a...),
	}
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        Environment
}

type Environment interface {
	NewEnclosedEnvironment() Environment
	Eval(ast.Node) Object
	Set(string, Object) Object
}

func (f *Function) Type() ObjectType { return FunctionObjectType }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	var params []string
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn(")
	out.WriteString(strings.Join(params, ","))
	out.WriteString("{\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

func ApplyFunction(fn Object, args []Object) Object {
	function, ok := fn.(*Function)
	if !ok {
		return NewError("not a function: %s", fn.Type().String())
	}

	extendedEnv := function.Env.NewEnclosedEnvironment()
	for i, param := range function.Parameters {
		extendedEnv.Set(param.Value, args[i])
	}

	evaluated := extendedEnv.Eval(function.Body)
	if retVal, ok := evaluated.(*ReturnValue); ok {
		return retVal.Value
	}
	return evaluated
}
