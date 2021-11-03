default:
	mdbook serve

build:
	-rm docs
	mdbook build
	-rm docs/.gitignore
	-rm docs/.nojekyll

clean:
