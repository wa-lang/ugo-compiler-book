declare i32 @printf(i8*,...)

@format = constant [4 x i8] c"%d\0A\00"

define i32 @main() {
	; 1 + 3 - 2
	%x1 = add i32 1, 3
	%x2 = sub i32 %x1, 2

	; printf("%d\n", x2)
	call i32(i8*,...) @printf(i8* getelementptr([4 x i8], [4 x i8]* @format, i32 0, i32 0), i32 %x2)

	ret i32 0
}
