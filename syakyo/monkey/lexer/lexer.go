package lexer

import "github.com/daichimukai/x/syakyo/monkey/token"

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

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

var byteToTokenTypeMap map[byte]token.TokenType = map[byte]token.TokenType{
	'=': token.TypeAssign,
	'+': token.TypePlus,
	'(': token.TypeLeftParen,
	')': token.TypeRightParen,
	'{': token.TypeLeftBrace,
	'}': token.TypeRightBrace,
	',': token.TypeComma,
	';': token.TypeSemicolon,
	0:   token.TypeEof,
}

// NextToken returns the next token from the input.
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	var ok bool
	if tok.Type, ok = byteToTokenTypeMap[l.ch]; !ok {
		return tok
	}
	if tok.Type != token.TypeEof {
		tok.Literal = string(l.ch)
	} else {
		tok.Literal = ""
	}

	l.readChar()
	return tok
}
