package runner

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/compiler"
	"github.com/chai2010/ugo/lexer"
	"github.com/chai2010/ugo/parser"
	"github.com/chai2010/ugo/token"
)

type CodeMode string

const (
	CodeMode_Expr CodeMode = "expr"
	CodeMode_File CodeMode = "file"
)

type App struct {
	filename string
	code     string
	codeMode CodeMode
	tokens   []token.Token
	node     ast.Node
	llir     string
}

func NewApp(filename, code string, codeMode CodeMode) *App {
	p := &App{
		filename: filename,
		code:     code,
		codeMode: codeMode,
		tokens:   lexer.Lex(filename, code, lexer.Option{}),
	}
	return p
}

func (p *App) GetTokens() []token.Token {
	return p.tokens
}

func (p *App) GetAST() (ast.Node, error) {
	if p.node != nil {
		return p.node, nil
	}
	switch p.codeMode {
	case CodeMode_Expr:
		node, err := parser.ParseExpr(p.filename, p.code, parser.Option{})
		if err != nil {
			return nil, err
		}
		return node, nil
	case CodeMode_File:
		node, err := parser.ParseFile(p.filename, p.code, parser.Option{})
		if err != nil {
			return nil, err
		}
		return node, nil
	default:
		panic("unknown codeMode: " + p.codeMode)
	}
}

func (p *App) GetLLIR() (string, error) {
	if p.llir != "" {
		return p.llir, nil
	}
	node, err := p.GetAST()
	if err != nil {
		return "", err
	}

	switch p.codeMode {
	case CodeMode_Expr:
		p.llir = new(compiler.Compiler).CompileExpr(node.(ast.Expr))
		return p.llir, nil
	case CodeMode_File:
		p.llir = new(compiler.Compiler).CompileFile(node.(*ast.File))
		return p.llir, nil
	default:
		panic("unknown codeMode: " + p.codeMode)
	}
}

func (p *App) Run() error {
	output, err := p.GetLLIR()
	if err != nil {
		return err
	}

	a_out_ll := "a.out.ll"
	os.WriteFile(a_out_ll, []byte(output), 0666)

	stdoutStderr, err := exec.Command("clang", "-Wno-override-module", "-o", "a.out", "a.out.ll").CombinedOutput()
	if err != nil {
		fmt.Print(string(stdoutStderr))
		return err
	}

	if stdoutStderr, err := exec.Command("./a.out").CombinedOutput(); err != nil {
		fmt.Print(string(stdoutStderr))
		fmt.Println("exit:", err.(*exec.ExitError).ExitCode())
		return err
	} else {
		fmt.Print(string(stdoutStderr))
		fmt.Println("exit:", 0)
		return nil
	}
}
