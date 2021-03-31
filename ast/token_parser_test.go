package ast_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nobishino/1go/ast"
)

func TestTParser(t *testing.T) {
	testcases := [...]struct {
		in     string
		expect *ast.Node
	}{
		{
			in: "1",
			expect: &ast.Node{
				Value: 1,
				Kind:  ast.Num,
			},
		},
		{
			in: "1+1",
			expect: &ast.Node{
				Kind: ast.Add,
				Lhs: &ast.Node{
					Value: 1,
					Kind:  ast.Num,
				},
				Rhs: &ast.Node{
					Value: 1,
					Kind:  ast.Num,
				},
			},
		},
		{
			in: "1+2+3",
			expect: &ast.Node{
				Kind: ast.Add,
				Lhs: &ast.Node{
					Kind: ast.Add,
					Lhs: &ast.Node{
						Value: 1,
						Kind:  ast.Num,
					},
					Rhs: &ast.Node{
						Value: 2,
						Kind:  ast.Num,
					},
				},
				Rhs: &ast.Node{
					Value: 3,
					Kind:  ast.Num,
				},
			},
		},
		{
			in: "3 - 2",
			expect: &ast.Node{
				Kind: ast.Sub,
				Lhs: &ast.Node{
					Kind:  ast.Num,
					Value: 3,
				},
				Rhs: &ast.Node{
					Kind:  ast.Num,
					Value: 2,
				},
			},
		},
		{
			in: "1 + 2*3",
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
			in: "5-4/2",
			expect: &ast.Node{
				Kind: ast.Sub,
				Lhs: &ast.Node{
					Kind:  ast.Num,
					Value: 5,
				},
				Rhs: &ast.Node{
					Kind: ast.Div,
					Lhs: &ast.Node{
						Kind:  ast.Num,
						Value: 4,
					},
					Rhs: &ast.Node{
						Kind:  ast.Num,
						Value: 2,
					},
				},
			},
		},
	}
	for _, tt := range testcases {
		t.Run(tt.in, func(t *testing.T) {
			p, err := ast.NewTParser(tt.in)
			if err != nil {
				t.Errorf("expect error to be nil but got %v while creating parser", err)
			}
			got, err := p.Parse()
			if err != nil {
				t.Errorf("expect error to be nil but got %+v while parsing source %q", err, tt.in)
			}
			if diff := cmp.Diff(got, tt.expect); diff != "" {
				t.Errorf("input: %s\ndiffers: (-got +expect)\n%s\n", tt.in, diff)
			}
		})
	}
}
