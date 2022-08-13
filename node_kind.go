package gocc

type NodeKind int

const (
	NdADD = iota
	NdSUB
	NdMUL
	NdDIV
	NdNUM
	NdEQ
	NdNEQ
	NdLES
	NdLESEQ
	NdASSIGN
)
