package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/logger"
	"github.com/chai2010/ugo/runner"
	"github.com/chai2010/ugo/token"
)

var (
	flagFile = flag.String("file", "", "set ugo file")
	flagCode = flag.String("code", "", "set ugo code")
	flagMode = flag.String("mode", "expr", "set ugo code mode(expr|file)")

	flagLex  = flag.Bool("lex", false, "show lex tokens")
	flagAst  = flag.Bool("ast", false, "show ast")
	flagLLIR = flag.Bool("llir", false, "show llvm ir")

	flagDebug = flag.Bool("debug", false, "set debug mode")
)

func init() {
	if strings.Contains(os.Args[0], "go-build") {
		os.Args[0] = "ugo"
	}
}

func main() {
	flag.Parse()
	logger.DebugMode = *flagDebug

	filename := *flagFile
	code := *flagCode

	if filename == "" && code == "" {
		fmt.Printf("ERR: no code")
		os.Exit(1)
	}
	if code == "" {
		data, err := os.ReadFile(filename)
		if err != nil {
			fmt.Printf("ERR: %v", err)
			os.Exit(1)
		}
		code = string(data)
	}
	if filename == "" {
		filename = "_a.out.ugo"
	}

	app := runner.NewApp(filename, code, runner.CodeMode(*flagMode))

	if *flagLex {
		if _, err := os.Lstat(filename); err != nil {
			os.WriteFile(filename, []byte(code), 0666)
		}

		fmt.Println("lex:")
		for i, x := range app.GetTokens() {
			fmt.Printf("\t%03d: %-20v # %v\n", i, x, token.PosString(filename, []byte(code), x.Pos))
		}
	}
	if *flagAst {
		node, _ := app.GetAST()
		fmt.Println("ast:")
		ast.Print(node)
	}

	if *flagLLIR {
		llir, _ := app.GetLLIR()
		fmt.Println("llir:")
		fmt.Println(llir)
	}

	if err := app.Run(); err != nil {
		fmt.Println("ERR:", err)
		os.Exit(1)
	}
}
