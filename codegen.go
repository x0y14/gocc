package gocc

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

var labelCounter int

func init() {
	labelCounter = 0
}

func label() string {
	l := fmt.Sprintf("__lable_%d", labelCounter)
	labelCounter++
	return l
}

func comment(c string) string {
	return fmt.Sprintf("/*%s*/", c)
}

func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}

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
	case NdIF:
		mark := randomString(5)
		fmt.Println(comment(fmt.Sprintf("start #%s", mark)))
		// gotoに使用するラベルを生成
		falseLabel := label()
		// 条件式を生成
		Gen(node.cond)
		// 条件式の結果を取り出す
		// x8 = [sp]; sp+=16
		fmt.Println("  ldr x8, [sp], #16")
		// is 0(false) ?
		fmt.Println("  cmp x8, 0")
		fmt.Printf("  b.eq %s\n", falseLabel)
		// trueだった場合のコードを生成
		Gen(node.lhs)
		// falseだった場合のジャンプ先
		fmt.Printf("%s:\n", falseLabel)
		fmt.Println(comment(fmt.Sprintf("end #%s", mark)))
		return
	case NdIFELSE:
		mark := randomString(5)
		fmt.Println(comment(fmt.Sprintf("start #%s", mark)))
		// * 上から順に実行される
		//  if (A == 0)
		//    goto els;
		//  B;
		//  goto end;
		//els:
		//  C;
		//end:

		// if (A) B else C
		// Aがfalseの場合、Cを実行するためジャンプするのに使用
		elseLabel := label()
		// Aがtrueだった場合、Bを実行した後終了地点に移動するために使用
		endLabel := label()
		// 条件式を生成
		Gen(node.cond)
		// 条件式の結果をスタックから取り出す
		// x8 = [sp]; sp+=16
		fmt.Println("  ldr x8, [sp], #16")
		// (x8 == 0)
		fmt.Println("  cmp x8, #0")
		// true: goto else C
		fmt.Printf("  b.eq %s\n", elseLabel)
		// false: B goto end
		Gen(node.lhs)
		fmt.Printf("  b %s\n", endLabel)
		fmt.Printf("%s:\n", elseLabel)
		Gen(node.rhs)
		fmt.Printf("%s:\n", endLabel)
		fmt.Println(comment(fmt.Sprintf("end #%s", mark)))
		return
	case NdWHILE:
		mark := randomString(5)
		fmt.Println(comment(fmt.Sprintf("start #%s", mark)))
		// while (A) B
		// begin:
		//   A
		//   pop x8
		//   cmp x8, 0
		//   b.eq end
		//   B
		//   b begin
		// end:
		beginLabel := label()
		endLabel := label()
		// begin:
		fmt.Printf("%s:\n", beginLabel)
		// A
		Gen(node.cond)
		// Aの結果を取り出す
		// x8 = [sp]; sp += 16
		fmt.Println("ldr x8, [sp], #16")
		// (x8 == 0)
		fmt.Println("  cmp x8, #0")
		// true(1), A == false, goto end
		fmt.Printf("  b.eq %s\n", endLabel)
		// B
		Gen(node.lhs)
		// goto begin, ture-loop
		fmt.Printf("  b %s\n", beginLabel)
		// end:
		fmt.Printf("%s:\n", endLabel)
		fmt.Println(comment(fmt.Sprintf("end #%s", mark)))
		return
	case NdFOR:
		// for (A;B;C) D
		//   A
		// begin:
		//   B
		//   pop x8
		//   cmp x8, 0
		//   b.eq end
		//   D
		//   C
		//   b begin
		// end:
		mark := randomString(5)
		beginLabel := label()
		endLabel := label()
		fmt.Println(comment(fmt.Sprintf("start #%s", mark)))
		// A
		Gen(node.init)
		// begin:
		fmt.Printf("%s:\n", beginLabel)
		// B
		Gen(node.cond)
		// pop x8
		// x8 = [sp]; sp+=16
		fmt.Println("  ldr x8, [sp], #16")
		// cmp x8, 0
		fmt.Println("  cmp x8, 0")
		// b.eq end
		fmt.Printf("  b.eq %s\n", endLabel)
		// D
		Gen(node.lhs)
		// C
		Gen(node.loop)
		// b begin
		fmt.Printf("  b %s\n", beginLabel)
		// end:
		fmt.Printf("%s:\n", endLabel)
		fmt.Println(comment(fmt.Sprintf("end #%s", mark)))
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
