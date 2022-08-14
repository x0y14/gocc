package main

import (
	"fmt"
	"github.com/x0y14/gocc"
	"log"
	"os"
	"runtime"
)

func main() {
	argv := os.Args
	if len(argv) != 2 {
		_, _ = fmt.Fprintf(os.Stderr, "引数の個数が正しくありません\n")
		os.Exit(1)
	}

	userInput := argv[1]

	token := gocc.Tokenize([]rune(userInput))
	node := gocc.Parse(token)

	fmt.Println(".text")
	fmt.Println(".align 2")

	// macos: _main
	// linux:  main
	var mainSymbol string
	if runtime.GOOS == "darwin" {
		mainSymbol = "_main"
	} else if runtime.GOOS == "linux" {
		mainSymbol = "main"
	} else {
		log.Fatalf("unsupported platform: %v", runtime.GOOS)
	}
	fmt.Printf(".global %s\n", mainSymbol)
	fmt.Printf("%s:\n", mainSymbol)
	// wzr(32bit zero register)をスタックに書き込む
	// 何をしてるの不明
	// https://developer.arm.com/documentation/den0024/a/ch05s01s03
	// コーディングテクニックと書いてある..?
	// もしかしたら現在地の保存をしているのかもしれない..? -> ただ、参照されている感じがしない
	// 16Nづつ動かす
	fmt.Println("  sub sp, sp, #16")
	fmt.Println("str wzr, [sp]")

	// nodeをたどってコードを生成する
	gocc.Gen(node)

	// 最終的な計算結果はスタックに保存されているので取り出す
	fmt.Println("  ldr w8, [sp]")
	// 最終的な計算結果の保存に使用したスタック分spを戻してあげる
	fmt.Println("  add sp, sp, #16")

	// w0はプログラムの結果として使用されるレジスタ
	// w0に最終結果を書き込んであげる
	// w0 = w8
	fmt.Println("  mov w0, w8")
	// 一番最初にwzrを書き込んだ分のspを戻す
	fmt.Println("  add sp, sp, #16")

	// 終了
	fmt.Println("  ret")
}
