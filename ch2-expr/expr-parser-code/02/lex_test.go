package main

import (
	"reflect"
	"testing"
)

func TestLex(t *testing.T) {
	var tests = []struct {
		input  string
		tokens []string
	}{
		{"1", []string{"1"}},
		{"1+22*333", []string{"1", "+", "22", "*", "333"}},
		{"1+2*(3+4)", []string{"1", "+", "2", "*", "(", "3", "+", "4", ")"}},
	}
	for i, tt := range tests {
		if got := Lex(tt.input); !reflect.DeepEqual(got, tt.tokens) {
			t.Fatalf("%d: expect = %v, got = %v", i, tt.tokens, got)
		}
	}
}
