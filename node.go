package gocc

type Node struct {
	kind   NodeKind
	lhs    *Node
	rhs    *Node
	val    int
	offset int
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
