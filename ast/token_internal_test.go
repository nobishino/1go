package ast

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewToken(t *testing.T) {
	testcases := [...]struct {
		title  string
		kind   TokenKind
		str    string
		expect *Token
	}{
		{
			title: "+",
			kind:  TKReserved,
			str:   "+",
			expect: &Token{
				kind: TKReserved,
				str:  "+",
				len:  1,
			},
		},
	}
	for _, tt := range testcases {
		t.Run(tt.title, func(t *testing.T) {
			token := Token{}
			got := newToken(tt.kind, &token, tt.str)
			if diff := cmp.Diff(got, tt.expect, cmp.AllowUnexported(Token{})); diff != "" {
				t.Errorf("[%s] differs: (-got +expect)\n%s", tt.title, diff)
			}
		})
	}
}

func TestTokenizeToLinkedList(t *testing.T) {
	testcases := [...]struct {
		title  string
		kind   TokenKind
		str    string
		expect *Token
	}{
		{
			title: "1+1(with space)",
			kind:  TKReserved,
			str:   " 1 + 1 ",
			expect: &Token{
				kind: TKNum,
				str:  "1",
				val:  1,
				next: &Token{
					kind: TKReserved,
					str:  "+",
					len:  1,
					next: &Token{
						kind: TKNum,
						str:  "1",
						val:  1,
						next: &Token{kind: TKEOF},
					},
				},
			},
		},
		{
			title: "1+1",
			kind:  TKReserved,
			str:   "1+1",
			expect: &Token{
				kind: TKNum,
				str:  "1",
				val:  1,
				next: &Token{
					kind: TKReserved,
					str:  "+",
					len:  1,
					next: &Token{
						kind: TKNum,
						str:  "1",
						val:  1,
						next: &Token{kind: TKEOF},
					},
				},
			},
		},
		// {
		// 	title: "1==1",
		// 	kind:  TKReserved,
		// 	str:   "1==1",
		// 	expect: &Token{
		// 		kind: TKNum,
		// 		str:  "1",
		// 		val:  1,
		// 		next: &Token{
		// 			kind: TKReserved,
		// 			str:  "==",
		// 			next: &Token{
		// 				kind: TKNum,
		// 				str:  "1",
		// 				val:  1,
		// 				next: &Token{kind: TKEOF},
		// 			},
		// 		},
		// 	},
		// },
	}
	for _, tt := range testcases {
		t.Run(tt.title, func(t *testing.T) {
			got, err := tokenize(tt.str)
			if err != nil {
				t.Errorf("expect error to be nil, but got %v", err)
			}
			if diff := cmp.Diff(got, tt.expect, cmp.AllowUnexported(Token{})); diff != "" {
				t.Errorf("[%s] differs: (-got +expect)\n%s", tt.title, diff)
			}
		})
	}
}
