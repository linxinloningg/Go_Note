package Advanced

import (
	"fmt"
	"regexp"
)

func Flags() {
	/*

		regexp 包知道以下标志 [引用自文档]：

			    i 不区分大小写（默认为 false）
			    m 多行模式：^ 和 $ 匹配开始/结束行以及开始/结束文本（默认为 false）
			    s 让.匹配\n （默认为false）
			    U ungreedy：交换 x* 和 x*?、x+ 和 x+? 等的含义（默认 false）

			标志语法是 xyz（设置）或 -xyz（清除）或 xy-z（设置 xy，清除 z）。
	*/

	//区分大小写
	//您可能已经知道某些字符在两种情况下存在：大写和小写。 [现在你可能会说：“我当然知道，每个人都知道！” 好吧，如果您认为这是微不足道的，那么考虑大写/小写问题是这几种情况：a、$、本、ß、Ω。 好吧，我们不要把事情复杂化，只考虑英语。]
	//如果你明确地想忽略这种情况，换句话说，如果你想允许正则表达式或它的一部分的两种情况，你使用'i'标志。
	// Do we have an 'N' or 'n' at the beginning?
	fmt.Printf("%v", regexp.MustCompile(`(?i)^n`).MatchString("Never say never."))
	// true, case insensitive

	//Greedy vs. Non-Greedy
	// 正如我们之前看到的，正则表达式可能包含重复符号。  在某些情况下，正则表达式匹配给定字符串实际上可能有不止一种解决方案。
	// 例如，给定正则表达式 '.*' （包括引号），这将如何匹配：
	// 'abc'，'def'，'ghi'
	// 您可能期望检索“abc”。  不是这样。
	//默认情况下，正则表达式是Greedy的。  它们将使用尽可能多的字符来匹配正则表达式。  因此答案是'abc','def','ghi'，因为中间的引号也匹配点“.”！  像这儿：
	fmt.Printf("<%v>", regexp.MustCompile(`'.*'`).FindString(" 'abc','def','ghi' "))
	// Will print: <'abc','def','ghi'>

	//没有简单的方法可以让您指定匹配“abc”、“def”的正则表达式。
	//您可以使用标志 U 恢复正则表达式的行为以使Non-Greedy成为默认值
	fmt.Printf("<%v>", regexp.MustCompile(`(?U)'.*'`).FindString(" 'abc','def','ghi' "))
	// Will print: <'abc'>


}
