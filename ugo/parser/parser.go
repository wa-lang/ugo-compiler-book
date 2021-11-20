package parser

import (
	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/errors"
	"github.com/chai2010/ugo/lexer"
	"github.com/chai2010/ugo/logger"
	"github.com/chai2010/ugo/token"
)

type Option struct {
	Debug bool
}

func ParseFile(filename, src string, opt Option) (*ast.File, error) {
	p := newParser(filename, src, opt)
	return p.ParseFile()
}

func ParseExpr(filename, src string, opt Option) (ast.Expr, error) {
	p := newParser(filename, src, opt)
	return p.ParseExpr()
}

type parser struct {
	opt      Option
	filename string
	src      string
	r        *token.Reader

	file *ast.File
	expr ast.Expr
	err  error
}

func newParser(filename, src string, opt Option) *parser {
	p := &parser{
		filename: filename,
		src:      src,
		opt:      opt,
		r:        token.NewReader(lexer.Lex(filename, string(src), lexer.Option{})),
	}
	return p
}

func (p *parser) ParseFile() (file *ast.File, err error) {
	logger.Debugln(string(p.src))

	defer func() {
		if !logger.DebugMode {
			if r := recover(); r != nil {
				if errx, ok := r.(*errors.Error); !ok {
					panic(errx)
				}
			}
		}
		file, err = p.file, p.err
	}()

	p.parseFile()
	return
}

func (p *parser) ParseExpr() (expr ast.Expr, err error) {
	logger.Debugln(string(p.src))

	defer func() {
		if r := recover(); r != nil {
			if errx, ok := r.(*errors.Error); !ok {
				panic(errx)
			}
		}
		expr, err = p.expr, p.err
	}()

	p.expr = p.parseExpr()
	return
}

func (p *parser) errorf(pos token.Pos, format string, args ...interface{}) {
	p.err = errors.Newf(
		token.PosString(p.filename, []byte(p.src), pos),
		format, args...,
	)
	panic(p.err)
}
