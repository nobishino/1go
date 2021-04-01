package ast

import (
	"strconv"
)

type TokenKind int

const (
	TKReserved TokenKind = iota
	TKNum
	TKEOF
)

func (tk TokenKind) String() string {
	switch tk {
	case TKReserved:
		return "RESERVED"
	case TKNum:
		return "NUM"
	case TKEOF:
		return "EOF"
	default:
		return "UNDEFINED"
	}
}

type Token struct {
	kind TokenKind
	next *Token
	val  int    // TKNumの場合の値
	str  string // トークン文字列
	len  int    // トークン文字列の長さ。TKReservedの場合のみ >0
}

// var token *Token // 現在着目しているトークン. 連結リスト構造を持つ

// 新しいTokenを作成してtokenにつなげる
func newNumToken(cur *Token, str string) (*Token, error) {
	new := &Token{
		kind: TKNum,
		str:  str,
	}
	val, err := strconv.Atoi(str)
	if err != nil {
		return nil, err
	}
	new.val = val
	cur.next = new
	return new, nil
}

func newToken(kind TokenKind, cur *Token, str string) *Token {
	if kind == TKNum { // Num tokenは扱えないので何もしない
		return cur
	}
	new := &Token{
		kind: kind,
		str:  str,
		len:  len(str),
	}
	cur.next = new
	return new
}

const (
	plus  rune = '+'
	minus rune = '-'
)

// one-char ops: +, -, *, /
// nums 1, 2, 3, 10
func tokenize(src string) (*Token, error) {
	head := new(Token)
	cur := head
	rs := []rune(src)
	for len(rs) > 0 {
		if isSpace(rs[0]) {
			rs = rs[1:]
			continue
		}
		if len(rs) > 1 {
			head := string(rs[:2])
			if head == "==" {
				cur = newToken(TKReserved, cur, head)
				rs = rs[2:]
				continue
			}
		}
		if rs[0] == '+' || rs[0] == '-' || rs[0] == '*' || rs[0] == '/' || rs[0] == '(' || rs[0] == ')' {
			cur = newToken(TKReserved, cur, string(rs[0]))
			rs = rs[1:]
			continue
		}

		if i := readDigit(rs); i > 0 {
			c, err := newNumToken(cur, string(rs[:i]))
			if err != nil {
				return nil, err
			}
			cur = c
			rs = rs[i:]
			continue
		}
		break // 予想しない文字が来た場合はその場でtokenizeを終了する
	}
	cur = newToken(TKEOF, cur, "")
	return head.next, nil
}

func isSpace(r rune) bool {
	return r == ' '
}

// 何桁目まで数値であるかを返す
func readDigit(rs []rune) int {
	var i int
	for i < len(rs) && digits[rs[i]] {
		i++
	}
	return i
}

var digits = map[rune]bool{
	'0': true,
	'1': true,
	'2': true,
	'3': true,
	'4': true,
	'5': true,
	'6': true,
	'7': true,
	'8': true,
	'9': true,
}

func isDigit(r rune) bool {
	_, ok := digits[r]
	return ok
}
