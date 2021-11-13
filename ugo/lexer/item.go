package lexer

import (
	"fmt"
	"strconv"

	"github.com/chai2010/ugo/token"
)

type Item struct {
	Token   token.Token
	Literal string
	Pos     token.Pos
	Value   interface{}
}

func (i Item) RuneValue() rune {
	if x, ok := i.Value.(rune); ok {
		return x
	}
	x, _, _, _ := strconv.UnquoteChar(i.Literal, '\'')
	i.Value = x
	return x
}
func (i Item) IntValue() int64 {
	if x, ok := i.Value.(int64); ok {
		return x
	}
	x, _ := strconv.ParseInt(i.Literal, 10, 64)
	i.Value = x
	return x
}
func (i Item) FloatValue() float64 {
	if x, ok := i.Value.(float64); ok {
		return x
	}
	x, _ := strconv.ParseFloat(i.Literal, 64)
	i.Value = x
	return x
}

func (i Item) StringValue() string {
	if x, ok := i.Value.(string); ok {
		return x
	}
	s, _ := strconv.Unquote(i.Literal)
	i.Value = s
	return s
}

func (a Item) equal(b Item) bool {
	if a.Token != b.Token {
		return false
	}
	if a.Literal != b.Literal {
		return false
	}
	if a.Pos != token.NoPos && b.Pos != token.NoPos {
		if a.Pos != b.Pos {
			return false
		}
	}

	return true
}

func (i Item) String() string {
	return fmt.Sprintf("{%v:%q}", i.Token, i.Literal)
}
