package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/compiler"
	"github.com/chai2010/ugo/parser"
)

func main() {
	app := cli.NewApp()
	app.Name = "ugo"
	app.Usage = "ugo is a tool for managing µGo source code."
	app.Version = "0.0.1"

	app.Commands = []*cli.Command{
		{
			Name:  "build",
			Usage: "compile µGo source code",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
		{
			Name:  "run",
			Usage: "compile and run µGo program",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
		{
			Name:  "ast",
			Usage: "parse µGo source code and print ast",
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					fmt.Fprintf(os.Stderr, "")
					os.Exit(1)
				}

				code := loadCode("./hello.ugo")
				f, err := parser.ParseFile("./hello.ugo", code)
				if err != nil {
					panic(err)
				}

				ast.Print(f)
				return nil
			},
		},
		{
			Name:  "lex",
			Usage: "lex µGo source code and print token list",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
	}

	app.Run(os.Args)
}

func main2() {
	code := loadCode("./hello.ugo")
	f, err := parser.ParseFile("./hello.ugo", code)
	if err != nil {
		panic(err)
	}

	ast.Print(f)

	ll := new(compiler.Compiler).Compile(f)
	//fmt.Print(ll)
	_ = ll
}

func loadCode(filename string) string {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return string(data)
}
