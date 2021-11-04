.intel_syntax noprefix
.globl _main

_main:
	// 1 + 3 - 2
	mov rax, 1
	add rax, 3
	sub rax, 2
	ret
