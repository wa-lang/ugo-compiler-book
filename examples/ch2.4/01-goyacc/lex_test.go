package main

import "testing"

func TestLex(t *testing.T) {
	for i, tt := range tLexTests {
		gotTokens := Lex(tt.Code)
		for j, got := range gotTokens {
			if got.Type != tt.Tokens[j].Type {
				t.Fatalf(
					"%d/%d: expect-type = %v, got = %+v, code = %q",
					i, j, tt.Tokens[j].Type, gotTokens, tt.Code,
				)
			}
		}
	}
}

var tLexTests = []struct {
	Code   string
	Tokens []Token
}{
	{
		Code: "1",
		Tokens: []Token{
			{Type: NUMBER},
			{Type: EOF},
		},
	},
	{
		Code: "1+2*(3-4/1)",
		Tokens: []Token{
			{Type: NUMBER},
			{Type: ADD},
			{Type: NUMBER},
			{Type: MUL},
			{Type: LPAREN},
			{Type: NUMBER},
			{Type: SUB},
			{Type: NUMBER},
			{Type: DIV},
			{Type: NUMBER},
			{Type: RPAREN},
			{Type: EOF},
		},
	},
}
