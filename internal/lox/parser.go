package lox

import "fmt"

type Parser struct {
	tokens []Token
	lox    *Lox

	current int
}

func NewParser(tokens []Token, lox *Lox) Parser {
	return Parser{tokens: tokens, lox: lox}
}

func (p *Parser) Parse() Expr {
	if p.isAtEnd() {
		return nil
	}
	return p.Expression()
}

func (p *Parser) Expression() Expr {
	return p.Equality()
}

func (p *Parser) Equality() Expr {
	expr := p.Comparison()

	for p.match(BangEqual, EqualEqual) {
		operator := p.previous()
		right := p.Comparison()
		expr = Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) Comparison() Expr {
	expr := p.Term()

	for p.match(Greater, GreaterEqual, Less, LessEqual) {
		operator := p.previous()
		right := p.Term()
		expr = Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr
}

func (p *Parser) Term() Expr {
	expr := p.Factor()

	for p.match(Minus, Plus) {
		operator := p.previous()
		right := p.Factor()
		expr = Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr
}

func (p *Parser) Factor() Expr {
	expr := p.Unary()

	for p.match(Slash, Star) {
		operator := p.previous()
		right := p.Unary()
		expr = Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr
}

func (p *Parser) Unary() Expr {
	if p.match(Bang, Minus) {
		operator := p.previous()
		right := p.Unary()
		return Unary{Operator: operator, Right: right}
	}
	return p.Primary()
}

func (p *Parser) Primary() Expr {
	var e Expr
	switch {
	case p.match(False, True, Nil, Number, String):
		e = Literal{Value: p.previous()}
	case p.match(LeftParen):
		expr := p.Expression()
		t, err := p.consume(RightParen, "Expect ')' after expression")
		if err != nil {
			p.error(t, err.Error())
			return e
		}
		e = Grouping{Expr: expr}
	default:
		p.error(p.peek(), "Expect expression")
		return nil
	}
	return e
}

//nolint:unused
func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().TokenType == Semicolon {
			return
		}

		switch p.peek().TokenType {
		case Class, Fun, Var, For, If, While, Print, Return:
			return
		}

		p.advance()
	}
}

func (p *Parser) error(token Token, message string) {
	if token.TokenType == EOF {
		p.lox.report(token.Line, " at end", message)
	} else {
		p.lox.report(token.Line, " at '"+token.Lexeme+"'", message)
	}
}

func (p *Parser) match(tokenTypes ...TokenType) bool {
	for _, tokenType := range tokenTypes {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) check(tokenType TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().TokenType == tokenType
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) consume(tokenType TokenType, message string) (Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}
	return Token{}, fmt.Errorf("%s", message)
}

func (p *Parser) isAtEnd() bool {
	return p.peek().TokenType == EOF
}
