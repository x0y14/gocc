package gocc

type TokenKind int

const (
	TkRESERVED = iota
	TkNUM
	TkIDENT
	TkRETURN
	TkIF
	TkEOF
)

func (tk TokenKind) String() string {
	return []string{
		"TkRESERVED",
		"TkNUM",
		"TkIDENT",
		"TkRETURN",
		"TkIF",
		"TkEOF",
	}[tk]
}
