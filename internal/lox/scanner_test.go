package lox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewScanner(t *testing.T) {
	source := "var x = 10;"
	lox := &Lox{}
	scanner := NewScanner(source, lox)

	assert.Equal(t, []rune(source), scanner.source)
	assert.Equal(t, lox, scanner.lox)
	assert.Equal(t, 0, scanner.start)
	assert.Equal(t, 0, scanner.current)
	assert.Equal(t, 1, scanner.line)
	assert.Empty(t, scanner.tokens)
}

func TestScanner_isAtEnd(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		current  int
		expected bool
	}{
		{
			name:     "empty source",
			source:   "",
			current:  0,
			expected: true,
		},
		{
			name:     "at beginning",
			source:   "hello",
			current:  0,
			expected: false,
		},
		{
			name:     "at middle",
			source:   "hello",
			current:  2,
			expected: false,
		},
		{
			name:     "at exact end",
			source:   "hello",
			current:  5,
			expected: true,
		},
		{
			name:     "past end",
			source:   "hello",
			current:  10,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := Scanner{
				source:  []rune(tt.source),
				current: tt.current,
			}

			result := scanner.isAtEnd()

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestScanner_scanTokens(t *testing.T) {
	tests := []struct {
		name           string
		source         string
		expectedTypes  []TokenType
		expectedCount  int
		expectHadError bool
	}{
		{
			name:           "empty source",
			source:         "",
			expectedTypes:  []TokenType{EOF},
			expectedCount:  1,
			expectHadError: false,
		},
		{
			name:           "single character tokens",
			source:         "(){},.-+;*",
			expectedTypes:  []TokenType{LeftParen, RightParen, LeftBrace, RightBrace, Comma, Dot, Minus, Plus, Semicolon, Star, EOF},
			expectedCount:  11,
			expectHadError: false,
		},
		{
			name:           "left paren only",
			source:         "(",
			expectedTypes:  []TokenType{LeftParen, EOF},
			expectedCount:  2,
			expectHadError: false,
		},
		{
			name:           "right paren only",
			source:         ")",
			expectedTypes:  []TokenType{RightParen, EOF},
			expectedCount:  2,
			expectHadError: false,
		},
		{
			name:           "braces",
			source:         "{}",
			expectedTypes:  []TokenType{LeftBrace, RightBrace, EOF},
			expectedCount:  3,
			expectHadError: false,
		},
		{
			name:           "arithmetic operators",
			source:         "+-*",
			expectedTypes:  []TokenType{Plus, Minus, Star, EOF},
			expectedCount:  4,
			expectHadError: false,
		},
		{
			name:           "punctuation",
			source:         ",;.",
			expectedTypes:  []TokenType{Comma, Semicolon, Dot, EOF},
			expectedCount:  4,
			expectHadError: false,
		},
		{
			name:           "unexpected character",
			source:         "@",
			expectedTypes:  []TokenType{EOF},
			expectedCount:  1,
			expectHadError: true,
		},
		{
			name:           "valid tokens with unexpected character",
			source:         "(+@)",
			expectedTypes:  []TokenType{LeftParen, Plus, RightParen, EOF},
			expectedCount:  4,
			expectHadError: true,
		},
		{
			name:           "multiple unexpected characters",
			source:         "@#$",
			expectedTypes:  []TokenType{EOF},
			expectedCount:  1,
			expectHadError: true,
		},
		{
			name:           "bang operator",
			source:         "!",
			expectedTypes:  []TokenType{Bang, EOF},
			expectedCount:  2,
			expectHadError: false,
		},
		{
			name:           "bang equal operator",
			source:         "!=",
			expectedTypes:  []TokenType{BangEqual, EOF},
			expectedCount:  2,
			expectHadError: false,
		},
		{
			name:           "equal operator",
			source:         "=",
			expectedTypes:  []TokenType{Equal, EOF},
			expectedCount:  2,
			expectHadError: false,
		},
		{
			name:           "equal equal operator",
			source:         "==",
			expectedTypes:  []TokenType{EqualEqual, EOF},
			expectedCount:  2,
			expectHadError: false,
		},
		{
			name:           "less than operator",
			source:         "<",
			expectedTypes:  []TokenType{Less, EOF},
			expectedCount:  2,
			expectHadError: false,
		},
		{
			name:           "less than or equal operator",
			source:         "<=",
			expectedTypes:  []TokenType{LessEqual, EOF},
			expectedCount:  2,
			expectHadError: false,
		},
		{
			name:           "greater than operator",
			source:         ">",
			expectedTypes:  []TokenType{Greater, EOF},
			expectedCount:  2,
			expectHadError: false,
		},
		{
			name:           "greater than or equal operator",
			source:         ">=",
			expectedTypes:  []TokenType{GreaterEqual, EOF},
			expectedCount:  2,
			expectHadError: false,
		},
		{
			name:           "comparison operators mixed",
			source:         "!===<=>=",
			expectedTypes:  []TokenType{BangEqual, EqualEqual, LessEqual, GreaterEqual, EOF},
			expectedCount:  5,
			expectHadError: false,
		},
		{
			name:           "comparison operators mixed with spaces",
			source:         "! == < = >=",
			expectedTypes:  []TokenType{Bang, EqualEqual, Less, Equal, GreaterEqual, EOF},
			expectedCount:  6,
			expectHadError: false,
		},
		{
			name:           "comparison with other tokens",
			source:         "(!=)",
			expectedTypes:  []TokenType{LeftParen, BangEqual, RightParen, EOF},
			expectedCount:  4,
			expectHadError: false,
		},
		{
			name:           "single and double character mix",
			source:         "!==<=>",
			expectedTypes:  []TokenType{BangEqual, Equal, LessEqual, Greater, EOF},
			expectedCount:  5,
			expectHadError: false,
		},
		{
			name:           "slash token",
			source:         "/",
			expectedTypes:  []TokenType{Slash, EOF},
			expectedCount:  2,
			expectHadError: false,
		},
		{
			name:           "comment at end of input",
			source:         "//",
			expectedTypes:  []TokenType{EOF},
			expectedCount:  1,
			expectHadError: false,
		},
		{
			name:           "comment with text",
			source:         "// this is a comment",
			expectedTypes:  []TokenType{EOF},
			expectedCount:  1,
			expectHadError: false,
		},
		{
			name:           "token followed by comment",
			source:         "+ // comment",
			expectedTypes:  []TokenType{Plus, EOF},
			expectedCount:  2,
			expectHadError: false,
		},
		{
			name:           "comment followed by newline and token",
			source:         "// comment\n+",
			expectedTypes:  []TokenType{Plus, EOF},
			expectedCount:  2,
			expectHadError: false,
		},
		{
			name:           "multiple tokens with comment",
			source:         "(+) // comment",
			expectedTypes:  []TokenType{LeftParen, Plus, RightParen, EOF},
			expectedCount:  4,
			expectHadError: false,
		},
		{
			name:           "comment between tokens",
			source:         "+\n// comment\n-",
			expectedTypes:  []TokenType{Plus, Minus, EOF},
			expectedCount:  3,
			expectHadError: false,
		},
		{
			name:           "whitespace is ignored",
			source:         " \t\r\n + \t - \n",
			expectedTypes:  []TokenType{Plus, Minus, EOF},
			expectedCount:  3,
			expectHadError: false,
		},
		{
			name:           "simple string",
			source:         "\"hello\"",
			expectedTypes:  []TokenType{String, EOF},
			expectedCount:  2,
			expectHadError: false,
		},
		{
			name:           "empty string",
			source:         "\"\"",
			expectedTypes:  []TokenType{String, EOF},
			expectedCount:  2,
			expectHadError: false,
		},
		{
			name:           "string with spaces",
			source:         "\"hello world\"",
			expectedTypes:  []TokenType{String, EOF},
			expectedCount:  2,
			expectHadError: false,
		},
		{
			name:           "multiline string",
			source:         "\"hello\nworld\"",
			expectedTypes:  []TokenType{String, EOF},
			expectedCount:  2,
			expectHadError: false,
		},
		{
			name:           "string between tokens",
			source:         "+ \"hello\" -",
			expectedTypes:  []TokenType{Plus, String, Minus, EOF},
			expectedCount:  4,
			expectHadError: false,
		},
		{
			name:           "multiple strings",
			source:         "\"foo\" \"bar\"",
			expectedTypes:  []TokenType{String, String, EOF},
			expectedCount:  3,
			expectHadError: false,
		},
		{
			name:           "unterminated string",
			source:         "\"hello",
			expectedTypes:  []TokenType{EOF},
			expectedCount:  1,
			expectHadError: true,
		},
		{
			name:           "unterminated multiline string",
			source:         "\"hello\nworld",
			expectedTypes:  []TokenType{EOF},
			expectedCount:  1,
			expectHadError: true,
		},
		{
			name:           "integer number",
			source:         "123",
			expectedTypes:  []TokenType{Number, EOF},
			expectedCount:  2,
			expectHadError: false,
		},
		{
			name:           "single digit",
			source:         "0",
			expectedTypes:  []TokenType{Number, EOF},
			expectedCount:  2,
			expectHadError: false,
		},
		{
			name:           "decimal number",
			source:         "123.456",
			expectedTypes:  []TokenType{Number, EOF},
			expectedCount:  2,
			expectHadError: false,
		},
		{
			name:           "decimal starting with zero",
			source:         "0.5",
			expectedTypes:  []TokenType{Number, EOF},
			expectedCount:  2,
			expectHadError: false,
		},
		{
			name:           "number between tokens",
			source:         "+ 42 -",
			expectedTypes:  []TokenType{Plus, Number, Minus, EOF},
			expectedCount:  4,
			expectHadError: false,
		},
		{
			name:           "multiple numbers",
			source:         "1 2.3 456",
			expectedTypes:  []TokenType{Number, Number, Number, EOF},
			expectedCount:  4,
			expectHadError: false,
		},
		{
			name:           "number followed by dot and non-digit",
			source:         "123.abc",
			expectedTypes:  []TokenType{Number, Dot, Identifier, EOF},
			expectedCount:  4,
			expectHadError: false,
		},
		{
			name:           "trailing decimal point",
			source:         "123.",
			expectedTypes:  []TokenType{Number, Dot, EOF},
			expectedCount:  3,
			expectHadError: false,
		},
		{
			name:           "simple identifier",
			source:         "foo",
			expectedTypes:  []TokenType{Identifier, EOF},
			expectedCount:  2,
			expectHadError: false,
		},
		{
			name:           "identifier with underscore",
			source:         "foo_bar",
			expectedTypes:  []TokenType{Identifier, EOF},
			expectedCount:  2,
			expectHadError: false,
		},
		{
			name:           "identifier starting with underscore",
			source:         "_foo",
			expectedTypes:  []TokenType{Identifier, EOF},
			expectedCount:  2,
			expectHadError: false,
		},
		{
			name:           "identifier with numbers",
			source:         "foo123",
			expectedTypes:  []TokenType{Identifier, EOF},
			expectedCount:  2,
			expectHadError: false,
		},
		{
			name:           "identifier between tokens",
			source:         "+ foo -",
			expectedTypes:  []TokenType{Plus, Identifier, Minus, EOF},
			expectedCount:  4,
			expectHadError: false,
		},
		{
			name:           "multiple identifiers",
			source:         "foo bar baz",
			expectedTypes:  []TokenType{Identifier, Identifier, Identifier, EOF},
			expectedCount:  4,
			expectHadError: false,
		},
		{
			name:           "keyword and",
			source:         "and",
			expectedTypes:  []TokenType{And, EOF},
			expectedCount:  2,
			expectHadError: false,
		},
		{
			name:           "keyword class",
			source:         "class",
			expectedTypes:  []TokenType{Class, EOF},
			expectedCount:  2,
			expectHadError: false,
		},
		{
			name:           "keyword else",
			source:         "else",
			expectedTypes:  []TokenType{Else, EOF},
			expectedCount:  2,
			expectHadError: false,
		},
		{
			name:           "keyword false",
			source:         "false",
			expectedTypes:  []TokenType{False, EOF},
			expectedCount:  2,
			expectHadError: false,
		},
		{
			name:           "keyword fun",
			source:         "fun",
			expectedTypes:  []TokenType{Fun, EOF},
			expectedCount:  2,
			expectHadError: false,
		},
		{
			name:           "keyword for",
			source:         "for",
			expectedTypes:  []TokenType{For, EOF},
			expectedCount:  2,
			expectHadError: false,
		},
		{
			name:           "keyword if",
			source:         "if",
			expectedTypes:  []TokenType{If, EOF},
			expectedCount:  2,
			expectHadError: false,
		},
		{
			name:           "keyword nil",
			source:         "nil",
			expectedTypes:  []TokenType{Nil, EOF},
			expectedCount:  2,
			expectHadError: false,
		},
		{
			name:           "keyword or",
			source:         "or",
			expectedTypes:  []TokenType{Or, EOF},
			expectedCount:  2,
			expectHadError: false,
		},
		{
			name:           "keyword print",
			source:         "print",
			expectedTypes:  []TokenType{Print, EOF},
			expectedCount:  2,
			expectHadError: false,
		},
		{
			name:           "keyword return",
			source:         "return",
			expectedTypes:  []TokenType{Return, EOF},
			expectedCount:  2,
			expectHadError: false,
		},
		{
			name:           "keyword super",
			source:         "super",
			expectedTypes:  []TokenType{Super, EOF},
			expectedCount:  2,
			expectHadError: false,
		},
		{
			name:           "keyword this",
			source:         "this",
			expectedTypes:  []TokenType{This, EOF},
			expectedCount:  2,
			expectHadError: false,
		},
		{
			name:           "keyword true",
			source:         "true",
			expectedTypes:  []TokenType{True, EOF},
			expectedCount:  2,
			expectHadError: false,
		},
		{
			name:           "keyword var",
			source:         "var",
			expectedTypes:  []TokenType{Var, EOF},
			expectedCount:  2,
			expectHadError: false,
		},
		{
			name:           "keyword while",
			source:         "while",
			expectedTypes:  []TokenType{While, EOF},
			expectedCount:  2,
			expectHadError: false,
		},
		{
			name:           "keywords and identifiers mixed",
			source:         "var x = 10",
			expectedTypes:  []TokenType{Var, Identifier, Equal, Number, EOF},
			expectedCount:  5,
			expectHadError: false,
		},
		{
			name:           "identifier that starts with keyword",
			source:         "variable",
			expectedTypes:  []TokenType{Identifier, EOF},
			expectedCount:  2,
			expectHadError: false,
		},
		{
			name:           "identifier that contains keyword",
			source:         "myvar",
			expectedTypes:  []TokenType{Identifier, EOF},
			expectedCount:  2,
			expectHadError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lox := &Lox{}
			scanner := NewScanner(tt.source, lox)
			tokens := scanner.ScanTokens()

			assert.Equal(t, tt.expectedCount, len(tokens))
			assert.Equal(t, tt.expectHadError, lox.hadError)

			for i, expectedType := range tt.expectedTypes {
				assert.Equal(t, expectedType, tokens[i].TokenType, "token %d type mismatch", i)
			}
		})
	}
}

func TestScanner_scanToken(t *testing.T) {
	tests := []struct {
		name           string
		source         string
		expectedType   TokenType
		expectHadError bool
		expectTokens   int
	}{
		{
			name:           "left paren",
			source:         "(",
			expectedType:   LeftParen,
			expectHadError: false,
			expectTokens:   1,
		},
		{
			name:           "right paren",
			source:         ")",
			expectedType:   RightParen,
			expectHadError: false,
			expectTokens:   1,
		},
		{
			name:           "left brace",
			source:         "{",
			expectedType:   LeftBrace,
			expectHadError: false,
			expectTokens:   1,
		},
		{
			name:           "right brace",
			source:         "}",
			expectedType:   RightBrace,
			expectHadError: false,
			expectTokens:   1,
		},
		{
			name:           "comma",
			source:         ",",
			expectedType:   Comma,
			expectHadError: false,
			expectTokens:   1,
		},
		{
			name:           "dot",
			source:         ".",
			expectedType:   Dot,
			expectHadError: false,
			expectTokens:   1,
		},
		{
			name:           "minus",
			source:         "-",
			expectedType:   Minus,
			expectHadError: false,
			expectTokens:   1,
		},
		{
			name:           "plus",
			source:         "+",
			expectedType:   Plus,
			expectHadError: false,
			expectTokens:   1,
		},
		{
			name:           "semicolon",
			source:         ";",
			expectedType:   Semicolon,
			expectHadError: false,
			expectTokens:   1,
		},
		{
			name:           "star",
			source:         "*",
			expectedType:   Star,
			expectHadError: false,
			expectTokens:   1,
		},
		{
			name:           "unexpected character",
			source:         "@",
			expectHadError: true,
			expectTokens:   0,
		},
		{
			name:           "another unexpected character",
			source:         "#",
			expectHadError: true,
			expectTokens:   0,
		},
		{
			name:           "bang",
			source:         "!",
			expectedType:   Bang,
			expectHadError: false,
			expectTokens:   1,
		},
		{
			name:           "bang equal",
			source:         "!=",
			expectedType:   BangEqual,
			expectHadError: false,
			expectTokens:   1,
		},
		{
			name:           "equal",
			source:         "=",
			expectedType:   Equal,
			expectHadError: false,
			expectTokens:   1,
		},
		{
			name:           "equal equal",
			source:         "==",
			expectedType:   EqualEqual,
			expectHadError: false,
			expectTokens:   1,
		},
		{
			name:           "less",
			source:         "<",
			expectedType:   Less,
			expectHadError: false,
			expectTokens:   1,
		},
		{
			name:           "less equal",
			source:         "<=",
			expectedType:   LessEqual,
			expectHadError: false,
			expectTokens:   1,
		},
		{
			name:           "greater",
			source:         ">",
			expectedType:   Greater,
			expectHadError: false,
			expectTokens:   1,
		},
		{
			name:           "greater equal",
			source:         ">=",
			expectedType:   GreaterEqual,
			expectHadError: false,
			expectTokens:   1,
		},
		{
			name:           "slash",
			source:         "/",
			expectedType:   Slash,
			expectHadError: false,
			expectTokens:   1,
		},
		{
			name:           "comment - no token",
			source:         "//",
			expectHadError: false,
			expectTokens:   0,
		},
		{
			name:           "simple string",
			source:         "\"hello\"",
			expectedType:   String,
			expectHadError: false,
			expectTokens:   1,
		},
		{
			name:           "empty string",
			source:         "\"\"",
			expectedType:   String,
			expectHadError: false,
			expectTokens:   1,
		},
		{
			name:           "string with spaces",
			source:         "\"hello world\"",
			expectedType:   String,
			expectHadError: false,
			expectTokens:   1,
		},
		{
			name:           "unterminated string",
			source:         "\"hello",
			expectHadError: true,
			expectTokens:   0,
		},
		{
			name:           "integer number",
			source:         "123",
			expectedType:   Number,
			expectHadError: false,
			expectTokens:   1,
		},
		{
			name:           "single digit",
			source:         "5",
			expectedType:   Number,
			expectHadError: false,
			expectTokens:   1,
		},
		{
			name:           "decimal number",
			source:         "123.456",
			expectedType:   Number,
			expectHadError: false,
			expectTokens:   1,
		},
		{
			name:           "zero",
			source:         "0",
			expectedType:   Number,
			expectHadError: false,
			expectTokens:   1,
		},
		{
			name:           "decimal with leading zero",
			source:         "0.123",
			expectedType:   Number,
			expectHadError: false,
			expectTokens:   1,
		},
		{
			name:           "simple identifier",
			source:         "foo",
			expectedType:   Identifier,
			expectHadError: false,
			expectTokens:   1,
		},
		{
			name:           "identifier with underscore",
			source:         "foo_bar",
			expectedType:   Identifier,
			expectHadError: false,
			expectTokens:   1,
		},
		{
			name:           "identifier starting with underscore",
			source:         "_variable",
			expectedType:   Identifier,
			expectHadError: false,
			expectTokens:   1,
		},
		{
			name:           "identifier with numbers",
			source:         "var123",
			expectedType:   Identifier,
			expectHadError: false,
			expectTokens:   1,
		},
		{
			name:           "keyword var",
			source:         "var",
			expectedType:   Var,
			expectHadError: false,
			expectTokens:   1,
		},
		{
			name:           "keyword if",
			source:         "if",
			expectedType:   If,
			expectHadError: false,
			expectTokens:   1,
		},
		{
			name:           "keyword while",
			source:         "while",
			expectedType:   While,
			expectHadError: false,
			expectTokens:   1,
		},
		{
			name:           "keyword true",
			source:         "true",
			expectedType:   True,
			expectHadError: false,
			expectTokens:   1,
		},
		{
			name:           "keyword false",
			source:         "false",
			expectedType:   False,
			expectHadError: false,
			expectTokens:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lox := &Lox{}
			scanner := NewScanner(tt.source, lox)
			scanner.scanToken()

			assert.Equal(t, tt.expectTokens, len(scanner.tokens))
			assert.Equal(t, tt.expectHadError, lox.hadError)

			if tt.expectTokens > 0 {
				assert.Equal(t, tt.expectedType, scanner.tokens[0].TokenType)
			}
		})
	}
}

