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
	-rm -rf docs/ugo

clean:
