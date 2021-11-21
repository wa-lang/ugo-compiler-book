# 3. 最小µGo程序

本章将尝试编译最小的最小µGo程序，代码如下：

```go
package main

func main() {
	exit(40+2) // 推出码 42
}
```

针对最小µGo程序，我们需要重新设计完善AST，然后编译main函数的唯一一个函数调用语句。exit则是作为builtin函数，完成退出程序的操作。
