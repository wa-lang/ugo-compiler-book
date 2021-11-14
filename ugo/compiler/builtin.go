package compiler

const builtin_llir = `
declare i32 @printf(i8*,...)

@format_int = constant [4 x i8] c"%d\0A\00"

define i32 @double(i32 %a) {
	%1 = add i32 %a, %a
	ret i32 %1
}

define i32 @println(i32 %x) {
	call i32(i8*,...) @printf(i8* getelementptr([4 x i8], [4 x i8]* @format_int, i32 0, i32 0), i32 %x)
	ret i32 0
}
`
