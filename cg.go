package gocc

import (
	"fmt"
	"runtime"
	"time"
)

var lines []*Line
var spRelativePos int
var localLabelCount int
var stringLabelCount int
var defStrings []*DefString

func init() {
	lines = []*Line{}
	spRelativePos = 0
	localLabelCount = 0
	stringLabelCount = 0
	defStrings = []*DefString{}
}

func addLine(l *Line) {
	lines = append(lines, l)
}

// eg “.text“
func text() {
	addLine(NewSrcLine(".text", ""))
}

// eg “.align 2“
func align(i int) {
	addLine(NewSrcLine(fmt.Sprintf(".align %d", i), ""))
}

// eg “.global _main“
func global(label string) {
	addLine(NewSrcLine(fmt.Sprintf(".global %s", label), ""))
}

// eg “stp x29, x30, [sp]“
func stpS(src1, src2 string) {
	addLine(NewSrcLine(fmt.Sprintf("  stp %s, %s, [sp]", src1, src2), ""))
}

// cf “ldp x29, x30, [sp]
func ldpS(dest1, dest2 string) {
	addLine(NewSrcLine(fmt.Sprintf("  ldp %s, %s, [sp]", dest1, dest2), ""))
}

// eg “mov x0, x8“
func movR(dest, src string) {
	addLine(NewSrcLine(fmt.Sprintf("  mov %s, %s", dest, src), ""))
}

// eg “mov x8, #16“
func movI(dest string, val int) {
	addLine(NewSrcLine(fmt.Sprintf("  mov %s, #%d", dest, val), ""))
}

// cg “add x8, x8, #1“
func addI(dest string, src string, val int) {
	addLine(NewSrcLine(fmt.Sprintf("  add %s, %s, #%d", dest, src, val), ""))
}

// cg “add x8, x8, x9“
func addR(dest string, src string, val string) {
	addLine(NewSrcLine(fmt.Sprintf("  add %s, %s, %s", dest, src, val), ""))
}

// cg “sub x8, x8, #1“
func subI(dest, src string, val int) {
	addLine(NewSrcLine(fmt.Sprintf("  sub %s, %s, #%d", dest, src, val), ""))
}

// cg “sub x8, x8, x9“
func subR(dest, src, val string) {
	addLine(NewSrcLine(fmt.Sprintf("  sub %s, %s, %s", dest, src, val), ""))
}

// cg “subs x8, x8, x9“
func subsR(dest, src, val string) {
	addLine(NewSrcLine(fmt.Sprintf("  subs %s, %s, %s", dest, src, val), ""))
}

// cg “mul x8, x8, x9“
func mulR(dest, src, val string) {
	addLine(NewSrcLine(fmt.Sprintf("  mul %s, %s, %s", dest, src, val), ""))
}

// cg “sdiv x8, x8, x9“
func sdivR(dest, src, val string) {
	addLine(NewSrcLine(fmt.Sprintf("  sdiv %s, %s, %s", dest, src, val), ""))
}

// cg “cset x8, eq“
func cset(dest, op string) {
	addLine(NewSrcLine(fmt.Sprintf("  cset %s, %s", dest, op), ""))
}

// cg “and x8, x8, x9
func andR(dest, val1, val2 string) {
	addLine(NewSrcLine(fmt.Sprintf("  and %s, %s, %s", dest, val1, val2), ""))
}

// cg “and x8, x8, #1
func andI(dest, val1 string, val2 int) {
	addLine(NewSrcLine(fmt.Sprintf("  and %s, %s, #%d", dest, val1, val2), ""))
}

// cg “orr x8, x8, x9“
func orrR(dest, val1, val2 string) {
	addLine(NewSrcLine(fmt.Sprintf("  orr %s, %s, %s", dest, val1, val2), ""))
}

// cg “cmp x8, #1
func cmpI(val1 string, val2 int) {
	addLine(NewSrcLine(fmt.Sprintf("  cmp %s, #%d", val1, val2), ""))
}

// cg “ldr x8, [sp]“
func ldrSt(dest string) {
	addLine(NewSrcLine(fmt.Sprintf("  ldr %s, [sp]", dest), ""))
}

// cg “ldr x8, [x8]“
func ldrArs(dest, address string) {
	addLine(NewSrcLine(fmt.Sprintf("  ldr %s, [%s]", dest, address), ""))
}

