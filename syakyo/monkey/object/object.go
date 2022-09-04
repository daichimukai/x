// Package object implements the object system of the monkey language.
package object

import "fmt"

type ObjectType int

const (
	_ ObjectType = iota
	IntegerObjectType
	BooleanObjectType
	NullObjectType
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

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BooleanObjectType }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

type Null struct{}

func (n *Null) Type() ObjectType { return NullObjectType }
func (n *Null) Inspect() string  { return "null" }
