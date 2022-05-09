package Basic_knowledge

import (
	"fmt"
	"regexp"
)

func LiteralSpecialCharacters() {
	//找到一个反斜杠''：它必须在正则表达式中转义两次，在字符串中转义一次。
	r1, err := regexp.Compile("C:\\\\")
	if err != nil {
		panic(err)
	}
	if r1.MatchString("Working on drive C:\\") == true {
		fmt.Printf("Matches.") // <---
	} else {
		fmt.Printf("No match.")
	}

	// 查找点：
	r2, err := regexp.Compile(`\.`)
	if err != nil {
		panic(err)
	}
	if r2.MatchString("Short.") == true {
		fmt.Printf("Has a dot.") // <---
	} else {
		fmt.Printf("Has no dot.")
	}
	//与构造正则表达式相关的其他特殊字符以类似的方式工作：.+?()|[]{}^

	//查找文字美元符号：
	r3, err := regexp.Compile(`\$`)
	if err != nil {
		panic(err)
	}
	if len(r3.FindString("He paid $150 for that software.")) != 0 {
		fmt.Printf("Found $-symbol.") // <-
	} else {
		fmt.Printf("No $$$.")
	}
}
