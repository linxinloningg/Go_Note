package main

import "io/ioutil"

/*
在dir目录里创建一个新的、使用prfix作为前缀的临时文件夹，并返回文件夹的路径。如果dir是空字符串，TempDir使用默认用于临时文件的目录（参见os.TempDir函数）。
不同程序同时调用该函数会创建不同的临时目录，调用本函数的程序有责任在不需要临时文件夹时摧毁它。
*/
func main() {
	/*
	第一个参数如果为空，表明在系统默认的临时目录（ os.TempDir ）中创建临时目录；
	第二个参数指定临时目录名的前缀，该函数返回临时目录的路径。
	 */
	name, _ := ioutil.TempDir("D:\\Go_Project\\src\\Standardlibrary\\ioutil", "go-build")
	print(name)
}
