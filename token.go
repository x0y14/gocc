package gocc

type Token struct {
	kind TokenKind
	next *Token
	val  int
	str  []rune
	len  int
}

func NewToken(kind TokenKind, cur *Token, str []rune, len int) *Token {
	tok := &Token{
		kind: kind,
		next: nil,
		val:  0,
		str:  str,
		len:  len,
	}
	cur.next = tok
	return tok
}
