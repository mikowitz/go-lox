package lox

import (
	"fmt"
	"slices"
)

type Parser struct {
	tokens []Token
	lox    *Lox

	current int
}

func NewParser(tokens []Token, lox *Lox) Parser {
	return Parser{
		tokens: tokens,
		lox:    lox,
	}
}

func (p *Parser) Parse() (Expr, error) {
	if p.isAtEnd() {
		return nil, nil
	}
	return p.Expression()
}

func (p *Parser) Expression() (Expr, error) {
	return p.Equality()
}

func (p *Parser) Equality() (Expr, error) {
	expr, err := p.Comparison()
	if err != nil {
		return nil, err
	}

	for p.match(BangEqual, EqualEqual) {
		operator := p.previous()
		right, err := p.Comparison()
		if err != nil {
			return nil, err
		}
		expr = Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr, nil
}

func (p *Parser) Comparison() (Expr, error) {
	expr, err := p.Term()
	if err != nil {
		return nil, err
	}

	for p.match(Greater, GreaterEqual, Less, LessEqual) {
		operator := p.previous()
		right, err := p.Term()
		if err != nil {
			return nil, err
		}
		expr = Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr, nil
}

func (p *Parser) Term() (Expr, error) {
	expr, err := p.Factor()
	if err != nil {
		return nil, err
	}

	for p.match(Minus, Plus) {
		operator := p.previous()
		right, err := p.Factor()
		if err != nil {
			return nil, err
		}
		expr = Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr, nil
}

func (p *Parser) Factor() (Expr, error) {
	expr, err := p.Unary()
	if err != nil {
		return nil, err
	}

	for p.match(Slash, Star) {
		operator := p.previous()
		right, err := p.Unary()
		if err != nil {
			return nil, err
		}
		expr = Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr, nil
}

func (p *Parser) Unary() (Expr, error) {
	if p.match(Bang, Minus) {
		operator := p.previous()
		right, err := p.Unary()
		if err != nil {
			return nil, err
		}
		return Unary{Operator: operator, Right: right}, nil
	}

	return p.Primary()
}

func (p *Parser) Primary() (Expr, error) {
	var expr Expr
	switch {
	case p.match(False, True, Nil, Number, String):
		expr = Literal{Value: p.previous()}
	case p.match(LeftParen):
		innerExpr, err := p.Expression()
		if err != nil {
			return nil, err
		}
		err = p.consume(RightParen, "expect ')' after expression")
		if err != nil {
			p.error(p.previous(), err.Error())
			return nil, err
		}
		expr = Grouping{Expr: innerExpr}
	default:
		token := p.peek()
		errStr := fmt.Sprintf("unexpected token '%s'", token.Lexeme)
		p.error(token, errStr)
		return nil, fmt.Errorf("%s", errStr)
	}

	return expr, nil
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

func (p *Parser) consume(tokenType TokenType, message string) error {
	if p.check(tokenType) {
		p.advance()
		return nil
	}
	return fmt.Errorf("%s", message)
}

func (p *Parser) match(tokenTypes ...TokenType) bool {
	if slices.ContainsFunc(tokenTypes, func(tokenType TokenType) bool {
		return p.check(tokenType)
	}) {
		p.advance()
		return true
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

func (p *Parser) isAtEnd() bool {
	return p.peek().TokenType == EOF
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}
