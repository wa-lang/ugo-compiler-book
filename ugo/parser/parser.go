package parser

import (
	"fmt"

	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/lexer"
	"github.com/chai2010/ugo/logger"
	"github.com/chai2010/ugo/token"
)

type Option struct {
	Debug bool
}

func ParseFile(filename, src string, opt Option) (*ast.File, error) {
	logger.Debugln("ParseFile:", string(src))

	p := &parser{
		filename: filename,
		src:      src,
		opt:      opt,
	}
	p.parseFile()
	return p.file, p.err
}

func ParseBlock(filename, src string, opt Option) (*ast.BlockStmt, error) {
	logger.Debugln("ParseBlock:", string(src))

	p := &parser{
		filename: filename,
		src:      src,
		opt:      opt,
		input:    lexer.Lex(filename, string(src), lexer.Option{}),
	}
	p.block = p.parseBlock()
	return p.block, p.err
}

func ParseExpr(filename, src string, opt Option) (ast.Expr, error) {
	logger.Debugln("ParseExpr:", string(src))

	p := &parser{
		filename: filename,
		src:      src,
		opt:      opt,
		input:    lexer.Lex(filename, string(src), lexer.Option{}),
	}
	p.expr = p.parseExpr()
	return p.expr, p.err
}

type parser struct {
	opt      Option
	filename string
	src      string

	input []token.Token // the tokens being parsed.
	start int           // start position of this item.
	pos   int           // current position in the input.
	width int           // width of last rune read from input.

	file  *ast.File
	block *ast.BlockStmt
	expr  ast.Expr
	err   error
}

func (p *parser) next() token.Token {
	if p.pos >= len(p.input) {
		p.width = 0
		return token.Token{Type: token.EOF}
	}
	tok := p.input[p.pos]
	p.width = 1
	p.pos += p.width
	return tok
}

func (p *parser) peek() token.Token {
	tok := p.next()
	p.backup()
	return tok
}

func (p *parser) peekTokenType() token.TokenType {
	return p.peek().Type
}

func (p *parser) backup() {
	p.pos -= p.width
}

func (p *parser) ignore() {
	p.start = p.pos
}

func (p *parser) accept(validTokens ...token.TokenType) bool {
	tok := p.next()
	for _, x := range validTokens {
		if tok.Type == x {
			return true
		}
	}
	p.backup()
	return false
}

func (p *parser) acceptRun(validTokens ...token.TokenType) {
	for p.accept(validTokens...) {
	}
	p.backup()
}

func (p *parser) errorf(format string, args ...interface{}) {
	p.err = fmt.Errorf(format, args...)
}
