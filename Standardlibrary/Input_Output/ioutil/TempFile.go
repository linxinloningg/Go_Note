package main

import "io/ioutil"

/*
在dir目录下创建一个新的、使用prefix为前缀的临时文件，以读写模式打开该文件并返回os.File指针。
如果dir是空字符串，TempFile使用默认用于临时文件的目录（参见os.TempDir函数）。
不同程序同时调用该函数会创建不同的临时文件，调用本函数的程序有责任在不需要临时文件时摧毁它。
 */
func main() {
	// tmp_dir, _ := ioutil.TempDir("D:\\Go_Project\\src\\Standardlibrary\\ioutil", "go-build")

	tmp_file, _ := ioutil.TempFile("", "gofmt")

	print(tmp_file)
}
