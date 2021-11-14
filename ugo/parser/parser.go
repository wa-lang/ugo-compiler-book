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

	input []token.Token // the tokens being parsed.
	start int           // start position of this item.
	pos   int           // current position in the input.
	width int           // width of last rune read from input.

	file *ast.File
	expr ast.Expr
	err  error
}

func newParser(filename, src string, opt Option) *parser {
	p := &parser{
		filename: filename,
		src:      src,
		opt:      opt,
		input:    lexer.Lex(filename, string(src), lexer.Option{}),
	}
	return p
}

func (p *parser) ParseFile() (file *ast.File, err error) {
	logger.Debugln(string(p.src))

	defer func() {
		if r := recover(); r != nil {
			if errx, ok := r.(*errors.Error); !ok {
				panic(errx)
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

func (p *parser) nextToken() token.Token {
	if p.pos >= len(p.input) {
		p.width = 0
		return token.Token{Type: token.EOF}
	}
	tok := p.input[p.pos]
	p.width = 1
	p.pos += p.width
	return tok
}

func (p *parser) peekToken() token.Token {
	tok := p.nextToken()
	p.backupToken()
	return tok
}

func (p *parser) peekTokenType() token.TokenType {
	return p.peekToken().Type
}

func (p *parser) backupToken() {
	p.pos -= p.width
}

func (p *parser) ignoreToken() {
	p.start = p.pos
}

func (p *parser) acceptToken(validTokens ...token.TokenType) (token.Token, bool) {
	tok := p.nextToken()
	for _, x := range validTokens {
		if tok.Type == x {
			return tok, true
		}
	}
	p.backupToken()
	return token.Token{}, false
}

func (p *parser) acceptTokenRun(validTokens ...token.TokenType) {
	for {
		if _, ok := p.acceptToken(validTokens...); !ok {
			break
		}
	}
	p.backupToken()
}

func (p *parser) errorf(format string, args ...interface{}) {
	pos := token.PosString(p.filename, []byte(p.src), token.Pos(p.start+1))
	p.err = errors.Newf(pos, format, args...)
	panic(p.err)
}
