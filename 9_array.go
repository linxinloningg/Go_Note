package main

/*
声明数组
Go 语言数组声明需要指定元素类型及元素个数，语法格式如下：

var variable_name [SIZE] variable_type
*/

/*
初始化数组
以下演示了数组初始化：

var balance = [5]float32{1000.0, 2.0, 3.4, 7.0, 50.0}
我们也可以通过字面量在声明数组的同时快速初始化数组：

balance := [5]float32{1000.0, 2.0, 3.4, 7.0, 50.0}
如果数组长度不确定，可以使用 ... 代替数组的长度，编译器会根据元素个数自行推断数组的长度：

var balance = [...]float32{1000.0, 2.0, 3.4, 7.0, 50.0}
或
balance := [...]float32{1000.0, 2.0, 3.4, 7.0, 50.0}
如果设置了数组的长度，我们还可以通过指定下标来初始化元素：

//  将索引为 1 和 3 的元素初始化
balance := [5]float32{1:2.0,3:7.0}
*/
func main() {
	var array = [...]int{}

	print(len(array))
	for i := range array {
		print(array[i])
	}

}
