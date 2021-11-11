package main

type ExprNode struct {
	Value string // +, -, *, /, 123
	Left  *ExprNode
	Right *ExprNode
}
