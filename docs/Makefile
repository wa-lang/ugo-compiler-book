default:
	mdbook serve

build:
	-rm docs
	mdbook build
	-rm docs/.gitignore
	-rm docs/.nojekyll
	-rm -rf docs/.git
	-rm -rf docs/docs
	-rm -rf docs/examples
	-rm -rf docs/talks
	-rm -rf docs/ugo

	make build-talks

build-talks:
	mkdir -p ./docs/talks
	cp -r ./talks/static ./docs/static
	cp -r ./talks/go-compiler-intro ./docs/talks/go-compiler-intro
	cp ./talks/go-compiler-intro.html ./docs/talks/go-compiler-intro.html

clean:
