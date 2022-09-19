package token

type TokenType uint

const (
	TypeIllegal TokenType = iota // token is illegal.
	TypeEof                      // the input reached to the end.

	TypeIdent  // identifier literal, e.g. x, foo.
	TypeInt    // integer literal e.g. 0, 100, -1.
	TypeString // string literal, e.g. "foo".

	TypeAssign   // =
	TypePlus     // +
	TypeMinus    // -
	TypeBang     // !
	TypeAsterisk // *
	TypeSlash    // /
	TypeLt       // <
	TypeGt       // >
	TypeEq       // ==
	TypeNotEq    // !=

	TypeComma       // ,
	TypeSemicolon   // ;
	TypeLeftParen   // (
	TypeRightParen  // )
	TypeLeftBrace   // {
	TypeRightBrace  // }
	TypeLeftBraket  // [
	TypeRightBraket // ]

	TypeFunction // keyword "funcion"
	TypeLet      // keyword "let"
	TypeTrue     // keyword "true"
	TypeFalse    // keyword "false"
	TypeIf       // keyword "if"
	TypeElse     // keyword "else"
	TypeReturn   // keywork "return"
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
