package main

import "testing"

func TestRun(t *testing.T) {
	for i, tt := range tests {
		if got := run(tt.code); got != tt.value {
			t.Fatalf("%d: expect = %v, got = %v", i, tt.value, got)
		}
	}
}

var tests = []struct {
	code  string
	value int
}{
	{code: "1", value: 1},
	{code: "1+1", value: 2},
	{code: "1 + 3 - 2  ", value: 2},
	{code: "1+2+3+4", value: 10},
}
