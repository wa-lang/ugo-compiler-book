# 《µGo语言实现——从头开发一个迷你Go语言编译器》

----

- [蚂蚁 - 可信原生技术部 - 云原生运维专家(杭州P7-8)](https://github.com/chai2010/chai2010/blob/master/jobs.md)
- [蚂蚁 - 可信原生技术部 - 专用编程语言设计研发(杭州P7-8)](https://github.com/chai2010/chai2010/blob/master/jobs.md)

----

本书尝试以实现 µGo 编译器为线索，以边学习边完善的自举方式实现一个玩具语言。

![](cover.png)

- [µGo 简介](ugo) - ugo 子目录是正在实现中的代码
- 在线阅读: https://chai2010.cn/ugo-compiler-book/
- µGo 输出C语言: https://github.com/3dgen/ugo-c-book

---

## Why: 凹(读音Wa)坑的起因

- 因为坑就在那里
- 不希望被Rxxx语言把脸摁在地上摩擦
- 凹坑的工具差不多齐全了
- ？

## What: µGo 例子

```go
package main

import "libc"
import m "libc.math"

const Pi = 3.14
const Pi_2 = Pi * 2

type MyInt int
type MyInt2 = int

var x = println(1 + 2*(3+4) + -10 + double(50))

func println() int

func main() int {}
```

## Output: 输出的目标格式

为了跨平台和方便测试，输出LLVM汇编代码，如果以后可能会增加WASM文件。

## License 版权

学习目的可自由使用。
