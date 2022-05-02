package main

import "io/ioutil"

func main() {

	content := []byte("hello world")

	ioutil.WriteFile("D:\\Go_Project\\src\\Standardlibrary\\ioutil\\test.txt", content, 0666)
}
