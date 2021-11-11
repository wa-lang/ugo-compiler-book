# ugo-compiler-book: 从头开发一个迷你Go语言

本书尝试以实现 µGo 编译器为线索，尝试以边学习边完善的自举方式开发一个玩具语言。

- 在线阅读: https://chai2010.cn/ugo-compiler-book/
- µGo 输出C语言: https://github.com/3dgen/ugo-c-book

## µGo 介绍

µGo 是迷你Go语言玩具版本，只保留最基本的int数据类型、变量定义和函数、分支和循环等最基本的特性。µGo 有以下的关键字：`var`、`func`、`if`、`for`、`return`。此外有一个`int`内置的数据类型，`func input() int` 函数读取一个整数，`println(...)` 打印函数。

比如计算1到100的和对应以下代码：

```go
func main() {
	var sum int
	for i := 0; i <= 100; i = i+1 {
		sum = sum+1
	}
	println(sum)
}
```

µGo 的具体定义会根据需要和实现的情况调整，目前可以将其作为Go的最小子集就可。

## 输出的目标格式

为了跨平台和方便测试，输出LLVM汇编代码，如果以后可能会增加WASM文件。

## 凹坑的起因

- 因为坑就在那里
- 不希望被Rxxx语言把脸摁在地上摩擦
- 凹坑的工具差不多齐全了
- ？
