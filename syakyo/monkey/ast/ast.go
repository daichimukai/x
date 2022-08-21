package ast

import "github.com/daichimukai/x/syakyo/monkey/token"

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	return ""
}

type LetStatement struct {
	Statement

	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

type Identifier struct {
	Expression

	Token token.Token
	Value string
}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
