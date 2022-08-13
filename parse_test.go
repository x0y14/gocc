package gocc

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParse(t *testing.T) {
	var tests = []struct {
		name   string
		token  *Token
		expect *Node
	}{
		{
			name:   "1-3",
			token:  Tokenize([]rune("1-3")),
			expect: NewNode(NdSUB, NewNodeNum(1), NewNodeNum(3)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expect, Parse(tt.token))
		})
	}
}
