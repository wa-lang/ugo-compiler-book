package token

// 记号类型
type TokenType int

// ch3中 µGo程序用到的记号类型
const (
	EOF TokenType = iota
	ERROR

	IDENT
	NUMBER

	PACKAGE
	FUNC

	ADD // +
	SUB // -
	MUL // *
	DIV // /

	LPAREN // (
	RPAREN // )
	LBRACE // {
	RBRACE // }

	SEMICOLON // ;
)

// 记号值
type Token struct {
	Type    TokenType   // 记号的类型
	Value   interface{} // 记号的值, 目前只有 int
	Pos     int         // 记号所在的位置(从1开始)
	Literal string      // 程序中原始的字符串
}

var keywords = map[string]TokenType{
	"package": PACKAGE,
	"func":    FUNC,
}

func Lookup(ident string) TokenType {
	if tok, is_keyword := keywords[ident]; is_keyword {
		return tok
	}
	return IDENT
}
