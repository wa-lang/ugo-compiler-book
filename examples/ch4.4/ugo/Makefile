default:
	go run main.go -goos=wasm -debug=false asm ./_examples/hello.ugo
	go run main.go \
		-goos=wasm -debug=true \
		-wasm-llc=/usr/local/Cellar/llvm/9.0.0/bin/llc \
		-wasm-ld=/usr/local/Cellar/llvm/9.0.0/bin/wasm-ld \
		build ./_examples/hello.ugo

	wasm2wat a.out.wasm
	node run_wasm.js

clean:
