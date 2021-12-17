package builtin

import _ "embed"

//go:embed _builtin.ll
var llBuiltin string

func GetBuiltinLL(goos, goarch string) string {
	switch goos {
	case "darwin":
	case "linux":
	case "windows":
	}
	return llBuiltin
}

const Header = `
declare i32 @ugo_builtin_println(i32)
declare i32 @ugo_builtin_exit(i32)
`

const MainMain = `
define i32 @main() {
	call i32() @ugo_main_init()
	call i32() @ugo_main_main()
	ret i32 0
}
`
