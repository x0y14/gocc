package gocc

import (
	"log"
	"strconv"
	"strings"
	"unicode"
)

// 現在着目しているuserInputのruneの場所
var userInput []rune
var p int

func getRunes(len int) []rune {
	var current []rune
	for i := 0; i < len; i++ {
		current = append(current, current[p+i])
	}
	return current
}

func startWith(q string) bool {
	qRunes := []rune(q)
	for i := 0; i < len(qRunes); i++ {
		if userInput[p+i] != qRunes[i] {
			return false
		}
	}
	return true
}

// isAlnum
// r in [0-9a-zA-Z_]
func isAlnum(r rune) bool {
	return ('a' <= r && r <= 'z') ||
		('A' <= r && r <= 'Z') ||
		('0' <= r && r <= '9') || r == '_'
}

// startWithAndAfterIsNotAlnum qで始まるか　かつ　その後が変数として使用可能な文字
func startWithAndAfterIsNotAlnum(q string) bool {
	return startWith(q) && !isAlnum(userInput[p+len(q)])
}

// userInput[p]が数字である限り取得し、繋げ、次へ進む
func strToI() int {
	var integer []rune
	for p < len(userInput) {
		if unicode.IsDigit(userInput[p]) {
			integer = append(integer, userInput[p])
			p++
		} else {
			break
		}
	}

	v, err := strconv.Atoi(string(integer))
	if err != nil {
		log.Panicf("整数をパースできませんでした%s", err)
	}
	return v
}

// userInput[p]がアルファベットであるかぎり取得し、繋げ、次へ進む
func strToIdent() []rune {
	var ident []rune
	for p < len(userInput) {
		if isAlnum(userInput[p]) {
			ident = append(ident, userInput[p])
			p++
		} else {
			break
		}
	}
	return ident
}

func errorAt(text string) {
	log.Fatalf("\n%s\n%s^ %s", string(userInput), strings.Repeat(" ", p), text)
}
func Tokenize(r []rune) *Token {
	userInput = r
	p = 0
	var head Token
	head.next = nil
	cur := &head

	// while (*p)
userInputLoop:
	for p < len(r) {
		if unicode.IsSpace(userInput[p]) {
			p++
			continue
		}

		for _, longOp := range []string{"==", "!=", "<=", ">=", "&&", "||"} {
			if startWith(longOp) {
				p += len(longOp)
				cur = NewToken(TkRESERVED, cur, []rune(longOp), len(longOp))
				continue userInputLoop
			}
		}

		if strings.ContainsRune("{}+-*/()=<>!;,", userInput[p]) {
			cur = NewToken(TkRESERVED, cur, []rune{userInput[p]}, 1)
			p++
			continue
		}

		if unicode.IsDigit(userInput[p]) {
			integer := strToI()
			strInt := strconv.Itoa(integer)
			cur = NewToken(TkNUM, cur, []rune(strInt), len([]rune(strInt)))
			cur.val = integer
			continue
		}

		if startWithAndAfterIsNotAlnum("return") {
			cur = NewToken(TkRETURN, cur, []rune("return"), 6)
			p += 6
			continue
		}

		if startWithAndAfterIsNotAlnum("if") {
			cur = NewToken(TkIF, cur, []rune("if"), 2)
			p += 2
			continue
		}

		if startWithAndAfterIsNotAlnum("else") {
			cur = NewToken(TkELSE, cur, []rune("else"), 4)
			p += 4
			continue
		}

		if startWithAndAfterIsNotAlnum("while") {
			cur = NewToken(TkWHILE, cur, []rune("while"), 5)
			p += 5
			continue
		}

		if startWithAndAfterIsNotAlnum("for") {
			cur = NewToken(TkFOR, cur, []rune("for"), 3)
			p += 3
			continue
		}

		// 数字から始まることはない
		if ('a' <= userInput[p] && userInput[p] <= 'z') ||
			('A' <= userInput[p] && userInput[p] <= 'Z') || '_' == userInput[p] {
			ident := strToIdent()
			cur = NewToken(TkIDENT, cur, ident, len(ident))
			continue
		}

		errorAt("トークナイズできません")

	}

	NewToken(TkEOF, cur, []rune{}, 0)
	return head.next
}
