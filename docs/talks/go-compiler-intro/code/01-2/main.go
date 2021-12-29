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
	compile(string(code)) // HL
}

func compile(code string) {
	output := fmt.Sprintf(tmpl, code) // HL
	os.WriteFile("a.out.ll", []byte(output), 0666)
	exec.Command("clang", "-o", "a.out", "a.out.ll").Run() // HL
}

const tmpl = `
define i32 @main() {
	ret i32 %v
}
`
