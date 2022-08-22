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
	NdSTRING
	NdFUNCTION
)

var nodeKinds = [...]string{
	NdADD:      "NdADD",
	NdSUB:      "NdSUB",
	NdMUL:      "NdMUL",
	NdDIV:      "NdDIV",
	NdEQ:       "NdEQ",
	NdNE:       "NdNE",
	NdLT:       "NdLT",
	NdLE:       "NdLE",
	NdAND:      "NdAND",
	NdOR:       "NdOR",
	NdASSIGN:   "NdASSIGN",
	NdNUM:      "NdNUM",
	NdLVAR:     "NdLVAR",
	NdCALL:     "NdCALL",
	NdRETURN:   "NdRETURN",
	NdIF:       "NdIF",
	NdIFELSE:   "NdIFELSE",
	NdWHILE:    "NdWHILE",
	NdFOR:      "NdFOR",
	NdBLOCK:    "NdBLOCK",
	NdSTRING:   "NdSTRING",
	NdFUNCTION: "NdFUNCTION",
}

func (n NodeKind) String() string {
	return nodeKinds[n]
}
