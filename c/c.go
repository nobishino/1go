package c

import (
	"bytes"
	"fmt"
)

// Compile compiles source code and returns assembly code
func Compile(src string) string {
	buf := new(bytes.Buffer)
	fmt.Fprintln(buf, ".intel_syntax noprefix")
	fmt.Fprintln(buf, ".globl main")
	fmt.Fprintln(buf, "")
	fmt.Fprintln(buf, "main:")
	op := "mov"
	tokens := Tokenize(src)
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
	return buf.String()
}

func Tokenize(src string) []string {
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
		result = append(result, token)
	}
	return result
}
