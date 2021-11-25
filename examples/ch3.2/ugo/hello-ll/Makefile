default:
	clang -S -emit-llvm builtin.c
	clang builtin.ll main.ll
	./a.out || echo $$?

clean:
