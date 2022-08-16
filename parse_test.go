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
				NewNode(NdIF,
					NewNode(NdLT, NewNodeNum(8), NewNodeNum(2)),
					NewNode(NdRETURN, NewNodeNum(10), nil)),
				NewNode(NdRETURN, NewNodeNum(20), nil),
				nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expect, Parse(tt.token))
		})
	}
}

func TestGen(t *testing.T) {
	token := Tokenize([]rune("a=1;"))
	code := Parse(token)
	for _, node := range code {
		Gen(node)
	}
}
