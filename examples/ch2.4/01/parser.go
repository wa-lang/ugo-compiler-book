package main

func ParseExpr(input []Token) *ExprNode {
	r := NewTokenReader(input)
	return parseExpr(r)
}

func parseExpr(r *TokenReader) *ExprNode {
	return parseExpr_binary(r, 1)
}

func parseExpr_binary(r *TokenReader, prec int) *ExprNode {
	x := parseExpr_unary(r)
	for {
		op := r.PeekToken()
		if op.Type.Precedence() < prec {
			return x
		}

		r.MustAcceptToken(op.Type)
		y := parseExpr_binary(r, op.Type.Precedence()+1)
		x = &ExprNode{Token: op, Left: x, Right: y}
	}
	return nil
}

func parseExpr_unary(r *TokenReader) *ExprNode {
	if _, ok := r.AcceptToken(ADD); ok {
		return parseExpr_primary(r)
	}
	if _, ok := r.AcceptToken(SUB); ok {
		return &ExprNode{
			Token: Token{Type: SUB},
			Left:  &ExprNode{Token: Token{Type: NUMBER, Val: "0"}},
			Right: parseExpr_primary(r),
		}
	}
	return parseExpr_primary(r)
}
func parseExpr_primary(r *TokenReader) *ExprNode {
	if _, ok := r.AcceptToken(LPAREN); ok {
		expr := parseExpr(r)
		r.MustAcceptToken(RPAREN)
		return expr
	}
	return &ExprNode{
		Token: r.MustAcceptToken(NUMBER),
	}
}
