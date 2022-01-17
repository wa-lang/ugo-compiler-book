package parser

import (
	"fmt"

	"github.com/wa-lang/ugo/ast"
	"github.com/wa-lang/ugo/lexer"
	"github.com/wa-lang/ugo/token"
)

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
		if r := recover(); r != p.err {
			panic(r)
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
	p.err = fmt.Errorf("%s: %s",
		pos.Position(p.filename, p.src),
		fmt.Sprintf(format, args...),
	)
	panic(p.err)
}
