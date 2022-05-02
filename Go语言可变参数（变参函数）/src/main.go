package main

import (
	"bytes"
	"fmt"
)

/*
函数 Func_int() 接受不定数量的参数，这些参数的类型全部是 int
*/
func Func_int(args ...int) {
	for _, arg := range args {
		fmt.Println(arg)
	}
}

/*
形如...type格式的类型只能作为函数的参数类型存在，并且必须是最后一个参数，它是一个语法糖（syntactic sugar），
即这种语法对语言的功能并没有影响，但是更方便程序员使用，通常来说，使用语法糖能够增加程序的可读性，从而减少程序出错的可能。

从内部实现机理上来说，类型...type本质上是一个数组切片，也就是[]type，
这也是为什么上面的参数 args 可以用 for 循环来获得每个传入的参数。

从函数的实现角度来看，这没有任何影响，该怎么写就怎么写，但从调用方来说，情形则完全不同：
Func_int_([]int{1, 3, 7, 13})

大家会发现，我们不得不加上[]int{}来构造一个数组切片实例，但是有了...type这个语法糖，我们就不用自己来处理了。
*/
func Func_int_(args []int) {
	for _, arg := range args {
		fmt.Println(arg)
	}
}

/*
任意类型的可变参数
*/

func Func_interface(args ...interface{}) {
	for _, arg := range args {
		switch arg.(type) {
		case int:
			fmt.Println(arg, "is an int value.")
		case string:
			fmt.Println(arg, "is a string value.")
		case int64:
			fmt.Println(arg, "is an int64 value.")
		default:
			fmt.Println(arg, "is an unknown type.")
		}
	}
}

/*
遍历可变参数列表——获取每一个参数的值
可变参数列表的数量不固定，传入的参数是一个切片，如果需要获得每一个参数的具体值时，可以对可变参数变量进行遍历
*/

// 定义一个函数, 参数数量为0~n, 类型约束为字符串
func Func_string(slist ...string) string {
	// 定义一个字节缓冲, 快速地连接字符串
	var b bytes.Buffer
	// 遍历可变参数列表slist, 类型为[]string
	for _, s := range slist {
		// 将遍历出的字符串连续写入字节数组
		b.WriteString(s)
	}
	// 将连接好的字节数组转换为字符串并输出
	return b.String()
}

func Func_printType(slist ...interface{}) string {
	// 字节缓冲作为快速字符串连接
	var b bytes.Buffer
	// 遍历参数
	for _, s := range slist {
		// 将interface{}类型格式化为字符串
		str := fmt.Sprintf("%v", s)
		// 类型的字符串描述
		var typeString string
		// 对s进行类型断言
		switch s.(type) {
		case bool: // 当s为布尔类型时
			typeString = "bool"
		case string: // 当s为字符串类型时
			typeString = "string"
		case int: // 当s为整型类型时
			typeString = "int"
		}
		// 写字符串前缀
		b.WriteString("value: ")
		// 写入值
		b.WriteString(str)
		// 写类型前缀
		b.WriteString(" type: ")
		// 写类型字符串
		b.WriteString(typeString)
		// 写入换行符
		b.WriteString("\n")
	}
	return b.String()
}

/*
在多个可变参数函数中传递参数
*/
// 实际打印的函数
func rawPrint(rawList ...interface{}) {
	// 遍历可变参数切片
	for _, a := range rawList {
		// 打印参数
		fmt.Println(a)
	}
}

// 打印函数封装
func Func_Print(slist ...interface{}) {
	// 将slist可变参数切片完整传递给下一个函数
	rawPrint(slist...)
}

func main() {
	Func_Print(1, 2, 3)
}
