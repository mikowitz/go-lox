package lox

import (
	"testing"
)

func TestInterpreter_Literals(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expr
		expected any
		wantErr  bool
	}{
		{
			name:     "number literal",
			expr:     Literal{Value: Token{TokenType: Number, Literal: 42.0}},
			expected: 42.0,
		},
		{
			name:     "string literal",
			expr:     Literal{Value: Token{TokenType: String, Literal: "hello"}},
			expected: "hello",
		},
		{
			name:     "true literal",
			expr:     Literal{Value: Token{TokenType: True}},
			expected: true,
		},
		{
			name:     "false literal",
			expr:     Literal{Value: Token{TokenType: False}},
			expected: false,
		},
		{
			name:     "nil literal",
			expr:     Literal{Value: Token{TokenType: Nil}},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lox := &Lox{}
			i := NewInterpreter(lox)
			result, err := i.Interpret(tt.expr)

			if (err != nil) != tt.wantErr {
				t.Errorf("Interpret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.expected {
				t.Errorf("Interpret() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestInterpreter_UnaryExpressions(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expr
		expected any
		wantErr  bool
	}{
		{
			name: "negation of number",
			expr: Unary{
				Operator: Token{TokenType: Minus, Lexeme: "-"},
				Right:    Literal{Value: Token{TokenType: Number, Literal: 5.0}},
			},
			expected: -5.0,
		},
		{
			name: "negation of negative number",
			expr: Unary{
				Operator: Token{TokenType: Minus, Lexeme: "-"},
				Right:    Literal{Value: Token{TokenType: Number, Literal: -3.0}},
			},
			expected: 3.0,
		},
		{
			name: "logical not of true",
			expr: Unary{
				Operator: Token{TokenType: Bang, Lexeme: "!"},
				Right:    Literal{Value: Token{TokenType: True}},
			},
			expected: false,
		},
		{
			name: "logical not of false",
			expr: Unary{
				Operator: Token{TokenType: Bang, Lexeme: "!"},
				Right:    Literal{Value: Token{TokenType: False}},
			},
			expected: true,
		},
		{
			name: "logical not of nil",
			expr: Unary{
				Operator: Token{TokenType: Bang, Lexeme: "!"},
				Right:    Literal{Value: Token{TokenType: Nil}},
			},
			expected: true,
		},
		{
			name: "logical not of number (truthy)",
			expr: Unary{
				Operator: Token{TokenType: Bang, Lexeme: "!"},
				Right:    Literal{Value: Token{TokenType: Number, Literal: 0.0}},
			},
			expected: false,
		},
		{
			name: "double negation",
			expr: Unary{
				Operator: Token{TokenType: Bang, Lexeme: "!"},
				Right: Unary{
					Operator: Token{TokenType: Bang, Lexeme: "!"},
					Right:    Literal{Value: Token{TokenType: True}},
				},
			},
			expected: true,
		},
		{
			name: "negation of non-number should error",
			expr: Unary{
				Operator: Token{TokenType: Minus, Lexeme: "-"},
				Right:    Literal{Value: Token{TokenType: String, Literal: "hello"}},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lox := &Lox{}
			i := NewInterpreter(lox)
			result, err := i.Interpret(tt.expr)

			if (err != nil) != tt.wantErr {
				t.Errorf("Interpret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("Interpret() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestInterpreter_ArithmeticExpressions(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expr
		expected any
		wantErr  bool
	}{
		{
			name: "addition",
			expr: Binary{
				Left:     Literal{Value: Token{TokenType: Number, Literal: 5.0}},
				Operator: Token{TokenType: Plus, Lexeme: "+"},
				Right:    Literal{Value: Token{TokenType: Number, Literal: 3.0}},
			},
			expected: 8.0,
		},
		{
			name: "subtraction",
			expr: Binary{
				Left:     Literal{Value: Token{TokenType: Number, Literal: 10.0}},
				Operator: Token{TokenType: Minus, Lexeme: "-"},
				Right:    Literal{Value: Token{TokenType: Number, Literal: 3.0}},
			},
			expected: 7.0,
		},
		{
			name: "multiplication",
			expr: Binary{
				Left:     Literal{Value: Token{TokenType: Number, Literal: 4.0}},
				Operator: Token{TokenType: Star, Lexeme: "*"},
				Right:    Literal{Value: Token{TokenType: Number, Literal: 5.0}},
			},
			expected: 20.0,
		},
		{
			name: "division",
			expr: Binary{
				Left:     Literal{Value: Token{TokenType: Number, Literal: 15.0}},
				Operator: Token{TokenType: Slash, Lexeme: "/"},
				Right:    Literal{Value: Token{TokenType: Number, Literal: 3.0}},
			},
			expected: 5.0,
		},
		{
			name: "string concatenation",
			expr: Binary{
				Left:     Literal{Value: Token{TokenType: String, Literal: "hello"}},
				Operator: Token{TokenType: Plus, Lexeme: "+"},
				Right:    Literal{Value: Token{TokenType: String, Literal: " world"}},
			},
			expected: "hello world",
		},
		{
			name: "complex arithmetic: (5 + 3) * 2",
			expr: Binary{
				Left: Grouping{
					Expr: Binary{
						Left:     Literal{Value: Token{TokenType: Number, Literal: 5.0}},
						Operator: Token{TokenType: Plus, Lexeme: "+"},
						Right:    Literal{Value: Token{TokenType: Number, Literal: 3.0}},
					},
				},
				Operator: Token{TokenType: Star, Lexeme: "*"},
				Right:    Literal{Value: Token{TokenType: Number, Literal: 2.0}},
			},
			expected: 16.0,
		},
		{
			name: "mixed types in addition should error",
			expr: Binary{
				Left:     Literal{Value: Token{TokenType: Number, Literal: 5.0}},
				Operator: Token{TokenType: Plus, Lexeme: "+"},
				Right:    Literal{Value: Token{TokenType: String, Literal: "hello"}},
			},
			wantErr: true,
		},
		{
			name: "non-numbers in subtraction should error",
			expr: Binary{
				Left:     Literal{Value: Token{TokenType: String, Literal: "hello"}},
				Operator: Token{TokenType: Minus, Lexeme: "-"},
				Right:    Literal{Value: Token{TokenType: String, Literal: "world"}},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lox := &Lox{}
			i := NewInterpreter(lox)
			result, err := i.Interpret(tt.expr)

			if (err != nil) != tt.wantErr {
				t.Errorf("Interpret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("Interpret() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestInterpreter_ComparisonExpressions(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expr
		expected any
		wantErr  bool
	}{
		{
			name: "less than: true",
			expr: Binary{
				Left:     Literal{Value: Token{TokenType: Number, Literal: 3.0}},
				Operator: Token{TokenType: Less, Lexeme: "<"},
				Right:    Literal{Value: Token{TokenType: Number, Literal: 5.0}},
			},
			expected: true,
		},
		{
			name: "less than: false",
			expr: Binary{
				Left:     Literal{Value: Token{TokenType: Number, Literal: 5.0}},
				Operator: Token{TokenType: Less, Lexeme: "<"},
				Right:    Literal{Value: Token{TokenType: Number, Literal: 3.0}},
			},
			expected: false,
		},
		{
			name: "less than or equal: true (less)",
			expr: Binary{
				Left:     Literal{Value: Token{TokenType: Number, Literal: 3.0}},
				Operator: Token{TokenType: LessEqual, Lexeme: "<="},
				Right:    Literal{Value: Token{TokenType: Number, Literal: 5.0}},
			},
			expected: true,
		},
		{
			name: "less than or equal: true (equal)",
			expr: Binary{
				Left:     Literal{Value: Token{TokenType: Number, Literal: 5.0}},
				Operator: Token{TokenType: LessEqual, Lexeme: "<="},
				Right:    Literal{Value: Token{TokenType: Number, Literal: 5.0}},
			},
			expected: true,
		},
		{
			name: "greater than: true",
			expr: Binary{
				Left:     Literal{Value: Token{TokenType: Number, Literal: 5.0}},
				Operator: Token{TokenType: Greater, Lexeme: ">"},
				Right:    Literal{Value: Token{TokenType: Number, Literal: 3.0}},
			},
			expected: true,
		},
		{
			name: "greater than or equal: true (equal)",
			expr: Binary{
				Left:     Literal{Value: Token{TokenType: Number, Literal: 5.0}},
				Operator: Token{TokenType: GreaterEqual, Lexeme: ">="},
				Right:    Literal{Value: Token{TokenType: Number, Literal: 5.0}},
			},
			expected: true,
		},
		{
			name: "equality: numbers equal",
			expr: Binary{
				Left:     Literal{Value: Token{TokenType: Number, Literal: 5.0}},
				Operator: Token{TokenType: EqualEqual, Lexeme: "=="},
				Right:    Literal{Value: Token{TokenType: Number, Literal: 5.0}},
			},
			expected: true,
		},
		{
			name: "equality: numbers not equal",
			expr: Binary{
				Left:     Literal{Value: Token{TokenType: Number, Literal: 5.0}},
				Operator: Token{TokenType: EqualEqual, Lexeme: "=="},
				Right:    Literal{Value: Token{TokenType: Number, Literal: 3.0}},
			},
			expected: false,
		},
		{
			name: "equality: strings equal",
			expr: Binary{
				Left:     Literal{Value: Token{TokenType: String, Literal: "hello"}},
				Operator: Token{TokenType: EqualEqual, Lexeme: "=="},
				Right:    Literal{Value: Token{TokenType: String, Literal: "hello"}},
			},
			expected: true,
		},
		{
			name: "equality: nil == nil",
			expr: Binary{
				Left:     Literal{Value: Token{TokenType: Nil}},
				Operator: Token{TokenType: EqualEqual, Lexeme: "=="},
				Right:    Literal{Value: Token{TokenType: Nil}},
			},
			expected: true,
		},
		{
			name: "inequality: numbers",
			expr: Binary{
				Left:     Literal{Value: Token{TokenType: Number, Literal: 5.0}},
				Operator: Token{TokenType: BangEqual, Lexeme: "!="},
				Right:    Literal{Value: Token{TokenType: Number, Literal: 3.0}},
			},
			expected: true,
		},
		{
			name: "inequality: same numbers",
			expr: Binary{
				Left:     Literal{Value: Token{TokenType: Number, Literal: 5.0}},
				Operator: Token{TokenType: BangEqual, Lexeme: "!="},
				Right:    Literal{Value: Token{TokenType: Number, Literal: 5.0}},
			},
			expected: false,
		},
		{
			name: "comparison of non-numbers should error",
			expr: Binary{
				Left:     Literal{Value: Token{TokenType: String, Literal: "hello"}},
				Operator: Token{TokenType: Less, Lexeme: "<"},
				Right:    Literal{Value: Token{TokenType: String, Literal: "world"}},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lox := &Lox{}
			i := NewInterpreter(lox)
			result, err := i.Interpret(tt.expr)

			if (err != nil) != tt.wantErr {
				t.Errorf("Interpret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("Interpret() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestInterpreter_ComplexNestedExpressions(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expr
		expected any
		wantErr  bool
	}{
		{
			name: "nested arithmetic: (10 + 5) / (3 - 1)",
			expr: Binary{
				Left: Grouping{
					Expr: Binary{
						Left:     Literal{Value: Token{TokenType: Number, Literal: 10.0}},
						Operator: Token{TokenType: Plus, Lexeme: "+"},
						Right:    Literal{Value: Token{TokenType: Number, Literal: 5.0}},
					},
				},
				Operator: Token{TokenType: Slash, Lexeme: "/"},
				Right: Grouping{
					Expr: Binary{
						Left:     Literal{Value: Token{TokenType: Number, Literal: 3.0}},
						Operator: Token{TokenType: Minus, Lexeme: "-"},
						Right:    Literal{Value: Token{TokenType: Number, Literal: 1.0}},
					},
				},
			},
			expected: 7.5,
		},
		{
			name: "comparison with arithmetic: 5 + 3 > 2 * 3",
			expr: Binary{
				Left: Binary{
					Left:     Literal{Value: Token{TokenType: Number, Literal: 5.0}},
					Operator: Token{TokenType: Plus, Lexeme: "+"},
					Right:    Literal{Value: Token{TokenType: Number, Literal: 3.0}},
				},
				Operator: Token{TokenType: Greater, Lexeme: ">"},
				Right: Binary{
					Left:     Literal{Value: Token{TokenType: Number, Literal: 2.0}},
					Operator: Token{TokenType: Star, Lexeme: "*"},
					Right:    Literal{Value: Token{TokenType: Number, Literal: 3.0}},
				},
			},
			expected: true,
		},
		{
			name: "negation of grouped expression: -(5 + 3)",
			expr: Unary{
				Operator: Token{TokenType: Minus, Lexeme: "-"},
				Right: Grouping{
					Expr: Binary{
						Left:     Literal{Value: Token{TokenType: Number, Literal: 5.0}},
						Operator: Token{TokenType: Plus, Lexeme: "+"},
						Right:    Literal{Value: Token{TokenType: Number, Literal: 3.0}},
					},
				},
			},
			expected: -8.0,
		},
		{
			name: "logical not of comparison: !(5 > 3)",
			expr: Unary{
				Operator: Token{TokenType: Bang, Lexeme: "!"},
				Right: Grouping{
					Expr: Binary{
						Left:     Literal{Value: Token{TokenType: Number, Literal: 5.0}},
						Operator: Token{TokenType: Greater, Lexeme: ">"},
						Right:    Literal{Value: Token{TokenType: Number, Literal: 3.0}},
					},
				},
			},
			expected: false,
		},
		{
			name: "deeply nested: ((10 - 5) * 2) + (3 / 3)",
			expr: Binary{
				Left: Grouping{
					Expr: Binary{
						Left: Grouping{
							Expr: Binary{
								Left:     Literal{Value: Token{TokenType: Number, Literal: 10.0}},
								Operator: Token{TokenType: Minus, Lexeme: "-"},
								Right:    Literal{Value: Token{TokenType: Number, Literal: 5.0}},
							},
						},
						Operator: Token{TokenType: Star, Lexeme: "*"},
						Right:    Literal{Value: Token{TokenType: Number, Literal: 2.0}},
					},
				},
				Operator: Token{TokenType: Plus, Lexeme: "+"},
				Right: Grouping{
					Expr: Binary{
						Left:     Literal{Value: Token{TokenType: Number, Literal: 3.0}},
						Operator: Token{TokenType: Slash, Lexeme: "/"},
						Right:    Literal{Value: Token{TokenType: Number, Literal: 3.0}},
					},
				},
			},
			expected: 11.0,
		},
		{
			name: "equality with different types",
			expr: Binary{
				Left:     Literal{Value: Token{TokenType: Number, Literal: 5.0}},
				Operator: Token{TokenType: EqualEqual, Lexeme: "=="},
				Right:    Literal{Value: Token{TokenType: String, Literal: "5"}},
			},
			expected: false,
		},
		{
			name: "chained comparisons: (5 > 3) == true",
			expr: Binary{
				Left: Grouping{
					Expr: Binary{
						Left:     Literal{Value: Token{TokenType: Number, Literal: 5.0}},
						Operator: Token{TokenType: Greater, Lexeme: ">"},
						Right:    Literal{Value: Token{TokenType: Number, Literal: 3.0}},
					},
				},
				Operator: Token{TokenType: EqualEqual, Lexeme: "=="},
				Right:    Literal{Value: Token{TokenType: True}},
			},
			expected: true,
		},
		{
			name: "multiple string concatenations",
			expr: Binary{
				Left: Binary{
					Left:     Literal{Value: Token{TokenType: String, Literal: "hello"}},
					Operator: Token{TokenType: Plus, Lexeme: "+"},
					Right:    Literal{Value: Token{TokenType: String, Literal: " "}},
				},
				Operator: Token{TokenType: Plus, Lexeme: "+"},
				Right:    Literal{Value: Token{TokenType: String, Literal: "world"}},
			},
			expected: "hello world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lox := &Lox{}
			i := NewInterpreter(lox)
			result, err := i.Interpret(tt.expr)

			if (err != nil) != tt.wantErr {
				t.Errorf("Interpret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("Interpret() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestInterpreter_Grouping(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expr
		expected any
		wantErr  bool
	}{
		{
			name: "simple grouping of number",
			expr: Grouping{
				Expr: Literal{Value: Token{TokenType: Number, Literal: 42.0}},
			},
			expected: 42.0,
		},
		{
			name: "grouping of string",
			expr: Grouping{
				Expr: Literal{Value: Token{TokenType: String, Literal: "test"}},
			},
			expected: "test",
		},
		{
			name: "nested grouping",
			expr: Grouping{
				Expr: Grouping{
					Expr: Literal{Value: Token{TokenType: Number, Literal: 10.0}},
				},
			},
			expected: 10.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lox := &Lox{}
			i := NewInterpreter(lox)
			result, err := i.Interpret(tt.expr)

			if (err != nil) != tt.wantErr {
				t.Errorf("Interpret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("Interpret() = %v, want %v", result, tt.expected)
			}
		})
	}
}
