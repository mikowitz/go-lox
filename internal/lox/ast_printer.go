package lox

import (
	"strings"
)

type AstPrinter struct {
	strings.Builder
}

func NewAstPrinter() AstPrinter {
	return AstPrinter{}
}

func (a *AstPrinter) Print(e Expr) string {
	a.Reset()
	e.Accept(a)
	return a.String()
}

func (a *AstPrinter) VisitBinary(b Binary) {
	a.parenthesize(b.Operator.Lexeme, b.Left, b.Right)
}

func (a *AstPrinter) VisitGrouping(g Grouping) {
	a.parenthesize("group", g.Expr)
}

func (a *AstPrinter) VisitLiteral(l Literal) {
	if l.Value.TokenType == Nil {
		a.WriteString("nil")
		return
	}
	a.WriteString(l.Value.Lexeme)
}

func (a *AstPrinter) VisitUnary(u Unary) {
	a.parenthesize(u.Operator.Lexeme, u.Right)
}

func (a *AstPrinter) parenthesize(label string, exprs ...Expr) {
	a.WriteString("(")
	a.WriteString(label)
	for _, e := range exprs {
		a.WriteString(" ")
		e.Accept(a)
	}
	a.WriteString(")")
}
