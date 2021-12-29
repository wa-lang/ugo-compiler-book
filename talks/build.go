package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func main() {
	go startServer()
	time.Sleep(time.Second * 3)

	err := download("http://127.0.0.1:3999/go-compiler-intro.slide", "go-compiler-intro.html")
	if err != nil {
		log.Fatal(err)
	}
}

func startServer() {
	// present -base=. -play=false
	cmd := exec.Command("present", "-base=.", "-play=false")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", stdoutStderr)
}

func download(url, filename string) (err error) {
	fmt.Println("Downloading ", url, " to ", filename)

	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	f, err := os.Create(filename)
	if err != nil {
		return
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return
}
