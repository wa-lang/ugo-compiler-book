# 可乘除的表达式

这节我们在加减表达式基础之上增加乘除和小括弧的支持。在开始之前我们先对比下加减表达式和可乘除的表达式的本质区别。

加减表达式`1+2+3`可以翻译为以下汇编：

```s
// 1 + 2 + 3
mov rax, 1 // rax = 1
add rax, 2 // rax = rax + 2
sub rax, 3 // rax = rax + 3
```

因为加减表达式只有一种优先级，可以从左到右依次在当前结果的基础上计算出每一步的结果，因此只需要一个 `rax` 寄存器表示中间结果即可。

但是乘除比加减法有着更高的优先级，如果载结合小括弧将出现层次结构，如果是从左到右依次计算则需要保存更多更多的中间结果。我们可以借助栈来保存中间结果：




比如 `1+2*(3+4)` 对应以下的树形结构：

```
  +
 / \
1   *
   / \
  2   +
     / \
    3   4
```

先简单采用以下的结构表示：

```go
type ExprNode struct {
	Left  interface{} // num, *Node
	Op    string      // +-*/
	Right interface{}
}

var expr = &ExprNode{
	Left: 1, Op: "+",
	Right: &ExprNode{
		Left: 2, Op: "*",
		Right: &ExprNode{
			Left: 3, Op: "+",
			Right: 4,
		},
	},
}
```

