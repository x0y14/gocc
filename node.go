package gocc

type Node struct {
	kind      NodeKind
	init      *Node  // for (A;B;C)のA
	cond      *Node  // for (A;B;C)、while (B)
	loop      *Node  // for (A;B;C;)のC
	lhs       *Node  // 左辺
	rhs       *Node  // 右辺
	val       int    // NdNUMの場合、数値が入る
	offset    int    // NdLVALの場合、FramePointerからどれだけ離れた位置のスタックにデータが格納されているかが入る
	data      string // NdSTRINGの場合、文字列が入る
	label     string // NdCALLの場合、関数の名前（ラベル）が入る
	code      []*Node
	arguments []*Node
}

func NewNode(kind NodeKind, lhs *Node, rhs *Node) *Node {
	return &Node{
		kind:   kind,
		lhs:    lhs,
		rhs:    rhs,
		val:    0,
		offset: 0,
	}
}

func NewNodeWithExpr(Kind NodeKind, init, cond, loop, lhs, rhs *Node) *Node {
	return &Node{
		kind:   Kind,
		init:   init,
		cond:   cond,
		loop:   loop,
		lhs:    lhs,
		rhs:    rhs,
		val:    0,
		offset: 0,
	}
}

func NewNodeBlock(code []*Node) *Node {
	return &Node{
		kind: NdBLOCK,
		code: code}
}

func NewNodeNum(val int) *Node {
	return &Node{
		kind: NdNUM,
		val:  val,
	}
}

func NewNodeString(data []rune) *Node {
	return &Node{
		kind: NdSTRING,
		data: string(data),
	}
}

func NewNodeLVar(offset int) *Node {
	return &Node{
		kind:   NdLVAR,
		offset: offset,
	}
}

func NewNodeCALL(functionName []rune, args []*Node) *Node {
	// 数のラベルはprefixとして"_"がつく。
	// なお外部から呼ばれる場合は .global ${label} する必要がある
	return &Node{
		kind:      NdCALL,
		label:     "_" + string(functionName),
		arguments: args,
	}
}
