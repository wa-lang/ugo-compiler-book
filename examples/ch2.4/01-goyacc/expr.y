%{
package main

import "fmt"

var _ = fmt.Sprint
%}

%union {
	node *ExprNode
	tok  Token
}

%token ILLEGAL

%token <tok> NUMBER

%token <tok> ADD // +
%token <tok> SUB // -
%token <tok> MUL // *
%token <tok> DIV // /

%token LPAREN // (
%token RPAREN // )

%type <node> expr mul primary

%%

// expr    = mul ("+" mul | "-" mul)*
// mul     = primary ("*" primary | "/" primary)*
// primary = num | "(" expr ")"

top: expr { yyrcvr.lval.node = $1 }

expr: mul { $$ = $1 }
	| expr ADD mul { $$ = NewExprNode($2, $1, $3) }
	| expr SUB mul { $$ = NewExprNode($2, $1, $3) }

mul: primary { $$ = $1 }
	| mul MUL primary { $$ = NewExprNode($2, $1, $3) }
	| mul DIV primary { $$ = NewExprNode($2, $1, $3) }

primary: NUMBER { $$ = NewExprNode($1, nil, nil) }
	| '(' expr ')' { $$ = $2 }

%%

const EOF = 0

type exprLex struct {
	tokens []Token
	pos    int
}

func (p *exprLex) Lex(yylval *yySymType) int {
	if p.pos >= len(p.tokens) {
		return EOF
	}
	
	yylval.tok = p.tokens[p.pos]
	p.pos++

	return int(yylval.tok.Type)
}

func (x *exprLex) Error(s string) {
	panic(s)
}

func ParseExpr(tokens []Token) *ExprNode {
	parser := yyNewParser().(*yyParserImpl)
	parser.Parse(&exprLex{tokens:tokens})
	return parser.lval.node
}