func ldrLabel(dest, label string) {
	addLine(NewSrcLine(fmt.Sprintf("  ldr %s, %s", dest, label), ""))
}

// cg “str x8, [sp]“
func strSt(src string) {
	addLine(NewSrcLine(fmt.Sprintf("  str %s, [sp]", src), ""))
}

// cg “str x9, [x8]“
func strArs(src, address string) {
	addLine(NewSrcLine(fmt.Sprintf("  str %s, [%s]", src, address), ""))
}

// cg “b _foo“
func b(label string) {
	addLine(NewSrcLine(fmt.Sprintf("  b %s", label), ""))
}

// cg “b.eq _foo
func bop(op, label string) {
	addLine(NewSrcLine(fmt.Sprintf("  b.%s %s", op, label), ""))
}

// cg “bl _foo“
func bl(label string) {
	addLine(NewSrcLine(fmt.Sprintf("  bl %s", label), ""))
}

func blr(reg string) {
	addLine(NewSrcLine(fmt.Sprintf("  blr %s", reg), ""))
}

// cg “ret“
func ret() {
	addLine(NewSrcLine("  ret", ""))
}

func adrReg(dest, label string) {
	addLine(NewSrcLine(fmt.Sprintf("  adr %s, %s", dest, label), ""))
}

// cg “add sp, sp, #16“
func addSp(val int) {
	addLine(NewSrcLine(fmt.Sprintf("  add sp, sp, #%d", val), fmt.Sprintf("%d -> %d", spRelativePos, spRelativePos+val)))
	spRelativePos += val
}

// cg “sub sp, sp, #16“
func subSp(val int) {
	addLine(NewSrcLine(fmt.Sprintf("  sub sp, sp, #%d", val), fmt.Sprintf("%d -> %d", spRelativePos, spRelativePos-val)))
	spRelativePos -= val
}

func adrpPage(reg string, label string) {
	//	adrp	x0, l_.str@PAGE
	addLine(NewSrcLine(fmt.Sprintf("  adrp %s, %s@PAGE", reg, label), ""))
}

func addPageOFF(dest, src, label string) {
	//add	x0, x0, l_.str@PAGEOFF
	addLine(NewSrcLine(fmt.Sprintf("  add %s, %s, %s@PAGEOFF", dest, src, label), ""))
}

func comment(text string) {
	addLine(NewComment(text))
}

func makeFuncLabel(l string) string {
	if runtime.GOOS == "darwin" {
		l = "_" + l
	}
	return l
}

func makeLocalLabel() string {
	l := fmt.Sprintf("L%d", localLabelCount)
	localLabelCount++
	return l
}

func makeStringLabel() string {
	var l string
	if stringLabelCount == 0 {
		l = fmt.Sprintf("l_.str")
	} else {
		l = fmt.Sprintf("l_.str.%d", stringLabelCount)
	}
	stringLabelCount++
	return l
}

func defLabel(l string) {
	addLine(NewSrcLine(l+":", ""))
}

// genLocalVariable
// 取得した変数のアドレスをsタックに積む、読み込むときはldrArs
func genLocalVariable(node *Node) {
	comment("変数のアドレスを取得する")
	movR("x8", "x29")
	subI("x8", "x8", node.offset)
	subSp(16)
	strSt("x8")
}

