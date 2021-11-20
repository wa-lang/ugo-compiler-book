package main

import "fmt"

// 词法记号类型
type lexType int

// 记号值
type Token struct {
	Type lexType // 记号类型
	Val  string  // 记号原始字面值
	Pos  int     // 开始位置
}

func (t Token) String() string {
	return fmt.Sprintf("Token(%v:%v)", t.Type, t.Val)
}
