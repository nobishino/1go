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
	fmt.Fprint(os.Stdout, c.Compile(os.Args[1]))
}
