package lexer

import (
	"strings"
	"unicode/utf8"
)

type Stream struct {
	name  string // 文件名
	input string // 输入的源代码
	start int    // 当前正解析中的记号的开始位置
	pos   int    // 当前读取的位置
	width int    // 最后一次读取utf8字符的字节宽度, 用于回退
}

func NewStream(name, src string) *Stream {
	return &Stream{name: name, input: src}
}

func (p *Stream) Name() string {
	return p.name
}
func (p *Stream) Input() string {
	return p.input
}

func (p *Stream) Pos() int {
	return p.pos
}

func (p *Stream) Peek() rune {
	x := p.Read()
	p.Unread()
	return x
}

func (p *Stream) Read() rune {
	if p.pos >= len(p.input) {
		p.width = 0
		return 0
	}

	r, size := utf8.DecodeRune([]byte(p.input[p.pos:]))
	p.width = size
	p.pos += p.width
	return r
}
func (p *Stream) Unread() {
	p.pos -= p.width
	return
}

func (p *Stream) Accept(valid string) bool {
	if strings.IndexRune(valid, rune(p.Read())) >= 0 {
		return true
	}
	return false
}

func (p *Stream) AcceptRun(valid string) (ok bool) {
	for p.Accept(valid) {
		ok = true
	}
	p.Unread()
	return
}

func (p *Stream) EmitToken() (lit string, pos int) {
	lit, pos = p.input[p.start:p.pos], p.start
	p.start = p.pos
	return
}

func (p *Stream) IgnoreToken() {
	_, _ = p.EmitToken()
}
