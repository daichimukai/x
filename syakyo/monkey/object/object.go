// Package object implements the object system of the monkey language.
package object

import "fmt"

//go:generate stringer -type ObjectType -linecomment
type ObjectType int

const (
	IntegerObjectType     ObjectType = iota // INTEGER
	BooleanObjectType                       // BOOLEAN
	NullObjectType                          // NULL
	ReturnValueObjectType                   // RETURN_VALUE
	ErrorObjectType                         // ERROR
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
