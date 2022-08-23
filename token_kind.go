package gocc

type TokenKind int

const (
	TkRESERVED = iota
	TkNUM
	TkIDENT
	TkRETURN
	TkIF
	TkELSE
	TkWHILE
	TkFOR
	TkEOF
)

func (tk TokenKind) String() string {
	return []string{
		"TkRESERVED",
		"TkNUM",
		"TkIDENT",
		"TkRETURN",
		"TkIF",
		"TkELSE",
		"TkWHILE",
		"TkFOR",
		"TkEOF",
	}[tk]
}
