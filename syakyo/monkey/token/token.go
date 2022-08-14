package token

type TokenType uint

// Token represents a token of the language.
// The zero value is a illegal token.
type Token struct {
	Type    TokenType
	Literal string
}

const (
	// TypeIllegal means that the token is illegal.
	TypeIllegal TokenType = iota
	// TypeEof appears if the input reached to the end.
	TypeEof

	// TypeIdent represents the identifier literal, e.g. x, foo.
	TypeIdent
	// TypeInt represents the integer literal, e.g. 0, 100, -1.
	TypeInt

	// TypePlus represents the operator "="
	TypeAssign
	// TypePlus represents the operator "+"
	TypePlus

	// TypeComma represents the delimiter ","
	TypeComma
	// TypeSemicolon represents the delimiter ";"
	TypeSemicolon
	// TypeLeftParen represents the token `(`
	TypeLeftParen
	// TypeRightParen represents the token `)`
	TypeRightParen
	// TypeLeftBrace represents the token `{`
	TypeLeftBrace
	// TypeRightBrace represents the token `}`
	TypeRightBrace

	// TypeFuncion represents the keyword "function"
	TypeFunction
	// TypeLeft represents the keyword "let"
	TypeLet
)
