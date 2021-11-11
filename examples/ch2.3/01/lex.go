package main

import "strings"

func Lex(code string) (tokens []string) {
	for code != "" {
		if idx := strings.IndexAny(code, "+-*/()"); idx >= 0 {
			if idx > 0 {
				tokens = append(tokens, strings.TrimSpace(code[:idx]))
			}
			tokens = append(tokens, code[idx:][:1])
			code = code[idx+1:]
			continue
		}

		tokens = append(tokens, strings.TrimSpace(code))
		return
	}
	return
}
