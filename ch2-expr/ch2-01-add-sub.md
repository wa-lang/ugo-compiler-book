# 2.1. 加减法表达式

在前一节我们通过最小编译器将一个整数编译为可以返回相同状态码的程序。现在我们尝试将加法和减法的表达式编译为同样的程序。

比如有 `1+3-2` 表达式，手工编写对应的LLVM汇编程序如下：

```ll
define i32 @main() {
	; 1 + 3 - 2
	%t0 = add i32 0, 1   ; t0 = 1
	%t1 = add i32 %t0, 3 ; t1 = t0 + 3
	%t2 = sub i32 %t1, 2 ; t2 = t1 - 2
	ret i32 %t2
}
```

如果将输入的`1+3-2`转化为`[]string{"1", "+", "3", "-", "2"}` 形式，我们则可以通过以下代码输出对应的汇编程序：

```go
func gen_asm(tokens []string) string {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, `define i32 @main() {`)

	var idx int
	for i, tok := range tokens {
		if i == 0 {
			fmt.Fprintf(&buf, "\t%%t%d = add i32 0, %v\n",
				idx, tokens[i],
			)
			continue
		}
		switch tok {
		case "+":
			idx++
			fmt.Fprintf(&buf, "\t%%t%d = add i32 %%t%d, %v\n",
				idx, idx-1, tokens[i+1],
			)
		case "-":
			idx++
			fmt.Fprintf(&buf, "\t%%t%d = sub i32 %%t%d, %v\n",
				idx, idx-1, tokens[i+1],
			)
		}
	}
	fmt.Fprintf(&buf, "\tret i32 %%t%d\n", idx)
	fmt.Fprintln(&buf, `}`)

	return buf.String()
}
```

而如何将输入的字符串拆分为记号数组本质上属于词法分析的问题。我们先以最简单的方式实现：

```go
func parse_tokens(code string) (tokens []string) {
	for code != "" {
		if idx := strings.IndexAny(code, "+-"); idx >= 0 {
			if idx > 0 {
				tokens = append(tokens, strings.TrimSpace(code[:idx]))
			}
			tokens = append(tokens, code[idx:][:1])
			code = code[idx+1:]
			continue
		}

		tokens = append(tokens, strings.TrimSpace(code))
		return
	}
	return
}
```

基本思路是通过 `strings.IndexAny(code, "+-")` 函数调用根据 `+-` 符号拆分，最终返回拆分后的词法列表。

然后对上个版本的compile函数稍加改造以支持加法和减法的运算表达式编译：

```go
func compile(code string) {
	tokens := parse_tokens(code)
	output := gen_asm(tokens)

	os.WriteFile("a.out.ll", []byte(output), 0666)
	exec.Command("clang", "-Wno-override-module", "-o", "a.out", "a.out.ll").Run()
}
```

为了便于测试，我们再包装一个run函数：

```go
func run(code string) int {
	compile(code)
	if err := exec.Command("./a.out").Run(); err != nil {
		return err.(*exec.ExitError).ExitCode()
	}
	return 0
}
```

run函数将输入的表达式程序编译并运行、最后返回状态码。然后构造单元测试：

```go
func TestRun(t *testing.T) {
	for i, tt := range tests {
		if got := run(tt.code); got != tt.value {
			t.Fatalf("%d: expect = %v, got = %v", i, tt.value, got)
		}
	}
}

var tests = []struct {
	code  string
	value int
}{
	{code: "1", value: 1},
	{code: "1+1", value: 2},
	{code: "1 + 3 - 2  ", value: 2},
	{code: "1+2+3+4", value: 10},
}
```

确认单元测试没有问题后，更新main函数：

```go
func main() {
	code, _ := io.ReadAll(os.Stdin)
	fmt.Println(run(string(code)))
}
```

通过以下命令执行：

```
$ echo "1+2+3" | go run main.go 
6
```
