package token

import (
	"fmt"
	"strconv"
)

// 记号值
type Token struct {
	Type    TokenType
	Literal string
	Pos     Pos
	Value   interface{}
}

func (i Token) EndPos() Pos {
	if i.Pos != NoPos {
		return i.Pos + Pos(len(i.Literal))
	}
	return NoPos
}

func (i Token) IdentName() string {
	return i.Literal
}

func (i Token) RuneValue() rune {
	if x, ok := i.Value.(rune); ok {
		return x
	}
	x, _, _, _ := strconv.UnquoteChar(i.Literal, '\'')
	i.Value = x
	return x
}
func (i Token) IntValue() int64 {
	if x, ok := i.Value.(int64); ok {
		return x
	}
	x, _ := strconv.ParseInt(i.Literal, 10, 64)
	i.Value = x
	return x
}
func (i Token) FloatValue() float64 {
	if x, ok := i.Value.(float64); ok {
		return x
	}
	x, _ := strconv.ParseFloat(i.Literal, 64)
	i.Value = x
	return x
}

func (i Token) StringValue() string {
	if x, ok := i.Value.(string); ok {
		return x
	}
	s, _ := strconv.Unquote(i.Literal)
	i.Value = s
	return s
}

func (a Token) equal(b Token) bool {
	if a.Type != b.Type {
		return false
	}
	if a.Literal != b.Literal {
		return false
	}
	if a.Pos != NoPos && b.Pos != NoPos {
		if a.Pos != b.Pos {
			return false
		}
	}

	return true
}

func (i Token) String() string {
	return fmt.Sprintf("{%v:%q}", i.Type, i.Literal)
}
