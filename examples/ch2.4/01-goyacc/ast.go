package main

type ExprNode struct {
	Token           // +, -, *, /, 123
	Left  *ExprNode `json:",omitempty"`
	Right *ExprNode `json:",omitempty"`
}

func NewExprNode(token Token, left, right *ExprNode) *ExprNode {
	return &ExprNode{
		Token: token,
		Left:  left,
		Right: right,
	}
}
