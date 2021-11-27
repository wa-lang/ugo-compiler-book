default:
	make -C ../builtin builtin-ll
	cp ../builtin/_builtin.ll .
	go run hello.go > a.out.ll

	clang -Wno-override-module _builtin.ll a.out.ll
	./a.out || echo $$?

clean:
