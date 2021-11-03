//go:build ignore

// echo 122 | go run main.go

package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

func main() {
	code, _ := io.ReadAll(os.Stdin)
	compile(string(code))
}

func compile(code string) {
	output := fmt.Sprintf(tmpl, code)
	os.WriteFile("_output_amd64.s", []byte(output), 0666)
	exec.Command("gcc", "-o", "a.out", "_output_amd64.s").Run()
}

const tmpl = `
.intel_syntax noprefix
.globl _main

_main:
	mov rax, %v
	ret
`