func gen(node *Node) {
	switch node.kind {
	case NdNUM:
		comment(node.kind.String())

		subSp(16)
		movI("x8", node.val)
		strSt("x8")
		return
	case NdSTRING:
		comment(node.kind.String())

		// 文字列につけるラベルを生成
		l := makeStringLabel()
		// コードジェネレータが最後に一覧を書き出せるように保存
		defStrings = append(defStrings, NewDefString(l, node.data))

		adrpPage("x0", l)
		addPageOFF("x0", "x0", l)
		subSp(16)
		strSt("x0")
		return
	case NdLVAR:
		comment(node.kind.String())
		genLocalVariable(node)

		comment("変数として与えられたアドレスから値を取り出す")
		ldrSt("x8")
		addSp(16)
		ldrArs("x8", "x8")
		subSp(16)
		strSt("x8")
		return
	case NdASSIGN:
		comment(node.kind.String())
		genLocalVariable(node.lhs)
		gen(node.rhs)

		comment("x9, x8の順にスタックからpopしてx8のデータをアドレスとしてx9のデータを割り当てる")
		ldrSt("x9")
		addSp(16)
		ldrSt("x8")
		addSp(16)
		strArs("x9", "x8")
		//subSp(16)
		//strSt("x9")
		return
	case NdRETURN:
		comment(node.kind.String())
		gen(node.lhs)

		comment("計算結果をpopしてx0に代入し、戻り値とする")
		ldrSt("x8")
		addSp(16)

		if spRelativePos < -32 {
			comment("[sp fixing]")
			comment(fmt.Sprintf("spRelativePos: %d", spRelativePos))
			// cg. now:-48, -48 + 32 = -16
			// add sp 16
			fixWidth := spRelativePos + 32
			addSp(-1 * fixWidth)
		}

		// epilogue
		movR("x0", "x8")
		ldpS("x29", "x30")
		addSp(32)
		ret()
		return
	case NdIF:
		comment(node.kind.String())

		// condの結果がfalseだった時、trueだった場合に実行されるコードを読み飛ばすために使用される
		falseLabel := makeLocalLabel()

		// 条件式の生成
		gen(node.cond)
		ldrSt("x8")
		addSp(16)
		cmpI("x8", 0)
		// cond == false then goto false case
		bop("eq", falseLabel)

		// true case
		gen(node.lhs)

		// false case
		defLabel(falseLabel)
		return
	case NdIFELSE:
		comment(node.kind.String())

		// 条件式の結果がfalseだった場合にtrueだった場合のコードを読み飛ばし、elseブロックに飛ぶために使用される
		elseLabel := makeLocalLabel()
		// 条件式の結果がtrueだった場合にtrueだった場合のコードを実行した後、elseブロックを読み飛ばすために使用される
		endLabel := makeLocalLabel()

		// 条件式の生成
		gen(node.cond)
		ldrSt("x8")
		addSp(16)
		cmpI("x8", 0)
		// cond == false then goto false case
		bop("eq", elseLabel)

		// true case
		gen(node.lhs)
		// elseブロックを読み飛ばす
		b(endLabel)

		// false case
		defLabel(elseLabel)
		gen(node.rhs)

		// end
		defLabel(endLabel)
		return
	case NdWHILE:
		comment(node.kind.String())

		beginLabel := makeLocalLabel()
		endLabel := makeLocalLabel()

		// ループブロック終了したら戻ってくる場所
		defLabel(beginLabel)
		// 条件式を生成
		gen(node.cond)
		ldrSt("x8")
		addSp(16)
		cmpI("x8", 0)
		// 条件式がfalseならループ外にとぶ
		bop("eq", endLabel)

		// true case
		gen(node.lhs)
		// ループ先頭へ
		b(beginLabel)

		defLabel(endLabel)
		return
	case NdFOR:
		comment(node.kind.String())

		beginLabel := makeLocalLabel()
		endLabel := makeLocalLabel()

		// 初期化式を生成
		if node.init != nil {
			comment("for.init")
			gen(node.init)
		}

		// ループブロック終了後に戻ってくる場所
		defLabel(beginLabel)

		// 条件式を生成
		if node.cond != nil {
			comment("for.cond")
			gen(node.cond)

			ldrSt("x8")
			addSp(16)
			cmpI("x8", 0)
			// 条件式の結果がfalseならループ外に飛ぶ
			bop("eq", endLabel)
		}

		// true case
		gen(node.lhs)

		// ループの最後に実行されるコードを生成
		if node.loop != nil {
			comment("for.loop")
			gen(node.loop)
		}

		// ループの最初に戻る
		b(beginLabel)

		// 条件式がfalseだった場合にループ外に出るために使用
		defLabel(endLabel)
		return
	case NdBLOCK:
		comment(node.kind.String())

		// 使用されなかった関数の戻り値をなかったことにするため、ブロックの最初と最後のspを合わせる処理
		// 今後エラー出たら多分ここが原因
		spBefore := spRelativePos
		for i, n := range node.code {
			comment(fmt.Sprintf("block.%d", i))
			gen(n)
		}
		spAfter := spRelativePos
		if spAfter < 0 && spBefore != spAfter {
			n16 := -1 * (spAfter - spBefore)
			comment(fmt.Sprintf("fix %d", n16))
			addSp(n16)
		}
		return
	case NdCALL:
		comment(node.kind.String())

		//adrReg("x8", node.label)
		//adrp x0, symbol@PAGE
		//ldr x1, [x0, symbol@PAGEOFF]

		//write(NewSrcLine(fmt.Sprintf("  adrp x0, %s@PAGE", node.label), ""))
		////write(NewSrcLine(fmt.Sprintf("  ldr x8, [x0, %s@PAGEOFF", node.label), ""))
		//write(NewSrcLine(fmt.Sprintf("  add x0, x0, %s@PAGEOFF", node.label), ""))
		//blr("x0")

		// arguments
		for i, arg := range node.arguments {
			gen(arg)
			ldrSt(fmt.Sprintf("x%d", i))
			addSp(16)
		}

		bl(node.label)

		subSp(16)
		strSt("x0")
		//movR("x8", "x0")

		// スタックトップに結果を自動的に保存するため、これが使われなかった場合、spだけがずれていく
		// 型を導入して、使用されなかったらvoid、、、？
		return
	case NdFUNCTION:
		// label
		defLabel(makeFuncLabel(node.data))

		// save data for restore
		subSp(32)
		stpS("x29", "x30")

		// parameters
		// ident(p1, p2, p3)
		// parameters are set on x0-x7 regs
		// x0 = p1
		// x1 = p2
		// x3 = p3
		// we need to convert received parameters to local variables
		for i, param := range node.arguments {
			// param = p1, p2, p3, p4, ...
			// 変数のアドレスを取得しスタックに積む
			genLocalVariable(param)

			subSp(16)
			// registerからデータをとってきてスタックに積む
			strSt(fmt.Sprintf("x%d", i))

			// 中身
			ldrSt("x9")
			addSp(16)
			// 変数のアドレス
			ldrSt("x8")
			addSp(16)
			// 代入
			strArs("x9", "x8")
		}

		gen(node.lhs)
		return
	}

	gen(node.lhs)
	gen(node.rhs)

	ldrSt("x9")
	addSp(16)

	ldrSt("x8")
	addSp(16)

	switch node.kind {
	case NdADD:
		addR("x8", "x8", "x9")
	case NdSUB:
		subR("x8", "x8", "x9")
	case NdMUL:
		mulR("x8", "x8", "x9")
	case NdDIV:
		sdivR("x8", "x8", "x9")
	case NdEQ:
		subsR("x8", "x8", "x9")
		cset("x8", "eq")
		andI("x8", "x8", 1)
	case NdNE:
		subsR("x8", "x8", "x9")
		cset("x8", "ne")
		andI("x8", "x8", 1)
	case NdLT:
		subsR("x8", "x8", "x9")
		cset("x8", "lt")
		andI("x8", "x8", 1)
	case NdLE:
		subsR("x8", "x8", "x9")
		cset("x8", "le")
		andI("x8", "x8", 1)
	case NdAND:
		andR("x8", "x8", "x9")
	case NdOR:
		orrR("x8", "x8", "x9")
	}

	subSp(16)
	strSt("x8")
}

