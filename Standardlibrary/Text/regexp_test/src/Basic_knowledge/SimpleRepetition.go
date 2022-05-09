package Basic_knowledge

import (
	"fmt"
	"regexp"
)

/*
FindAllString 函数返回一个包含所有匹配字符串的数组。
FindAllString 接受两个参数，一个字符串和应返回的最大匹配数。
如果您绝对想要所有匹配项，请使用“-1”。
*/

func SimpleRepetition() {
	//找词。单词是 \w 类型的字符序列。加号“+”表示重复：
	s1 := "Eenie meenie miny moe."
	r1, err := regexp.Compile(`\w+`)
	if err != nil {
		panic(err)
	}
	// Prints [Eenie meenie miny moe]
	fmt.Printf("\n%v", r1.FindAllString(s1, -1))

	//与命令行中用于文件名匹配的通配符相比，'' 不代表“任何字符”，而是前一个字符（或组）的重复。
	//虽然“+”要求其前一个符号至少出现一次，但“”也满足于 0 次出现。这可能会导致奇怪的结果。
	//一个空格
	s2 := "Firstname Lastname"
	r2, err := regexp.Compile(`\w+\s\w+`)
	if err != nil {
		panic(err)
	}
	// Prints Firstname Lastname
	fmt.Printf("\n%v", r2.FindString(s2))

	//但如果这是一些用户提供的输入，则可能有两个空格：
	//两个空格
	s3 := "Firstname  Lastname"
	r3, err := regexp.Compile(`\w+\s\w+`)
	if err != nil {
		panic(err)
	}
	// Prints nothing (the empty string=no match)
	fmt.Printf("\n%v", r3.FindString(s3))

	//我们允许任意数量（但至少一个）带有 '\s+' 的空格：
	s4 := "Firstname  Lastname"
	r4, err := regexp.Compile(`\w+\s+\w+`)
	if err != nil {
		panic(err)
	}
	// Prints Firstname  Lastname
	fmt.Printf("\n%v", r4.FindString(s4))

	//如果您阅读 INI 样式的文本文件，您可能希望允许等号周围的空格。
	s5 := "Key=Value"
	r5, err := regexp.Compile(`\w+=\w+`)
	if err != nil {
		panic(err)
	}
	// OK, prints Key=Value
	fmt.Printf("\n%v", r5.FindAllString(s5, -1))
	// fmt.Printf("\n%v", regexp.MustCompile("\\w+=\\w+").FindAllString("Key=Value", -1))

	//现在让我们在等号周围添加一些空格。
	s6 := "Key = Value"
	r6, err := regexp.Compile(`\w+=\w+`)
	if err != nil {
		panic(err)
	}
	// FAIL, prints nothing, the \w does not match the space.
	fmt.Printf("\n%v", r6.FindAllString(s6, -1))

	//因此，我们允许使用 '\s' 使用多个空格（可能包括 0）：
	s7 := "Key = Value"
	r7, err := regexp.Compile(`\w+\s*=\s*\w+`)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n%v", r7.FindAllString(s7, -1))

	//Go-regexp 模式支持更多用 '?' 构造的模式。

}
