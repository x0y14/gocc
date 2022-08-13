package gocc

type TokenKind int

const (
	TkRESERVED = iota
	TkNUM
	TkEOF
)

func (tk TokenKind) String() string {
	return []string{
		"TkRESERVED",
		"TkNUM",
		"TkEOF",
	}[tk]
}
