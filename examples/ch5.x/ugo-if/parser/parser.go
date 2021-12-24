package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/lexer"
	"github.com/chai2010/ugo/token"
)

var DebugMode = false

func ParseFile(filename, src string) (*ast.File, error) {
	p := NewParser(filename, src)
	return p.ParseFile()
}

type Parser struct {
	filename string
	src      string

	*TokenStream
	file *ast.File
	err  error
}

func NewParser(filename, src string) *Parser {
	return &Parser{filename: filename, src: src}
}

func (p *Parser) ParseFile() (file *ast.File, err error) {
	defer func() {
		if !DebugMode {
			if r := recover(); r != p.err {
				panic(r)
			}
		}
		file, err = p.file, p.err
	}()

	tokens, comments := lexer.Lex(p.filename, p.src)
	for _, tok := range tokens {
		if tok.Type == token.ERROR {
			p.errorf(tok.Pos, "invalid token: %s", tok.Literal)
		}
	}

	p.TokenStream = NewTokenStream(p.filename, p.src, tokens, comments)
	p.parseFile()

	return
}

func (p *Parser) errorf(pos token.Pos, format string, args ...interface{}) {
	_, filename, line := p.callerInfo(1)
	p.err = fmt.Errorf("%s: %s (%s:%d)",
		pos.Position(p.filename, p.src),
		fmt.Sprintf(format, args...),
		filename, line,
	)
	panic(p.err)
}

func (p *Parser) callerInfo(skip int) (fn, filename string, line int) {
	pc, filename, line, _ := runtime.Caller(skip + 1)
	fn = runtime.FuncForPC(pc).Name()
	if idx := strings.LastIndex(fn, "/"); idx >= 0 {
		fn = fn[idx+1:]
	}
	if wd, _ := os.Getwd(); wd != "" {
		if rel, err := filepath.Rel(wd, filename); err == nil {
			filename = rel
		}
	}
	return
}
