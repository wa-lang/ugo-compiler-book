package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"

	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/compiler"
	lexpkg "github.com/chai2010/ugo/lexer"
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
			Name:  "run",
			Usage: "compile and run µGo program",
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					fmt.Fprintf(os.Stderr, "no input file")
					os.Exit(1)
				}
				run(c.Args().First())
				return nil
			},
		},
		{
			Name:  "build",
			Usage: "compile µGo source code",
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					fmt.Fprintf(os.Stderr, "no input file")
					os.Exit(1)
				}
				build(c.Args().First())
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

				filename := c.Args().First()

				code := loadCode(filename)
				f, err := parser.ParseFile(filename, code)
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

				filename := c.Args().First()

				code := loadCode(filename)
				lexer := lexpkg.NewLexer(filename, code)

				for i, tok := range lexer.Tokens() {
					fmt.Printf(
						"%02d: %-12v: %-20q // %s\n",
						i, tok.Type, tok.Literal,
						tok.Pos.Position(filename, code),
					)
				}

				fmt.Println("----")

				for i, tok := range lexer.Comments() {
					fmt.Printf(
						"%02d: %-12v: %-20q // %s\n",
						i, tok.Type, tok.Literal,
						tok.Pos.Position(filename, code),
					)
				}
				return nil
			},
		},
	}

	app.Run(os.Args)
}

func run(filename string) {
	build(filename)
	data, err := exec.Command("./a.out").CombinedOutput()
	if len(data) != 0 {
		fmt.Print(string(data))
	}
	if errx, ok := err.(*exec.ExitError); ok {
		os.Exit(errx.ExitCode())
	}
}

func build(filename string) {
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
}

func loadCode(filename string) string {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return string(data)
}
