// Package parser implements a parser for Monkey source codes.
package parser

import (
	"fmt"
	"strconv"

	"github.com/daichimukai/x/syakyo/monkey/ast"
	"github.com/daichimukai/x/syakyo/monkey/lexer"
	"github.com/daichimukai/x/syakyo/monkey/token"
)

const (
	priorityLowest      int = iota
	priorityEquals          // ==
	priorityLessGreater     // > or <
	prioritySum             // +
	priorityProduct         // *
	priorityPrefix          // -X or !X
	priorityCall            // X(Y)
)

var precedences = map[token.TokenType]int{
	token.TypeEq:       priorityEquals,
	token.TypeNotEq:    priorityEquals,
	token.TypeLt:       priorityLessGreater,
	token.TypeGt:       priorityLessGreater,
	token.TypePlus:     prioritySum,
	token.TypeMinus:    prioritySum,
	token.TypeAsterisk: priorityProduct,
	token.TypeSlash:    priorityProduct,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:              l,
		prefixParseFns: map[token.TokenType]prefixParseFn{},
		infixParseFns:  map[token.TokenType]infixParseFn{},
	}

	p.registerPrefix(token.TypeTrue, p.parseBoolean)
	p.registerPrefix(token.TypeFalse, p.parseBoolean)
	p.registerPrefix(token.TypeIdent, p.parseIdentifier)
	p.registerPrefix(token.TypeInt, p.parseIntegerLiteral)
	p.registerPrefix(token.TypeMinus, p.parsePrefixExpression)
	p.registerPrefix(token.TypeBang, p.parsePrefixExpression)
	p.registerPrefix(token.TypeLeftParen, p.parseGroupedExpression)
	p.registerPrefix(token.TypeIf, p.parseIfExpression)
	p.registerPrefix(token.TypeFunction, p.parseFunctionLiteral)

	p.registerInfix(token.TypePlus, p.parseInfixExpression)
	p.registerInfix(token.TypeMinus, p.parseInfixExpression)
	p.registerInfix(token.TypeSlash, p.parseInfixExpression)
	p.registerInfix(token.TypeAsterisk, p.parseInfixExpression)
	p.registerInfix(token.TypeEq, p.parseInfixExpression)
	p.registerInfix(token.TypeNotEq, p.parseInfixExpression)
	p.registerInfix(token.TypeLt, p.parseInfixExpression)
	p.registerInfix(token.TypeGt, p.parseInfixExpression)

	// Set curToken and peekToken
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return priorityLowest
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return priorityLowest
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
	default:
		return p.parseExpressionStatement()
	}
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

func (p *Parser) parseExpressionStatement() (*ast.ExpressionStatement, error) {
	stmt := &ast.ExpressionStatement{
		Token:      p.curToken,
		Expression: p.parseExpression(priorityLowest),
	}

	if p.peekToken.Type == token.TypeSemicolon {
		p.nextToken()
	}

	return stmt, nil
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		return nil
	}
	expr := prefix()

	for p.peekToken.Type != token.TypeSemicolon && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			break
		}
		p.nextToken()
		expr = infix(expr)
	}

	return expr
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{
		Token: p.curToken,
		Value: p.curToken.Type == token.TypeTrue,
	}
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		return nil
	}

	return &ast.IntegerLiteral{
		Token: p.curToken,
		Value: value,
	}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	expression.Right = p.parseExpression(priorityPrefix)

	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(priorityLowest)
	if !p.expectPeek(token.TypeRightParen) {
		return nil
	}

	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	expr := &ast.IfExpression{
		Token: p.curToken,
	}

	if !p.expectPeek(token.TypeLeftParen) {
		return nil
	}
	p.nextToken()
	expr.Condition = p.parseExpression(priorityLowest)
	if !p.expectPeek(token.TypeRightParen) {
		return nil
	}

	if !p.expectPeek(token.TypeLeftBrace) {
		return nil
	}
	expr.Consequence = p.parseBlockStatement()
	if p.peekToken.Type == token.TypeElse {
		p.nextToken()
		if !p.expectPeek(token.TypeLeftBrace) {
			return nil
		}
		expr.Alternative = p.parseBlockStatement()
	}

	return expr
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{
		Token: p.curToken,
	}
	p.nextToken()

	for p.curToken.Type != token.TypeRightBrace && p.curToken.Type != token.TypeEof {
		if stmt, err := p.parseStatement(); err == nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{
		Token: p.curToken,
	}
	if !p.expectPeek(token.TypeLeftParen) {
		return nil
	}
	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.TypeLeftBrace) {
		return nil
	}
	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	var identifiers []*ast.Identifier

	if p.peekToken.Type == token.TypeRightParen {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
	identifiers = append(identifiers, ident)

	for p.peekToken.Type == token.TypeComma {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{
			Token: p.curToken,
			Value: p.curToken.Literal,
		}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.TypeRightParen) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) expectPeek(typ token.TokenType) bool {
	if p.peekToken.Type != typ {
		return false
	}
	p.nextToken()
	return true
}
