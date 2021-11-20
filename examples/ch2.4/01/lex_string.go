package main

func (t TokenType) String() string {
	switch t {
	case EOF:
		return "EOF"
	case ILLEGAL:
		return "ILLEGAL"

	case NUMBER:
		return "NUMBER"

	case ADD:
		return "ADD"
	case SUB:
		return "SUB"
	case MUL:
		return "MUL"
	case DIV:
		return "DIV"

	case LPAREN:
		return "LPAREN"
	case RPAREN:
		return "RPAREN"

	default:
		return "UNKNWON"
	}
}
