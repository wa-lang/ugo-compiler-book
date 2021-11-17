package parser

import (
	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/logger"
	"github.com/chai2010/ugo/token"
)

// func main() {}
// func name(a int, b int) int
// func (p *Type) method() int
// func() (int, int) {}

func (p *parser) parseFunc() *ast.Func {
	tokFunc := p.mustAcceptToken(token.FUNC)

	var funcSpec = &ast.Func{
		FuncPos: tokFunc.Pos,
	}

	switch p.peekTokenType() {
	case token.IDENT:
		p.backupToken()
		p.parseFunc_func(funcSpec)
	case token.LPAREN:
		p.parseFunc_method(funcSpec)
	default:
		p.errorf("invalid token = %v", token.IDENT, p.peekToken())
	}

	return funcSpec
}

func (p *parser) parseFunc_func(fn *ast.Func) {
	p.mustAcceptToken(token.FUNC)

	p.parseFunc_sig_name(fn)
	p.parseFunc_sig_args(fn)
	p.parseFunc_sig_returns(fn)

	// 函数声明
	if p.peekTokenType() == token.LBRACE {
		fn.Body = p.parseStmt_block()
	}

	p.acceptTokenRun(token.SEMICOLON)
}

func (p *parser) parseFunc_method(fn *ast.Func) {
	logger.Debugln("peek =", p.peekToken())

	p.parseFunc_sig_self(fn)
	p.parseFunc_sig_name(fn)
	p.parseFunc_sig_args(fn)
	p.parseFunc_sig_returns(fn)

	// 函数声明
	if p.peekTokenType() == token.LBRACE {
		fn.Body = p.parseStmt_block()
	}

	p.acceptTokenRun(token.SEMICOLON)
}

func (p *parser) parseFunc_closure(fn *ast.Func) {
	logger.Debugln("peek =", p.peekToken())

	p.parseFunc_sig_args(fn)
	p.parseFunc_sig_returns(fn)

	fn.Body = p.parseStmt_block()

	p.acceptTokenRun(token.SEMICOLON)
}

func (p *parser) parseFunc_sig_self(fn *ast.Func) {
	logger.Debugln("peek =", p.peekToken())

	fn.Self = p.parseFunc_sig_field()
}

func (p *parser) parseFunc_sig_name(fn *ast.Func) {
	tokIdent := p.mustAcceptToken(token.IDENT)

	fn.Name = &ast.Ident{
		NamePos: tokIdent.Pos,
		Name:    tokIdent.IdentName(),
	}
}

func (p *parser) parseFunc_sig_args(fn *ast.Func) {
	logger.Debugln("peek =", p.peekToken())

	if _, ok := p.acceptToken(token.LPAREN); !ok {
		p.errorf("invalid token = %v", token.IDENT, p.peekToken())
	}

	for {
		switch p.peekTokenType() {
		case token.RPAREN:
			p.nextToken()
			return
		default:
			if field := p.parseFunc_sig_field(); field != nil {
				fn.Args = append(fn.Args, field)
			} else {
				return
			}
		}
	}
}

func (p *parser) parseFunc_sig_returns(fn *ast.Func) {
	logger.Debugln("peek =", p.peekToken())

	p.acceptToken(token.LPAREN)
	for {
		switch p.peekTokenType() {
		case token.RPAREN, token.RBRACE, token.SEMICOLON:
			p.nextToken()
			return
		default:
			if field := p.parseFunc_sig_field(); field != nil {
				fn.Returns = append(fn.Returns, field)
			} else {
				return
			}
		}
	}
}

func (p *parser) parseFunc_sig_field() (field *ast.Field) {
	logger.Debugln("peek =", p.peekToken())

	switch p.peekTokenType() {
	case token.RPAREN, token.SEMICOLON:
		p.nextToken()
		return nil
	case token.LBRACE:
		return nil
	}

	tokIdent, ok := p.acceptToken(token.IDENT)
	if !ok {
		return nil
	}

	field = &ast.Field{
		Name: &ast.Ident{
			NamePos: tokIdent.Pos,
			Name:    tokIdent.IdentName(),
		},
	}

	if _, ok := p.acceptToken(token.COMMA); ok {
		return field
	}

	switch p.peekTokenType() {
	case token.RPAREN, token.LBRACE, token.SEMICOLON:
		return field
	}

	field.Type = p.parseExpr()
	return field
}
