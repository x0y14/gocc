package gocc

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTokenize(t *testing.T) {
	var tests = []struct {
		name   string
		in     string
		expect *Token
	}{
		{
			"plus",
			"1+1",
			nil,
		},
		{
			"plus2",
			"12+30",
			&Token{
				kind: TkNUM,
				next: &Token{
					kind: TkRESERVED,
					next: &Token{
						kind: TkNUM,
						next: &Token{TkEOF, nil, 0, []rune{}, 0},
						val:  30,
						str:  []rune("30"),
						len:  2,
					},
					val: 0,
					str: []rune("+"),
					len: 1,
				},
				val: 12,
				str: []rune("12"),
				len: 2,
			},
		},
		{
			"error",
			"1#",
			nil,
		},
		{
			"error2",
			"if ( 1 == 1 ) return 1;",
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tok := Tokenize([]rune(tt.in))
			assert.Equal(t, tt.expect, tok)
		})
	}

}
