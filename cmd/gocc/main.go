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
	code := gocc.Parse(token)

	// prologue
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

	fmt.Println("  stp x29, x30, [sp, #-32]!; prologue")

	fmt.Println()

	//fmt.Println("  sub sp, sp, #16")
	//fmt.Println("  str x29, [sp]")
	//
	//fmt.Println("  mov x29, sp")
	//
	//fmt.Printf("  sub sp, sp, #%d\n", 26*16)

	// nodeをたどってコードを生成する
	for _, node := range code {
		if node != nil {
			gocc.Gen(node)
			//// 最終的な計算結果はスタックに保存されているので取り出す
			//fmt.Println("  ldr x8, [sp]")
			//// 最終的な計算結果の保存に使用したスタック分spを戻してあげる
			//fmt.Println("  add sp, sp, #16")
		}
	}

	for _, asm := range gocc.Assembly {
		fmt.Println(asm.String())
	}

	// epilogue

	// w0はプログラムの結果として使用されるレジスタ
	// w0に最終結果を書き込んであげる
	// w0 =x8
	//fmt.Println("  mov x0, x8")
	//// 一番最初にxzrを書き込んだ分のspを戻す
	//fmt.Println("  add sp, sp, #16")
	//
	//// 終了
	//fmt.Println("  ret")
}
