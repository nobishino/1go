package ast

import (
	"strconv"

	"golang.org/x/xerrors"
)

type TokenKind int

const (
	TKReserved TokenKind = iota
	TKNum
	TKEOF
	TKIDENT
)

func (tk TokenKind) String() string {
	switch tk {
	case TKReserved:
		return "RESERVED"
	case TKNum:
		return "NUM"
	case TKIDENT:
		return "IDENTIFIER"
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

// 新しいIDENT Tokenを作成してcurにつなげる
func newIdentToken(cur *Token, str string) (*Token, error) {
	// validation
	for i, r := range str {
		if !isLatin(r) {
			return nil, xerrors.Errorf("%q is not a latin character which is illegal as variable name, %dth charachter of %q",
				r, i+1, str)
		}
	}
	new := &Token{
		kind: TKIDENT,
		str:  str,
		len:  len(str),
		// TODO: lenは設定する必要があるか？
	}
	cur.next = new
	return new, nil
}

// 新しい数値Tokenを作成してcurにつなげる
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

var reserved = map[int]map[string]bool{
	1: {
		"+": true,
		"-": true,
		"*": true,
		"/": true,
		"(": true,
		")": true,
		"<": true,
		">": true,
		"=": true,
		";": true,
	},
	2: {
		"==": true,
		"!=": true,
		"<=": true,
		">=": true,
	},
}

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
		reservedWord := func() string {
			if len(rs) > 1 && reserved[2][string(rs[:2])] {
				return string(rs[:2])
			}
			if reserved[1][string(rs[:1])] {
				return string(rs[:1])
			}
			return ""
		}()
		if len(reservedWord) > 0 { // 何らかの予約語トークンにマッチした場合
			cur = newToken(TKReserved, cur, reservedWord)
			rs = rs[len(reservedWord):]
			continue
		}

		if i := readLatin(rs); i > 0 {
			c, err := newIdentToken(cur, string(rs[:i]))
			if err != nil {
				return nil, xerrors.Errorf("failed to read IDENT Token. cause: %w", err)
			}
			cur = c
			rs = rs[i:]
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
	for i < len(rs) && isDigit(rs[i]) {
		i++
	}
	return i
}

// 何桁目までがラテン文字 'a'-'z'及び'A'-'Z'であるかを返す
func readLatin(rs []rune) int {
	var i int
	for i < len(rs) && isLatin(rs[i]) {
		i++
	}
	return i
}

// isLatin は、rがラテン文字である時にtrueを返す
func isLatin(r rune) bool {
	if 'a' <= r && r <= 'z' {
		return true
	}
	if 'A' <= r && r <= 'Z' {
		return true
	}
	return false
}

func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}
