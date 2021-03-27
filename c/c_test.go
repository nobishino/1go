package c_test

import (
	"io"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/nobishino/1go/ast"
	"github.com/nobishino/1go/c"
)

func TestTokenize(t *testing.T) {
	testcases := [...]struct {
		title  string
		in     string
		expect []string
	}{
		{
			title:  "1 term",
			in:     "5",
			expect: []string{"5"},
		},
		{
			title:  "2 terms#1",
			in:     "5+3",
			expect: []string{"5", "+", "3"},
		},
		{
			title:  "2 terms#2",
			in:     "5-31",
			expect: []string{"5", "-", "31"},
		},
		{
			title:  "3 terms",
			in:     "5+13-4",
			expect: []string{"5", "+", "13", "-", "4"},
		},
		{
			title:  "3 terms#2",
			in:     " 5 + 13 - 4",
			expect: []string{"5", "+", "13", "-", "4"},
		},
	}
	for _, tt := range testcases {
		t.Run(tt.title, func(t *testing.T) {
			got, err := c.Tokenize(tt.in)
			if err != nil {
				t.Errorf("error should be nil but got %v", err)
			}
			if diff := cmp.Diff(got, tt.expect); diff != "" {
				t.Errorf("differs: (-got +expect)\n%s", diff)
			}
		})
	}
	t.Run("Invalid Input", func(t *testing.T) {
		testcases := [...]struct {
			title string
			in    string
		}{
			{
				title: "invalid character",
				in:    "f",
			},
		}
		for _, tt := range testcases {
			t.Run(tt.title, func(t *testing.T) {
				_, err := c.Tokenize(tt.in)
				if err == nil {
					t.Fatal("error should not be nil but got nil")
				}
			})
		}
	})
}

func TestValidate(t *testing.T) {
	t.Run("Valid Input", func(t *testing.T) {
		testcases := [...]struct {
			title string
			in    string
		}{
			{
				title: "only digits",
				in:    "456",
			},
			{
				title: "has space",
				in:    "4 5 6",
			},
			{
				title: "mixture",
				in:    "4 + 50 - 6",
			},
		}
		for _, tt := range testcases {
			t.Run(tt.title, func(t *testing.T) {
				err := c.Validate(tt.in)
				if err != nil {
					t.Errorf("expect err to be nil, but got %v", err)
				}
			})
		}
	})
	t.Run("Invalid Input", func(t *testing.T) {
		testcases := [...]struct {
			title  string
			in     string
			errMsg string
		}{
			{
				title: "latin",
				in:    "a",
				errMsg: `invalid character at 1
a
^`,
			},
		}
		for _, tt := range testcases {
			t.Run(tt.title, func(t *testing.T) {
				err := c.Validate(tt.in)
				if err == nil {
					t.Fatal("expect err not to be nil, but got nil")
				}
				gotErrMsg := err.Error()
				if gotErrMsg != tt.errMsg {
					t.Errorf("expect error message %q, but got %q", tt.errMsg, gotErrMsg)
				}
			})
		}
	})
}

func TestAddSub(t *testing.T) {
	testcases := [...]struct {
		in     string
		expect string
	}{
		{
			in:     "5+20-4",
			expect: "./testdata/1.s",
		},
	}
	for _, tt := range testcases {
		got, err := c.NaiveAddSub(tt.in)
		if err != nil {
			t.Errorf("error should be nil but got %v", err)
		}
		expect := readTestFile(t, tt.expect)
		if got != expect {
			t.Errorf("\n[expect]\n%s\n[got]\n%s", expect, got)
		}
	}
}

func readTestFile(t *testing.T, path string) string {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	s, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	return string(s)
}

func TestCompileAST(t *testing.T) {
	testcases := [...]struct {
		title  string
		in     *ast.Node
		expect []string
	}{
		{
			title:  "nil",
			in:     nil,
			expect: nil,
		},
		{
			title: "1+2",
			in: &ast.Node{
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
			expect: []string{
				".intel_syntax noprefix",
				".globl main",
				"",
				"main:",
				"    push 1",
				"    push 2",
				"    pop rdi",
				"    pop rax",
				"    add rax, rdi",
				"    push rax",
				"    pop rax",
				"    ret",
				"",
			},
		},
		{
			title: "2-1",
			in: &ast.Node{
				Kind: ast.Sub,
				Lhs: &ast.Node{
					Kind:  ast.Num,
					Value: 2,
				},
				Rhs: &ast.Node{
					Kind:  ast.Num,
					Value: 1,
				},
			},
			expect: []string{
				".intel_syntax noprefix",
				".globl main",
				"",
				"main:",
				"    push 2",
				"    push 1",
				"    pop rdi",
				"    pop rax",
				"    sub rax, rdi",
				"    push rax",
				"    pop rax",
				"    ret",
				"",
			},
		},
	}
	for _, tt := range testcases {
		t.Run(tt.title, func(t *testing.T) {
			got := c.Gen(tt.in)
			if diff := cmp.Diff(got, tt.expect); diff != "" {
				t.Errorf("differs: (-got +expect)\n%s", diff)
			}
		})
	}
}
