package ast

import (
	"fmt"

	"golang.org/x/xerrors"
)

type TParser struct {
	token *Token
	pos   int
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
	p.pos++
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
	node, err := p.assign()
	if err != nil {
		return nil, xerrors.Errorf("failed to parse expr %w", err)
	}
	return node, nil
}

func (p *TParser) assign() (*Node, error) {
	node, err := p.equality()
	if err != nil {
		return nil, xerrors.Errorf("failed to parse left hand side of =. caused by %w", err)
	}
	if p.token.kind != TKEOF && p.consume("=") {
		rhs, err := p.assign()
		if err != nil {
			return nil, xerrors.Errorf("failed to parse right hand side of =. caused by %w", err)
		}
		node = NewNode(Assign, node, rhs)
	}
	return node, nil
}

func (p *TParser) debug() {
	fmt.Printf("DEBUG: current pos = %v, kind = %q, label = %q\n", p.pos, p.token.kind, p.token.str)
}

func (p *TParser) equality() (*Node, error) {
	node, err := p.relational()
	if err != nil {
		return nil, xerrors.Errorf("failed to parse equality, because of %w", err)
	}
	for p.token.kind != TKEOF {
		if p.consume("==") {
			rhs, err := p.relational()
			if err != nil {
				return nil, xerrors.Errorf("failed to parse right-hand side of ==, because of %w", err)
			}
			node = NewNode(Eq, node, rhs)
			continue
		}
		if p.consume("!=") {
			rhs, err := p.relational()
			if err != nil {
				return nil, xerrors.Errorf("failed to parse right-hand side of !=, because of %w", err)
			}
			node = NewNode(Neq, node, rhs)
		}
		break
	}
	return node, nil
}

func (p *TParser) relational() (*Node, error) {
	node, err := p.add()
	if err != nil {
		return nil, xerrors.Errorf("failed to parse leftmost part of relational. cause: %w", err)
	}
	for p.token.kind != TKEOF {
		if p.consume("<") {
			rhs, err := p.add()
			if err != nil {
				return nil, xerrors.Errorf("failed to parse right hand side of <. cause: %w", err)
			}
			node = NewNode(LT, node, rhs)
		}
		if p.consume("<=") {
			rhs, err := p.add()
			if err != nil {
				return nil, xerrors.Errorf("failed to parse right hand side of <=. cause: %w", err)
			}
			node = NewNode(LE, node, rhs)
		}
		if p.consume(">") {
			rhs, err := p.add()
			if err != nil {
				return nil, xerrors.Errorf("failed to parse right hand side of <. cause: %w", err)
			}
			node = NewNode(LT, rhs, node) // 逆向きの < としてparseする
		}
		if p.consume(">=") {
			rhs, err := p.add()
			if err != nil {
				return nil, xerrors.Errorf("failed to parse right hand side of >=. cause: %w", err)
			}
			node = NewNode(LE, rhs, node) // 逆向きの <= としてparseする
		}
		break
	}
	return node, nil
}

func (p *TParser) add() (*Node, error) {
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
		break
	}
	return node, nil
}

// func (p *TParser) equality() (*Node,  error) {

// }

func (p *TParser) mul() (*Node, error) {
	node, err := p.unary()
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

func (p *TParser) unary() (*Node, error) {
	if p.consume("+") {
		node, err := p.primary()
		if err != nil {
			xerrors.Errorf("failed to parse unary: %w", err)
		}
		return node, nil
	}
	if p.consume("-") {
		node, err := p.primary()
		if err != nil {
			xerrors.Errorf("failed to parse unary: %w", err)
		}
		zero := &Node{
			Kind:  Num,
			Value: 0,
		}
		return NewNode(Sub, zero, node), nil
	}
	node, err := p.primary()
	if err != nil {
		xerrors.Errorf("failed to parse unary: %w", err)
	}
	return node, nil
}

func (p *TParser) primary() (*Node, error) {
	if p.consume("(") {
		e, err := p.add()
		if err != nil {
			return nil, err
		}
		if p.consume(")") {
			return e, nil
		}
		return nil, xerrors.Errorf("token ')' is missing in (expr), got %q", p.token.str)
	}
	if node, ok := p.parseIfIdentifier(); ok {
		return node, nil
	}
	node, err := p.expectNumber()
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (p *TParser) parseIfIdentifier() (*Node, bool) {
	if p.token.kind != TKIDENT {
		return nil, false
	}
	name := p.token.str
	p.token = p.token.next
	return &Node{
		Kind: Ident,
		Name: name,
	}, true
}

// func (p *TParser) consumeNumberToken() (int, bool) {
// 	if p.token.kind != TKNum {
// 		return 0, false
// 	}
// }

// これだと使いづらい
// tokenを直接触らないようにする
func (p *TParser) ident() (*Node, error) {
	if p.token.kind != TKIDENT {
		return nil, xerrors.Errorf("expect identifier token but got %+v", *p.token)
	}
	node := &Node{
		Kind: Ident,
		Name: p.token.str,
	}
	p.token = p.token.next
	p.pos++
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
	p.pos++
	return node, nil
}
