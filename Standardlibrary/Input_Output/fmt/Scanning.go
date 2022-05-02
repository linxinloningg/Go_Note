package main

import "fmt"

/*
一系列类似的函数可以扫描格式化文本以生成值。

Scan、Scanf和Scanln从标准输入os.Stdin读取文本；Fscan、Fscanf、Fscanln从指定的io.Reader接口读取文本；Sscan、Sscanf、Sscanln从一个参数字符串读取文本。

Scanln、Fscanln、Sscanln会在读取到换行时停止，并要求一次提供一行所有条目；Scanf、Fscanf、Sscanf只有在格式化文本末端有换行时会读取到换行为止；其他函数会将换行视为空白。

Scanf、Fscanf、Sscanf会根据格式字符串解析参数，类似Printf。例如%x会读取一个十六进制的整数，%v会按对应值的默认格式读取。
*/

func main() {

	//var value int

	/*
		Scan、Scanf 和 Scanln 从 os.Stdin 中读取；
	*/
	// Scan从标准输入扫描文本，将成功读取的空白分隔的值保存进成功传递给本函数的参数。换行视为空白。
	//返回成功扫描的条目个数和遇到的任何错误。如果读取的条目比提供的参数少，会返回一个错误报告原因。
	//fmt.Scan(&value)
	//print(value)

	// Scanf从标准输入扫描文本，根据format 参数指定的格式将成功读取的空白分隔的值保存进成功传递给本函数的参数。
	//返回成功扫描的条目个数和遇到的任何错误。
	//fmt.Scanf("%d", &value)
	//print(value)

	//Scanln类似Scan，但会在换行时才停止扫描。最后一个条目后必须有换行或者到达结束位置。
	//fmt.Scanln(&value)
	//print(value)

	/*
		Fscan、Fscanf 和 Fscanln 从指定的 io.Reader 中读取；
	*/
	//Fscan从r扫描文本，将成功读取的空白分隔的值保存进成功传递给本函数的参数。换行视为空白。
	//返回成功扫描的条目个数和遇到的任何错误。如果读取的条目比提供的参数少，会返回一个错误报告原因。
	//file, _ := os.Open("D:\\Go_Project\\src\\Standardlibrary\\io\\test.txt")
	//fmt.Fscan(file)

	//Fscanf从r扫描文本，根据format 参数指定的格式将成功读取的空白分隔的值保存进成功传递给本函数的参数。
	//返回成功扫描的条目个数和遇到的任何错误。
	//file, _ := os.Open("D:\\Go_Project\\src\\Standardlibrary\\io\\test.txt")
	//fmt.Fscanf(file, "%s")

	//Fscanln类似Fscan，但会在换行时才停止扫描。最后一个条目后必须有换行或者到达结束位置。
	//file, _ := os.Open("D:\\Go_Project\\src\\Standardlibrary\\io\\test.txt")
	//fmt.Fscanln(file)

	/*
		Sscan、Sscanf 和 Sscanln 从实参字符串中读取
	*/
	//Sscan从字符串str扫描文本，将成功读取的空白分隔的值保存进成功传递给本函数的参数。换行视为空白。
	//返回成功扫描的条目个数和遇到的任何错误。如果读取的条目比提供的参数少，会返回一个错误报告原因。
	/*
		var (
			test1 string
			test2 string
		)
		var str string = "test1 test2"
		data, _ := fmt.Sscan(str, &test1, &test2)
		fmt.Printf("返回成功扫描的条目个数:%d,本函数的参数1:%s,本函数的参数2:%s", data, test1, test2)
	*/
	//Sscanf从字符串str扫描文本，根据format 参数指定的格式将成功读取的空白分隔的值保存进成功传递给本函数的参数。
	//返回成功扫描的条目个数和遇到的任何错误。
	/*
		var (
			test1 string
			test2 string
		)
		var str string = "test1 test2"
		data, _ := fmt.Sscanf(str, "%s%s", &test1, &test2)
		fmt.Printf("返回成功扫描的条目个数:%d,本函数的参数1:%s,本函数的参数2:%s", data, test1, test2)
	*/

	//Sscanln类似Sscan，但会在换行时才停止扫描。最后一个条目后必须有换行或者到达结束位置。
	var (
		test1 string
		test2 string
	)
	var str string = "test1\n test2"
	data, _ := fmt.Sscanln(str, &test1, &test2)
	fmt.Printf("返回成功扫描的条目个数:%d,本函数的参数1:%s,本函数的参数2:%s", data, test1, test2)

	// Scanln、Fscanln 和 Sscanln 在换行符处停止扫描，且需要条目紧随换行符之后；
	// Scanf、Fscanf 和 Sscanf 需要输入换行符来匹配格式中的换行符；其它函数则将换行符视为空格。

}
