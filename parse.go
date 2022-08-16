package gocc

import (
	"fmt"
)

// 現在着目しているトークン
var token *Token

// 着目しているローカル変数
var locals *LocalVariable

// 構文解析により得られたノード群
var code []*Node

func consume(op string) bool {
	if token.kind != TkRESERVED ||
		len([]rune(op)) != token.len ||
		!runeCmp(token.str, []rune(op)) {
		return false
	}
	token = token.next
	return true
}

func consumeIdent() *Token {
	if token.kind != TkIDENT {
		return nil
	}
	tok := token
	token = token.next
	return tok
}

func consumeKeyword(kind TokenKind) bool {
	if token.kind != kind {
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

func findLocalVariable(tok *Token) *LocalVariable {
	for variable := locals; variable != nil; variable = variable.next {
		if variable != nil && runeCmp(variable.name, tok.str) {
			return variable
		}
	}
	return nil
}

func program() {
	for !atEof() {
		//fmt.Println("no eof")
		code = append(code, stmt())
	}
	code = append(code, nil)
}

func stmt() *Node {
	var node *Node
	if consumeKeyword(TkRETURN) {
		node = NewNode(NdRETURN, expr(), nil)
	} else {
		node = expr()
	}
	expect(";")
	return node
}

func expr() *Node {
	return assign()
}

func assign() *Node {
	node := equality()
	if consume("=") {
		return NewNode(NdASSIGN, node, equality())
	}
	return node
}

func equality() *Node {
	node := relational()
	for {
		if consume("==") {
			node = NewNode(NdEQ, node, relational())
		} else if consume("!=") {
			node = NewNode(NdNE, node, relational())
		} else {
			return node
		}
	}
}

func relational() *Node {
	node := add()
	for {
		if consume("<") {
			node = NewNode(NdLT, node, add())
		} else if consume("<=") {
			node = NewNode(NdLE, node, add())
		} else if consume(">") {
			node = NewNode(NdLT, add(), node)
		} else if consume(">=") {
			node = NewNode(NdLE, add(), node)
		} else {
			return node
		}
	}
}

func add() *Node {
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
	} else if consume("-") {
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
	if identToken := consumeIdent(); identToken != nil {
		var node *Node
		lVar := findLocalVariable(identToken)
		if lVar != nil {
			// すでに宣言されているので、どこにデータが入っているのかというデータをもつノードを
			// 作成する
			node = NewNodeLVar(lVar.offset)
		} else {
			// 新しく宣言
			lVar = &LocalVariable{
				next:   locals,
				name:   identToken.str,
				len:    identToken.len,
				offset: locals.offset + 16, // 最後に宣言されたものから16離れた位置にデータを入れてあげる
			}
			node = NewNodeLVar(lVar.offset)
			// 新しく宣言したのでそれを一覧につなげてあげる
			locals = lVar
		}
		return node
	}
	return NewNodeNum(expectNumber())
}

func Parse(tok *Token) []*Node {
	// init
	token = tok
	locals = &LocalVariable{}
	code = []*Node{}

	program()
	return code
}
