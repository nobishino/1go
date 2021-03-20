package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		log.Println("引数の個数が不正です")
		os.Exit(1)
	}
	fmt.Println(".intel_syntax noprefix")
	fmt.Println(".globl main")
	fmt.Println("main:")
	digit, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}
	fmt.Printf("    mov rax, %d\n", digit)
	fmt.Println("    ret")
}
