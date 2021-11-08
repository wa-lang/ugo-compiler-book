# 解析表达式语法树

要生成语法树需要将代码的字符序列转化为词法序列，然后再将词法序列解析为结构化的语法树。本节将开发一个简易版本的词法解析器，然后在此基础之上开发一个语法解析器，最终产生前一节的语法树。

## 词法解析

在加法和减法表达式的例子我们已经实现了一个词法解析器，现在我们继续添加对`*/()`的解析支持。词法解析器一般叫lexer，我们新的函数改名为lex：

```go
func Lex(code string) (tokens []string) {
	for code != "" {
		if idx := strings.IndexAny(code, "+-*/()"); idx >= 0 {
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

其中`strings.IndexAny`增加了乘除法和小括弧的支持。目前我们暂时忽略错误的输入。开发调试的同时可以添加测试代码，如下：

```go
func TestLex(t *testing.T) {
	var tests = []struct {
		input  string
		tokens []string
	}{
		{"1", []string{"1"}},
		{"1+22*333", []string{"1", "+", "22", "*", "333"}},
		{"1+22*(3+4)", []string{"1", "+", "22", "*", "(", "3", "+", "4", ")"}},
	}
	for i, tt := range tests {
		if got := Lex(tt.input); !reflect.DeepEqual(got, tt.tokens) {
			t.Fatalf("%d: expect = %v, got = %v", i, tt.tokens, got)
		}
	}
}
```

目前的词法解析器虽然简陋，有了单元测试后面就可以放心重构和优化。

## 语法解析

语法解析和词法解析输入类似：前者是字符序列、后者是Token序列。输出的结果少有差异：词法解析器产生的是Token扁平的序列，而语法解析产生的是结构化的语法树。

目前依然复用之前的语法树结构：

```go
type ExprNode struct {
	Value string // +, -, *, /, 123
	Left  *ExprNode
	Right *ExprNode
}
```

然后定义语法解析函数：

```go
func ParseExpr(code string) *ExprNode {
	panic("TODO")
}
```

TODO