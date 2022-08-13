package main

import (
	"fmt"
	"os"
)

func main() {
	argv := os.Args
	if len(argv) != 2 {
		_, _ = fmt.Fprintf(os.Stderr, "引数の個数が正しくありません\n")
		os.Exit(1)
	}

	userInput := argv[1]
}
