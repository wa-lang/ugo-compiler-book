package lexer

import (
	"strings"
	"unicode/utf8"
)

type SourceStream struct {
	name  string // 文件名
	input string // 输入的源代码
	start int    // 当前正解析中的记号的开始位置
	pos   int    // 当前读取的位置
	width int    // 最后一次读取utf8字符的字节宽度, 用于回退
}

func NewSourceStream(name, src string) *SourceStream {
	return &SourceStream{name: name, input: src}
}

func (p *SourceStream) Name() string {
	return p.name
}
func (p *SourceStream) Input() string {
	return p.input
}

func (p *SourceStream) Pos() int {
	return p.pos
}

func (p *SourceStream) Peek() rune {
	x := p.Read()
	p.Unread()
	return x
}

func (p *SourceStream) Read() rune {
	if p.pos >= len(p.input) {
		p.width = 0
		return 0
	}

	r, size := utf8.DecodeRune([]byte(p.input[p.pos:]))
	p.width = size
	p.pos += p.width
	return r
}
func (p *SourceStream) Unread() {
	p.pos -= p.width
	return
}

func (p *SourceStream) Accept(valid string) bool {
	if strings.IndexRune(valid, rune(p.Read())) >= 0 {
		return true
	}
	return false
}

func (p *SourceStream) AcceptRun(valid string) (ok bool) {
	for p.Accept(valid) {
		ok = true
	}
	p.Unread()
	return
}

func (p *SourceStream) EmitToken() (lit string, pos int) {
	lit, pos = p.input[p.start:p.pos], p.start
	p.start = p.pos
	return
}

func (p *SourceStream) IgnoreToken() {
	_, _ = p.EmitToken()
}
