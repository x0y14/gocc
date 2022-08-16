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
	NdIF     //
	NdIFELSE //
	NdWHILE
	NdFOR
)
