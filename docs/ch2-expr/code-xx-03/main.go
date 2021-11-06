package main

func main() {}

//   +
//  / \
// 1   *
//    / \
//   2   +
//      / \
//     3   4
type ExprNode struct {
	Left  interface{} // num, *Node
	Op    string      // +-*/
	Right interface{}
}

var expr = &ExprNode{
	Left: 1, Op: "+",
	Right: &ExprNode{
		Left: 2, Op: "*",
		Right: &ExprNode{
			Left: 3, Op: "+",
			Right: 4,
		},
	},
}
