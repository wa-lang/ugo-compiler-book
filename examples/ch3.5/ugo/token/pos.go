package token

import (
	"fmt"
)

// Pos 类似一个指针, 表示文件中的位置.
type Pos int

// NoPos 类似指针的 nil 值, 表示一个无效的位置.
const NoPos Pos = 0

func (p Pos) IsValid() bool {
	return p != NoPos
}

// PosString 格式化 pos 为字符串
func PosString(filename string, src []byte, pos Pos) string {
	var p = &Position{
		Filename: filename,
	}

	if pos.IsValid() {
		p.Offset = int(pos) - 1
	}

	for _, c := range string(src) {
		if c == '\n' {
			p.Line++
		}
	}

	return p.String()

}

type Position struct {
	Filename string // filename, if any
	Offset   int    // offset, starting at 0
	Line     int    // line number, starting at 1
	Column   int    // column number, starting at 1 (byte count)
}

// IsValid reports whether the position is valid.
func (pos *Position) IsValid() bool { return pos.Line > 0 }

// String returns a string in one of several forms:
//
//	file:line:column    valid position with file name
//	file:line           valid position with file name but no column (column == 0)
//	line:column         valid position without file name
//	line                valid position without file name and no column (column == 0)
//	file                invalid position with file name
//	-                   invalid position without file name
//
func (pos Position) String() string {
	s := pos.Filename
	if pos.IsValid() {
		if s != "" {
			s += ":"
		}
		s += fmt.Sprintf("%d", pos.Line)
		if pos.Column != 0 {
			s += fmt.Sprintf(":%d", pos.Column)
		}
	}
	if s == "" {
		s = "-"
	}
	return s
}
