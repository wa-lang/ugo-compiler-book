package main

import (
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
	var expr = &ExprNode{
		Value: "+",
		Left: &ExprNode{
			Value: "1",
		},
		Right: &ExprNode{
			Value: "*",
			Left: &ExprNode{
				Value: "2",
			},
			Right: &ExprNode{
				Value: "+",
				Left: &ExprNode{
					Value: "3",
				},
				Right: &ExprNode{
					Value: "4",
				},
			},
		},
	}

	result := run(expr)
	fmt.Println(result)
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
