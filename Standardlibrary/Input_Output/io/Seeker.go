package main

import (
	"fmt"
	"io"
	"strings"
)

/*
Seeker接口用于包装基本的移位方法。

Seek方法设定下一次读写的位置：偏移量为offset，校准点由whence确定：0表示相对于文件起始；1表示相对于当前位置；2表示相对于文件结尾。Seek方法返回新的位置以及可能遇到的错误。

移动到一个绝对偏移量为负数的位置会导致错误。移动到任何偏移量为正数的位置都是合法的，但其下一次I/O操作的具体行为则要看底层的实现。
*/
func main() {
	reader := strings.NewReader("Go语言中文网")
	reader.Seek(-6, io.SeekEnd)
	r, _, _ := reader.ReadRune()
	fmt.Printf("%c\n", r)
}
