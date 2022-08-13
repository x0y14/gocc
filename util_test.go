package gocc

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRuneCmp(t *testing.T) {
	var tests = []struct {
		name   string
		v1     string
		v2     string
		expect bool
	}{
		{
			"same",
			"hello",
			"hello",
			true,
		},
		{
			"different",
			"1",
			"2",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expect, runeCmp([]rune(tt.v1), []rune(tt.v2)))
		})
	}
}
