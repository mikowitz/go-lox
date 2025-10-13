package lox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAstPrinter(t *testing.T) {
	printer := NewAstPrinter()

	e := Binary{
		Left: Unary{
			Operator: NewToken(Minus, "-", nil, 1),
			Right: Literal{
				Value: NewToken(Number, "123", 123.0, 1),
			},
		},
		Operator: NewToken(Star, "*", nil, 1),
		Right: Grouping{
			Expr: Literal{
				Value: NewToken(Number, "45.67", 45.67, 1),
			},
		},
	}

	assert.Equal(t, "(* (- 123) (group 45.67))", printer.Print(e))
}
