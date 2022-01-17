package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/urfave/cli/v2"

	"github.com/wa-lang/ugo/build"
)

func main() {
	app := cli.NewApp()
	app.Name = "ugo"
	app.Usage = "ugo is a tool for managing µGo source code."
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		&cli.StringFlag{Name: "goos", Usage: "set GOOS", Value: runtime.GOOS},
		&cli.StringFlag{Name: "goarch", Usage: "set GOARCH", Value: runtime.GOARCH},
		&cli.StringFlag{Name: "clang", Value: "", Usage: "set clang"},
		&cli.StringFlag{Name: "wasm-llc", Value: "", Usage: "set wasm-llc"},
		&cli.StringFlag{Name: "wasm-ld", Value: "", Usage: "set wasm-ld"},
		&cli.BoolFlag{Name: "debug", Aliases: []string{"d"}, Usage: "set debug mode"},
	}

	app.Action = func(c *cli.Context) error {
		if c.NArg() == 0 {
			fmt.Fprintln(os.Stderr, "no input file")
			os.Exit(1)
		}

		ctx := build.NewContext(build_Options(c))
		data, err := ctx.Run(c.Args().First(), nil)
		if len(data) != 0 {
			fmt.Print(string(data))
		}
		if errx, ok := err.(*exec.ExitError); ok {
			os.Exit(errx.ExitCode())
		}
		if err != nil {
			fmt.Println(err)
		}
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

				ctx := build.NewContext(build_Options(c))
				data, err := ctx.Run(c.Args().First(), nil)
				if len(data) != 0 {
					fmt.Print(string(data))
				}
				if errx, ok := err.(*exec.ExitError); ok {
					os.Exit(errx.ExitCode())
				}
				if err != nil {
					fmt.Println(err)
				}
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

				ctx := build.NewContext(build_Options(c))
				data, err := ctx.Build(c.Args().First(), nil, "")
				if len(data) != 0 {
					fmt.Print(string(data))
				}
				if errx, ok := err.(*exec.ExitError); ok {
					os.Exit(errx.ExitCode())
				}
				if err != nil {
					fmt.Println(err)
				}
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
				code, err := os.ReadFile(filename)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				ctx := build.NewContext(build_Options(c))
				tokens, comments, err := ctx.Lex(c.Args().First(), nil)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				for i, tok := range tokens {
					fmt.Printf(
						"%02d: %-12v: %-20q // %s\n",
						i, tok.Type, tok.Literal,
						tok.Pos.Position(filename, string(code)),
					)
				}

				if len(comments) != 0 {
					fmt.Println("----")
				}

				for i, tok := range comments {
					fmt.Printf(
						"%02d: %-12v: %-20q // %s\n",
						i, tok.Type, tok.Literal,
						tok.Pos.Position(filename, string(code)),
					)
				}
				return nil
			},
		},
		{
			Name:  "ast",
			Usage: "parse µGo source code and print ast",
			Flags: []cli.Flag{
				&cli.BoolFlag{Name: "json", Usage: "output json format"},
			},
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					fmt.Fprintf(os.Stderr, "no input file")
					os.Exit(1)
				}

				ctx := build.NewContext(build_Options(c))
				f, err := ctx.AST(c.Args().First(), nil)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				if c.Bool("json") {
					fmt.Println(f.JSONString())
				} else {
					fmt.Println(f.String())
				}
				return nil
			},
		},
		{
			Name:  "asm",
			Usage: "parse µGo source code and print llvm-ir",
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					fmt.Fprintf(os.Stderr, "no input file")
					os.Exit(1)
				}

				ctx := build.NewContext(build_Options(c))
				ll, err := ctx.ASM(c.Args().First(), nil)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Println(ll)
				return nil
			},
		},
	}

	app.Run(os.Args)
}

func build_Options(c *cli.Context) *build.Option {
	return &build.Option{
		Debug:   c.Bool("debug"),
		GOOS:    c.String("goos"),
		GOARCH:  c.String("goarch"),
		Clang:   c.String("clang"),
		WasmLLC: c.String("wasm-llc"),
		WasmLD:  c.String("wasm-ld"),
	}
}
