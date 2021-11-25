package token

// 记号类型
type TokenType int

// ch3中 µGo程序用到的记号类型
const (
	EOF TokenType = iota

	PACKAGE
	FUNC

	ADD // +
	SUB // -
	MUL // *
	DIV // /
)

// 记号值
type Token struct {
	Type    TokenType   // 记号的类型
	Value   interface{} // 记号的值, 目前只有 int
	Pos     int         // 记号所在的位置(从1开始)
	Literal string      // 程序中原始的字符串
}
