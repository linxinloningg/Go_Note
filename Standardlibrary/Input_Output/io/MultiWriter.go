package main

import (
	"io"
	"os"
)

func main() {
	file, err := os.Create("tmp.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 定义多个输出位置
	writers := []io.Writer{
		file,
		os.Stdout,
	}
	writer := io.MultiWriter(writers...)
	writer.Write([]byte("Go语言中文网"))
}
