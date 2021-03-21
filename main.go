package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nobishino/1go/c"
)

func main() {
	if len(os.Args) != 2 {
		log.Println("引数の個数が不正です")
		os.Exit(1)
	}
	asm, err := c.Compile(os.Args[1])
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	fmt.Fprint(os.Stdout, asm)
}
