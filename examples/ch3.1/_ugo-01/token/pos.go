package token

import (
	"fmt"
	gotoken "go/token"
)

// Pos 类似一个指针, 表示文件中的位置.
type Pos = gotoken.Pos

// NoPos 类似指针的 nil 值, 表示一个无效的位置.
const NoPos Pos = 0

// PosString 格式化 pos 为字符串
func PosString(filename string, src []byte, pos Pos) string {
	fset := gotoken.NewFileSet()
	fset.AddFile(filename, 1, len(src)).SetLinesForContent(src)
	return fmt.Sprintf("%v", fset.Position(pos))
}
