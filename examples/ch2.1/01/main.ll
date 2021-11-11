define i32 @main() {
	; 1 + 3 - 2
	%t0 = add i32 0, 1   ; t0 = 1
	%t1 = add i32 %t0, 3 ; t1 = t0 + 3
	%t2 = sub i32 %t1, 2 ; t2 = t1 - 2
	ret i32 %t2
}
