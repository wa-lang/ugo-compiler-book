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

deploy:
	-@make clean
	mnbook build
	-rm book/.gitignore
	-rm -rf book/.git
	-rm -rf book/examples

	cd book && git init
	cd book && git add .
	cd book && git commit -m "first commit"
	cd book && git branch -M gh-pages
	cd book && git remote add origin git@github.com:wa-lang/ugo-compiler-book.git
	cd book && git push -f origin gh-pages

clean:
	-rm -rf ./book

