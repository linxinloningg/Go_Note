**包用于组织 Go 源代码以获得更好的可重用性和可读性。 包是位于同一目录中的 Go 源文件的集合。 包提供了代码划分，因此很容易维护 Go 项目。**

### 主函数和主包

每个可执行的 Go 应用程序都必须包含 main 函数。 该函数是执行的入口点。 主要功能应该驻留在主包中。

**package package_name**指定特定源文件属于包**package_name**. 这应该是每个 go 源文件的第一行。

一般来说，一个文件夹可以作为 package，同一个 package 内部变量、类型、方法等定义可以相互看到。

比如我们新建一个文件 `calc.go`， `main.go` 平级，分别定义 add 和 main 方法。

```
// calc.go
package main

func add(num1 int, num2 int) int {
	return num1 + num2
}
// main.go
package main

import "fmt"

func main() {
	fmt.Println(add(3, 5)) // 8
}
```

运行 `go run main.go`，会报错，add 未定义：

```
./main.go:6:14: undefined: add
```

因为 `go run main.go` 仅编译 main.go 一个文件，所以命令需要换成 

```
$ go run main.go calc.go
8
```

或 

```
$ go run .
8
```

Go 语言也有 Public 和 Private 的概念，粒度是包。如果类型/接口/方法/函数/字段的首字母大写，则是 Public 的，对其他 package 可见，如果首字母小写，则是 Private 的，对其他 package 不可见。

### Go Module

**Go Module 只不过是 Go 包的集合。** 现在你可能会想到这个问题。 为什么我们需要 Go 模块来创建自定义包？ 答案是 **我们创建的自定义包的导入路径来源于 go 模块的名称** 。 除此之外，我们的应用程序使用的所有其他第三方包（例如来自 github 的源代码）都将出现在 go.mod文件连同版本。 这 go.mod文件是在我们创建新模块时创建的。

[Go Modules](https://github.com/golang/go/wiki/Modules) 是 Go 1.11 版本之后引入的，Go 1.11 之前使用 $GOPATH 机制。Go Modules 可以算作是较为完善的包管理工具。同时支持代理，国内也能享受高速的第三方包镜像服务。接下来简单介绍 `go mod` 的使用。Go Modules 在 1.13 版本仍是可选使用的，环境变量 GO111MODULE 的值默认为 AUTO，强制使用 Go Modules 进行依赖管理，可以将 GO111MODULE 设置为 ON。



### 创建一个 Go 模块

在一个空文件夹下，初始化一个 Module

```
$ go mod init example
go: creating new go.mod: module example
```

此时，在当前文件夹下生成了`go.mod`，这个文件记录当前模块的模块名以及所有依赖包的版本。

接着，我们在当前目录下新建文件 `main.go`，添加如下代码：

```
package main

import (
	"fmt"

	"rsc.io/quote"
)

func main() {
	fmt.Println(quote.Hello())  // Ahoy, world!
}
```

运行 `go run .`，将会自动触发第三方包 `rsc.io/quote`的下载，具体的版本信息也记录在了`go.mod`中：

```
module example

go 1.13

require rsc.io/quote v3.1.0+incompatible
```

我们在当前目录，添加一个子 package calc，代码目录如下：

```
demo/
   |--calc/
      |--calc.go
   |--main.go
```

在 `calc.go` 中写入

```
package calc

func Add(num1 int, num2 int) int {
	return num1 + num2
}
```

在 package main 中如何使用 package cal 中的 Add 函数呢？`import 模块名/子目录名` 即可，修改后的 main 函数如下：

```go
package main

import (
	"fmt"
	"example/calc"

	"rsc.io/quote"
)

func main() {
	fmt.Println(quote.Hello())
	fmt.Println(calc.Add(10, 3))
}
$ go run .
Ahoy, world!
13
```

问题描述：A connection attempt failed because the connected party did not properly respond after a period of time, or established connection failed because connected host has failed to respo nd.

解决：

goproxy.cn
在最go1.11发布后,使用go modules管理包依赖，同时还发布一个goproxy提供代理服务，github地址： https://github.com/goproxy，goproxy.cn是专门服务于中国的，依赖于七牛云。

思路：

1. 设置env proxy:

```bash
go env -w GOPROXY=https://goproxy.cn,direct
```

2. 正常使用go modules管理包:

问题描述：$GOPATH/go.mod exists but should not

产生原因：开启模块支持后，并不能与$GOPATH共存,所以把项目从$GOPATH中移出即可



