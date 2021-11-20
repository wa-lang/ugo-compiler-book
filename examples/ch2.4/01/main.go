package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func main() {
	code, _ := io.ReadAll(os.Stdin)
	tokens := Lex(string(code))
	fmt.Println(tokens)

	ast := ParseExpr(tokens)
	fmt.Println(JSONString(ast))

	fmt.Println(run(ast))
}

func run(node *ExprNode) int {
	compile(node)
	if data, err := exec.Command("./a.out").CombinedOutput(); err != nil {
		fmt.Print(string(data))
		return err.(*exec.ExitError).ExitCode()
	}
	return 0
}

func compile(node *ExprNode) {
	output := new(Compiler).GenLLIR(node)

	os.WriteFile("a.out.ll", []byte(output), 0666)
	data, err := exec.Command("clang", "-Wno-override-module", "-o", "a.out", "a.out.ll").CombinedOutput()
	if err != nil {
		fmt.Print(string(data))
		os.Exit(1)
	}
}

func JSONString(x interface{}) string {
	d, _ := json.MarshalIndent(x, "", "    ")
	return string(d)
}
