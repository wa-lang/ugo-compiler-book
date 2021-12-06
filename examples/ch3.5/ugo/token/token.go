package token

import (
	"strconv"
)

// 记号类型
type TokenType int

// ch3中 µGo程序用到的记号类型
const (
	EOF TokenType = iota // = 0
	ERROR
	COMMENT

	IDENT
	NUMBER

	PACKAGE
	FUNC

	ADD // +
	SUB // -
	MUL // *
	DIV // /

	LPAREN // (
	RPAREN // )
	LBRACE // {
	RBRACE // }

	SEMICOLON // ;
)

// 记号值
type Token struct {
	Pos     Pos       // 记号所在的位置(从1开始)
	Type    TokenType // 记号的类型
	Literal string    // 程序中原始的字符串
}

var tokens = [...]string{
	EOF:     "EOF",
	ERROR:   "ERROR",
	COMMENT: "COMMENT",

	IDENT:  "IDENT",
	NUMBER: "NUMBER",

	PACKAGE: "PACKAGE",
	FUNC:    "FUNC",

	ADD: "+",
	SUB: "-",
	MUL: "*",
	DIV: "/",

	LPAREN: "(",
	RPAREN: ")",
	LBRACE: "{",
	RBRACE: "}",

	SEMICOLON: ";",
}

func (tokType TokenType) String() string {
	s := ""
	if 0 <= tokType && tokType < TokenType(len(tokens)) {
		s = tokens[tokType]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tokType)) + ")"
	}
	return s
}

var keywords = map[string]TokenType{
	"package": PACKAGE,
	"func":    FUNC,
}

func Lookup(ident string) TokenType {
	if tok, is_keyword := keywords[ident]; is_keyword {
		return tok
	}
	return IDENT
}

func (op TokenType) Precedence() int {
	switch op {
	case ADD, SUB:
		return 1
	case MUL, DIV:
		return 2
	}
	return 0
}
