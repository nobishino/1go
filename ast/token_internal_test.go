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
		source string
		expect *Token
	}{
		{
			title:  "1+1(with space)",
			source: " 1 + 1 ",
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
			title:  "1+1",
			source: "1+1",
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
			title:  "1==1",
			source: "1==1",
			expect: &Token{
				kind: TKNum,
				str:  "1",
				val:  1,
				next: &Token{
					kind: TKReserved,
					str:  "==",
					len:  2,
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
			title:  "a=1",
			source: "a=1",
			expect: &Token{
				kind: TKIDENT,
				str:  "a",
				len:  1,
				next: &Token{
					kind: TKReserved,
					str:  "=",
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
			title:  "a=z",
			source: "a=z",
			expect: &Token{
				kind: TKIDENT,
				str:  "a",
				len:  1,
				next: &Token{
					kind: TKReserved,
					str:  "=",
					len:  1,
					next: &Token{
						kind: TKIDENT,
						str:  "z",
						len:  1,
						next: &Token{kind: TKEOF},
					},
				},
			},
		},
		{
			title:  "1;1",
			source: "1;1",
			expect: &Token{
				kind: TKNum,
				val:  1,
				str:  "1",
				next: &Token{
					kind: TKReserved,
					str:  ";",
					len:  1,
					next: &Token{
						kind: TKNum,
						val:  1,
						str:  "1",
						next: &Token{kind: TKEOF},
					},
				},
			},
		},
	}
	for _, tt := range testcases {
		t.Run(tt.title, func(t *testing.T) {
			got, err := tokenize(tt.source)
			if err != nil {
				t.Errorf("expect error to be nil, but got %v", err)
			}
			if diff := cmp.Diff(got, tt.expect, cmp.AllowUnexported(Token{})); diff != "" {
				t.Errorf("[%s] differs: (-got +expect)\n%s", tt.title, diff)
			}
		})
	}
}
