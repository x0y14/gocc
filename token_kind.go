package gocc

type TokenKind int

const (
	TkRESERVED = iota
	TkNUM
	TkIDENT
	TkEOF
)

func (tk TokenKind) String() string {
	return []string{
		"TkRESERVED",
		"TkNUM",
		"TkIDENT",
		"TkEOF",
	}[tk]
}