func Generate(nodes []*Node) string {
	// prologue
	addLine(NewSrcLine("    .section  __TEXT,__text,regular,pure_instructions", ""))
	comment(fmt.Sprintf("compiled at %s", time.Now().String()))
	text()
	align(2)

	//mainLabel := makeFuncLabel("main")
	//global(mainLabel)
	//defLabel(mainLabel)

	global(makeFuncLabel("main"))

	//subSp(32)
	//stpS("x29", "x30")

	// generate assemblies from nodes
	var src string
	for _, n := range nodes {
		if n != nil {
			gen(n)
		}
	}

	addLine(NewSrcLine("    .section    __TEXT,__cstring,cstring_literals", ""))
	// define string bottom
	for _, l := range defStrings {
		//src += fmt.Sprintf("%s:\n", l.label)
		//src += fmt.Sprintf(".asciz %s\n", l.data)
		defLabel(l.label)
		addLine(NewSrcLine(fmt.Sprintf("  .asciz \"%s\"", l.data), ""))
	}

	addLine(NewSrcLine(".subsections_via_symbols", ""))

	// print source code lines
	for _, l := range lines {
		src += l.String() + "\n"
	}

	// print warn top
	if spRelativePos != 0 {
		fmt.Printf("; warn: stack pointer is %d at the end.\n", spRelativePos)
	}

	return src
}
