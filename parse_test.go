package gocc

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParse(t *testing.T) {
	var tests = []struct {
		name   string
		token  *Token
		expect []*Node
	}{
		//{
		//	name:   "1-3;",
		//	token:  Tokenize([]rune("1-3;")),
		//	expect: NewNode(NdSUB, NewNodeNum(1), NewNodeNum(3)),
		//},
		{
			name:  "five * 5",
			token: Tokenize([]rune("five = 5;")),
			expect: []*Node{
				NewNode(NdASSIGN, NewNodeLVar(16), NewNodeNum(5)),
				nil,
			},
		},
		{
			name:  "return 0;",
			token: Tokenize([]rune("return 0;")),
			expect: []*Node{
				NewNode(NdRETURN, NewNodeNum(0), nil),
				nil,
			},
		},
		{
			name:  "expect 20",
			token: Tokenize([]rune("if (8 < 2) return 10; return 20;")),
			expect: []*Node{
				NewNodeWithExpr(NdIF, nil, NewNode(NdLT, NewNodeNum(8), NewNodeNum(2)), nil, NewNode(NdRETURN, NewNodeNum(10), nil), nil),
				NewNode(NdRETURN, NewNodeNum(20), nil),
				nil,
			},
		},
		{
			name:  "block",
			token: Tokenize([]rune("{return 10;}")),
			expect: []*Node{
				NewNodeBlock([]*Node{
					NewNode(NdRETURN, NewNodeNum(10), nil),
				}),
				nil,
			},
		},
		{
			name:  "call func without argument",
			token: Tokenize([]rune("add1();")),
			expect: []*Node{
				NewNodeCALL([]rune("add1"), nil),
				nil,
			},
		},
		{
			name:  "call func 1 argument",
			token: Tokenize([]rune("add(1);")),
			expect: []*Node{
				NewNodeCALL([]rune("add"), []*Node{NewNodeNum(1)}),
				nil,
			},
		},
		{
			name:  "call func with 2 argument (num, call)",
			token: Tokenize([]rune("add(1, add(2));")),
			expect: []*Node{
				NewNodeCALL([]rune("add"), []*Node{NewNodeNum(1), NewNodeCALL([]rune("add"), []*Node{NewNodeNum(2)})}),
				nil,
			},
		},
		{
			name:  "call func with 2 argument(num, num)",
			token: Tokenize([]rune("add(1, 2);")),
			expect: []*Node{
				NewNodeCALL([]rune("add"), []*Node{NewNodeNum(1), NewNodeNum(2)}),
				nil,
			},
		},
		{
			name:  "call func 3 argument(ident, call, eq)",
			token: Tokenize([]rune("add(one, add(two), one==1);")),
			expect: []*Node{
				NewNodeCALL([]rune("add"), []*Node{NewNodeLVar(16), NewNodeCALL([]rune("add"), []*Node{NewNodeLVar(32)}), NewNode(NdEQ, NewNodeLVar(16), NewNodeNum(1))}),
				nil,
			},
		},
		//{
		//	name:  "string",
		//	token: Tokenize([]rune("name=\"john\"; return 45;")),
		//	expect: []*Node{
		//		NewNode(NdASSIGN, NewNodeLVar(16), NewNodeString([]rune("john"))),
		//		NewNode(NdRETURN, NewNodeNum(45), nil),
		//		nil,
		//	},
		//},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Parse(tt.token)
			assert.Equal(t, tt.expect, p)
		})
	}
}
