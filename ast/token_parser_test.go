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
					Value: 1,
					Kind:  ast.Num,
				},
				Rhs: &ast.Node{
					Kind: ast.Add,
					Lhs: &ast.Node{
						Value: 2,
						Kind:  ast.Num,
					},
					Rhs: &ast.Node{
						Value: 3,
						Kind:  ast.Num,
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
