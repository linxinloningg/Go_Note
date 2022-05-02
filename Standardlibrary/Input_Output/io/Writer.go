package main

import (
	"bytes"
	"os"
)
/*
Writer接口用于包装基本的写入方法。

Write方法len(p) 字节数据从p写入底层的数据流。它会返回写入的字节数(0 <= n <= len(p))和遇到的任何导致写入提取结束的错误。Write必须返回非nil的错误，如果它返回的 n < len(p)。Write不能修改切片p中的数据，即使临时修改也不行。
 */
func WriteTo(content string) {
	file, err := os.Create("writetest.txt")
	if err != nil {
		panic(err)
	}
	reader := bytes.NewReader([]byte(content))
	reader.WriteTo(file)

}
func main() {
	WriteTo("text")
}
