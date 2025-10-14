package lox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser_Expression(t *testing.T) {
	tests := []struct {
		name           string
		tokens         []Token
		expectedOutput string
		expectError    bool
	}{
		// Literals
		{
			name: "single number literal",
			tokens: []Token{
				NewToken(Number, "42", 42.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedOutput: "42",
			expectError:    false,
		},
		{
			name: "single string literal",
			tokens: []Token{
				NewToken(String, "\"hello\"", "hello", 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedOutput: "\"hello\"",
			expectError:    false,
		},
		{
			name: "true literal",
			tokens: []Token{
				NewToken(True, "true", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedOutput: "true",
			expectError:    false,
		},
		{
			name: "false literal",
			tokens: []Token{
				NewToken(False, "false", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedOutput: "false",
			expectError:    false,
		},
		{
			name: "nil literal",
			tokens: []Token{
				NewToken(Nil, "nil", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedOutput: "nil",
			expectError:    false,
		},

		// Grouping (1 level of nesting)
		{
			name: "simple grouping",
			tokens: []Token{
				NewToken(LeftParen, "(", nil, 1),
				NewToken(Number, "42", 42.0, 1),
				NewToken(RightParen, ")", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedOutput: "(group 42)",
			expectError:    false,
		},

		// Unary expressions (1 level)
		{
			name: "unary minus",
			tokens: []Token{
				NewToken(Minus, "-", nil, 1),
				NewToken(Number, "42", 42.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedOutput: "(- 42)",
			expectError:    false,
		},
		{
			name: "unary bang",
			tokens: []Token{
				NewToken(Bang, "!", nil, 1),
				NewToken(True, "true", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedOutput: "(! true)",
			expectError:    false,
		},

		// Binary expressions - addition/subtraction (1 level)
		{
			name: "simple addition",
			tokens: []Token{
				NewToken(Number, "1", 1.0, 1),
				NewToken(Plus, "+", nil, 1),
				NewToken(Number, "2", 2.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedOutput: "(+ 1 2)",
			expectError:    false,
		},
		{
			name: "simple subtraction",
			tokens: []Token{
				NewToken(Number, "5", 5.0, 1),
				NewToken(Minus, "-", nil, 1),
				NewToken(Number, "3", 3.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedOutput: "(- 5 3)",
			expectError:    false,
		},

		// Binary expressions - multiplication/division (1 level)
		{
			name: "simple multiplication",
			tokens: []Token{
				NewToken(Number, "3", 3.0, 1),
				NewToken(Star, "*", nil, 1),
				NewToken(Number, "4", 4.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedOutput: "(* 3 4)",
			expectError:    false,
		},
		{
			name: "simple division",
			tokens: []Token{
				NewToken(Number, "10", 10.0, 1),
				NewToken(Slash, "/", nil, 1),
				NewToken(Number, "2", 2.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedOutput: "(/ 10 2)",
			expectError:    false,
		},

		// Comparison expressions (1 level)
		{
			name: "greater than",
			tokens: []Token{
				NewToken(Number, "5", 5.0, 1),
				NewToken(Greater, ">", nil, 1),
				NewToken(Number, "3", 3.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedOutput: "(> 5 3)",
			expectError:    false,
		},
		{
			name: "greater than or equal",
			tokens: []Token{
				NewToken(Number, "5", 5.0, 1),
				NewToken(GreaterEqual, ">=", nil, 1),
				NewToken(Number, "3", 3.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedOutput: "(>= 5 3)",
			expectError:    false,
		},
		{
			name: "less than",
			tokens: []Token{
				NewToken(Number, "3", 3.0, 1),
				NewToken(Less, "<", nil, 1),
				NewToken(Number, "5", 5.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedOutput: "(< 3 5)",
			expectError:    false,
		},
		{
			name: "less than or equal",
			tokens: []Token{
				NewToken(Number, "3", 3.0, 1),
				NewToken(LessEqual, "<=", nil, 1),
				NewToken(Number, "5", 5.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedOutput: "(<= 3 5)",
			expectError:    false,
		},

		// Equality expressions (1 level)
		{
			name: "equal",
			tokens: []Token{
				NewToken(Number, "5", 5.0, 1),
				NewToken(EqualEqual, "==", nil, 1),
				NewToken(Number, "5", 5.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedOutput: "(== 5 5)",
			expectError:    false,
		},
		{
			name: "not equal",
			tokens: []Token{
				NewToken(Number, "5", 5.0, 1),
				NewToken(BangEqual, "!=", nil, 1),
				NewToken(Number, "3", 3.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedOutput: "(!= 5 3)",
			expectError:    false,
		},

		// 2 levels of nesting
		{
			name: "addition with unary",
			tokens: []Token{
				NewToken(Minus, "-", nil, 1),
				NewToken(Number, "1", 1.0, 1),
				NewToken(Plus, "+", nil, 1),
				NewToken(Number, "2", 2.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedOutput: "(+ (- 1) 2)",
			expectError:    false,
		},
		{
			name: "multiplication with grouping",
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
			expectedOutput: "(* (group (+ 1 2)) 3)",
			expectError:    false,
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
			expectedOutput: "(group (group 42))",
			expectError:    false,
		},
		{
			name: "comparison with addition",
			tokens: []Token{
				NewToken(Number, "1", 1.0, 1),
				NewToken(Plus, "+", nil, 1),
				NewToken(Number, "2", 2.0, 1),
				NewToken(Greater, ">", nil, 1),
				NewToken(Number, "3", 3.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedOutput: "(> (+ 1 2) 3)",
			expectError:    false,
		},
		{
			name: "equality with multiplication",
			tokens: []Token{
				NewToken(Number, "2", 2.0, 1),
				NewToken(Star, "*", nil, 1),
				NewToken(Number, "3", 3.0, 1),
				NewToken(EqualEqual, "==", nil, 1),
				NewToken(Number, "6", 6.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedOutput: "(== (* 2 3) 6)",
			expectError:    false,
		},
		{
			name: "chained addition",
			tokens: []Token{
				NewToken(Number, "1", 1.0, 1),
				NewToken(Plus, "+", nil, 1),
				NewToken(Number, "2", 2.0, 1),
				NewToken(Plus, "+", nil, 1),
				NewToken(Number, "3", 3.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedOutput: "(+ (+ 1 2) 3)",
			expectError:    false,
		},
		{
			name: "chained multiplication",
			tokens: []Token{
				NewToken(Number, "2", 2.0, 1),
				NewToken(Star, "*", nil, 1),
				NewToken(Number, "3", 3.0, 1),
				NewToken(Star, "*", nil, 1),
				NewToken(Number, "4", 4.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedOutput: "(* (* 2 3) 4)",
			expectError:    false,
		},

		// 3 levels of nesting
		{
			name: "complex arithmetic expression",
			tokens: []Token{
				NewToken(Number, "1", 1.0, 1),
				NewToken(Plus, "+", nil, 1),
				NewToken(Number, "2", 2.0, 1),
				NewToken(Star, "*", nil, 1),
				NewToken(Number, "3", 3.0, 1),
				NewToken(Minus, "-", nil, 1),
				NewToken(Number, "4", 4.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedOutput: "(- (+ 1 (* 2 3)) 4)",
			expectError:    false,
		},
		{
			name: "unary with grouped expression",
			tokens: []Token{
				NewToken(Minus, "-", nil, 1),
				NewToken(LeftParen, "(", nil, 1),
				NewToken(Number, "1", 1.0, 1),
				NewToken(Plus, "+", nil, 1),
				NewToken(Number, "2", 2.0, 1),
				NewToken(RightParen, ")", nil, 1),
				NewToken(Star, "*", nil, 1),
				NewToken(Number, "3", 3.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedOutput: "(* (- (group (+ 1 2))) 3)",
			expectError:    false,
		},
		{
			name: "nested comparison",
			tokens: []Token{
				NewToken(Number, "1", 1.0, 1),
				NewToken(Plus, "+", nil, 1),
				NewToken(Number, "2", 2.0, 1),
				NewToken(Greater, ">", nil, 1),
				NewToken(Number, "3", 3.0, 1),
				NewToken(Minus, "-", nil, 1),
				NewToken(Number, "1", 1.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedOutput: "(> (+ 1 2) (- 3 1))",
			expectError:    false,
		},
		{
			name: "triple nested grouping",
			tokens: []Token{
				NewToken(LeftParen, "(", nil, 1),
				NewToken(LeftParen, "(", nil, 1),
				NewToken(LeftParen, "(", nil, 1),
				NewToken(Number, "42", 42.0, 1),
				NewToken(RightParen, ")", nil, 1),
				NewToken(RightParen, ")", nil, 1),
				NewToken(RightParen, ")", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedOutput: "(group (group (group 42)))",
			expectError:    false,
		},
		{
			name: "complex equality with multiple operations",
			tokens: []Token{
				NewToken(LeftParen, "(", nil, 1),
				NewToken(Number, "1", 1.0, 1),
				NewToken(Plus, "+", nil, 1),
				NewToken(Number, "2", 2.0, 1),
				NewToken(RightParen, ")", nil, 1),
				NewToken(Star, "*", nil, 1),
				NewToken(Number, "3", 3.0, 1),
				NewToken(EqualEqual, "==", nil, 1),
				NewToken(Number, "9", 9.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedOutput: "(== (* (group (+ 1 2)) 3) 9)",
			expectError:    false,
		},
		{
			name: "chained operations with three levels",
			tokens: []Token{
				NewToken(Number, "2", 2.0, 1),
				NewToken(Star, "*", nil, 1),
				NewToken(Number, "3", 3.0, 1),
				NewToken(Plus, "+", nil, 1),
				NewToken(Number, "4", 4.0, 1),
				NewToken(Star, "*", nil, 1),
				NewToken(Number, "5", 5.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectedOutput: "(+ (* 2 3) (* 4 5))",
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lox := &Lox{}
			parser := NewParser(tt.tokens, lox)

			expr, err := parser.Parse()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if expr == nil {
					t.Fatalf("expected non-nil expression, got nil")
				}

				printer := NewAstPrinter()
				output := printer.Print(expr)
				assert.Equal(t, tt.expectedOutput, output)
			}
		})
	}
}

func TestParser_Expression_Errors(t *testing.T) {
	tests := []struct {
		name             string
		tokens           []Token
		expectError      bool
		expectNil        bool
		expectedErrorMsg string
	}{
		{
			name: "unclosed parenthesis at end",
			tokens: []Token{
				NewToken(LeftParen, "(", nil, 1),
				NewToken(Number, "42", 42.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectError:      true,
			expectedErrorMsg: "expect ')' after expression",
		},
		{
			name: "unclosed parenthesis in middle",
			tokens: []Token{
				NewToken(LeftParen, "(", nil, 1),
				NewToken(Number, "1", 1.0, 1),
				NewToken(Plus, "+", nil, 1),
				NewToken(Number, "2", 2.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectError:      true,
			expectedErrorMsg: "expect ')' after expression",
		},
		{
			name: "unclosed nested parenthesis",
			tokens: []Token{
				NewToken(LeftParen, "(", nil, 1),
				NewToken(LeftParen, "(", nil, 1),
				NewToken(Number, "42", 42.0, 1),
				NewToken(RightParen, ")", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectError:      true,
			expectedErrorMsg: "expect ')' after expression",
		},
		{
			name: "invalid character at beginning - plus operator",
			tokens: []Token{
				NewToken(Plus, "+", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectError: true,
			expectNil:   true, // Parser returns nil without error for unexpected tokens
		},
		{
			name: "invalid character at beginning - star operator",
			tokens: []Token{
				NewToken(Star, "*", nil, 1),
				NewToken(Number, "5", 5.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectError: false,
			expectNil:   true,
		},
		{
			name: "invalid character at beginning - right paren",
			tokens: []Token{
				NewToken(RightParen, ")", nil, 1),
				NewToken(Number, "42", 42.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectError: false,
			expectNil:   true,
		},
		{
			name: "invalid character in middle - consecutive operators",
			tokens: []Token{
				NewToken(Number, "1", 1.0, 1),
				NewToken(Plus, "+", nil, 1),
				NewToken(Star, "*", nil, 1),
				NewToken(Number, "2", 2.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectError: false,
			expectNil:   true,
		},
		{
			name: "binary operator without right operand",
			tokens: []Token{
				NewToken(Number, "5", 5.0, 1),
				NewToken(Plus, "+", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectError: false,
			expectNil:   true,
		},
		{
			name: "multiple unclosed parentheses",
			tokens: []Token{
				NewToken(LeftParen, "(", nil, 1),
				NewToken(LeftParen, "(", nil, 1),
				NewToken(Number, "42", 42.0, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectError:      true,
			expectedErrorMsg: "expect ')' after expression",
		},
		{
			name: "empty grouping",
			tokens: []Token{
				NewToken(LeftParen, "(", nil, 1),
				NewToken(RightParen, ")", nil, 1),
				NewToken(EOF, "", nil, 1),
			},
			expectError: false,
			expectNil:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lox := &Lox{}
			parser := NewParser(tt.tokens, lox)

			var expr Expr
			output, err := captureOutput(func() error {
				expr, _ = parser.Parse()
				return nil
			})

			if tt.expectError {
				// Either the parser returns an error or lox.hadError is set
				if err == nil && !lox.hadError {
					t.Errorf("expected an error but got none")
				}

				if tt.expectedErrorMsg != "" {
					if err != nil {
						assert.Contains(t, err.Error(), tt.expectedErrorMsg)
					} else {
						assert.Contains(t, string(output), tt.expectedErrorMsg)
					}
				}

				// For most errors, expr should be nil
				if err != nil {
					assert.Nil(t, expr)
				}
			}
		})
	}
}

func TestParser_Expression_EOFOnly(t *testing.T) {
	tests := []struct {
		name        string
		tokens      []Token
		expectNil   bool
		expectError bool
	}{
		{
			name: "only EOF token",
			tokens: []Token{
				NewToken(EOF, "", nil, 1),
			},
			expectNil:   true,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lox := &Lox{}
			parser := NewParser(tt.tokens, lox)

			expr, err := parser.Parse()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.expectNil {
				assert.Nil(t, expr, "expected nil expression for EOF-only input")
			}

			// Verify that hadError is false for EOF-only case
			assert.False(t, lox.hadError, "EOF-only input should not be an error")
		})
	}
}
