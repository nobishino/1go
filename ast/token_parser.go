package ast

import (
	"fmt"

	"golang.org/x/xerrors"
)

type TParser struct {
	token *Token
	pos   int
	lvar  *LVar
}

func NewTParser(src string) (*TParser, error) {
	t, err := tokenize(src)
	if err != nil {
		return nil, err
	}
	return &TParser{
		token: t,
		lvar:  &LVar{}, // offset = 0 で name == ""のダミーローカル変数を設定しておく
	}, nil
}

// 仮実装。全体で1つの関数＝ローカル変数空間しか存在しないという前提でParserがoffsetを返すようにしておく
// FIXME
func (p *TParser) GetOffset() int {
	return p.lvar.offset
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
	node, err := p.stmt()
	if err != nil {
		return nil, err
	}
	return node, nil
}

// Program は、複数のStatementを含むプログラムソースコードをparseし、1つのStatementを1つのNodeとするスライスとして返す.
func (p *TParser) Program() ([]*Node, error) {
	var result []*Node
	for p.token.kind != TKEOF {
		node, err := p.stmt()
		if err != nil {
			return result, xerrors.Errorf("failed to parse program. cause: %w", err)
		}
		result = append(result, node)
	}
	return result, nil
}

func (p *TParser) stmt() (*Node, error) {
	node, err := p.expr()
	if err != nil {
		return nil, xerrors.Errorf("failed to parse statement. cause: %w", err)
	}
	if err := p.expect(";"); err != nil {
		return nil, xerrors.Errorf("failed to parse statement. cause: %w", err)
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
	if len(name) == 0 {
		return nil, false
	}
	// 初めて現れたローカル変数名である場合は登録する
	lvar := p.findLVar(p.token)
	if lvar == nil {
		newLVar := &LVar{
			name:   name,
			len:    len(name),
			next:   p.lvar,
			offset: p.lvar.offset + 8,
		}
		p.lvar = newLVar
		lvar = newLVar
	}
	p.token = p.token.next
	return &Node{
		Kind:   LocalVar,
		Name:   name,
		Offset: lvar.offset,
	}, true
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

func (p *TParser) expect(s string) error {
	if !p.consume(s) {
		return xerrors.Errorf("expect %q but got %v", s, p.token)
	}
	return nil
}

// 指定されたtokenに合致するローカル変数をすでに定義されたローカル変数から検索する。
// 存在しなければnilを返す。
func (p *TParser) findLVar(token *Token) *LVar {
	if token.kind != TKIDENT {
		return nil
	}
	lvar := p.lvar
	for lvar != nil {
		if lvar.name == token.str {
			return lvar
		}
		lvar = lvar.next
	}
	return nil
}

// LVar は、ローカル変数の集まりを管理するための連結リスト
type LVar struct {
	name   string // 変数名
	len    int    // nameの長さ
	offset int    // 変数に割り当てるスタック領域のBase Pointerからのoffset
	next   *LVar  // 1つ前に定義されたLVarへのポインタ
}
