# 2.3. 解析表达式语法树

要生成语法树需要将代码的字符序列转化为词法序列，然后再将词法序列解析为结构化的语法树。本节将开发一个简易版本的词法解析器，然后在此基础之上开发一个语法解析器，最终产生前一节的语法树。

## 2.3.1 词法解析

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

目前的词法解析器虽然简陋，有了单元测试后面就可以放心重构和优化。词法解析可以参考Rob Pike的报告：https://talks.golang.org/2011/lex.slide。

## 2.3.2 语法定义

语法解析和词法解析输入类似：前者是字符序列、后者是Token序列。输出的结果少有差异：词法解析器产生的是Token扁平的序列，而语法解析产生的是结构化的语法树。

目前依然复用之前的语法树结构：

```go
type ExprNode struct {
	Value string    // +, -, *, /, 123
	Left  *ExprNode `json:",omitempty"`
	Right *ExprNode `json:",omitempty"`
}

func NewExprNode(value string, left, right *ExprNode) *ExprNode {
	return &ExprNode{
		Value: value,
		Left:  left,
		Right: right,
	}
}
```

为了便于JSON打印，我们忽略了空指针。同时增加了NewExprNode构造函数。

在解析语法之前需要明确定义语法规则，下面是EBNF表示的四则运算规则：

```bnf
expr    = mul ("+" mul | "-" mul)*
mul     = primary ("*" primary | "/" primary)*
primary = num | "(" expr ")"
```

可以将EBNF看作是正则表达式的增强版本，其中`|`表示或、`()`表示组合、`*`表示0或多个。比如 expr 规则表示了由 mul 表示的乘法元素再次通过加法或减法组合（隐含了乘法有更高的优先级）。mul 则定义了如何通过 primary 实现乘法或除法组合。而 primary 则表示更小更不可分的数字或小括弧包含的表达式元素。

## 2.3.3 手工递归下降解析

有了EBNF参考之后我们就可以很容易手写一个递归下降的解析程序。手写定义一个parser解析器对象，其中包含词法序列和当前处理的pos位置。

```go
type parser struct {
	tokens []string
	pos    int
}

func (p *parser) peekToken() string {
	if p.pos >= len(p.tokens) {
		return ""
	}
	return p.tokens[p.pos]
}
func (p *parser) nextToken() {
	if p.pos < len(p.tokens) {
		p.pos++
	}
}
```

同时也定义了2个辅助方法：peekToken预取下个元素，nextToken则是移动到下个元素。

然后参考3个规则定义3个有着相同结构的递归函数，每个函数递归构造出语法树：

```go
func (p *parser) build_expr() *ExprNode {
	node := p.build_mul()
	for {
		switch p.peekToken() {
		case "+":
			p.nextToken()
			node = NewExprNode("+", node, p.build_mul())
		case "-":
			p.nextToken()
			node = NewExprNode("-", node, p.build_mul())
		default:
			return node
		}
	}
}
func (p *parser) build_mul() *ExprNode {
	node := p.build_primary()
	for {
		switch p.peekToken() {
		case "*":
			p.nextToken()
			node = NewExprNode("*", node, p.build_primary())
		case "/":
			p.nextToken()
			node = NewExprNode("/", node, p.build_primary())
		default:
			return node
		}
	}
}
func (p *parser) build_primary() *ExprNode {
	if tok := p.peekToken(); tok == "(" {
		p.nextToken()
		node := p.build_expr()
		p.nextToken() // skip ')'
		return node
	} else {
		p.nextToken()
		return NewExprNode(tok, nil, nil)
	}
}
```

然后再包装一个ParseExpr函数：

```go
func ParseExpr(tokens []string) *ExprNode {
	p := &parser{tokens: tokens}
	return p.build_expr()
}
```

然后可以构造一个例子测试语法树的解析：

```go
func main() {
	// 1+2*(3+4)
	expr_tokens := []string{"1", "+", "2", "*", "(", "3", "+", "4", ")"}

	ast := ParseExpr(expr_tokens)
	fmt.Println(JSONString(ast))
}

func JSONString(x interface{}) string {
	d, _ := json.MarshalIndent(x, "", "    ")
	return string(d)
}
```

输出的JSON如下：

```json
{
    "Value": "+",
    "Left": {
        "Value": "1"
    },
    "Right": {
        "Value": "*",
        "Left": {
            "Value": "2"
        },
        "Right": {
            "Value": "+",
            "Left": {
                "Value": "3"
            },
            "Right": {
                "Value": "4"
            }
        }
    }
}
```

然后结合前一节的AST到LLVM的编译函数就可以实现表达式到可执行程序的翻译了。

## 2.3.4 goyacc等工具

其实在Go1.5之前都是基于goyacc工具来产生编译器，在最初的版本我们也提供了基于goyacc的例子（代码还在仓库）。但是对于新手来说，并不推荐goyacc和AntLR等自动生成解析器代码的工具，因此删除了这部分内容。首先是手写解析器对于Go这种语法比较规则的语言并不困难，手写代码不仅仅可以熟悉解析器的工作模式，也可以为错误处理带来更大的灵活性。正如Rob Pike所言，我们也不建议通过goyacc自动生成代码的迂回战术、而是要手写解析器的方式迎头而上解决问题。

