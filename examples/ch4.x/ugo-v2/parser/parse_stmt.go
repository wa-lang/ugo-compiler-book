package parser

import (
	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/token"
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

		case token.LBRACE: // {}
			block.List = append(block.List, p.parseStmt_block())
		case token.RBRACE: // }
			break Loop

		case token.VAR:
			block.List = append(block.List, p.parseStmt_var())

		default:
			// exprList ;
			// exprList := exprList;
			// exprList = exprList;
			expr := p.parseExpr()
			switch tok := p.PeekToken(); tok.Type {
			case token.SEMICOLON:
				block.List = append(block.List, &ast.ExprStmt{
					X: expr,
				})
			case token.ASSIGN:
				p.ReadToken()
				exprValue := p.parseExpr()
				block.List = append(block.List, &ast.AssignStmt{
					Target: expr.(*ast.Ident),
					OpPos:  tok.Pos,
					Op:     tok.Type,
					Value:  exprValue,
				})

			default:
				panic("aa")
			}

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
