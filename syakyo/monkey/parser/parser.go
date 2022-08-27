// Package parser implements a parser for Monkey source codes.
package parser

import (
	"fmt"

	"github.com/daichimukai/x/syakyo/monkey/ast"
	"github.com/daichimukai/x/syakyo/monkey/lexer"
	"github.com/daichimukai/x/syakyo/monkey/token"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	// Set curToken and peekToken
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() (*ast.Program, error) {
	program := &ast.Program{}

	for p.curToken.Type != token.TypeEof {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		program.Statements = append(program.Statements, stmt)
		p.nextToken()
	}
	return program, nil
}

func (p *Parser) parseStatement() (ast.Statement, error) {
	switch p.curToken.Type {
	case token.TypeLet:
		return p.parseLetStatement()
	case token.TypeReturn:
		return p.parseReturnStatement()
	}
	return nil, fmt.Errorf("unexpected token: %s", p.curToken.Literal)
}

func (p *Parser) parseLetStatement() (*ast.LetStatement, error) {
	stmt := &ast.LetStatement{
		Token: p.curToken,
	}
	if ok := p.expectPeek(token.TypeIdent); !ok {
		return nil, fmt.Errorf("expected identifier, got %s", p.peekToken.Literal)
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if ok := p.expectPeek(token.TypeAssign); !ok {
		return nil, fmt.Errorf("expected =, got %s", p.peekToken.Literal)
	}

	// TODO: for now, skip while got semicolon
	for p.curToken.Type != token.TypeSemicolon {
		p.nextToken()
	}

	return stmt, nil
}

func (p *Parser) parseReturnStatement() (*ast.ReturnStatement, error) {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	// TODO: for now, skip until semicolon
	for p.curToken.Type != token.TypeSemicolon {
		p.nextToken()
	}

	return stmt, nil
}

func (p *Parser) expectPeek(typ token.TokenType) bool {
	if p.peekToken.Type != typ {
		return false
	}
	p.nextToken()
	return true
}
