# 5. if分支和for循环

for循环的例子：

```go
package main

func main() {
	for i := 5; i < 10; i = i + 1 {
		println(i)
	}
}
```

执行：

```
$ go run main.go run ./_examples/hello.ugo
5
6
7
8
9
```