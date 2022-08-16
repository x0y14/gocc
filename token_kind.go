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
		"WHILE",
		"TkEOF",
	}[tk]
}
