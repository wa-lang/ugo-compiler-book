package main

import "fmt"

// 词法记号类型
type TokenType int

const (
	EOF TokenType = iota
	ILLEGAL

	NUMBER

	ADD // +
	SUB // -
	MUL // *
	DIV // /

	LPAREN // (
	RPAREN // )
)

func (op TokenType) Precedence() int {
	switch op {
	case ADD, SUB:
		return 1
	case MUL, DIV:
		return 2
	}
	return 0
}

// 记号值
type Token struct {
	Type TokenType // 记号类型
	Val  string    // 记号原始字面值
	Pos  int       // 开始位置
}

func (t Token) String() string {
	return fmt.Sprintf("Token(%v:%v)", t.Type, t.Val)
}
