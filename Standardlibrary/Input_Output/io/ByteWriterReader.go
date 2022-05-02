package main

import (
	"bytes"
	"fmt"
)
/*
ByteReader是基本的ReadByte方法的包装。

ReadByte读取输入中的单个字节并返回。如果没有字节可读取，会返回错误。
 */

/*
ByteWriter是基本的WriteByte方法的包装。
 */
func main() {
	var char byte
	// fmt.Scanf("%c\n", &char)
	char = 'h'
	buffer := new(bytes.Buffer)
	// 写入一个字节
	buffer.WriteByte(char)
	// 读取一个字节
	newCh, _ := buffer.ReadByte()
	fmt.Printf("%c", newCh)
}
