package gocc

import (
	"github.com/stretchr/testify/assert"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"testing"
	"time"
)

const (
	compiler = "./bin/gocc"
)

func TestSingleIdent(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10; i++ {
		want, statement := GenSingleIdentStmt()
		log.Println("in:", statement)
		log.Println("expect:", want)

		asm, err := exec.Command(compiler, statement).Output()
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
		assert.Equal(t, want, exitCode)
	}

}