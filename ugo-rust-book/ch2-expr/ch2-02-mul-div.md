# 2.2. 乘除法表达式

这节我们在加减表达式基础之上增加乘除和小括弧的支持。在开始之前我们先对比下加减表达式和可乘除的表达式的本质区别。

因为加减表达式只有一种优先级，可以从左到右依次在当前结果的基础上计算出每一步的结果，因此只需要一个寄存器表示中间结果即可。如果基于有限数量的寄存器进行任意的加减乘除运算则可能需要借助栈保存中间结果。不过LLVM是一个SSA静态单赋值的抽象汇编语言，其虚拟的寄存器数量是无限的，因此只要处理好乘除法的优先级就可以轻松完成任意的加减乘除表达式运算。

## 2.2.1 测试数据

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

这个表达式语法树类似普通的二叉树，节点中的值对应表达式的运算符或整数。

## 2.2.2 定义语法树

我们先在`src/ast.rs`文件定义表达式树结构：

```rust,noplayground
// src/ast.rs
pub struct ExprNode<'a> {
    pub value: &'a str, // +, -, *, /, 123
    pub left: Option<Box<ExprNode<'a>>>,
    pub right: Option<Box<ExprNode<'a>>>,
}
```

然后在 `src/main.rs` 文件的 main 函数构造以上表达式树：

```rust,noplayground
#![allow(unused)]

mod ast;

fn main() {
    // 1+2*(3+4)
    let node = crate::ast::ExprNode {
        value: "+",
        left: Some(Box::new(crate::ast::ExprNode {
            value: "1",
            left: None,
            right: None,
        })),
        right: Some(Box::new(crate::ast::ExprNode {
            value: "*",
            left: Some(Box::new(crate::ast::ExprNode {
                value: "2",
                left: None,
                right: None,
            })),
            right: Some(Box::new(crate::ast::ExprNode {
                value: "+",
                left: Some(Box::new(crate::ast::ExprNode {
                    value: "3",
                    left: None,
                    right: None,
                })),
                right: Some(Box::new(crate::ast::ExprNode {
                    value: "4",
                    left: None,
                    right: None,
                })),
            })),
        })),
    };

    ...
}
```

## 2.2.3 构造 Compiler 对象

现在可以在 `sc/compiler.rs` 文件中构造针对 `ExprNode` 的 Compiler 对象：

```rust,noplayground
// sc/compiler.rs

pub struct Compiler {
    next_id: i32,
}

impl Compiler {
    pub fn new() -> Self {
        Compiler { next_id: 0 }
    }
}
```

然后构造 `Compiler.gen_llir` 方法，用于将语法树编译为LLIR汇编程序：

```rust,noplayground
impl Compiler {
    pub fn gen_llir(&mut self, node: &crate::ast::ExprNode) -> String {
        let mut result = String::new();
        result.push_str("define i32 @main() {\n");
        let v = self.gen_value(&mut result, node);
        result.push_str(&format!("\tret i32 {}\n", v));
        result.push_str("}\n");
        result
    }
}
```

其中 `Compiler.gen_value` 方法用于遍历每个内部节点，并返回结果：

```rust,noplayground
impl Compiler {
    fn gen_value(&mut self, output: &mut String, node: &crate::ast::ExprNode) -> String {
        let id = self.gen_id();
        match &node.value {
            &"+" => {
                let x = self.gen_value(output, node.left.as_ref().unwrap());
                let y = self.gen_value(output, node.right.as_ref().unwrap());
                output.push_str(&format!("\t{} = add i32 {}, {}\n", id, x, y));
                return id;
            }
            &"-" => {
                let x = self.gen_value(output, node.left.as_ref().unwrap());
                let y = self.gen_value(output, node.right.as_ref().unwrap());
                output.push_str(&format!("\t{} = sub i32 {}, {}\n", id, x, y));
                return id;
            }
            &"*" => {
                let x = self.gen_value(output, node.left.as_ref().unwrap());
                let y = self.gen_value(output, node.right.as_ref().unwrap());
                output.push_str(&format!("\t{} = mul i32 {}, {}\n", id, x, y));
                return id;
            }
            &"/" => {
                let x = self.gen_value(output, node.left.as_ref().unwrap());
                let y = self.gen_value(output, node.right.as_ref().unwrap());
                output.push_str(&format!("\t{} = sdiv i32 {}, {}\n", id, x, y));
                return id;
            }
            _ => {
                output.push_str(&format!("\t{} = add i32 0, {}\n", id, node.value));
                return id;
            }
        }
    }
}
```

如果`node.value`是加减乘除运算符，则递归编译左右子树并产生新的结果，如果不是运算符则作为数值直接返回（通过将数值和0相加产生一个值）。


最后是辅助产生内部变量名的辅助方法：

```rust,noplayground
impl Compiler {
    fn gen_id(&mut self) -> String {
        let s = format!("%t{}", self.next_id);
        self.next_id += 1;
        s
    }
}
```

genId辅助方法用于生成唯一的局部虚拟寄存器名字。

## 2.2.4 编译表达式树

包装main函数执行对表达式的翻译：

```rust,noplayground
fn main() {
    // 1+2*(3+4)
    let node = crate::ast::ExprNode { ... }
    println!("{}", run(&node));
}

fn run(node: &crate::ast::ExprNode) -> i32 {
    compile(node);

    let status = std::process::Command::new("./a.out").status().unwrap();

    match status.code() {
        Some(code) => code,
        None => -1,
    }
}

fn compile(node: &crate::ast::ExprNode) {
    let mut c = crate::compiler::Compiler::new();
    let output = c.gen_llir(node);

    std::fs::write("a.out.ll", output).unwrap();

    std::process::Command::new("clang")
        .arg("-Wno-override-module")
        .arg("-o")
        .arg("a.out")
        .arg("a.out.ll")
        .output()
        .unwrap();
}

```

运行代码将得到以下的LLVM汇编：

```ll
define i32 @main() {
	%t1 = add i32 0, 1
	%t3 = add i32 0, 2
	%t5 = add i32 0, 3
	%t6 = add i32 0, 4
	%t4 = add i32 %t5, %t6
	%t2 = mul i32 %t3, %t4
	%t0 = add i32 %t1, %t2
	ret i32 %t0
}
```

这样就完成了表达式树到LLVM汇编的翻译。
