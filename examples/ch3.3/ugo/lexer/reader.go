package lexer

import (
	"strings"
	"unicode/utf8"
)

type Reader interface {
	Input() string

	Peek() rune
	Read() rune
	Unread()

	Accept(valid string) (lit string, pos int, ok bool)
	AcceptRun(valid string) (litList []string, posList []int, ok bool)

	EmitToken() (lit string, pos int)
	IgnoreToken()
}

func NewReader(src string) Reader {
	return &srcReader{input: src}
}

type srcReader struct {
	input string // 输入的源代码
	start int    // 当前正解析中的记号的开始位置
	pos   int    // 当前读取的位置
	width int    // 最后一次读取utf8字符的字节宽度, 用于回退
}

func (p *srcReader) Input() string {
	return p.input
}

func (p *srcReader) Peek() rune {
	x := p.Read()
	p.Unread()
	return x
}

func (p *srcReader) Read() rune {
	if p.pos >= len(p.input) {
		p.width = 0
		return 0
	}
	r, size := utf8.DecodeRune([]byte(p.input[p.pos:]))
	p.width = size
	p.pos += p.width
	return r
}
func (p *srcReader) Unread() {
	p.pos -= p.width
	return
}

func (p *srcReader) EmitToken() (lit string, pos int) {
	lit, pos = p.input[p.start:p.pos], p.start
	p.start = p.pos
	return
}

func (p *srcReader) IgnoreToken() {
	_, _ = p.EmitToken()
}

func (p *srcReader) Accept(valid string) (lit string, pos int, ok bool) {
	if strings.IndexRune(valid, rune(p.Read())) >= 0 {
		lit, pos = p.EmitToken()
		ok = true
		return
	}
	p.Unread()
	ok = false
	return
}

func (p *srcReader) AcceptRun(valid string) (litList []string, posList []int, ok bool) {
	for {
		lit, pos, ok := p.Accept(valid)
		if !ok {
			break
		}
		litList = append(litList, lit)
		posList = append(posList, pos)
	}
	ok = len(litList) > 0
	return
}
