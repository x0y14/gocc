package gocc

import (
	"fmt"
)

// 現在着目しているトークン
var token *Token

func consume(op string) bool {
	if token.kind != TkRESERVED ||
		len([]rune(op)) != token.len ||
		!runeCmp(token.str, []rune(op)) {
		return false
	}
	token = token.next
	return true
}

func expect(op string) {
	if token.kind != TkRESERVED ||
		!runeCmp(token.str, []rune(op)) {
		errorAt(fmt.Sprintf("%sではありません", op))
	}
	token = token.next
}

func expectNumber() int {
	if token.kind != TkNUM {
		errorAt("数字ではありません")
	}
	val := token.val
	token = token.next
	return val
}

func atEof() bool {
	return token.kind == TkEOF
}

func expr() *Node {
	node := mul()
	for {
		if consume("+") {
			node = NewNode(NdADD, node, mul())
		} else if consume("-") {
			node = NewNode(NdSUB, node, mul())
		} else {
			return node
		}
	}
}

func mul() *Node {
	node := unary()
	for {
		if consume("*") {
			node = NewNode(NdMUL, node, unary())
		} else if consume("/") {
			node = NewNode(NdDIV, node, unary())
		} else {
			return node
		}
	}
}

func unary() *Node {
	if consume("+") {
		return primary()
	}
	if consume("-") {
		return NewNode(NdSUB, NewNodeNum(0), primary())
	}
	return primary()
}

func primary() *Node {
	if consume("(") {
		node := expr()
		expect(")")
		return node
	}
	return NewNodeNum(expectNumber())
}

func parse(tok *Token) *Node {
	token = tok
	return expr()
}
