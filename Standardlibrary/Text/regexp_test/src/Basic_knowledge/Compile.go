package Basic_knowledge

/*
Compile 是 regexp 包的核心。
每个正则表达式都必须使用 Compile 或其类似函数 MustCompile 进行准备。
MustCompile 函数的行为几乎与 Compile 类似，但如果无法编译正则表达式，则会引发恐慌。
由于 MustCompile 中的任何错误都会导致恐慌，因此无需返回错误代码作为第二个返回值。
这使得将 MustCompile 调用与您选择的匹配函数链接起来更容易，如下所示：
（但出于性能原因，您应该避免在循环中重复编译正则表达式。）
 */
import (
	"fmt"
	"regexp"
)

func Compile() {
	r, err := regexp.Compile(`Hello`)

	if err != nil {
		fmt.Printf("There is a problem with your regexp.\n")
		return
	}

	// Will print 'Match'
	if r.MatchString("Hello Regular Expression.") == true {
		fmt.Printf("Match ")
	} else {
		fmt.Printf("No match ")
	}
}
