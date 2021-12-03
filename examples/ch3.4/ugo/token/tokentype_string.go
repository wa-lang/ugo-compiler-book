// Code generated by "stringer -type=TokenType"; DO NOT EDIT.

package token

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[EOF-0]
	_ = x[ERROR-1]
	_ = x[COMMENT-2]
	_ = x[IDENT-3]
	_ = x[NUMBER-4]
	_ = x[PACKAGE-5]
	_ = x[FUNC-6]
	_ = x[ADD-7]
	_ = x[SUB-8]
	_ = x[MUL-9]
	_ = x[DIV-10]
	_ = x[LPAREN-11]
	_ = x[RPAREN-12]
	_ = x[LBRACE-13]
	_ = x[RBRACE-14]
	_ = x[SEMICOLON-15]
}

const _TokenType_name = "EOFERRORCOMMENTIDENTNUMBERPACKAGEFUNCADDSUBMULDIVLPARENRPARENLBRACERBRACESEMICOLON"

var _TokenType_index = [...]uint8{0, 3, 8, 15, 20, 26, 33, 37, 40, 43, 46, 49, 55, 61, 67, 73, 82}

func (i TokenType) String() string {
	if i < 0 || i >= TokenType(len(_TokenType_index)-1) {
		return "TokenType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}