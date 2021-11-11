package main

// expr    = mul ("+" mul | "-" mul)*
// mul     = primary ("*" primary | "/" primary)*
// primary = num | "(" expr ")"

func ParseExpr(tokens []string) *ExprNode {
	p := &parser{tokens: tokens}
	return p.build_expr()
}

type parser struct {
	tokens []string
	pos    int
}

func (p *parser) build_expr() *ExprNode {
	node := p.build_mul()
	for {
		switch p.peekToken() {
		case "+":
			p.nextToken()
			node = NewExprNode("+", node, p.build_mul())
		case "-":
			p.nextToken()
			node = NewExprNode("-", node, p.build_mul())
		default:
			return node
		}
	}
}
func (p *parser) build_mul() *ExprNode {
	node := p.build_primary()
	for {
		switch p.peekToken() {
		case "*":
			p.nextToken()
			node = NewExprNode("*", node, p.build_primary())
		case "/":
			p.nextToken()
			node = NewExprNode("/", node, p.build_primary())
		default:
			return node
		}
	}
}
func (p *parser) build_primary() *ExprNode {
	if tok := p.peekToken(); tok == "(" {
		p.nextToken()
		node := p.build_expr()
		p.nextToken() // skip ')'
		return node
	} else {
		p.nextToken()
		return NewExprNode(tok, nil, nil)
	}
}

func (p *parser) peekToken() string {
	if p.pos >= len(p.tokens) {
		return ""
	}
	return p.tokens[p.pos]
}
func (p *parser) nextToken() {
	if p.pos < len(p.tokens) {
		p.pos++
	}
}
