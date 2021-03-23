package ast_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nobishino/1go/ast"
)

func TestParser(t *testing.T) {
	testcases := [...]struct {
		title  string
		tokens []string
		expect *ast.Node
	}{
		{
			title:  "single digit",
			tokens: []string{"1"},
			expect: &ast.Node{
				Value: 1,
				Kind:  ast.Num,
			},
		},
		{
			title:  "sum of 2 numbers",
			tokens: []string{"1", "+", "2"},
			expect: &ast.Node{
				Kind: ast.Add,
				Lhs: &ast.Node{
					Kind:  ast.Num,
					Value: 1,
				},
				Rhs: &ast.Node{
					Kind:  ast.Num,
					Value: 2,
				},
			},
		},
	}
	for _, tt := range testcases {
		t.Run(tt.title, func(t *testing.T) {
			got, err := ast.NewParser(tt.tokens).Parse()
			if err != nil {
				t.Fatalf("expect error to be nil but got %v", err)
			}
			if diff := cmp.Diff(got, tt.expect); diff != "" {
				t.Errorf("differs: (-got +expect)\n%s", diff)
			}
		})
	}
}
