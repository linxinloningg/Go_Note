package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

/*
ReadAll从r读取数据直到EOF或遇到error，返回读取的数据和遇到的错误。成功的调用返回的err为nil而非EOF。
因为本函数定义为读取r直到EOF，它不会将读取返回的EOF视为应报告的错误。
*/
func main() {
	file, _ := os.Open("D:\\Go_Project\\src\\Standardlibrary\\io\\test.txt")
	data, _ := ioutil.ReadAll(file)
	fmt.Printf("读取到的数据是：%s\n", data)
}
