package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	//   +
	//  / \
	// 1   *
	//    / \
	//   2   +
	//      / \
	//     3   4
	expr_tokens := []string{"1", "+", "2", "*", "(", "3", "+", "4", ")"}

	ast := ParseExpr(expr_tokens)
	fmt.Println(JSONString(ast))

	fmt.Println(run(ast))
}

func run(node *ExprNode) int {
	compile(node)
	if err := exec.Command("./a.out").Run(); err != nil {
		return err.(*exec.ExitError).ExitCode()
	}
	return 0
}

func compile(node *ExprNode) {
	output := new(Compiler).GenLLIR(node)

	os.WriteFile("a.out.ll", []byte(output), 0666)
	exec.Command("clang", "-Wno-override-module", "-o", "a.out", "a.out.ll").Run()
}

func JSONString(x interface{}) string {
	d, _ := json.MarshalIndent(x, "", "    ")
	return string(d)
}
