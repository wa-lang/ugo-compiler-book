package main

import (
	"bytes"
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

	<-make(chan bool)
}

func startServer() {
	// present -base=. -play=false
	cmd := exec.Command("present", "-base=.", "-play=false")
	go func() {
		time.Sleep(time.Second * 6)
		cmd.Process.Kill()
		os.Exit(0)
	}()
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

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	data = bytes.TrimSpace(data)
	data = bytes.Replace(data,
		[]byte(`src='/static/slides.js'`),
		[]byte(`src='static/slides.js'`),
		-1,
	)

	err = os.WriteFile(filename, data, 0666)
	if err != nil {
		return
	}
	return
}
