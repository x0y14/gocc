package gocc

type NodeKind int

const (
	NdADD = iota // +
	NdSUB        // -
	NdMUL        // *
	NdDIV        // /

	NdEQ // ==
	NdNE //!=
	NdLT // <
	NdLE // <=

	NdASSIGN // =

	NdNUM  // 数字
	NdLVAR // Local VAR

	NdRETURN
	NdIF     // 条件, trueの場合の式...?
	NdIFELSE // 条件, trueの場合, falseの場合

)
