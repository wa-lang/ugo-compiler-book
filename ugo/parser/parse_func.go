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
	tokFunc := p.r.MustAcceptToken(token.FUNC)

	var funcSpec = &ast.Func{
		FuncPos: tokFunc.Pos,
	}

	switch tok := p.r.PeekToken(); tok.Type {
	case token.IDENT:
		p.r.UnreadToken()
		p.parseFunc_func(funcSpec)
	case token.LPAREN:
		p.parseFunc_method(funcSpec)
	default:
		p.errorf(tok.Pos, "invalid token = %v", tok)
	}

	return funcSpec
}

func (p *parser) parseFunc_func(fn *ast.Func) {
	p.r.MustAcceptToken(token.FUNC)

	p.parseFunc_sig_name(fn)
	p.parseFunc_sig_args(fn)
	p.parseFunc_sig_returns(fn)

	// 函数声明
	if p.r.PeekToken().Type == token.LBRACE {
		fn.Body = p.parseStmt_block()
	}

	p.r.AcceptTokenList(token.SEMICOLON)
}

func (p *parser) parseFunc_method(fn *ast.Func) {
	logger.Debugln("peek =", p.r.PeekToken())

	p.parseFunc_sig_self(fn)
	p.parseFunc_sig_name(fn)
	p.parseFunc_sig_args(fn)
	p.parseFunc_sig_returns(fn)

	// 函数声明
	if p.r.PeekToken().Type == token.LBRACE {
		fn.Body = p.parseStmt_block()
	}

	p.r.AcceptTokenList(token.SEMICOLON)
}

func (p *parser) parseFunc_closure(fn *ast.Func) {
	logger.Debugln("peek =", p.r.PeekToken())

	p.parseFunc_sig_args(fn)
	p.parseFunc_sig_returns(fn)

	fn.Body = p.parseStmt_block()

	p.r.AcceptTokenList(token.SEMICOLON)
}

func (p *parser) parseFunc_sig_self(fn *ast.Func) {
	logger.Debugln("peek =", p.r.PeekToken())

	fn.Self = p.parseFunc_sig_field()
}

func (p *parser) parseFunc_sig_name(fn *ast.Func) {
	tokIdent := p.r.MustAcceptToken(token.IDENT)

	fn.Name = &ast.Ident{
		NamePos: tokIdent.Pos,
		Name:    tokIdent.IdentName(),
	}
}

func (p *parser) parseFunc_sig_args(fn *ast.Func) {
	logger.Debugln("peek =", p.r.PeekToken())

	if tok, ok := p.r.AcceptToken(token.LPAREN); !ok {
		p.errorf(tok.Pos, "invalid token = %v", tok)
	}

	for {
		switch p.r.PeekToken().Type {
		case token.RPAREN:
			p.r.ReadToken()
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
	logger.Debugln("peek =", p.r.PeekToken())

	p.r.AcceptToken(token.LPAREN)
	for {
		switch p.r.PeekToken().Type {
		case token.RPAREN, token.RBRACE, token.SEMICOLON:
			p.r.ReadToken()
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
	logger.Debugln("peek =", p.r.PeekToken())

	switch p.r.PeekToken().Type {
	case token.RPAREN, token.SEMICOLON:
		p.r.ReadToken()
		return nil
	case token.LBRACE:
		return nil
	}

	tokIdent, ok := p.r.AcceptToken(token.IDENT)
	if !ok {
		return nil
	}

	field = &ast.Field{
		Name: &ast.Ident{
			NamePos: tokIdent.Pos,
			Name:    tokIdent.IdentName(),
		},
	}

	if _, ok := p.r.AcceptToken(token.COMMA); ok {
		return field
	}

	switch p.r.PeekToken().Type {
	case token.RPAREN, token.LBRACE, token.SEMICOLON:
		return field
	}

	field.Type = p.parseExpr()
	return field
}
