package main

import (
	"fmt"
)

func main() {
	//Fprint采用默认格式将其参数格式化并写入w。如果两个相邻的参数都不是字符串，会在它们的输出之间添加空格。
	//返回写入的字节数和遇到的任何错误。
	//n, _ := fmt.Fprint(os.Stdout, "hello", "world\n")
	//print(n)

	//Fprintf根据format参数生成格式化的字符串并写入w。返回写入的字节数和遇到的任何错误。
	//n, _ := fmt.Fprintf(os.Stdout, "%s%s", "hello", "world\n")
	//print(n)

	//Fprintln采用默认格式将其参数格式化并写入w。总是会在相邻参数的输出之间添加空格并在输出结束后添加换行符。
	//返回写入的字节数和遇到的任何错误。
	//n, _ := fmt.Fprintln(os.Stdout, "hello", "world\n")
	//print(n)

	//Sprint采用默认格式将其参数格式化，串联所有输出生成并返回一个字符串。
	//如果两个相邻的参数都不是字符串，会在它们的输出之间添加空格。
	//data := fmt.Sprint("hello", "world\n")
	//print(data)

	//Sprintf根据format参数生成格式化的字符串并返回该字符串。
	//data := fmt.Sprintf("%s%s", "hello", "world\n")
	//print(data)

	//Sprintln采用默认格式将其参数格式化，串联所有输出生成并返回一个字符串。
	//总是会在相邻参数的输出之间添加空格并在输出结束后添加换行符。
	//data := fmt.Sprintln("hello", "world\n")
	//print(data)

	//Print采用默认格式将其参数格式化并写入标准输出。
	//如果两个相邻的参数都不是字符串，会在它们的输出之间添加空格。返回写入的字节数和遇到的任何错误。
	//n, _ := fmt.Print("hello", "world\n")
	//print(n)

	//Printf根据format参数生成格式化的字符串并写入标准输出。
	//返回写入的字节数和遇到的任何错误。
	//n, _ := fmt.Printf("%s%s", "hello", "world\n")
	//print(n)

	//Println采用默认格式将其参数格式化并写入标准输出。总是会在相邻参数的输出之间添加空格并在输出结束后添加换行符。
	//返回写入的字节数和遇到的任何错误。
	n, _ := fmt.Println("hello", "world\n")
	print(n)

}
