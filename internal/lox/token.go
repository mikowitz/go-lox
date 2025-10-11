//go:generate stringer -type TokenType
package lox

// Token represents a lexical token in the Lox language.
type Token struct {
	Lexeme    string    // The raw text of the token
	Literal   any       // The literal value (for numbers, strings, etc.)
	TokenType TokenType // The type of token
	Line      int       // The line number where the token appears
}

// NewToken creates a new Token with the given properties.
func NewToken(tokenType TokenType, lexeme string, literal any, line int) Token {
	return Token{
		TokenType: tokenType,
		Lexeme:    lexeme,
		Literal:   literal,
		Line:      line,
	}
}

// TokenType represents the type of a lexical token.
type TokenType int

const (
	// Single character tokens
	LeftParen TokenType = iota
	RightParen
	LeftBrace
	RightBrace
	Comma
	Dot
	Minus
	Plus
	Semicolon
	Slash
	Star
	// One or two character tokens
	Bang
	BangEqual
	Equal
	EqualEqual
	Greater
	GreaterEqual
	Less
	LessEqual
	// Literals
	Identifier
	String
	Number
	// Keywords
	And
	Class
	Else
	False
	Fun
	For
	If
	Nil
	Or
	Print
	Return
	Super
	This
	True
	Var
	While
	// EOF
	EOF
)

var keywords = map[string]TokenType{
	"and":    And,
	"class":  Class,
	"else":   Else,
	"false":  False,
	"fun":    Fun,
	"for":    For,
	"if":     If,
	"nil":    Nil,
	"or":     Or,
	"print":  Print,
	"return": Return,
	"super":  Super,
	"this":   This,
	"true":   True,
	"var":    Var,
	"while":  While,
}
