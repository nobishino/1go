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
		{
			title:  "difference of 2 numbers",
			tokens: []string{"10", "-", "8"},
			expect: &ast.Node{
				Kind: ast.Sub,
				Lhs: &ast.Node{
					Kind:  ast.Num,
					Value: 10,
				},
				Rhs: &ast.Node{
					Kind:  ast.Num,
					Value: 8,
				},
			},
		},
		{
			title:  "4 + 21 - 12",
			tokens: []string{"4", "+", "21", "-", "12"},
			expect: &ast.Node{
				Kind: ast.Sub,
				Lhs: &ast.Node{
					Kind: ast.Add,
					Lhs: &ast.Node{
						Kind:  ast.Num,
						Value: 4,
					},
					Rhs: &ast.Node{
						Kind:  ast.Num,
						Value: 21,
					},
				},
				Rhs: &ast.Node{
					Kind:  ast.Num,
					Value: 12,
				},
			},
		},
		{
			title:  "1 + 2 * 3",
			tokens: []string{"1", "+", "2", "*", "3"},
			expect: &ast.Node{
				Kind: ast.Add,
				Lhs: &ast.Node{
					Kind:  ast.Num,
					Value: 1,
				},
				Rhs: &ast.Node{
					Kind: ast.Mul,
					Lhs: &ast.Node{
						Kind:  ast.Num,
						Value: 2,
					},
					Rhs: &ast.Node{
						Kind:  ast.Num,
						Value: 3,
					},
				},
			},
		},
		{
			title:  "1 + 2 / 3",
			tokens: []string{"1", "+", "2", "/", "3"},
			expect: &ast.Node{
				Kind: ast.Add,
				Lhs: &ast.Node{
					Kind:  ast.Num,
					Value: 1,
				},
				Rhs: &ast.Node{
					Kind: ast.Div,
					Lhs: &ast.Node{
						Kind:  ast.Num,
						Value: 2,
					},
					Rhs: &ast.Node{
						Kind:  ast.Num,
						Value: 3,
					},
				},
			},
		},
		{
			title:  "5 * (9 - 6)",
			tokens: []string{"5", "*", "(", "9", "-", "6", ")"},
			expect: &ast.Node{
				Kind: ast.Mul,
				Lhs: &ast.Node{
					Kind:  ast.Num,
					Value: 5,
				},
				Rhs: &ast.Node{
					Kind: ast.Sub,
					Lhs: &ast.Node{
						Kind:  ast.Num,
						Value: 9,
					},
					Rhs: &ast.Node{
						Kind:  ast.Num,
						Value: 6,
					},
				},
			},
		},
		{
			title:  "-1(unary)",
			tokens: []string{"-", "1"},
			expect: &ast.Node{
				Kind: ast.Sub,
				Lhs: &ast.Node{
					Kind:  ast.Num,
					Value: 0,
				},
				Rhs: &ast.Node{
					Kind:  ast.Num,
					Value: 1,
				},
			},
		},
		{
			title:  "+1(unary)",
			tokens: []string{"+", "1"},
			expect: &ast.Node{
				Kind:  ast.Num,
				Value: 1,
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
				t.Errorf("differs: (-got +expect)\n%s\ninput = %v", diff, tt.tokens)
			}
		})
	}
}
