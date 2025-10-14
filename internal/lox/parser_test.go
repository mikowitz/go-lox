package lox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser_Literal(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []Token
		expected Expr
	}{
		{
			name: "number literal",
			tokens: []Token{
				NewToken(Number, "123", 123.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expected: Literal{Value: NewToken(Number, "123", 123.0, 1)},
		},
		{
			name: "string literal",
			tokens: []Token{
				NewToken(String, "hello", "hello", 1),
				NewToken(EOF, "", nil, 1),
			},
			expected: Literal{Value: NewToken(String, "hello", "hello", 1)},
		},
		{
			name: "true literal",
			tokens: []Token{
				NewToken(True, "true", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
			expected: Literal{Value: NewToken(True, "true", nil, 1)},
		},
		{
			name: "false literal",
			tokens: []Token{
				NewToken(False, "false", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
			expected: Literal{Value: NewToken(False, "false", nil, 1)},
		},
		{
			name: "nil literal",
			tokens: []Token{
				NewToken(Nil, "nil", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
			expected: Literal{Value: NewToken(Nil, "nil", nil, 1)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lox := &Lox{}
			parser := NewParser(tt.tokens, lox)
			result, err := parser.Parse()
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParser_Grouping(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []Token
		expected Expr
	}{
		{
			name: "grouped number",
			tokens: []Token{
				NewToken(LeftParen, "(", nil, 1),
				NewToken(Number, "123", 123.0, 1),
				NewToken(RightParen, ")", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
			expected: Grouping{
				Expr: Literal{Value: NewToken(Number, "123", 123.0, 1)},
			},
		},
		{
			name: "nested grouping",
			tokens: []Token{
				NewToken(LeftParen, "(", nil, 1),
				NewToken(LeftParen, "(", nil, 1),
				NewToken(Number, "42", 42.0, 1),
				NewToken(RightParen, ")", nil, 1),
				NewToken(RightParen, ")", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
			expected: Grouping{
				Expr: Grouping{
					Expr: Literal{Value: NewToken(Number, "42", 42.0, 1)},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lox := &Lox{}
			parser := NewParser(tt.tokens, lox)
			result, err := parser.Parse()
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParser_Unary(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []Token
		expected Expr
	}{
		{
			name: "negation",
			tokens: []Token{
				NewToken(Minus, "-", nil, 1),
				NewToken(Number, "123", 123.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expected: Unary{
				Operator: NewToken(Minus, "-", nil, 1),
				Right:    Literal{Value: NewToken(Number, "123", 123.0, 1)},
			},
		},
		{
			name: "logical not",
			tokens: []Token{
				NewToken(Bang, "!", nil, 1),
				NewToken(True, "true", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
			expected: Unary{
				Operator: NewToken(Bang, "!", nil, 1),
				Right:    Literal{Value: NewToken(True, "true", nil, 1)},
			},
		},
		{
			name: "double negation",
			tokens: []Token{
				NewToken(Minus, "-", nil, 1),
				NewToken(Minus, "-", nil, 1),
				NewToken(Number, "5", 5.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expected: Unary{
				Operator: NewToken(Minus, "-", nil, 1),
				Right: Unary{
					Operator: NewToken(Minus, "-", nil, 1),
					Right:    Literal{Value: NewToken(Number, "5", 5.0, 1)},
				},
			},
		},
		{
			name: "negation of grouped expression",
			tokens: []Token{
				NewToken(Minus, "-", nil, 1),
				NewToken(LeftParen, "(", nil, 1),
				NewToken(Number, "10", 10.0, 1),
				NewToken(RightParen, ")", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
			expected: Unary{
				Operator: NewToken(Minus, "-", nil, 1),
				Right: Grouping{
					Expr: Literal{Value: NewToken(Number, "10", 10.0, 1)},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lox := &Lox{}
			parser := NewParser(tt.tokens, lox)
			result, err := parser.Parse()
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParser_Binary(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []Token
		expected Expr
	}{
		{
			name: "addition",
			tokens: []Token{
				NewToken(Number, "1", 1.0, 1),
				NewToken(Plus, "+", nil, 1),
				NewToken(Number, "2", 2.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expected: Binary{
				Left:     Literal{Value: NewToken(Number, "1", 1.0, 1)},
				Operator: NewToken(Plus, "+", nil, 1),
				Right:    Literal{Value: NewToken(Number, "2", 2.0, 1)},
			},
		},
		{
			name: "subtraction",
			tokens: []Token{
				NewToken(Number, "10", 10.0, 1),
				NewToken(Minus, "-", nil, 1),
				NewToken(Number, "5", 5.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expected: Binary{
				Left:     Literal{Value: NewToken(Number, "10", 10.0, 1)},
				Operator: NewToken(Minus, "-", nil, 1),
				Right:    Literal{Value: NewToken(Number, "5", 5.0, 1)},
			},
		},
		{
			name: "multiplication",
			tokens: []Token{
				NewToken(Number, "3", 3.0, 1),
				NewToken(Star, "*", nil, 1),
				NewToken(Number, "4", 4.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expected: Binary{
				Left:     Literal{Value: NewToken(Number, "3", 3.0, 1)},
				Operator: NewToken(Star, "*", nil, 1),
				Right:    Literal{Value: NewToken(Number, "4", 4.0, 1)},
			},
		},
		{
			name: "division",
			tokens: []Token{
				NewToken(Number, "8", 8.0, 1),
				NewToken(Slash, "/", nil, 1),
				NewToken(Number, "2", 2.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expected: Binary{
				Left:     Literal{Value: NewToken(Number, "8", 8.0, 1)},
				Operator: NewToken(Slash, "/", nil, 1),
				Right:    Literal{Value: NewToken(Number, "2", 2.0, 1)},
			},
		},
		{
			name: "equality",
			tokens: []Token{
				NewToken(Number, "5", 5.0, 1),
				NewToken(EqualEqual, "==", nil, 1),
				NewToken(Number, "5", 5.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expected: Binary{
				Left:     Literal{Value: NewToken(Number, "5", 5.0, 1)},
				Operator: NewToken(EqualEqual, "==", nil, 1),
				Right:    Literal{Value: NewToken(Number, "5", 5.0, 1)},
			},
		},
		{
			name: "inequality",
			tokens: []Token{
				NewToken(Number, "5", 5.0, 1),
				NewToken(BangEqual, "!=", nil, 1),
				NewToken(Number, "3", 3.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expected: Binary{
				Left:     Literal{Value: NewToken(Number, "5", 5.0, 1)},
				Operator: NewToken(BangEqual, "!=", nil, 1),
				Right:    Literal{Value: NewToken(Number, "3", 3.0, 1)},
			},
		},
		{
			name: "less than",
			tokens: []Token{
				NewToken(Number, "3", 3.0, 1),
				NewToken(Less, "<", nil, 1),
				NewToken(Number, "5", 5.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expected: Binary{
				Left:     Literal{Value: NewToken(Number, "3", 3.0, 1)},
				Operator: NewToken(Less, "<", nil, 1),
				Right:    Literal{Value: NewToken(Number, "5", 5.0, 1)},
			},
		},
		{
			name: "greater than or equal",
			tokens: []Token{
				NewToken(Number, "7", 7.0, 1),
				NewToken(GreaterEqual, ">=", nil, 1),
				NewToken(Number, "5", 5.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expected: Binary{
				Left:     Literal{Value: NewToken(Number, "7", 7.0, 1)},
				Operator: NewToken(GreaterEqual, ">=", nil, 1),
				Right:    Literal{Value: NewToken(Number, "5", 5.0, 1)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lox := &Lox{}
			parser := NewParser(tt.tokens, lox)
			result, err := parser.Parse()
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParser_Precedence(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []Token
		expected Expr
	}{
		{
			name: "multiplication before addition",
			tokens: []Token{
				NewToken(Number, "1", 1.0, 1),
				NewToken(Plus, "+", nil, 1),
				NewToken(Number, "2", 2.0, 1),
				NewToken(Star, "*", nil, 1),
				NewToken(Number, "3", 3.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expected: Binary{
				Left:     Literal{Value: NewToken(Number, "1", 1.0, 1)},
				Operator: NewToken(Plus, "+", nil, 1),
				Right: Binary{
					Left:     Literal{Value: NewToken(Number, "2", 2.0, 1)},
					Operator: NewToken(Star, "*", nil, 1),
					Right:    Literal{Value: NewToken(Number, "3", 3.0, 1)},
				},
			},
		},
		{
			name: "unary before multiplication",
			tokens: []Token{
				NewToken(Minus, "-", nil, 1),
				NewToken(Number, "2", 2.0, 1),
				NewToken(Star, "*", nil, 1),
				NewToken(Number, "3", 3.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expected: Binary{
				Left: Unary{
					Operator: NewToken(Minus, "-", nil, 1),
					Right:    Literal{Value: NewToken(Number, "2", 2.0, 1)},
				},
				Operator: NewToken(Star, "*", nil, 1),
				Right:    Literal{Value: NewToken(Number, "3", 3.0, 1)},
			},
		},
		{
			name: "grouping overrides precedence",
			tokens: []Token{
				NewToken(LeftParen, "(", nil, 1),
				NewToken(Number, "1", 1.0, 1),
				NewToken(Plus, "+", nil, 1),
				NewToken(Number, "2", 2.0, 1),
				NewToken(RightParen, ")", nil, 1),
				NewToken(Star, "*", nil, 1),
				NewToken(Number, "3", 3.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expected: Binary{
				Left: Grouping{
					Expr: Binary{
						Left:     Literal{Value: NewToken(Number, "1", 1.0, 1)},
						Operator: NewToken(Plus, "+", nil, 1),
						Right:    Literal{Value: NewToken(Number, "2", 2.0, 1)},
					},
				},
				Operator: NewToken(Star, "*", nil, 1),
				Right:    Literal{Value: NewToken(Number, "3", 3.0, 1)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lox := &Lox{}
			parser := NewParser(tt.tokens, lox)
			result, err := parser.Parse()
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParser_Complex(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []Token
		expected Expr
	}{
		{
			name: "complex expression with all types",
			tokens: []Token{
				// -2 * (3 + 4)
				NewToken(Minus, "-", nil, 1),
				NewToken(Number, "2", 2.0, 1),
				NewToken(Star, "*", nil, 1),
				NewToken(LeftParen, "(", nil, 1),
				NewToken(Number, "3", 3.0, 1),
				NewToken(Plus, "+", nil, 1),
				NewToken(Number, "4", 4.0, 1),
				NewToken(RightParen, ")", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
			expected: Binary{
				Left: Unary{
					Operator: NewToken(Minus, "-", nil, 1),
					Right:    Literal{Value: NewToken(Number, "2", 2.0, 1)},
				},
				Operator: NewToken(Star, "*", nil, 1),
				Right: Grouping{
					Expr: Binary{
						Left:     Literal{Value: NewToken(Number, "3", 3.0, 1)},
						Operator: NewToken(Plus, "+", nil, 1),
						Right:    Literal{Value: NewToken(Number, "4", 4.0, 1)},
					},
				},
			},
		},
		{
			name: "comparison with arithmetic",
			tokens: []Token{
				// 5 + 3 > 2 * 4
				NewToken(Number, "5", 5.0, 1),
				NewToken(Plus, "+", nil, 1),
				NewToken(Number, "3", 3.0, 1),
				NewToken(Greater, ">", nil, 1),
				NewToken(Number, "2", 2.0, 1),
				NewToken(Star, "*", nil, 1),
				NewToken(Number, "4", 4.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expected: Binary{
				Left: Binary{
					Left:     Literal{Value: NewToken(Number, "5", 5.0, 1)},
					Operator: NewToken(Plus, "+", nil, 1),
					Right:    Literal{Value: NewToken(Number, "3", 3.0, 1)},
				},
				Operator: NewToken(Greater, ">", nil, 1),
				Right: Binary{
					Left:     Literal{Value: NewToken(Number, "2", 2.0, 1)},
					Operator: NewToken(Star, "*", nil, 1),
					Right:    Literal{Value: NewToken(Number, "4", 4.0, 1)},
				},
			},
		},
		{
			name: "equality with comparison",
			tokens: []Token{
				// 5 < 10 == true
				NewToken(Number, "5", 5.0, 1),
				NewToken(Less, "<", nil, 1),
				NewToken(Number, "10", 10.0, 1),
				NewToken(EqualEqual, "==", nil, 1),
				NewToken(True, "true", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
			expected: Binary{
				Left: Binary{
					Left:     Literal{Value: NewToken(Number, "5", 5.0, 1)},
					Operator: NewToken(Less, "<", nil, 1),
					Right:    Literal{Value: NewToken(Number, "10", 10.0, 1)},
				},
				Operator: NewToken(EqualEqual, "==", nil, 1),
				Right:    Literal{Value: NewToken(True, "true", nil, 1)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lox := &Lox{}
			parser := NewParser(tt.tokens, lox)
			result, err := parser.Parse()
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParser_Errors(t *testing.T) {
	tests := []struct {
		name           string
		tokens         []Token
		expectHadError bool
	}{
		{
			name: "empty input (only EOF)",
			tokens: []Token{
				NewToken(EOF, "", nil, 1),
			},
			expectHadError: false,
		},
		{
			name: "unclosed left parenthesis",
			tokens: []Token{
				NewToken(LeftParen, "(", nil, 1),
				NewToken(Number, "5", 5.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectHadError: true,
		},
		{
			name: "unclosed nested parenthesis",
			tokens: []Token{
				NewToken(LeftParen, "(", nil, 1),
				NewToken(Number, "1", 1.0, 1),
				NewToken(Plus, "+", nil, 1),
				NewToken(LeftParen, "(", nil, 1),
				NewToken(Number, "2", 2.0, 1),
				NewToken(RightParen, ")", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectHadError: true,
		},
		{
			name: "unexpected token - binary operator at start",
			tokens: []Token{
				NewToken(Plus, "+", nil, 1),
				NewToken(Number, "5", 5.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectHadError: true,
		},
		{
			name: "unexpected token - star at start",
			tokens: []Token{
				NewToken(Star, "*", nil, 1),
				NewToken(Number, "10", 10.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectHadError: true,
		},
		{
			name: "unexpected token - comparison operator at start",
			tokens: []Token{
				NewToken(Greater, ">", nil, 1),
				NewToken(Number, "3", 3.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectHadError: true,
		},
		{
			name: "unexpected token - equality operator at start",
			tokens: []Token{
				NewToken(EqualEqual, "==", nil, 1),
				NewToken(Number, "7", 7.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectHadError: true,
		},
		{
			name: "unexpected token - right paren without left",
			tokens: []Token{
				NewToken(Number, "5", 5.0, 1),
				NewToken(RightParen, ")", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectHadError: true,
		},
		{
			name: "unexpected token - semicolon",
			tokens: []Token{
				NewToken(Semicolon, ";", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectHadError: true,
		},
		{
			name: "empty parentheses",
			tokens: []Token{
				NewToken(LeftParen, "(", nil, 1),
				NewToken(RightParen, ")", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectHadError: true,
		},
		{
			name: "multiple operators in a row",
			tokens: []Token{
				NewToken(Number, "5", 5.0, 1),
				NewToken(Plus, "+", nil, 1),
				NewToken(Star, "*", nil, 1),
				NewToken(Number, "3", 3.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectHadError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lox := &Lox{}
			parser := NewParser(tt.tokens, lox)
			result, _ := parser.Parse()

			assert.Equal(t, tt.expectHadError, lox.hadError, "hadError flag mismatch")
			assert.Nil(t, result, "result mismatch")
		})
	}
}
