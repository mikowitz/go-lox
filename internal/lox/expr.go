package lox

type Visitor interface {
	VisitBinary(b Binary)
	VisitGrouping(g Grouping)
	VisitLiteral(l Literal)
	VisitUnary(u Unary)
}

type Expr interface {
	Accept(v Visitor)
}

type Binary struct {
	Left, Right Expr
	Operator    Token
}

func (b Binary) Accept(v Visitor) {
	v.VisitBinary(b)
}

type Grouping struct {
	Expr Expr
}

func (g Grouping) Accept(v Visitor) {
	v.VisitGrouping(g)
}

type Literal struct {
	Value Token
}

func (l Literal) Accept(v Visitor) {
	v.VisitLiteral(l)
}

type Unary struct {
	Operator Token
	Right    Expr
}

func (u Unary) Accept(v Visitor) {
	v.VisitUnary(u)
}
