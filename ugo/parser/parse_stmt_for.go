package parser

import (
	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/token"
)

func (p *parser) parseStmt_for() *ast.ForStmt {
	tokFor := p.r.MustAcceptToken(token.FOR)

	forStmt := &ast.ForStmt{
		For: tokFor.Pos,
	}

	// for {}
	if _, ok := p.r.AcceptToken(token.LBRACE); ok {
		p.r.UnreadToken()
		forStmt.Body = p.parseStmt_block()
		return forStmt
	}

	// for Cond {}
	// for Init?; Cond?; Post? {}

	// for ; ...
	if _, ok := p.r.AcceptToken(token.SEMICOLON); ok {
		forStmt.Init = nil

		// for ;; ...
		if _, ok := p.r.AcceptToken(token.SEMICOLON); ok {
			if _, ok := p.r.AcceptToken(token.LBRACE); ok {
				// for ;; {}
				p.r.UnreadToken()
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
			p.r.MustAcceptToken(token.SEMICOLON)
			if _, ok := p.r.AcceptToken(token.LBRACE); ok {
				// for ; cond ; {}
				p.r.UnreadToken()
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

		if _, ok := p.r.AcceptToken(token.LBRACE); ok {
			// for cond {}
			p.r.UnreadToken()
			if expr, ok := stmt.(ast.Expr); ok {
				forStmt.Cond = expr
			}
			forStmt.Body = p.parseStmt_block()
			return forStmt
		} else {
			// for init;
			p.r.MustAcceptToken(token.SEMICOLON)

			// for ;; ...
			if _, ok := p.r.AcceptToken(token.SEMICOLON); ok {
				if _, ok := p.r.AcceptToken(token.LBRACE); ok {
					// for ;; {}
					p.r.UnreadToken()
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
				p.r.MustAcceptToken(token.SEMICOLON)
				if _, ok := p.r.AcceptToken(token.LBRACE); ok {
					// for ; cond ; {}
					p.r.UnreadToken()
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
