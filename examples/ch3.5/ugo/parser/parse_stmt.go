package parser

import (
	"github.com/wa-lang/ugo/ast"
	"github.com/wa-lang/ugo/token"
)

func (p *Parser) parseStmt_block() *ast.BlockStmt {
	block := &ast.BlockStmt{}

	tokBegin := p.MustAcceptToken(token.LBRACE) // {

Loop:
	for {
		switch tok := p.PeekToken(); tok.Type {
		case token.EOF:
			break Loop
		case token.ERROR:
			p.errorf(tok.Pos, "invalid token: %s", tok.Literal)
		case token.SEMICOLON:
			p.AcceptTokenList(token.SEMICOLON)

		case token.RBRACE: // }
			break Loop

		default:
			block.List = append(block.List, p.parseStmt_expr())
		}
	}

	tokEnd := p.MustAcceptToken(token.RBRACE) // }

	block.Lbrace = tokBegin.Pos
	block.Rbrace = tokEnd.Pos

	return block
}

func (p *Parser) parseStmt_expr() *ast.ExprStmt {
	return &ast.ExprStmt{
		X: p.parseExpr(),
	}
}
