default:
	make -C ./builtin builtin-ll
	go run main.go > a.out.ll

	clang -Wno-override-module ./a.out.ll ./builtin/_builtin.ll
	./a.out || echo $$?

clean:
