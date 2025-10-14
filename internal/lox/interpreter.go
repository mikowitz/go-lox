package lox

import (
	"fmt"
	"log"
)

type Interpreter struct {
	value any
	err   error
	lox   *Lox
}

func NewInterpreter(lox *Lox) Interpreter {
	return Interpreter{lox: lox}
}

func (i *Interpreter) Interpret(e Expr) (any, error) {
	i.value = nil
	e.Accept(i)
	return i.value, i.err
}

func (i *Interpreter) VisitBinary(b Binary) {
	b.Left.Accept(i)
	left := i.value
	b.Right.Accept(i)
	right := i.value

	switch b.Operator.TokenType {
	case BangEqual:
		i.value = left != right
	case EqualEqual:
		i.value = left == right
	case Greater:
		l, r, err := checkNumbers(b.Operator, left, right)
		if err != nil {
			i.err = err
			i.lox.runtimeError(err, b.Operator.Line)
			return
		}
		i.value = l > r
	case GreaterEqual:
		l, r, err := checkNumbers(b.Operator, left, right)
		if err != nil {
			i.err = err
			i.lox.runtimeError(err, b.Operator.Line)
			return
		}
		i.value = l >= r
	case Less:
		l, r, err := checkNumbers(b.Operator, left, right)
		if err != nil {
			i.err = err
			i.lox.runtimeError(err, b.Operator.Line)
			return
		}
		i.value = l < r
	case LessEqual:
		l, r, err := checkNumbers(b.Operator, left, right)
		if err != nil {
			i.err = err
			i.lox.runtimeError(err, b.Operator.Line)
			return
		}
		i.value = l <= r
	case Minus:
		l, r, err := checkNumbers(b.Operator, left, right)
		if err != nil {
			i.err = err
			i.lox.runtimeError(err, b.Operator.Line)
			return
		}
		i.value = l - r
	case Slash:
		l, r, err := checkNumbers(b.Operator, left, right)
		if err != nil {
			i.err = err
			i.lox.runtimeError(err, b.Operator.Line)
			return
		}
		i.value = l / r
	case Star:
		l, r, err := checkNumbers(b.Operator, left, right)
		if err != nil {
			i.err = err
			i.lox.runtimeError(err, b.Operator.Line)
			return
		}
		i.value = l * r
	case Plus:
		l, r, err := checkNumbers(b.Operator, left, right)
		log.Printf("got %v, %v, %v", l, r, err)
		if err == nil {
			i.value = l + r
			return
		}
		lStr, rStr, errStr := checkStrings(b.Operator, left, right)
		if errStr == nil {
			i.value = lStr + rStr
			return
		}
		i.err = fmt.Errorf("operands to + must be two numbers or two strings")
		i.lox.runtimeError(i.err, b.Operator.Line)
	}
}

func (i *Interpreter) VisitGrouping(g Grouping) {
	g.Expr.Accept(i)
}

func (i *Interpreter) VisitLiteral(l Literal) {
	switch l.Value.TokenType {
	case False:
		i.value = false
	case True:
		i.value = true
	case Nil:
		i.value = nil
	case Number, String:
		i.value = l.Value.Literal
	}
}

func (i *Interpreter) VisitUnary(u Unary) {
	u.Right.Accept(i)

	switch u.Operator.TokenType {
	case Minus:
		f, err := checkNumber(u.Operator, i.value)
		if err != nil {
			i.err = err
			return
		}
		i.value = -f
		return
	case Bang:
		i.value = !isTruthy(i.value)
		return
	default:
		i.value = nil
	}
}

func isTruthy(b any) bool {
	if b == nil {
		return false
	}
	if b, ok := b.(bool); ok {
		return b
	}
	return true
}

func checkNumber(operator Token, operand any) (float64, error) {
	n, ok := operand.(float64)
	if ok {
		return n, nil
	}
	return n, fmt.Errorf(
		"operand to %s must be a number, got %v",
		operator.Lexeme, operand,
	)
}

func checkNumbers(operator Token, left, right any) (float64, float64, error) {
	l, lok := left.(float64)
	r, rok := right.(float64)
	if lok && rok {
		return l, r, nil
	}

	return l, r, fmt.Errorf(
		"operands to %s must be numbers, got %v, %v",
		operator.Lexeme, l, r,
	)
}

func checkStrings(operator Token, left, right any) (string, string, error) {
	l, lok := left.(string)
	r, rok := right.(string)
	if lok && rok {
		return l, r, nil
	}

	return l, r, fmt.Errorf(
		"operands to %s must be strings, got %v, %v",
		operator.Lexeme, l, r,
	)
}
