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
			},
		},
	}
	for _, tt := range testcases {
		t.Run(tt.title, func(t *testing.T) {
			token := Token{}
			got, err := newToken(tt.kind, &token, tt.str)
			if err != nil {
				t.Errorf("expect error to be nil, but got %v", err)
			}
			if diff := cmp.Diff(got, tt.expect, cmp.AllowUnexported(Token{})); diff != "" {
				t.Errorf("[%s] differs: (-got +expect)\n%s", tt.title, diff)
			}
		})
	}
}
