package ast

import "golang.org/x/xerrors"

type TParser struct {
	token *Token
}

func NewTParser(src string) (*TParser, error) {
	t, err := tokenize(src)
	if err != nil {
		return nil, err
	}
	return &TParser{
		token: t,
	}, nil
}

func (p *TParser) consume(s string) bool {
	if p.token.kind != TKReserved || p.token.str != s {
		return false
	}
	p.token = p.token.next
	return true
}

func (p *TParser) Parse() (*Node, error) {
	node, err := p.expr()
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (p *TParser) expr() (*Node, error) {
	node, err := p.mul()
	if err != nil {
		return nil, err
	}
	for p.token.kind != TKEOF {
		if p.consume("+") {
			rhs, err := p.mul()
			if err != nil {
				return nil, err
			}
			node = NewNode(Add, node, rhs)
			continue
		}
		if p.consume("-") {
			rhs, err := p.mul()
			if err != nil {
				return nil, err
			}
			node = NewNode(Sub, node, rhs)
			continue
		}
		return nil, xerrors.Errorf("unexpected token %+v (kind = %s)", *p.token, p.token.kind)
	}
	return node, nil
}

func (p *TParser) mul() (*Node, error) {
	node, err := p.primary()
	if err != nil {
		return nil, err
	}
	if p.consume("*") {
		rhs, err := p.primary()
		if err != nil {
			return nil, err
		}
		node = NewNode(Mul, node, rhs)
	}
	if p.consume("/") {
		rhs, err := p.primary()
		if err != nil {
			return nil, err
		}
		node = NewNode(Div, node, rhs)
	}
	return node, nil
}

func (p *TParser) primary() (*Node, error) {
	if p.consume("(") {
		e, err := p.expr()
		if err != nil {
			return nil, err
		}
		if p.consume(")") {
			return e, nil
		}
		return nil, xerrors.Errorf("token ')' is missing in (expr), got %q", p.token.str)
	}
	node, err := p.expectNumber()
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (p *TParser) expectNumber() (*Node, error) {
	if p.token.kind != TKNum {
		return nil, xerrors.Errorf("expect number but token %+v", *p.token)
	}
	node := &Node{
		Kind:  Num,
		Value: p.token.val,
	}
	p.token = p.token.next
	return node, nil
}
