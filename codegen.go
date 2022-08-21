package gocc

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

var labelCounter int
var nest int
var Assembly []*Line

const separator = " - "

func init() {
	labelCounter = 0
	nest = 0
	Assembly = []*Line{}
}

func write(l *Line) {
	Assembly = append(Assembly, l)
}

func label() string {
	l := fmt.Sprintf("__lable_%d", labelCounter)
	labelCounter++
	return l
}

func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}

func genLVal(node *Node) {
	// ブロック表示
	mark := randomString(10)
	nest++
	write(NewSeparator(mark, separator, nest, false))
	defer func(s string) {
		write(NewSeparator(mark, separator, nest, true))
		nest--
	}(mark)
	write(NewComment("generated at genLVal()"))

	if node.kind != NdLVAR {
		log.Fatalf("代入の左辺値が変数ではありません")
	}

	write(NewComment("変数のアドレスを取得する"))

	// x8にフレームポインタを保存
	write(NewSrcLine("  mov x8, x29", "x8 = x29"))
	// １文字変数の位置分スタックを移動
	write(NewSrcLine(fmt.Sprintf("  sub x8, x8, #%d", node.offset), "x8-=offset"))
	// スタックに変数のデータはどこから始まるのかを保存
	write(NewSrcLine("  str x8, [sp, #-16]!", "sp-=16, [sp] = x8"))
}

