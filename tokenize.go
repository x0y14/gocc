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
	for p < len(r) {
		if unicode.IsSpace(userInput[p]) {
			p++
			continue
		}

		for _, longOp := range []string{"==", "!=", "<=", ">="} {
			if startWith(longOp) {
				p += 2
				cur = NewToken(TkRESERVED, cur, []rune(longOp), 2)
				continue
			}
		}

		if strings.ContainsRune("+-*/()=<>!;", userInput[p]) {
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

		if 'a' <= userInput[p] && userInput[p] <= 'z' {
			cur = NewToken(TkIDENT, cur, []rune{userInput[p]}, 1)
			p++
			continue
		}

		errorAt("パースできません")

	}

	NewToken(TkEOF, cur, []rune{}, 0)
	return head.next
}
