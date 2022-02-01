default:
	mdbook serve

build:
	-rm book
	mdbook build
	-rm book/.gitignore
	-rm book/.nojekyll
	-rm -rf book/.git
	-rm -rf book/talks

	make build-talks

build-talks:
	mkdir -p ./book/talks
	cp -r ./talks/static-fix-prefix ./book/talks/static
	cp -r ./talks/go-compiler-intro ./book/talks/go-compiler-intro
	cp ./talks/go-compiler-intro.html ./book/talks/go-compiler-intro.html

clean:
