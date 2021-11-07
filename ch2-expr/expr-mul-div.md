# 乘除法表达式

这节我们在加减表达式基础之上增加乘除和小括弧的支持。在开始之前我们先对比下加减表达式和可乘除的表达式的本质区别。

因为加减表达式只有一种优先级，可以从左到右依次在当前结果的基础上计算出每一步的结果，因此只需要一个寄存器表示中间结果即可。如果基于有限数量的寄存器进行任意的加减乘除运算则可能需要借助栈保存中间结果。不过LLVM是一个字SSA静态单赋值的抽象汇编语言，其虚拟的寄存器数量是无限的，因此只要处理好乘除法的优先级就可以轻松完成任意的加减乘除表达式运算。

为了简化，我们先假设输入的表达式已经根据优先级被转化为树形结构。比如 `1+2*(3+4)` 对应以下的树形结构：

```
  +
 / \
1   *
   / \
  2   +
     / \
    3   4
```

这个表达式语法树类似普通的二叉树，节点中的值对应表达式的运算符或整数。我们可以单采用以下的结构表示：

```go
type ExprNode struct {
	Value string // +, -, *, /, 123
	Left  *ExprNode
	Right *ExprNode
}

var expr = &ExprNode{
	Value: "+",
	Left: &ExprNode{
		Value: "1",
	},
	Right: &ExprNode{
		Value: "*",
		Left: &ExprNode{
			Value: "2",
		},
		Right: &ExprNode{
			Value: "+",
			Left:  &ExprNode{
				Value: "3",
			},
			Right:  &ExprNode{
				Value: "4",
			},
		},
	},
}
```

现在可以构造针对ExprNode的 Compiler 对象：

```go
type Compiler struct {
	nextId int
}

func (p *Compiler) GenLLIR(node *ExprNode) string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "define i32 @main() {\n")
	fmt.Fprintf(&buf, "    ret i32 %s\n", p.genValue(&buf, node))
	fmt.Fprintf(&buf, "}\n")

	return buf.String()
}

func (p *Compiler) genId() string {
	id := fmt.Sprintf("%%t%d", p.nextId)
	p.nextId++
	return id
}
```

其中 GenLLIR 方法用于将 node 翻译为一个LLVM汇编语言，表达式的终极节点通过`p.genValue(&buf, node)`完成编译。此外还有一个genId辅助方法用于生成唯一的局部虚拟寄存器名字。

genValue 的实现如下：

```go
func (p *Compiler) genValue(w io.Writer, node *ExprNode) (id string) {
	if node == nil {
		return ""
	}
	id = p.genId()
	switch node.Value {
	case "+":
		fmt.Fprintf(w, "\t%s = add i32 %s, %s\n",
			id, p.genValue(w, node.Left), p.genValue(w, node.Right),
		)
	case "-":
		fmt.Fprintf(w, "\t%s = sub i32 %s, %s\n",
			id, p.genValue(w, node.Left), p.genValue(w, node.Right),
		)
	case "*":
		fmt.Fprintf(w, "\t%s = mul i32 %s, %s\n",
			id, p.genValue(w, node.Left), p.genValue(w, node.Right),
		)
	case "/":
		fmt.Fprintf(w, "\t%s = div i32 %s, %s\n",
			id, p.genValue(w, node.Left), p.genValue(w, node.Right),
		)
	default:
		fmt.Fprintf(w, "\t%[1]s = add i32 0, %[2]s; %[1]s = %[2]s\n",
			id, node.Value,
		)
	}
	return
}
```

如果`node.Value`是加减乘除运算符，则递归编译左右子树并产生新的结果，如果不是运算符则作为数值直接返回（通过将数值和0相加产生一个值）。

包装main函数执行对表达式的翻译：

```go
func main() {
	result := run(expr) // 1+2*(3+4)
	fmt.Println(result)
}
```

运行代码将得到以下的LLVM汇编：

```ll
define i32 @main() {
	%t1 = add i32 0, 1; %t1 = 1
	%t3 = add i32 0, 2; %t3 = 2
	%t5 = add i32 0, 3; %t5 = 3
	%t6 = add i32 0, 4; %t6 = 4
	%t4 = add i32 %t5, %t6
	%t2 = mul i32 %t3, %t4
	%t0 = add i32 %t1, %t2
	ret i32 %t0
}
```

这样就完成了表达式树到LLVM汇编的翻译。
