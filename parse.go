package gocc

import (
	"fmt"
	"log"
)

// 現在着目しているトークン
var token *Token

// 着目しているローカル変数
var locals *LocalVariable

// 構文解析により得られたノード群
var code []*Node

var toplevelFunctionName string

func init() {
	toplevelFunctionName = ""
}

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

//func expectString() []rune {
//	if token.kind != TkSTRING {
//		errorAt("文字列ではありません")
//	}
//	data := token.str
//	token = token.next
//	return data
//}

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
		code = append(code, toplevel())
	}
	code = append(code, nil)
}

// toplevel   = ident "(" (ident ","?)* ")" stmt
func toplevel() *Node {
	// function define
	if ident := consumeIdent(); ident != nil {
		expect("(")

		toplevelFunctionName = string(ident.str)
		var args []*Node
		for !consume(")") {
			args = append(args, expr())
			if !consume(",") {
				expect(")")
				break
			}
		}
		//toplevelFunctionName = ""
		return NewNodeFunction(ident.str, args, stmt())
	}

	log.Fatalf("unsupported token on toplevel: %s", token.kind.String())
	return nil
}

func stmt() *Node {
	var node *Node
	switch {
	case consumeKeyword(TkRETURN):
		node = NewNode(NdRETURN, expr(), nil)
		expect(";")
	case consumeKeyword(TkIF):
		expect("(")
		// 条件式
		// A
		cond := expr()
		expect(")")
		b := stmt()
		if consumeKeyword(TkELSE) {
			// if (A) B else C
			c := stmt()
			node = NewNodeWithExpr(NdIFELSE, nil, cond, nil, b, c)
		} else {
			// if (A) B
			node = NewNodeWithExpr(NdIF, nil, cond, nil, b, nil)
		}
	case consumeKeyword(TkWHILE):
		// while (A) B
		expect("(")
		cond := expr()
		expect(")")
		node = NewNodeWithExpr(NdWHILE, nil, cond, nil, stmt(), nil)
	case consumeKeyword(TkFOR):
		// for (A;B;C) D
		expect("(")
		var init, cond, loop *Node
		if !consume(";") {
			init = expr()
			expect(";")
		}
		if !consume(";") {
			cond = expr()
			expect(";")
		}
		if !consume(")") {
			loop = expr()
			expect(")")
		}
		node = NewNodeWithExpr(NdFOR, init, cond, loop, stmt(), nil)
	case consume("{"):
		var nodesInBlock []*Node
	eofLoop:
		for !atEof() {
			if consume("}") {
				break eofLoop
			}
			nodesInBlock = append(nodesInBlock, stmt())
		}
		//nodesInBlock = append(nodesInBlock, nil)
		node = NewNodeBlock(nodesInBlock)
	default:
		node = expr()
		expect(";")
	}
	return node
}

func expr() *Node {
	return assign()
}

func assign() *Node {
	node := andor()
	if consume("=") {
		return NewNode(NdASSIGN, node, andor())
	}
	return node
}

func andor() *Node {
	node := equality()
	for {
		if consume("&&") {
			node = NewNode(NdAND, node, equality())
		} else if consume("||") {
			node = NewNode(NdOR, node, equality())
		} else {
			return node
		}
	}
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
		// f() : call function
		// f   : ident
		// のパターンがある
		var node *Node

		// 関数
		if consume("(") {
			//node = NewNodeCALL(identToken.str)
			var args []*Node
			for !consume(")") {
				args = append(args, expr())
				if !consume(",") {
					expect(")")
					break
				}
			}
			node = NewNodeCALL(identToken.str, args)
			return node
		}

		// 変数
		// 関数間で変数名が被らないようにするスコープ機能
		// 変数の先頭に関数の名前をつけてあげる
		identToken.str = append([]rune(toplevelFunctionName+"_"), identToken.str...)
		identToken.len = len(identToken.str)

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

	// string
	//if consume("\"") {
	//	node := NewNodeString(expectString())
	//	expect("\"")
	//	return node
	//}

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