func TestScanner_scanToken_whitespace(t *testing.T) {
	tests := []struct {
		name         string
		source       string
		expectedLine int
		expectTokens int
	}{
		{
			name:         "space",
			source:       " ",
			expectedLine: 1,
			expectTokens: 0,
		},
		{
			name:         "tab",
			source:       "\t",
			expectedLine: 1,
			expectTokens: 0,
		},
		{
			name:         "carriage return",
			source:       "\r",
			expectedLine: 1,
			expectTokens: 0,
		},
		{
			name:         "newline",
			source:       "\n",
			expectedLine: 2,
			expectTokens: 0,
		},
		{
			name:         "multiple spaces",
			source:       "   ",
			expectedLine: 1,
			expectTokens: 0,
		},
		{
			name:         "multiline string",
			source:       "\"hello\nworld\"",
			expectedLine: 2,
			expectTokens: 1,
		},
		{
			name:         "multiline string with multiple newlines",
			source:       "\"line1\nline2\nline3\"",
			expectedLine: 3,
			expectTokens: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lox := &Lox{}
			scanner := NewScanner(tt.source, lox)
			scanner.scanToken()

			assert.Equal(t, tt.expectTokens, len(scanner.tokens))
			assert.Equal(t, tt.expectedLine, scanner.line)
		})
	}
}
