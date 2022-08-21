package gocc

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	compiler = "./bin/gocc"
)

func TestGenStmt(t *testing.T) {
	for i := 0; i < 20; i++ {
		result, stmts, values := GenStmt()
		expression, err := govaluate.NewEvaluableExpression(stmts[1])
		if err != nil {
			panic(err)
		}
		result_, err := expression.Evaluate(values)
		assert.Equal(t, result, int((result_).(float64)))

		statement := fmt.Sprintf("%s return %s;", stmts[0], stmts[1])

		exit := ExecAndGetExitCode(compiler, statement)
		ok := assert.Equal(t, result, exit)
		if ok {
			fmt.Printf("[OK] %s => %d\n", statement, result)
		} else {
			fmt.Printf("[FAIL] %s\n", statement)
			fmt.Printf("\t%d expected, but got %d\n", result, exit)
		}
	}
}

func TestGenFizzBuzz(t *testing.T) {
	fizzBuzz := func(i int) int {
		if i%15 == 0 {
			return 0
		}
		if i%3 == 0 {
			return 1
		}
		if i%5 == 0 {
			return 2
		}
		return 3
	}
	for i := 0; i < 100; i++ {
		ok := assert.Equal(t, fizzBuzz(i), ExecAndGetExitCode(compiler, GenFizzBuzz(i)))
		if ok {
			fmt.Printf("[OK] %d\n", i)
		} else {
			fmt.Printf("[FAIL] %d\n", i)
		}
	}
}

func Test1(t *testing.T) {
	var tests = []struct {
		name   string
		stmt   string
		expect int
	}{
		{
			"1",
			"return 1*5;",
			5,
		},
		{
			"2",
			"a=1;b=1;c=1;d=1;e=1;f=1;g=1;h=1;i=1;j=1;k=1;l=1;m=1;n=1;o=1;p=1;q=1;r=1;s=1;t=1; return a+b+c+d+e+f+g+h+i+j+k+l+m+n+o+p+q+r+s+t+a+b+c+d+e+f+g+h+i+j+k+l+m+n+o+p+q+r+s+t;",
			40,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exit := ExecAndGetExitCode(compiler, tt.stmt)
			ok := assert.Equal(t, tt.expect, exit)
			if ok {
				fmt.Printf("[OK] %s => %d\n", tt.stmt, tt.expect)
			} else {
				fmt.Printf("[FAIL] %s\n", tt.stmt)
				fmt.Printf("\t%d expected, but got %d\n", tt.expect, exit)
			}
		})
	}
}
