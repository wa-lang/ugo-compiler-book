package parser

import (
	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/token"
)

func (p *Parser) parseStmt_block() *ast.BlockStmt {
	p.MustAcceptToken(token.LBRACE)       // {
	defer p.MustAcceptToken(token.RBRACE) // }

Loop:
	for {
		switch tok := p.PeekToken(); tok.Type {
		case token.EOF:
			break Loop
		case token.ERROR:
			p.errorf(tok.Pos, "invalid token: %s", tok.Literal)
		case token.SEMICOLON:
			p.AcceptTokenList(token.SEMICOLON)

		default:
			// p.parseStmt_expr(block)
		}
	}

	panic("TODO")
}

func (p *Parser) parseStmt_expr(block *ast.BlockStmt) {
	// todo
}

/*

func (p *parser) parseStmt_block() (block *ast.BlockStmt) {
	logger.Debugln("peek =", p.r.PeekToken())

	block = new(ast.BlockStmt)
	p.r.MustAcceptToken(token.LBRACE)

Loop:
	for {
		switch tok := p.r.PeekToken(); tok.Type {
		case token.EOF:
			break Loop
		case token.ILLEGAL:
			panic(tok)
		case token.SEMICOLON:
			p.r.AcceptTokenList(token.SEMICOLON)

		case token.RBRACE:
			break Loop

		case token.CONST:
			block.List = append(block.List, p.parseStmt_const())
		case token.TYPE:
			block.List = append(block.List, p.parseStmt_type())
		case token.VAR:
			block.List = append(block.List, p.parseStmt_var())

		case token.DEFER:
			block.List = append(block.List, p.parseStmt_defer())
		case token.IF:
			block.List = append(block.List, p.parseStmt_if())
		case token.FOR:
			block.List = append(block.List, p.parseStmt_for())
		case token.RETURN:
			block.List = append(block.List, p.parseStmt_return())

		default:
			p.parseStmt_assign(block)
		}
	}

	// parse stmt list

	p.r.AcceptToken(token.RBRACE)
	return
}

*/
