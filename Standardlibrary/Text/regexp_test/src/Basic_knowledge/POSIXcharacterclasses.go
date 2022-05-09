package Basic_knowledge

import (
	"fmt"
	"regexp"
)

func POSIXcharacter() {
	/*
		Golang 正则表达式库实现了 POSIX 字符类。这些只是被赋予了更易读名称的常用类的别名。
		The classes are: (https://github.com/google/re2/blob/master/doc/syntax.txt)
	*/

	/*
		[:alnum:]	alphanumeric (≡ [0-9A-Za-z])
		[:alpha:]	alphabetic (≡ [A-Za-z])
		[:ascii:]	ASCII (≡ [\x00-\x7F])
		[:blank:]	blank (≡ [\t ])
		[:cntrl:]	control (≡ [\x00-\x1F\x7F])
		[:digit:]	digits (≡ [0-9])
		[:graph:]	graphical (≡ [!-~] == [A-Za-z0-9!"#$%&'()*+,\-./:;<=>?@[\\\]^_`{|}~])
		[:lower:]	lower case (≡ [a-z])
		[:print:]	printable (≡ [ -~] == [ [:graph:]])
		[:punct:]	punctuation (≡ [!-/:-@[-`{-~])
		[:space:]	whitespace (≡ [\t\n\v\f\r ])
		[:upper:]	upper case (≡ [A-Z])
		[:word:]	word characters (≡ [0-9A-Za-z_])
		[:xdigit:]	hex digit (≡ [0-9A-Fa-f])
	*/

	/*
		请注意，您必须在 [] 中包装一个 ASCII 字符类。
		此外请注意，每当我们谈论字母表时，我们只谈论 ASCII 范围 65-90 中的 26 个字母，不包括带有变音符号的字母。
		示例：查找由小写字母、标点符号、空格（空白）和数字组成的序列：
	*/

	if regexp.MustCompile(`[[:lower:]][[:punct:]][[:blank:]][[:digit:]]`).MatchString("Fred: 12345769") == true {

		fmt.Printf("Match ")
	} else {
		fmt.Printf("No match ")
	}

}
