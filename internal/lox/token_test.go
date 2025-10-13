package lox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenType_String(t *testing.T) {
	tests := []struct {
		tokenType    TokenType
		expectedName string
	}{
		// Single character tokens
		{LeftParen, "LeftParen"},
		{RightParen, "RightParen"},
		{LeftBrace, "LeftBrace"},
		{RightBrace, "RightBrace"},
		{Comma, "Comma"},
		{Dot, "Dot"},
		{Minus, "Minus"},
		{Plus, "Plus"},
		{Semicolon, "Semicolon"},
		{Slash, "Slash"},
		{Star, "Star"},
		// One or two character tokens
		{Bang, "Bang"},
		{BangEqual, "BangEqual"},
		{Equal, "Equal"},
		{EqualEqual, "EqualEqual"},
		{Greater, "Greater"},
		{GreaterEqual, "GreaterEqual"},
		{Less, "Less"},
		{LessEqual, "LessEqual"},
		// Literals
		{Identifier, "Identifier"},
		{String, "String"},
		{Number, "Number"},
		// Keywords
		{And, "And"},
		{Class, "Class"},
		{Else, "Else"},
		{False, "False"},
		{Fun, "Fun"},
		{For, "For"},
		{If, "If"},
		{Nil, "Nil"},
		{Or, "Or"},
		{Print, "Print"},
		{Return, "Return"},
		{Super, "Super"},
		{This, "This"},
		{True, "True"},
		{Var, "Var"},
		{While, "While"},
		// EOF
		{EOF, "EOF"},
	}

	for _, tt := range tests {
		t.Run(tt.expectedName, func(t *testing.T) {
			result := tt.tokenType.String()
			assert.Equal(t, tt.expectedName, result)
		})
	}
}

func TestNewToken(t *testing.T) {
	tests := []struct {
		name      string
		tokenType TokenType
		lexeme    string
		literal   any
		line      int
	}{
		{
			name:      "simple token without literal",
			tokenType: LeftParen,
			lexeme:    "(",
			literal:   nil,
			line:      1,
		},
		{
			name:      "string token with literal",
			tokenType: String,
			lexeme:    "\"hello\"",
			literal:   "hello",
			line:      1,
		},
		{
			name:      "number token with literal",
			tokenType: Number,
			lexeme:    "123.456",
			literal:   123.456,
			line:      5,
		},
		{
			name:      "identifier token",
			tokenType: Identifier,
			lexeme:    "myVar",
			literal:   nil,
			line:      10,
		},
		{
			name:      "keyword token",
			tokenType: Var,
			lexeme:    "var",
			literal:   nil,
			line:      3,
		},
		{
			name:      "EOF token",
			tokenType: EOF,
			lexeme:    "",
			literal:   nil,
			line:      100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := NewToken(tt.tokenType, tt.lexeme, tt.literal, tt.line)

			assert.Equal(t, tt.tokenType, token.TokenType)
			assert.Equal(t, tt.lexeme, token.Lexeme)
			assert.Equal(t, tt.literal, token.Literal)
			assert.Equal(t, tt.line, token.Line)
		})
	}
}
