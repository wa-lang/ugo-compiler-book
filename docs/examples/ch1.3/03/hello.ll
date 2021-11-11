declare i32 @printf(i8*,...)

@format = constant [4 x i8] c"%d\0A\00"

define i32 @main() {
entry:
	%x = alloca i32 ; var x int32 = 0
	store i32 0, i32* %x

	br label %loop
loop:
	%t1 = load i32, i32* %x ; t1 = x
	%t2 = add i32 %t1, 1    ; t2 = t1+1
	store i32 %t2, i32* %x  ; x = t2

	call i32(i8*,...) @printf(i8* getelementptr([4 x i8], [4 x i8]* @format, i32 0, i32 0), i32 %t2)
	br label %loop

	ret i32 0
}
