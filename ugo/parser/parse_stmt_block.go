package parser

import (
	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/logger"
	"github.com/chai2010/ugo/token"
)

// var x ...
// x := ...
// x, y = ...
// x, y := fn()

// if expr {} else {}
// for {}
// for expr {}
// for stmt; expr; stmt {}
// return expr?

func (p *parser) parseStmt_block() (block *ast.BlockStmt) {
	logger.Debugln("peek =", p.peekToken())

	block = new(ast.BlockStmt)
	p.mustAcceptToken(token.LBRACE)

Loop:
	for {
		switch tok := p.peekToken(); tok.Type {
		case token.EOF:
			break Loop
		case token.ILLEGAL:
			panic(tok)
		case token.SEMICOLON:
			p.acceptTokenRun(token.SEMICOLON)

		case token.RBRACE:
			break Loop

		case token.CONST:
			_ = p.parseConst()
		case token.TYPE:
			_ = p.parseType()
		case token.VAR:
			_ = p.parseVar()

		default:
			exprs := p.parseExprList()
			if len(exprs) > 1 {
				p.mustAcceptToken(token.ASSIGN, token.DEFINE)
				exprsRight := p.parseExprList()

				block.List = append(block.List, &ast.AssignStmt{
					Target: exprs[0],
					Value:  exprsRight[0],
				})
			} else {
				if _, ok := p.acceptToken(token.ASSIGN, token.DEFINE); ok {
					exprsRight := p.parseExprList()
					block.List = append(block.List, &ast.AssignStmt{
						Target: exprs[0],
						Value:  exprsRight[0],
					})
				} else {
					block.List = append(block.List, &ast.AssignStmt{
						Value: exprs[0],
					})
				}
			}

			logger.Debugln("peek =", p.peekTokenType())
		}
	}

	// parse stmt list

	p.acceptToken(token.RBRACE)
	return
}
