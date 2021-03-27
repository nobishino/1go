package ast

import (
	"strconv"

	"golang.org/x/xerrors"
)

// Parser parses given tokens as *ast.Node
type Parser struct {
	tokens []string
	loc    int
}

func (p *Parser) Parse() (*Node, error) {
	return p.expr()
}

func (p *Parser) expr() (*Node, error) {
	node, err := p.mul()
	if err != nil {
		return nil, err
	}
	if p.consume("+") {
		rhs, err := p.expr()
		if err != nil {
			return nil, err
		}
		return NewNode(Add, node, rhs), nil
	} else if p.consume("-") {
		rhs, err := p.expr()
		if err != nil {
			return nil, err
		}
		return NewNode(Sub, node, rhs), nil
	} else {
		return node, nil
	}
}

func (p *Parser) mul() (*Node, error) {
	node, err := p.primary()
	if err != nil {
		return nil, err
	}
	if p.consume("*") {
		rhs, err := p.mul()
		if err != nil {
			return nil, err
		}
		return NewNode(Mul, node, rhs), nil
	} else if p.consume("/") {
		rhs, err := p.mul()
		if err != nil {
			return nil, err
		}
		return NewNode(Div, node, rhs), nil
	} else {
		return node, nil
	}
}

func (p *Parser) primary() (*Node, error) {
	node, err := p.expectNumber()
	if err != nil {
		return nil, err
	}
	if p.consume("(") {
		exp, err := p.expr()
		if err != nil {
			return nil, err
		}
		if p.consume(")") {
			return exp, nil
		} else {
			return nil, xerrors.Errorf("parse failed at token #%v: %s", p.loc+1, p.current())
		}
	} else {
		return node, nil
	}
}

func (p *Parser) expectNumber() (*Node, error) {
	if !p.hasNext() {
		return nil, xerrors.New("no more token")
	}
	val, err := strconv.Atoi(p.current())
	if err != nil {
		return nil, xerrors.Errorf("failed to parse %s as number, err = %v", p.current(), err)
	}
	p.loc++
	return &Node{
		Kind:  Num,
		Value: val,
	}, nil
}

func (p *Parser) current() string {
	return p.tokens[p.loc]
}

func (p *Parser) hasNext() bool {
	return p.loc < len(p.tokens)
}

func (p *Parser) consume(token string) bool {
	if p.loc >= len(p.tokens) {
		return false
	}
	if token == p.tokens[p.loc] {
		p.loc++
		return true
	}
	return false
}

func NewParser(tokens []string) *Parser {
	return &Parser{
		tokens: tokens,
		loc:    0,
	}
}
