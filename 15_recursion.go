package main

import "fmt"

/*
递归，就是在运行的过程中调用自己。

语法格式如下：

func recursion() {
	//函数调用自身
   recursion()
}

func main() {
	recursion()
}
*/

func Factorial(n uint64) (result uint64) {
	if n > 0 {
		result = n * Factorial(n-1)
		return result
	}
	return 1
}

func main() {
	var i int = 15
	fmt.Printf("%d 的阶乘是 %d\n", i, Factorial(uint64(i)))
}