func Gen(node *Node) {
	mark := randomString(10)
	nest++
	write(NewSeparator(mark, separator, nest, false))
	defer func(s string) {
		write(NewSeparator(mark, separator, nest, true))
		nest--
	}(mark)
	write(NewComment("generated at Gen()"))

	switch node.kind {
	case NdNUM:
		write(NewComment(node.kind.String()))
		write(NewComment("即値を一度レジスタに格納してからスタックにpushする"))
		write(NewSrcLine("  sub sp, sp, #16", "sp-=16"))
		write(NewSrcLine(fmt.Sprintf("  mov x8, #%d", node.val), "x8 = val"))
		// x8をspから始まるスタック領域に書き込む
		write(NewSrcLine("  str x8, [sp]", "[sp] = x8"))
		return

	case NdLVAR:
		write(NewComment(node.kind.String()))
		genLVal(node)

		write(NewComment("変数として場所を与えられたアドレスから値を取り出す"))
		// 変数のデータはスタックのどこに書き込まれるのか取得
		write(NewSrcLine("  ldr x8, [sp], #16", "x8 = [sp], sp+=16"))
		// どこに書き込まれるのかというのをアドレスとして読み込みx8に書き込まれた内容を読み込む
		write(NewSrcLine("  ldr x8, [x8]", "x8 = [x8]"))
		// スタックに読み込んだ内容を書き込む
		write(NewSrcLine("  str x8, [sp, #-16]!", "sp-=16, [sp] = x8"))
		return

	case NdASSIGN:
		write(NewComment(node.kind.String()))
		genLVal(node.lhs)
		Gen(node.rhs)

		write(NewComment("x9, x8の順にスタックからpopしてx8のデータをアドレスとしてx9のデータを割り当てる"))
		// スタックからデータを読み込みx9へ。変数に代入するデータ。
		write(NewSrcLine("  ldr x9, [sp], #16", "x9 = [sp], sp+=16"))
		// スタックからデータを読み込みx8へ。変数のオフセット。
		write(NewSrcLine("  ldr x8, [sp], #16", "x8 = [sp], sp+=16"))
		// 変数アドレスにデータを書き込む
		write(NewSrcLine("  str x9, [x8]", "[x8] = x9"))
		// スタックトップに変数の内容を保存
		write(NewSrcLine("  str x9, [sp, #-16]!", "sp-=16, [sp] = x9"))
		return

	case NdRETURN:
		write(NewComment(node.kind.String()))
		Gen(node.lhs)

		write(NewComment("計算結果をpopしてx0に代入し、戻り値とする"))
		// 最終的な計算結果はスタックに保存されているので取り出す
		write(NewSrcLine("  ldr x8, [sp], #16", "x8 = [sp], #16"))
		write(NewSrcLine("  mov x0, x8", "x0 = x8"))
		write(NewSrcLine("  ldp x29, x30, [sp], #32; epilogue", ""))
		// 終了
		write(NewSrcLine("  ret", ""))
		return

	case NdIF:
		write(NewComment(node.kind.String()))

		// 疑似コード "IF (A) B"
		//  if (A == 0)
		//    goto end;
		//  B;
		//  end:

		// gotoに使用するラベルを生成
		falseLabel := label()

		// 条件式を生成
		Gen(node.cond)

		// 条件式の結果を取り出す
		write(NewSrcLine("  ldr x8, [sp], #16", "x8 = [sp], sp+=16"))
		// 結果が0と一致するか(条件式の結果がfalseかどうか)、一致すればZに1が自動的に設定される
		write(NewSrcLine("  cmp x8, 0", "if x8 == 0 then Z = 1"))
		// もしZが1なら、trueだった場合に実行するコードを読み飛ばす
		write(NewSrcLine(fmt.Sprintf("  b.eq %s", falseLabel), "if Z == 1 then goto false_case"))

		// trueだった場合のコードを生成
		write(NewComment("true_case"))
		Gen(node.lhs)

		// falseだった場合のジャンプ先
		write(NewComment("false_case"))
		write(NewSrcLine(fmt.Sprintf("%s:", falseLabel), "false_label"))
		return

	case NdIFELSE:
		write(NewComment(node.kind.String()))

		// 疑似コード "IF (A) B ELSE C"
		//  if (A == 0)
		//    goto els;
		//  B;
		//  goto end;
		//els:
		//  C;
		//end:

		// Aがfalseの場合、Cを実行するためジャンプするのに使用
		elseLabel := label()
		// Aがtrueだった場合、Bを実行した後終了地点に移動するために使用
		endLabel := label()

		// 条件式を生成
		Gen(node.cond)

		// 条件式の結果を取り出す
		write(NewSrcLine("  ldr x8, [sp], #16", "x8 = [sp], sp+=16"))
		// 結果が0と一致するか(条件式の結果がfalseかどうか)、一致すればZに1が自動的に設定される
		write(NewSrcLine("  cmp x8, #0", "if x8 == 0 then Z = 1"))
		// もしZが1なら、trueだった場合に実行するコードを読み飛ばし、elseの場合に実行されるコードの場所へ移動する
		write(NewSrcLine(fmt.Sprintf("  b.eq %s", elseLabel), "if Z == 1 then goto else_label"))

		// 条件式の結果がtrueになった場合実行されるコード
		write(NewComment("true_case"))
		Gen(node.lhs)
		// elseの場合に実行されるコードを読み飛ばす
		write(NewSrcLine(fmt.Sprintf("  b %s", endLabel), "goto end_label"))

		// 条件式の結果がfalseになった場合実行されるコード
		write(NewComment("else_case"))
		write(NewSrcLine(fmt.Sprintf("%s:", elseLabel), "else_label"))
		Gen(node.rhs)

		// trueだった場合、falseの場合のコードを読み飛ばすのに使用される
		write(NewSrcLine(fmt.Sprintf("%s:", endLabel), "end_label"))
		return

	case NdWHILE:
		write(NewComment(node.kind.String()))

		// 疑似コード "while (A) B"
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

		// ループブロック終了時に戻ってくるために使用する
		write(NewSrcLine(fmt.Sprintf("%s:", beginLabel), "begin_label"))

		// 条件式を生成
		Gen(node.cond)
		// 条件式の結果を取り出す
		write(NewSrcLine("  ldr x8, [sp], #16", "x8 = [sp], sp+=16"))
		// 条件式の答えがfalseであれば、Zに1が自動的に設定される
		write(NewSrcLine("  cmp x8, #0", "if x8 == 0 then Z = 1"))
		// 条件式の結果がfalseならループを抜ける必要がある
		write(NewSrcLine(fmt.Sprintf("  b.eq %s", endLabel), "if Z == 1 then goto end_label"))

		// 条件式の結果がtrueだった場合に実行されるコード
		write(NewComment("true_case"))
		Gen(node.lhs)
		// ループの先頭に戻る
		write(NewSrcLine(fmt.Sprintf("  b %s", beginLabel), "goto begin_label"))

		// 条件式がfalseだった場合にtrueだった場合に実行されるコードを読み飛ばすために使用される
		write(NewSrcLine(fmt.Sprintf("%s:", endLabel), "end_label"))
		return

	case NdFOR:
		write(NewComment(node.kind.String()))

		// 疑似コード "for (A;B;C) D"
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

		beginLabel := label()
		endLabel := label()

		// 初期化式を生成
		if node.init != nil {
			write(NewComment("for.init"))
			// 初期化式
			Gen(node.init)
		}

		// ループ内容終了後に一番再度ループを実行するために使用される
		write(NewSrcLine(fmt.Sprintf("%s:", beginLabel), "begin_label"))

		// 条件式を生成
		if node.cond != nil {
			write(NewComment("for.cond"))
			// 条件式
			Gen(node.cond)
			// 条件式の結果を取得
			write(NewSrcLine("  ldr x8, [sp], #16", "x8 = [sp], sp+=16"))
			// 条件式の結果がfalseであれば自動的にZに1が設定される
			write(NewSrcLine("  cmp x8, 0", "if x8 == 0 then Z = 1"))
			// もし条件式がfalseだったら、ループ外に移動する
			write(NewSrcLine(fmt.Sprintf("  b.eq %s", endLabel), "if Z == 1 then goto end_label"))
		}

		// 条件式がtrueだった場合実行されるコード
		write(NewComment("true_case"))
		Gen(node.lhs)

		// ループの最後に実行されるコードを生成
		if node.loop != nil {
			write(NewComment("for.loop"))
			Gen(node.loop)
		}

		// ループの最初に戻る
		write(NewSrcLine(fmt.Sprintf("  b %s", beginLabel), "goto begin_label"))

		// ループの条件式の結果がfalseだった場合、ループを離脱するために使用
		write(NewSrcLine(fmt.Sprintf("%s:", endLabel), "end_label"))
		return

	case NdBLOCK:
		write(NewComment(node.kind.String()))

		// {
		//   stmt
		//   stmt
		// }

		for i, c := range node.code {
			write(NewComment(fmt.Sprintf("(%d) the child of %s", i, mark)))
			Gen(c)
		}
		return

	case NdCALL:
		write(NewComment(node.kind.String()))

		write(NewComment("関数を呼び出す"))
		// 関数ラベルにジャンプする
		write(NewSrcLine(fmt.Sprintf("  bl %s", node.label), "goto function"))

		// w0:x0に結果が入っているので、取り出してスタックに保存
		// なぜか結果を取得することができない現象が発生しているのでチェックが必要
		write(NewSrcLine("  sub sp, sp, #16", "sp-=16"))
		write(NewSrcLine("  str x0, [sp]", "[sp] = x0"))
		write(NewSrcLine("  mov x0, xzr", "x0 = xzr"))
		return
	}

	// 左辺を生成
	write(NewComment("node.lhs"))
	Gen(node.lhs)

	// 右辺を生成
	write(NewComment("node.rhs"))
	Gen(node.rhs)

	// 右辺の結果を読み込む
	write(NewSrcLine("  ldr x9, [sp], #16", "x9 = [sp], sp+=16"))

	// 左辺を読み込む
	write(NewSrcLine("  ldr x8, [sp], #16", "x8 = [sp], sp+=16"))

	switch node.kind {
	case NdADD:
		write(NewSrcLine("  add x8, x8, x9", "x8 = x8 + x9"))

	case NdSUB:
		write(NewSrcLine("  sub x8, x8, x9", "x8 = x8 - x9"))

	case NdMUL:
		write(NewSrcLine("  mul x8, x8, x9", "x8 = x8 * x9"))

	case NdDIV:
		write(NewSrcLine("  sdiv x8, x8, x9", "x8 = x8 / x9"))

	case NdEQ:
		// x8 == x9
		// 結果によってZを変更する引き算を使用してx8 - x9がゼロになるか（等しいか）を計算、等しければZは1に、そうでなければZは0になる
		write(NewSrcLine("  subs x8, x8, x9", "x8 = x8 - x9, if x8 == 0 then Z = 1 else Z = 0"))
		// x8にZ == 1の結果を代入、trueであれば1、falseであれば0
		write(NewSrcLine("  cset x8, eq", "if Z == 1 then x8 = 1 else x8 = 0"))
		// x8が1であることを論理演算andを使用して確認
		write(NewSrcLine("  and x8, x8, 1", "x8 = x8 && 1"))

	case NdNE:
		// x8 != x9
		// 結果によってZを変更する引き算を使用してx8 - x9がゼロになるか（等しいか）を計算、Zフラグを更新、等しければZは1に、そうでなければZは0になる
		write(NewSrcLine("  subs x8, x8, x9", "x8 = x8 - x9, if x8 == 0 then Z = 1 else Z = 0"))
		// x8にZ != 1の結果を代入、trueであれば1、falseであれば0
		write(NewSrcLine("  cset x8, ne", "if Z == 0 then x8 = 1 else x8 = 0"))
		// x8が1であることを論理演算andを使用して確認
		write(NewSrcLine("  and x8, x8, 1", "x8 = x8 && 1"))

	case NdLT:
		// x8 < x9
		// 結果によってN、Vを変更する引き算こ使用して、x8 - x9が結果がオーバーフローするか、負になるかのフラグを更新
		write(NewSrcLine("  subs x8, x8, x9", "x8 = x8 - x9, update N, V"))
		// N != V : 結果が負でなく、かつ、オーバーフローしていなければtrue、そうでなければfalseをx8に代入
		write(NewSrcLine("  cset x8, lt", "if N != V then x8 = 1 else x8 = 0"))
		// x8が1であることを論理演算andを使用して確認
		write(NewSrcLine("  and x8, x8, 1", "x8 = x8 && 1"))

	case NdLE:
		// x8 <= x9
		// 結果によってN、Vを変更する引き算こ使用して、x8 - x9が結果がオーバーフローするか、負になるかあるいは等しいか、フラグを更新
		// EQ || LT
		write(NewSrcLine("  subs x8, x8, x9", "x8 = x8 - x9, update N, V"))
		// x8 = (N != V) || (Z==1)
		write(NewSrcLine("  cset x8, le", "if (N != V) || (Z ==1) then x8 = 1 else x8 = 0"))
		// x8が1であることを論理演算andを使用して確認
		write(NewSrcLine("  and x8, x8, 1", "x8 = x8 && 1"))

	case NdAND:
		// x8 && x9
		// x8に左辺の結果、x9に右辺の結果が入っている
		// x8 = x8 && x9
		write(NewSrcLine("  and x8, x8, x9", "x8 = x8 && x9"))

	case NdOR:
		// x8 || x9
		// x8に左辺の結果、x9に右辺の結果が入っている
		// x8 = x8 || x9
		write(NewSrcLine("  orr x8, x8, x9", "x8 = x8 || x9"))

	case NdNUM:
		log.Fatalf("記号ではなく数字トークンを発見しました")
	}

	// 計算結果をスタックに書き込む
	write(NewSrcLine("  str x8, [sp, #-16]!", "sp-=16, [sp] = x8"))
}
