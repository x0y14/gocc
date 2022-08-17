package gocc

type NodeKind int

const (
	NdADD = iota // +
	NdSUB        // -
	NdMUL        // *
	NdDIV        // /

	NdEQ  // ==
	NdNE  //!=
	NdLT  // <
	NdLE  // <=
	NdAND // &&
	NdOR  // ||

	NdASSIGN // =

	NdNUM  // 数字
	NdLVAR // Local VAR
	NdCALL

	NdRETURN
	NdIF     //
	NdIFELSE //
	NdWHILE
	NdFOR

	NdBLOCK // "{" stmt* "}"
)
