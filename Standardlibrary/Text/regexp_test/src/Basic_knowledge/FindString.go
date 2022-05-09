package Basic_knowledge

import (
	"fmt"
	"regexp"
)

func FindString() {
	//FindString 函数查找字符串。当您使用文字字符串时，结果显然是字符串本身。
	//只有当您开始使用模式和类时，结果才会更有趣。
	r1, err := regexp.Compile(`Hello`)
	if err != nil {
		panic(err)
	}
	// Will print 'Hello'
	fmt.Printf(r1.FindString("Hello Regular Expression. Hullo again."))

	// 当 FindString 没有找到匹配正则表达式的字符串时，它将返回空字符串。请注意，空字符串也可能是有效匹配的结果。
	r2, err := regexp.Compile(`Hxllo`)
	if err != nil {
		panic(err)
	}
	// Will print nothing (=the empty string)
	fmt.Printf(r2.FindString("Hello Regular Expression."))
	//FindString 在第一次匹配后返回。如果您对更多可能的匹配感兴趣，您可以使用 FindAllString()，见下文。
}
