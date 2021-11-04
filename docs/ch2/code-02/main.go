package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func main() {
	code, _ := io.ReadAll(os.Stdin)
	fmt.Println(run(string(code)))
}

func run(code string) int {
	compile(code)
	if err := exec.Command("./a.out").Run(); err != nil {
		return err.(*exec.ExitError).ExitCode()
	}
	return 0
}

func compile(code string) {
	tokens := parse_tokens(code)
	output := gen_asm(tokens)

	os.WriteFile("_output_amd64.s", []byte(output), 0666)
	exec.Command("gcc", "-o", "a.out", "_output_amd64.s").Run()
}

func parse_tokens(code string) (tokens []string) {
	for code != "" {
		if idx := strings.IndexAny(code, "+-"); idx >= 0 {
			if idx > 0 {
				tokens = append(tokens, strings.TrimSpace(code[:idx]))
			}
			tokens = append(tokens, code[idx:][:1])
			code = code[idx+1:]
			continue
		}

		tokens = append(tokens, strings.TrimSpace(code))
		return
	}
	return
}

func gen_asm(tokens []string) string {
	var buf bytes.Buffer

	fmt.Fprintln(&buf, `.intel_syntax noprefix`)
	fmt.Fprintln(&buf, `.globl _main`)
	fmt.Fprintln(&buf)

	fmt.Fprintln(&buf, `_main:`)
	for i, tok := range tokens {
		if i == 0 {
			fmt.Fprintln(&buf, `    mov rax,`, tokens[i])
			continue
		}
		switch tok {
		case "+":
			fmt.Fprintln(&buf, `    add rax,`, tokens[i+1])
		case "-":
			fmt.Fprintln(&buf, `    sub rax,`, tokens[i+1])
		}
	}
	fmt.Fprintln(&buf, `    ret`)

	return buf.String()
}
