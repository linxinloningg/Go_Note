package Advanced

import (
	"fmt"
	"regexp"
)

func Exactly() {
	s := "11110010101111100101001001110101"
	fmt.Printf("<%v>", regexp.MustCompile(`1{4}`).FindAllStringSubmatch(s, -1))
	// <[[1111] [1111]]>

	fmt.Printf("<%v>", regexp.MustCompile(`1{4}`).FindAllStringIndex(s, -1))
	// <[[0 4] [10 14]]>

	/*
	{} 语法很少使用。原因之一是在许多（所有？）情况下，您可以通过简单地写出重复次数来重写正则表达式。
	[I can see that you might not want to do that for, say, 120.]
	只有当您有非常具体的要求（例如 {123,130}）时，您才会想要使用 {}。
	{} 的一般模式是 x{n,m}。  “n”是最小出现次数，“m”是最大出现次数。
	 Go-regexp 包支持 {} 系列中的更多模式。
	 */

}
