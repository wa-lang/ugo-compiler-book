# 1.3. LLVM汇编简介

LLVM是低级虚拟机，其对应的指令可以看着是一种低级跨平台汇编语言（LLVM IR 是一种 SSA静态单赋值语言）。本节我们简单介绍LLVM汇编语言。


## 最简汇编程序

最小编译器的例子我们见识过最简LLVM汇编程序：

```ll
; hello.ll
define i32 @main() {
	ret i32 42
}
```

`;`开始的是行注释，注明了汇编程序文件为 `hello.ll`。define 定义一个 `@main` 函数，函数返回值是 i32 类型。`@main`函数的实现只有一个ret返回语句，返回 i32 类型的 42。

通过以下命令编译为可执行程序、并执行和查看返回值：

```shell
$ clang hello.ll
warning: overriding the module target triple with x86_64-apple-macosx10.15.4 [-Woverride-module]
1 warning generated.
$ ./a.out
$ echo $?
42
```

其中clang命令输出了一个警告信息：表示该LLVM程序没有指定目标平台，因此用本地环境覆盖了。警告错误中的 `x86_64-apple-macosx10.15.4` 为 LLVM 到目标三元组，第一个 `x86_64` 表示 CPU 类型、第二个 `apple` 表示操作系统类型、第三个 `macosx10.15.4` 操作系统版本信息。我们可以通过给 clang 添加 `-Wno-override-module` 命令行参数关闭该警告信息（在汇编程序中添加 `target triple = "x86_64-apple-macosx10.15.4"`也可以，不过汇编程序就绑死了目标平台）。

## 打印加减法结果

要实现加减法运算，可以通过LLVM的add和sub指令完成：

```ll
define i32 @main() {
	%x1 = add i32 1, 3
	%x2 = sub i32 %x1, 2
	ret i32 %x2 ; 1+3-2
}
```

其中 add 和 sub 分别做减法和减法，指令后面跟着的3个参数分别是 类型、二元操作数。返回的结果依然可以通过 `echo $?` 方式查看，不过需要注意到是 shell 只支持 0-255 范围内的返回值。

如果要直接输出运算结果，可以借助 C语言的 `@printf` 函数完成：

```ll
declare i32 @printf(i8*,...)

@format = constant [4 x i8] c"%d\0A\00"

define i32 @main() {
	; 1 + 3 - 2
	%x1 = add i32 1, 3
	%x2 = sub i32 %x1, 2

	; printf("%d\n", x2)
	call i32(i8*,...) @printf(i8* getelementptr([4 x i8], [4 x i8]* @format, i32 0, i32 0), i32 %x2)

	ret i32 0
}
```

首先通过 declare 指令从外部导入 `@printf` 函数。然后定义 `@format` 字符串常量 `"%d\n"`，用于 printf 函数的第一个参数（其中`%d`表示输出一个整数）。call 指令调用`@printf`打印函数输出`%x2`的值，其中 `getelementptr` 是将 `@format` 转为 `i8*`类型的指针传入第一个参数，第二个参数是 i32 类型的 `%x2`。

## 程序结构

LLVM汇编程序的结构大同小异，最开始是目标三元组（可省略）：

```ll
target triple = "x86_64-pc-linux-gnu"
```

然后是导入的和本地定义的函数、类型、常量、变量等：

```ll
; 声明 puts 函数
declare i32 @puts(i8*)

; 定义常量字符串
@msg = constant [14 x i8] c"Hello, world!\00"
```

最后是全局函数的定义：

```ll
define i32 @main() {
entry:
	call i32(i8*) @puts(i8* getelementptr([14 x i8], [14 x i8]* @msg, i32 0, i32 0))
	ret i32 0
}
```

这里的main函数通过外部的C语言puts函数输出一个“Hello, world!”字符串。

## 小结

LLVM IR 是一种 SSA静态单赋值语言，因此每个名称在您的程序中只能被分配一次。全局名称以 @开头，通常用于全局常量或函数名称，例如 @main。在每个函数中使用局部寄存器名称，这些名称以%开头（这些并不是真正的寄存器，因此数量并无限制）。

LLVM的常见类型：

- i1：1位整数，可以用于分支判断的条件
- i8：8为整数，对应 byte 类型
- i32：32位整数
- i64：64位整数
- `[N x type]`：长度为N，类似是type的数组，比如 `@format = constant [4 x i8] c"%d\0A\00"` 对应 `[4]int8` 类型数组。
- `type*`：指向type类型的指针
- `type(types)`：函数类型，type是函数返回值类型，types是参数的类型。比如 main 函数可能是 `i32()` 或 `i32(i32,i8**)` 类型。

每个 LLVM IR 函数内的指令被分组为基本块或普通指令。每个基本块都可以选择以一个标签开头（如果只有一个块则可以省略）。并且每个基本块都必须以特殊的“终止符指令”结尾（比如 br 、ret等指令）。每个基本块标签也对应一个虚拟寄存器，可以通过`%label`访问其对应的地址。

LLVM 语言参考在这里：https://llvm.org/docs/LangRef.html

