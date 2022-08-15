package token

type TokenType uint

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
	// TypeMinus represents the operator "-"
	TypeMinus
	// TypeBang represents the operator "!"
	TypeBang
	// TypeAsterisk represents the operator "*"
	TypeAsterisk
	// TypeSlash represents the operator "/"
	TypeSlash
	// TypeLt represents the operator "<"
	TypeLt
	// TypeGt represents the operator ">"
	TypeGt

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

// Token represents a token of the language.
// The zero value is a illegal token.
type Token struct {
	Type    TokenType
	Literal string
}

var keywords map[string]TokenType = map[string]TokenType{
	"fn":  TypeFunction,
	"let": TypeLet,
}

func LookupIdent(ident string) TokenType {
	if typ, ok := keywords[ident]; ok {
		return typ
	}
	return TypeIdent
}
