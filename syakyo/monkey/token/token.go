package token

type TokenType uint

const (
	// TypeIllegal means that the token is illegal.
	TypeIllegal TokenType = iota
	// TypeEof appears if the input reached to the end.
	TypeEof

	TypeIdent  // identifier literal, e.g. x, foo.
	TypeInt    // integer literal e.g. 0, 100, -1.
	TypeString // string literal, e.g. "foo".

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
	// TypeEq represents the operator "=="
	TypeEq
	// TypeNotEq represents the operator "!="
	TypeNotEq

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
	// TypeTrue represents the keyword "true"
	TypeTrue
	// TypeFalse represents the keyword "false"
	TypeFalse
	// TypeIf represents the keyword "if"
	TypeIf
	// TypeElse represents the keyword "else"
	TypeElse
	// TypeReturn represents the keyword "return"
	TypeReturn
)

// Token represents a token of the language.
// The zero value is a illegal token.
type Token struct {
	Type    TokenType
	Literal string
}

var keywords map[string]TokenType = map[string]TokenType{
	"fn":     TypeFunction,
	"let":    TypeLet,
	"true":   TypeTrue,
	"false":  TypeFalse,
	"if":     TypeIf,
	"else":   TypeElse,
	"return": TypeReturn,
}

func LookupIdent(ident string) TokenType {
	if typ, ok := keywords[ident]; ok {
		return typ
	}
	return TypeIdent
}
