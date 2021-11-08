package main

type ExprNode struct {
	Value string    // +, -, *, /, 123
	Left  *ExprNode `json:",omitempty"`
	Right *ExprNode `json:",omitempty"`
}

func NewExprNode(value string, left, right *ExprNode) *ExprNode {
	return &ExprNode{
		Value: value,
		Left:  left,
		Right: right,
	}
}
