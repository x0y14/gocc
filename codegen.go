package gocc

import (
	"fmt"
	"log"
)

func genLVal(node *Node) {
	if node.kind != NdLVAR {
		log.Fatalf("代入の左辺値が変数ではありません")
	}
	// x8にフレームポインタを保存
	// x8 = x29
	fmt.Println("  mov x8, x29")
	// １文字変数の位置分スタックを移動
	// x8 = x8 - offset
	fmt.Printf("  sub x8, x8, #%d\n", node.offset)
	// スタックに変数のデータはどこから始まるのかを保存
	// [sp] = x8
	fmt.Println("  str x8, [sp, #-16]!")
}

func Gen(node *Node) {
	switch node.kind {
	case NdNUM:
		// spは16Nづつ動かす cf. [ARM ABI](url)
		fmt.Println("  sub sp, sp, #16")
		// x8 = val
		fmt.Printf("  mov x8, #%d\n", node.val)
		// x8をspから始まるスタック領域に書き込む
		// str: store
		fmt.Println("  str x8, [sp]")
		return
	case NdLVAR:
		genLVal(node)
		// 変数のデータはスタックのどこに書き込まれるのか取得
		// x8 = [sp]
		fmt.Println("  ldr x8, [sp], #16")
		//fmt.Println("  add sp, sp, #16")
		// x8 = [x8]
		// どこに書き込まれるのかというのをアドレスとして読み込みx8に書き込まれた内容を読み込む
		fmt.Println("  ldr x8, [x8]")
		//
		fmt.Println("  sub sp, sp, #16")
		// 一歳に右辺値として読み取った値をスタックに書き込む
		fmt.Println("  str x8, [sp]")
		return
	case NdASSIGN:
		genLVal(node.lhs)
		Gen(node.rhs)

		fmt.Println("  ldr x9, [sp], #16")
		//fmt.Println("  add sp, sp, #16")
		fmt.Println("  ldr x8, [sp], #16")
		//fmt.Println("  add sp, sp, #16")
		fmt.Println("  str x9, [x8]")
		fmt.Println("  sub sp, sp, #16")
		fmt.Println("  str x9, [sp]")
		return
	case NdRETURN:
		Gen(node.lhs)

		// 最終的な計算結果はスタックに保存されているので取り出す
		fmt.Println("  ldr x8, [sp]")
		// 最終的な計算結果の保存に使用したスタック分spを戻してあげる
		fmt.Println("  add sp, sp, #16")

		fmt.Println("  mov x0, x8")
		// 一番最初にxzrを書き込んだ分のspを戻す
		fmt.Println("  add sp, sp, #16")

		// 終了
		fmt.Println("  ret")
		return
	}

	Gen(node.lhs)
	Gen(node.rhs)

	// RHS(右辺)の結果をx9に読み込む
	// x9にspの領域を読み込む(x9は32bitレジスタなのでspから4byte分読み込むと思う)
	// ldr : load register?
	fmt.Println("  ldr x9, [sp], #16") // ldr x9, [sp]; then sp += 16
	// スタックを進め、左辺の結果を読み込む準備をする
	// 16Nで動かさなければならない
	//fmt.Println("  add sp, sp, #16")
	// LHS(左辺)の結果をx8に読み込む
	// x8にspから始まる32bit=4byte分のデータを読み込む
	fmt.Println("  ldr x8, [sp], #16")
	//fmt.Println("  add sp, sp, #16")

	switch node.kind {
	case NdADD:
		// x8 = x8 +x9
		fmt.Println("  add x8, x8, x9")
	case NdSUB:
		// x8 = x8 -x9
		fmt.Println("  sub x8, x8, x9")
	case NdMUL:
		// x8 = x8 *x9
		fmt.Println("  mul x8, x8, x9")
	case NdDIV:
		// x8 = x8 /x9
		fmt.Println("  sdiv x8, x8, x9")
	case NdEQ:
		// SUB{S} Rd, Rn, OP2: キャリーなしで引きます。
		//   S: オプションのサフィックスです。が指定されている場合S、操作の結果で条件コード フラグが更新されます。
		//  Rd: Dest
		//  Rn: Operand1
		// OP2: Operand2
		// CSINC
		// 32bit(sf==0): CSINC Wd, Wn, Wm, cond
		// 64bit(sf==1): W->X

		// bits(datasize) result;
		//if ConditionHolds(cond) then
		//    result = X[n];
		//else
		//    result = X[m];
		//    result = result + 1;
		//
		//X[d] = result;

		// Conditional Select Increment は、条件が TRUE の場合、
		// デスティネーション レジスタに最初のソース レジスタの値を返し、それ以外の場合は、1 だけインクリメントされた 2 番目のソース レジスタの値を返します。
		//        subs    x8, x8, x9 -> x9 = x8 - x9; xzr = {0|1}?
		//        cset    x8, eq -> eq:z==1, if (xzr==1) x8 = 1 else x8 = 0
		//        and     x8, x8, #0x1 x8 = x8 == 1
		fmt.Println("  subs x8, x8, x9")
		fmt.Println("  cset x8, eq")
		// x8が1(0x01)か否か ==: 1, !=: 0
		// #0x01でも#1でも良いっぽい。
		// もしかしたらalignが関係してくるのかも
		fmt.Println("  and x8, x8, #0x01")
	case NdNE:
		//        subs    x8, x8,x9
		//        cset    x8, ne
		//        and     x8, x8, #0x1
		// https://www.mztn.org/dragon/arm6408cond.html#suffix
		fmt.Println("  subs x8, x8, x9")
		fmt.Println("  cset x8, ne")
		fmt.Println("  and x8, x8, #0x01")
	case NdLT:
		//        subs    x8, x8,x9
		//        cset    x8, lt
		//        and     x8, x8, #0x1
		fmt.Println("  subs x8, x8, x9")
		fmt.Println("  cset x8, lt")
		fmt.Println("  and x8, x8, #0x01")
	case NdLE:
		//        subs    x8, x8,x9
		//        cset    x8, le
		//        and     x8, x8, #0x1
		fmt.Println("  subs x8, x8, x9")
		fmt.Println("  cset x8, le")
		fmt.Println("  and x8, x8, #0x01")
	case NdNUM:
		log.Fatalf("記号ではなく数字トークンを発見しました")
	}

	// 計算結果x8のデータを戻り値として書き込むためspをずらす
	//fmt.Println("  sub sp, sp, #16")
	// spから始まる4byte分のスタック領域に戻り値を書き込む
	fmt.Println("  str x8, [sp, #-16]!")
}
