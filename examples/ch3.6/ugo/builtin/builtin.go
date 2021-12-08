package builtin

const Header = `
declare i32 @ugo_builtin_exit(i32)
`

const MainMain = `
define i32 @main() {
	call i32() @ugo_main_main()
	ret i32 0
}
`
