package main

import (
	"fmt"
	"strings"
	"text/scanner"
)

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

func Lex(code string) (tokens []Token) {
	var s scanner.Scanner
	s.Init(strings.NewReader(code))
	for x := s.Scan(); x != scanner.EOF; x = s.Scan() {
		var tok = Token{
			Val: s.TokenText(),
			Pos: s.Pos().Offset,
		}
		switch x {
		case scanner.Int:
			tok.Type = NUMBER
		default:
			switch s.TokenText() {
			case "+":
				tok.Type = ADD
			case "-":
				tok.Type = SUB
			case "*":
				tok.Type = MUL
			case "/":
				tok.Type = DIV
			case "(":
				tok.Type = LPAREN
			case ")":
				tok.Type = RPAREN
			default:
				tok.Type = ILLEGAL
				tokens = append(tokens, tok)
				return
			}
		}

		tokens = append(tokens, tok)
	}

	tokens = append(tokens, Token{Type: EOF})
	return
}
