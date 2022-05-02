//package main
//
///*
// //定义接口
//type interface_name interface {
//	method_name1 [return_type]
//	method_name2 [return_type]
//	method_name3 [return_type]
//	...
//	method_namen [return_type]
//}
//
// //定义结构体
//type struct_name struct {
//	 //variables
//}
//
// //实现接口方法
//func (struct_name_variable struct_name) method_name1() [return_type] {
//	 //方法实现
//}
//...
//func (struct_name_variable struct_name) method_namen() [return_type] {
//	 //方法实现
//}
//*/
//
//import (
//	"fmt"
//)
//
////定义
///*
//定义接口
//method_namen [return_type]
//call()
//*/
//type Phone interface {
//	call()
//}
//
///*
//定义结构体
//*/
//type NokiaPhone struct {
//}
//
///*
// 实现接口方法
//*/
//func (nokiaPhone NokiaPhone) call() {
//	fmt.Println("I am Nokia, I can call you!")
//}
//
//func main() {
//
//	//实现
//	var phone Phone
//	phone = new(NokiaPhone)
//	phone.call()
//
//
//}
package main

import (
	"fmt"
)

// 定义一个数据写入器
type DataWriter interface {
	WriteData(data interface{}) error
}

// 定义文件结构，用于实现DataWriter
type file struct {
}

// 实现DataWriter接口的WriteData方法
func (d *file) WriteData(data interface{}) error {

	// 模拟写入数据
	fmt.Println("WriteData:", data)
	return nil
}

func main() {

	// 实例化file
	f := new(file)

	// 声明一个DataWriter的接口
	var writer DataWriter

	// 将接口赋值f，也就是*file类型
	writer = f

	// 使用DataWriter接口进行数据写入
	writer.WriteData("data")
}