package gocc

import (
	"fmt"
	"log"
)

func Gen(node *Node) {
	if node.kind == NdNUM {
		// spは16Nづつ動かす cf. [ARM ABI](url)
		fmt.Println("  sub sp, sp, #16")
		// w9 = val
		fmt.Printf("  mov w8, #%d\n", node.val)
		// w8をspから始まるスタsy区領域に書き込む
		// str: store
		fmt.Println("  str w8, [sp]")
		return
	}

	Gen(node.lhs)
	Gen(node.rhs)

	// RHS(右辺)の結果をw9に読み込む
	// w9にspの領域を読み込む(w9は32bitレジスタなのでspから4byte分読み込むと思う)
	// ldr : load register?
	fmt.Println("  ldr w9, [sp]")
	// スタックを進め、左辺の結果を読み込む準備をする
	// 16Nで動かさなければならない
	fmt.Println("  add sp, sp, #16")
	// LHS(左辺)の結果をw8に読み込む
	// w8にspから始まる32bit=4byte分のデータを読み込む
	fmt.Println("  ldr w8, [sp]")
	fmt.Println("  add sp, sp, #16")

	switch node.kind {
	case NdADD:
		// w8 = w8 + w9
		fmt.Println("  add w8, w8, w9")
	case NdSUB:
		// w8 = w8 - w9
		fmt.Println("  sub w8, w8, w9")
	case NdMUL:
		// w8 = w8 * w9
		fmt.Println("  mul w8, w8, w9")
	case NdDIV:
		// w8 = w8 / w9
		fmt.Println("  sdiv w8, w8, w9")
	case NdNUM:
		log.Fatalf("記号ではなく数字トークンを発見しました")
	}

	// 計算結果w8のデータを戻り値として書き込むためspをずらす
	fmt.Println("  sub sp, sp, #16")
	// spから始まる4byte分のスタック領域に戻り値を書き込む
	fmt.Println("  str w8, [sp]")
}
