package main

import (
	"fmt"
	"github.com/x0y14/gocc"
	"os"
)

func main() {
	argv := os.Args
	if len(argv) != 2 {
		_, _ = fmt.Fprintf(os.Stderr, "引数の個数が正しくありません\n")
		os.Exit(1)
	}

	userInput := argv[1]

	token := gocc.Tokenize([]rune(userInput))
	code := gocc.Parse(token)

	asm := gocc.Generate(code)
	fmt.Println(asm)
}
