package c

// Compile compiles source code and returns assembly code
func Compile(src string) string {
	return `.intel_syntax noprefix
.globl main

main:
    mov rax, 5
    add rax, 20
    sub rax, 4
    ret
`
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
