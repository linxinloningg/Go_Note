package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	data, _ := ioutil.ReadFile("D:\\Go_Project\\src\\Standardlibrary\\ioutil\\test.txt")
	fmt.Printf("%s", data)
}
