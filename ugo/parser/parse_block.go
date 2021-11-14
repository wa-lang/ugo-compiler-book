package parser

import (
	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/logger"
	"github.com/chai2010/ugo/token"
)

func (p *parser) parseBlock() *ast.BlockStmt {
	logger.Debugln("peek =", p.peekToken())

	// parse stmt list

	p.acceptToken(token.RBRACE)
	return &ast.BlockStmt{}
}
