package parser

import (
	"github.com/wa-lang/ugo/ast"
	"github.com/wa-lang/ugo/token"
)

func (p *Parser) parseStmt_for() *ast.ForStmt {
	tokFor := p.MustAcceptToken(token.FOR)

	forStmt := &ast.ForStmt{
		For: tokFor.Pos,
	}

	// for {}
	if _, ok := p.AcceptToken(token.LBRACE); ok {
		p.UnreadToken()
		forStmt.Body = p.parseStmt_block()
		return forStmt
	}

	// for Cond {}
	// for Init?; Cond?; Post? {}

	// for ; ...
	if _, ok := p.AcceptToken(token.SEMICOLON); ok {
		forStmt.Init = nil

		// for ;; ...
		if _, ok := p.AcceptToken(token.SEMICOLON); ok {
			if _, ok := p.AcceptToken(token.LBRACE); ok {
				// for ;; {}
				p.UnreadToken()
				forStmt.Body = p.parseStmt_block()
				return forStmt
			} else {
				// for ; ; postStmt {}
				forStmt.Post = p.parseStmt()
				forStmt.Body = p.parseStmt_block()
				return forStmt
			}
		} else {
			// for ; cond ; ... {}
			forStmt.Cond = p.parseExpr()
			p.MustAcceptToken(token.SEMICOLON)
			if _, ok := p.AcceptToken(token.LBRACE); ok {
				// for ; cond ; {}
				p.UnreadToken()
				forStmt.Body = p.parseStmt_block()
				return forStmt
			} else {
				// for ; cond ; postStmt {}
				forStmt.Post = p.parseStmt()
				forStmt.Body = p.parseStmt_block()
				return forStmt
			}
		}
	} else {
		stmt := p.parseStmt()

		if _, ok := p.AcceptToken(token.LBRACE); ok {
			// for cond {}
			p.UnreadToken()
			if expr, ok := stmt.(ast.Expr); ok {
				forStmt.Cond = expr
			}
			forStmt.Body = p.parseStmt_block()
			return forStmt
		} else {
			// for init;
			p.MustAcceptToken(token.SEMICOLON)
			forStmt.Init = stmt

			// for ;; ...
			if _, ok := p.AcceptToken(token.SEMICOLON); ok {
				if _, ok := p.AcceptToken(token.LBRACE); ok {
					// for ;; {}
					p.UnreadToken()
					forStmt.Body = p.parseStmt_block()
					return forStmt
				} else {
					// for ; ; postStmt {}
					forStmt.Post = p.parseStmt()
					forStmt.Body = p.parseStmt_block()
					return forStmt
				}
			} else {
				// for ; cond ; ... {}
				forStmt.Cond = p.parseExpr()
				p.MustAcceptToken(token.SEMICOLON)
				if _, ok := p.AcceptToken(token.LBRACE); ok {
					// for ; cond ; {}
					p.UnreadToken()
					forStmt.Body = p.parseStmt_block()
					return forStmt
				} else {
					// for ; cond ; postStmt {}
					forStmt.Post = p.parseStmt()
					forStmt.Body = p.parseStmt_block()
					return forStmt
				}
			}
		}
	}
}
