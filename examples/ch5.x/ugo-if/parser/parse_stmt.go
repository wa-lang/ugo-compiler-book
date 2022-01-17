package parser

import (
	"github.com/wa-lang/ugo/ast"
	"github.com/wa-lang/ugo/token"
)

func (p *Parser) parseStmt() ast.Stmt {
	switch tok := p.PeekToken(); tok.Type {
	case token.EOF:
		return nil
	case token.ERROR:
		p.errorf(tok.Pos, "invalid token: %s", tok.Literal)
	case token.SEMICOLON:
		p.AcceptTokenList(token.SEMICOLON)
		return nil

	case token.LBRACE: // {}
		return p.parseStmt_block()

	case token.VAR:
		return p.parseStmt_var()

	default:
		// exprList ;
		// exprList := exprList;
		// exprList = exprList;
		exprList := p.parseExprList()
		switch tok := p.PeekToken(); tok.Type {
		case token.SEMICOLON, token.LBRACE:
			if len(exprList) != 1 {
				p.errorf(tok.Pos, "unknown token: %v", tok.Type)
			}
			return &ast.ExprStmt{
				X: exprList[0],
			}
		case token.DEFINE, token.ASSIGN:
			p.ReadToken()
			exprValueList := p.parseExprList()
			if len(exprList) != len(exprValueList) {
				p.errorf(tok.Pos, "unknown token: %v", tok)
			}
			var assignStmt = &ast.AssignStmt{
				Target: make([]*ast.Ident, len(exprList)),
				OpPos:  tok.Pos,
				Op:     tok.Type,
				Value:  make([]ast.Expr, len(exprList)),
			}
			for i, target := range exprList {
				assignStmt.Target[i] = target.(*ast.Ident)
				assignStmt.Value[i] = exprValueList[i]
			}
			return assignStmt
		default:
			p.errorf(tok.Pos, "unknown token: %v", tok)
		}
	}

	panic("unreachable")
}

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
		case token.IF:
			block.List = append(block.List, p.parseStmt_if())
		case token.FOR:
			block.List = append(block.List, p.parseStmt_for())

		default:
			// exprList ;
			// exprList := exprList;
			// exprList = exprList;
			exprList := p.parseExprList()
			switch tok := p.PeekToken(); tok.Type {
			case token.SEMICOLON:
				if len(exprList) != 1 {
					p.errorf(tok.Pos, "unknown token: %v", tok.Type)
				}
				block.List = append(block.List, &ast.ExprStmt{
					X: exprList[0],
				})
			case token.DEFINE, token.ASSIGN:
				p.ReadToken()
				exprValueList := p.parseExprList()
				if len(exprList) != len(exprValueList) {
					p.errorf(tok.Pos, "unknown token: %v", tok)
				}
				var assignStmt = &ast.AssignStmt{
					Target: make([]*ast.Ident, len(exprList)),
					OpPos:  tok.Pos,
					Op:     tok.Type,
					Value:  make([]ast.Expr, len(exprList)),
				}
				for i, target := range exprList {
					assignStmt.Target[i] = target.(*ast.Ident)
					assignStmt.Value[i] = exprValueList[i]
				}
				block.List = append(block.List, assignStmt)
			default:
				p.errorf(tok.Pos, "unknown token: %v", tok)
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
