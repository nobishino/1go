package c

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/nobishino/1go/ast"
)

// Compile compiles source code and returns assembly code
func Compile(src string) (string, error) {
	buf := new(bytes.Buffer)
	fmt.Fprintln(buf, ".intel_syntax noprefix")
	fmt.Fprintln(buf, ".globl main")
	fmt.Fprintln(buf, "")
	fmt.Fprintln(buf, "main:")
	op := "mov"
	tokens, err := Tokenize(src)
	if err != nil {
		return "", err
	}
	for _, token := range tokens {
		switch token {
		case "+":
			op = "add"
		case "-":
			op = "sub"
		default: // digit token
			fmt.Fprintf(buf, "    %s rax, %s\n", op, token)
		}
	}
	fmt.Fprintln(buf, "    ret")
	return buf.String(), nil
}

func Tokenize(src string) ([]string, error) {
	if err := Validate(src); err != nil {
		return nil, err
	}
	isOp := func(r rune) bool {
		operators := []rune{'+', '-'}
		for _, op := range operators {
			if r == op {
				return true
			}
		}
		return false
	}
	var i, j int
	var result []string
	s := []rune(src)
	for i < len(src) {
		for j < len(src) && !isOp(s[j]) {
			j++
		}
		var token string
		if i == j {
			token = string(s[i : i+1])
			j++
		} else {
			token = string(s[i:j])
		}
		i = j
		token = strings.Trim(token, " ")
		result = append(result, token)
	}
	return result, nil
}

func Validate(src string) error {
	isValid := func(r rune) bool {
		validChars := []rune{'+', '-', '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', ' '}
		for _, c := range validChars {
			if c == r {
				return true
			}
		}
		return false
	}
	for i, r := range src {
		if !isValid(r) {
			markerLine := strings.Repeat(" ", i) + "^"
			return errors.New(fmt.Sprintf("invalid character at %d\n%s\n%s", i+1, src, markerLine))
		}
	}
	return nil
}

func CompileAST(node *ast.Node) []string {
	if node == nil {
		return nil
	}
	result := []string{
		".intel_syntax noprefix",
		".globl main",
		"",
		"main:",
	}
	result = append(result, genAST(node)...)
	result = append(result, "    ret")
	return result
}

func genAST(node *ast.Node) []string {
	if node == nil {
		return nil
	}
	var result []string
	result = append(result, genAST(node.Lhs)...)
	result = append(result, genAST(node.Rhs)...)

	switch node.Kind {
	case ast.Add:
		result = append(result, add...)
	case ast.Sub:
		result = append(result, sub...)
	case ast.Num:
		result = append(result, fmt.Sprintf("    push %d", node.Value))
	}
	return result
}

var headers = []string{
	".intel_syntax noprefix",
	".globl main",
	"",
	"main:",
}

var add = []string{
	"    pop rdi",
	"    pop rax",
	"    add rax rdx",
	"    push rax",
}
var sub = []string{
	"    pop rdi",
	"    pop rax",
	"    sub rax rdx",
	"    push rax",
}
