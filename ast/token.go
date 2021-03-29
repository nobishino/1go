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

type Token struct {
	kind TokenKind
	next *Token
	val  int    // TKNumの場合の値
	str  string // トークン文字列
}

// var token *Token // 現在着目しているトークン. 連結リスト構造を持つ

// 新しいTokenを作成してtokenにつなげる
func newToken(kind TokenKind, cur *Token, str string) (*Token, error) {
	new := &Token{
		kind: kind,
		str:  str,
	}
	if kind == TKNum {
		val, err := strconv.Atoi(str)
		if err != nil {
			return nil, err
		}
		new.val = val
	}
	cur.next = new
	return new, nil
}

const (
	plus  rune = '+'
	minus rune = '-'
)

// one-char ops: +, -, *, /
// nums 1, 2, 3, 10
// func tokenize(src string) ([]string, error) {
// 	var i int
// 	rs := []rune(src)
// 	read := func(r ...rune) bool {
// 		if !(i+len(r) <= len(rs)) {
// 			return false
// 		}
// 		for k, v := range r {
// 			if rs[i+k] != v {
// 				return false
// 			}
// 		}
// 		i += len(r)
// 		return true
// 	}
// // readDigit := func() bool {
// // 	for k:=i;k < len(rs); k++ {
// // 		if !isDigit(rs[k]) { return false}
// // 	}
// // 	return true
// // }

// 	head := new(Token)
// 	cur := head
// 	for i < len(rs) {
// 		switch {
// 		case read(' '):
// 		case read('+'):
// 			c,err := newToken(TKReserved, cur,"+")
// 			if err != nil {
// 				return nil,err
// 			}
// 			cur = c
// 		case read('-'):
// 			newToken(TKReserved, cur,"-")
// 			c,err := newToken(TKReserved, cur,"+")
// 			if err != nil {
// 				return nil,err
// 			}
// 			cur = c
// 		case read('*'):
// 			newToken(TKReserved, cur,"*")
// 			c,err := newToken(TKReserved, cur,"+")
// 			if err != nil {
// 				return nil,err
// 			}
// 			cur = c
// 		case read('/'):
// 			newToken(TKReserved, cur,"/")
// 			c,err := newToken(TKReserved, cur,"+")
// 			if err != nil {
// 				return nil,err
// 			}
// 			cur = c
// 		}
// 	}
// 	cur = newToken(TKEOF, cur, "")
// }

var digits = map[rune]struct{}{
	'0': {},
	'1': {},
	'2': {},
	'3': {},
	'4': {},
	'5': {},
	'6': {},
	'7': {},
	'8': {},
	'9': {},
}

func isDigit(r rune) bool {
	_, ok := digits[r]
	return ok
}
