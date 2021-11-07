
default:
	clang -Wno-override-module main.ll
	-./a.out || echo $$?

	-@rm *.out

clean:
