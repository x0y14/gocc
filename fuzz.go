package gocc

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"math/rand"
)

func GenSingleIdentStmt() (int, string) {
	var s string
	var r string
	values := map[string]interface{}{}
	var result int

	for {

		s = ""
		r = ""

		for i := 0; i < 26; i++ {
			op := rand.Intn(3)
			val := rand.Intn(11)

			values[string(rune('a'+i))] = val

			s += fmt.Sprintf("%s=%d;", string(rune('a'+i)), val)
			if i == 0 {
				r += string(rune('a' + i))
				continue
			}

			switch op {
			case 0:
				r += "+" + string(rune('a'+i))
			case 1:
				r += "-" + string(rune('a'+i))
			case 2:
				r += "*" + string(rune('a'+i))
			}
		}

		expression, err := govaluate.NewEvaluableExpression(r)
		if err != nil {
			panic(err)
		}
		result_, err := expression.Evaluate(values)

		if 0 <= (result_).(float64) && (result_).(float64) < 255 {
			result = int((result_).(float64))
			break
		}

	}

	return result, fmt.Sprintf("%s%s;", s, r)
}
