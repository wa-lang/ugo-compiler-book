# 最小编译器

我们先从最小的整数开始，每个整数表示一个返回该值状态码的程序。

比如 0 表示 `os.Exit(0)`。它对应以下的Go程序：

```go
package main;

import "os"

func main() {
	os.Exit(0)
}
```

对应以下的AMD64汇编程序 main_amd64.s：

```s
.intel_syntax noprefix
.globl _main

_main:
	mov rax, 0
	ret
```

入口是`_main`函数，`mov`指令将要返回的0放入rax寄存器，然后调用`ret`指令返回。

可以通过以下命令编译并执行这个汇编程序：

```
$ gcc -o a.out main_amd64.s
$ ./a.out
$ echo $?
0
```

gcc将汇编程序编译为本地可执行程序，然后执行 a.out 程序，最后通过shell的 `echo $?` 命令查看 a.out 的推出状态码。

最小编译器就是将输入的整数翻译为可执行程序的返回该状态码的本地程序：

```go
func main() {
	code, _ := io.ReadAll(os.Stdin)
	compile(string(code))
}

func compile(code string) {
	output := fmt.Sprintf(tmpl, code)
	os.WriteFile("_output_amd64.s", []byte(output), 0666)
	exec.Command("gcc", "-o", "a.out", "_output_amd64.s").Run()
}

const tmpl = `
.intel_syntax noprefix
.globl _main

_main:
	mov rax, %v
	ret
`
```

其中 compile 是编译函数，将从stdin输入的代码先编译为汇编程序，然后调用gcc将汇编程序编译为本地可执行程序（`tmpl`是输出汇编的模板）。

通过以下命令将输入的状态码编译为一个对应的可执行程序：

```
$ echo 123 | go run main.go
$ ./a.out
$ echo $?
123
```

这样我们就实现了一个只能编译整数到本地可执行程序的最小编译器。
