package main

import (
	"fmt"
	"io"
	"os"
)

/*
Reader接口用于包装基本的读取方法。

Read方法读取len(p)字节数据写入p。它返回写入的字节数和遇到的任何错误。即使Read方法返回值n < len(p)，本方法在被调用时仍可能使用p的全部长度作为暂存空间。如果有部分可用数据，但不够len(p)字节，Read按惯例会返回可以读取到的数据，而不是等待更多数据。

当Read在读取n > 0个字节后遭遇错误或者到达文件结尾时，会返回读取的字节数。它可能会在该次调用返回一个非nil的错误，或者在下一次调用时返回0和该错误。一个常见的例子，Reader接口会在输入流的结尾返回非0的字节数，返回值err == EOF或err == nil。但不管怎样，下一次Read调用必然返回(0, EOF)。调用者应该总是先处理读取的n > 0字节再处理错误值。这么做可以正确的处理发生在读取部分数据后的I/O错误，也能正确处理EOF事件。

如果Read的某个实现返回0字节数和nil错误值，表示被阻碍；调用者应该将这种情况视为未进行操作
 */
func ReadFrom(reader io.Reader, num int) ([]byte, error) {
	p := make([]byte, num)
	n, err := reader.Read(p)
	if n > 0 {
		return p[:n], nil
	}
	return p, err
}

func main() {
	var (
		data []byte
		err  error
	)

	//// 从标准输入读取
	//data, err = ReadFrom(os.Stdin, 11)
	//fmt.Printf("读取到的数据是：%s\n", data)
	//fmt.Println("errorMsg is: ", err)
	//
	//// 从字符串读取
	//data, err = ReadFrom(strings.NewReader("from string"), 12)
	//fmt.Printf("读取到的数据是：%s\n", data)
	//fmt.Println("errorMsg is: ", err)

	// 从普通文件读取
	file, err := os.Open("D:\\Go_Project\\src\\Standardlibrary\\io\\test.txt")
	data, err = ReadFrom(file, 12)
	fmt.Printf("读取到的数据是：%s\n", data)
	fmt.Println("errorMsg is: ", err)

}
