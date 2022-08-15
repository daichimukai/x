package lexer

import (
	"github.com/daichimukai/x/syakyo/monkey/token"
)

type Lexer struct {
	input        string
	readPosition int  // position which we will read next
	position     int  // position which we had read
	ch           byte // char at `position`
}

// New returns a new lexer of `input`.
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) skipWhitespaces() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

var twoByteTokens map[string]token.TokenType = map[string]token.TokenType{
	"==": token.TypeEq,
	"!=": token.TypeNotEq,
}

var byteToTokenTypeMap map[byte]token.TokenType = map[byte]token.TokenType{
	'=': token.TypeAssign,
	'+': token.TypePlus,
	'-': token.TypeMinus,
	'!': token.TypeBang,
	'*': token.TypeAsterisk,
	'/': token.TypeSlash,
	'<': token.TypeLt,
	'>': token.TypeGt,
	'(': token.TypeLeftParen,
	')': token.TypeRightParen,
	'{': token.TypeLeftBrace,
	'}': token.TypeRightBrace,
	',': token.TypeComma,
	';': token.TypeSemicolon,
}

// NextToken returns the next token from the input.
func (l *Lexer) NextToken() token.Token {
	l.skipWhitespaces()

	var typ token.TokenType
	var literal string
	if l.ch == 0 {
		literal = ""
		typ = token.TypeEof
	} else if v, ok := twoByteTokens[string(l.ch)+string(l.peekChar())]; ok {
		literal = string(l.ch) + string(l.peekChar())
		typ = v
		l.readChar()
		l.readChar()
	} else if v, ok := byteToTokenTypeMap[l.ch]; ok {
		literal = string(l.ch)
		typ = v
		l.readChar()
	} else if isLetter(l.ch) {
		literal = l.readIdentifier()
		typ = token.LookupIdent(literal)
	} else if isDigit(l.ch) {
		literal = l.readNumber()
		typ = token.TypeInt
	}

	return token.Token{
		Type:    typ,
		Literal: literal,
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
