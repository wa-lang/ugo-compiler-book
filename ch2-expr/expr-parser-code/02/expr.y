%{
package main

import "fmt"

var _ = fmt.Sprint
%}

%union {
	node *ExprNode
}

%token <node> NUM
%token '+' '-' '*' '/' '(' ')'

%type <node> expr mul primary

%%

// expr    = mul ("+" mul | "-" mul)*
// mul     = primary ("*" primary | "/" primary)*
// primary = num | "(" expr ")"

top: expr { yyrcvr.lval.node = $1 }

expr: mul { $$ = $1 }
	| expr '+' mul { $$ = NewExprNode("+", $1, $3) }
	| expr '-' mul { $$ = NewExprNode("-", $1, $3) }

mul: primary { $$ = $1 }
	| mul '*' primary { $$ = NewExprNode("*", $1, $3) }
	| mul '/' primary { $$ = NewExprNode("/", $1, $3) }

primary: NUM { $$ = $1 }
	| '(' expr ')' { $$ = $2 }

%%

type exprLex struct {
	tokens []string
	pos    int
}

func (p *exprLex) read() (tok string) {
	if p.pos >= len(p.tokens) {
		return ""
	}
	tok = p.tokens[p.pos]
	p.pos++
	return
}

func (p *exprLex) Lex(yylval *yySymType) int {
	switch s := p.read(); s {
	case "+", "-", "*", "/", "(", ")":
		return int(s[0])
	default:
		if s != "" {
			yylval.node = NewExprNode(s, nil, nil)
			return NUM
		}
		return 0
	}
}

func (x *exprLex) Error(s string) {
	panic(s)
}

func ParseExpr(tokens []string) *ExprNode {
	parser := yyNewParser().(*yyParserImpl)
	parser.Parse(&exprLex{tokens:tokens})
	return parser.lval.node
}
