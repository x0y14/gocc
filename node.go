package gocc

type Node struct {
	kind   NodeKind
	init   *Node // for (A;B;C)のA
	cond   *Node // for (A;B;C)、while (B)
	loop   *Node // for (A;B;C;)のC
	lhs    *Node // 左辺
	rhs    *Node // 右辺
	val    int   // NdNUMの場合、数値が入る
	offset int   // NdLVALの場合、FramePointerからどれだけ離れた位置のスタックにデータが格納されているかが入る
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

func NewNodeNum(val int) *Node {
	return &Node{
		kind: NdNUM,
		val:  val,
	}
}

func NewNodeLVar(offset int) *Node {
	return &Node{
		kind:   NdLVAR,
		offset: offset,
	}
}
