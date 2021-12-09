package compileopts

import (
	"runtime"
)

func TargetTriple(goos, goarch string) string {
	if goos == "" {
		goos = runtime.GOOS
	}
	if goarch == "" {
		goarch = runtime.GOARCH
	}

	var llvmarch string
	var llvmos string

	switch goarch {
	case "386":
		llvmarch = "i386"
	case "amd64":
		llvmarch = "x86_64"
	case "arm64":
		llvmarch = "aarch64"
	default:
		panic("unsupport: " + goos + "/" + goarch)
	}

	switch goos {
	case "darwin":
		llvmos = "macosx10.12.0"
		if goarch == "arm64" {
			llvmarch = "arm64"
		}
	default:
		llvmos = goos
	}

	target := llvmarch + "-unknown-" + llvmos
	if goos == "windows" {
		target += "-gnu"
	}

	return target
}
