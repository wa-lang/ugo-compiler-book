package main

import (
	"fmt"
	"os"
	"os/exec"

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

	app.Action = func(c *cli.Context) error {
		if c.NArg() == 0 {
			fmt.Fprintf(os.Stderr, "no input file")
			os.Exit(1)
		}

		run(c.Args().First())
		return nil
	}

	app.Commands = []*cli.Command{
		{
			Name:  "build",
			Usage: "compile µGo source code",
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					fmt.Fprintf(os.Stderr, "no input file")
					os.Exit(1)
				}
				return nil
			},
		},
		{
			Name:  "run",
			Usage: "compile and run µGo program",
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					fmt.Fprintf(os.Stderr, "no input file")
					os.Exit(1)
				}
				return nil
			},
		},
		{
			Name:  "ast",
			Usage: "parse µGo source code and print ast",
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					fmt.Fprintf(os.Stderr, "no input file")
					os.Exit(1)
				}

				code := loadCode(c.Args().First())
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
				if c.NArg() == 0 {
					fmt.Fprintf(os.Stderr, "no input file")
					os.Exit(1)
				}
				return nil
			},
		},
	}

	app.Run(os.Args)
}

func run(filename string) {
	code := loadCode(filename)
	f, err := parser.ParseFile(filename, code)
	if err != nil {
		panic(err)
	}

	ll := new(compiler.Compiler).Compile(f)
	if err := os.WriteFile("a.out.ll", []byte(ll), 0666); err != nil {
		panic(err)
	}

	exec.Command("clang", "-Wno-override-module", "-o", "a.out", "a.out.ll", "./builtin/_builtin.ll").Run()
	if err := exec.Command("./a.out").Run(); err != nil {
		panic(err)
	}
}

func loadCode(filename string) string {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return string(data)
}
