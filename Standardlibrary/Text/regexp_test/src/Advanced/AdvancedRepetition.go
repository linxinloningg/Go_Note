package Advanced

import (
	"fmt"
	"regexp"
)

func AdvancedRepetition() {
	// Non-matching capture/group repetition
	/*
		如果一个复杂的正则表达式有多个组，您可能会遇到这样一种情况：我们使用括号进行分组，但对捕获的字符串并不是最不感兴趣的。
		要丢弃组的匹配项，您可以使用 (?:regex) 将其设为“非捕获组”。
		问号和冒号告诉编译器使用模式进行匹配但不存储它。
	*/
	result1 := regexp.MustCompile(`Mr(s)?\. (\w+) (\w+)`).FindStringSubmatch("Mrs. Leonora Spock")
	for i, value := range result1 {
		fmt.Printf("%d. %s\n", i, value)
	}
	// 0. Mrs. Leonora Spock
	// 1. s
	// 2. Leonora
	// 3. Spock

	//使用非捕获组：

	result2 := regexp.MustCompile(`Mr(?:s)?\. (\w+) (\w+)`).FindStringSubmatch("Mrs. Leonora Spock")
	for i, value := range result2 {
		fmt.Printf("%d. %s\n", i, value)
	}
	// 0. Mrs. Leonora Spock
	// 1. Leonora
	// 2. Spock
}
