package main

import "os"

/*
Closer接口用于包装基本的关闭方法。

在第一次调用之后再次被调用时，Close方法的的行为是未定义的。某些实现可能会说明他们自己的行为。
 */
func main() {
	file, err := os.Open("test.txt")
	defer file.Close()
	if err != nil {
		panic(err)
	}
}
