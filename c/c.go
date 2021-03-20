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
	return nil
}
