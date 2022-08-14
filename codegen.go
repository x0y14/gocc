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
		//        subs    w8, w8, w9 -> w9 = w8 - w9; wzr = {0|1}?
		//        cset    w8, eq -> eq:z==1, if (wzr==1) w8 = 1 else w8 = 0
		//        and     w8, w8, #0x1 w8 = w8 == 1
		fmt.Println("  subs w8, w8, w9")
		fmt.Println("  cset w8, eq")
		// w8が1(0x01)か否か ==: 1, !=: 0
		// #0x01でも#1でも良いっぽい。
		// もしかしたらalignが関係してくるのかも
		fmt.Println("  and w8, w8, #0x01")
	case NdNE:
		//        subs    w8, w8, w9
		//        cset    w8, ne
		//        and     w8, w8, #0x1
		// https://www.mztn.org/dragon/arm6408cond.html#suffix
		fmt.Println("subs w8, w8, w9")
		fmt.Println("  cset w8, ne")
		fmt.Println("  and w8, w8, #0x01")
	case NdLT:
		//        subs    w8, w8, w9
		//        cset    w8, lt
		//        and     w8, w8, #0x1
		fmt.Println("  subs w8, w8, w9")
		fmt.Println("  cset w8, lt")
		fmt.Println("  and w8, w8, #0x01")
	case NdLE:
		//        subs    w8, w8, w9
		//        cset    w8, le
		//        and     w8, w8, #0x1
		fmt.Println("  subs w8, w8, w9")
		fmt.Println("  cset w8, le")
		fmt.Println("  and w8, w8, #0x01")
	case NdNUM:
		log.Fatalf("記号ではなく数字トークンを発見しました")
	}

	// 計算結果w8のデータを戻り値として書き込むためspをずらす
	fmt.Println("  sub sp, sp, #16")
	// spから始まる4byte分のスタック領域に戻り値を書き込む
	fmt.Println("  str w8, [sp]")
}
