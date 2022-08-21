package gocc

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"log"
	"math/rand"
	"os"
	"os/exec"
)

func GenSingleIdentStmt() (int, string) {
	var s string
	var r string
	values := map[string]interface{}{}
	var result int

	for {

		s = ""
		r = ""

		for i := 0; i < 22; i++ {
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

	return result, fmt.Sprintf("{%s return %s;}", s, r)
}

// GenFizzBuzz return 15x: 0, 3x: 1, 5x: 2, else: 3
func GenFizzBuzz(in int) string {
	src := fmt.Sprintf("i = %d;", in)
	// 15
	src += "if ( ((i/15)*15) == i ) { return 0; }"
	// 3
	src += "else if ( ((i/3)*3) == i) { return 1; }"
	// 5
	src += "else if ( ((i/5)*5) == i ) {return 2; }"
	// else
	src += "else { return 3; }"
	return src
}

func ExecAndGetExitCode(compilerPath string, in string) int {
	asm, err := exec.Command(compilerPath, in).Output()
	if err != nil {
		log.Fatalf("failed to compile statement: %v", err)
	}

	err = os.WriteFile("./bin/tmp.s", asm, 0644)
	if err != nil {
		log.Fatalf("failed to write asm to tmp.s: %v", err)
	}

	err = exec.Command("cc", "-o", "./bin/tmp", "./bin/tmp.s").Run()
	if err != nil {
		log.Fatalf("failed to compile asm: %v", err)
	}

	cmd := exec.Command("./bin/tmp")
	_ = cmd.Run()
	exitCode := cmd.ProcessState.ExitCode()
	return exitCode
}
