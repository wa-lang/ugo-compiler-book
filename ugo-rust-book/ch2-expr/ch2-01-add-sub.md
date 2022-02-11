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

如果将输入的`1+3-2`转化为`vec!["1", "+", "3", "-", "2"]` 形式，我们则可以通过以下代码输出对应的汇编程序：

```rust,noplayground
fn gen_asm(tokens: &[&str]) -> String {
    let mut result = String::new();

    result.push_str("define i32 @main() {\n");

    let mut idx = 0;
    for (i, tok) in tokens.iter().enumerate() {
        if i == 0 {
            result.push_str(&format!("\t%t{} = add i32 0, {}\n", idx, tok));
            continue;
        }
        match tok {
            &"+" => {
                idx = idx + 1;
                result.push_str(&format!(
                    "\t%t{} = add i32 %t{}, {}\n",
                    idx,
                    idx - 1,
                    tokens[i + 1]
                ));
            }
            &"-" => {
                idx = idx + 1;
                result.push_str(&format!(
                    "\t%t{} = sub i32 %t{}, {}\n",
                    idx,
                    idx - 1,
                    tokens[i + 1]
                ));
            }
            _ => {}
        }
    }

    result.push_str(&format!("\tret i32 %t{}\n", idx));
    result.push_str("}\n");

    result
}
```

而如何将输入的字符串拆分为记号数组本质上属于词法分析的问题。我们先以最简单的方式实现：

```rust,noplayground
fn parse_tokens(code: &str) -> Vec<&str> {
    let mut tokens = Vec::new();
    let mut pos = 0;

    loop {
        if let Some(i) = code[pos..].find('+') {
            tokens.push(&code[pos..][..i]);
            tokens.push(&code[pos..][i..][..1]);
            pos = pos + i + 1;
            continue;
        }
        if let Some(i) = code[pos..].find('-') {
            tokens.push(&code[pos..][..i]);
            tokens.push(&code[pos..][i..][..1]);
            pos = pos + i + 1;
            continue;
        }

        tokens.push(&code[pos..]);
        return tokens;
    }
}
```

基本思路是通过遍历输入字符串，然后根据 `+-` 符号拆分，最终返回拆分后的词法列表。

然后对上个版本的compile函数稍加改造以支持加法和减法的运算表达式编译：

```rust,noplayground
fn compile(code: &str) {
    let tokens = parse_tokens(code);
    let output = gen_asm(&tokens);

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

为了便于测试，我们再包装一个run函数：

```rust,noplayground
fn run(code: &str) -> i32 {
    compile(code);

    let status = std::process::Command::new("./a.out").status().unwrap();

    match status.code() {
        Some(code) => code,
        None => -1,
    }
}
```

run函数将输入的表达式程序编译并运行、最后返回状态码。然后构造单元测试：

```rust,noplayground
#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn run_works() {
        assert_eq!(run("1"), 1);
        assert_eq!(run("1+1"), 2);
        assert_eq!(run("1 + 3 - 2"), 2);
        assert_eq!(run("1+2+3+4"), 10);
    }
}
```

运行`cargo test`测试命令。确认单元测试没有问题后，更新main函数：

```rust,noplayground
fn main() {
    let mut buffer = String::new();
    std::io::stdin().read_line(&mut buffer).unwrap();

    println!("{}", run(buffer.as_ref()));
}
```

通过以下命令执行：

```
$ echo "1+2+3" | cargo run
6
```
