package main

import "fmt"

func main() {
	fmt.Printf("%s%c%s%c\n", q, 0x60, q, 0x60) // HL
}

var q = `package main

import "fmt"

func main() {
	fmt.Printf("%s%c%s%c\n", q, 0x60, q, 0x60) // HL
}

var q = `
