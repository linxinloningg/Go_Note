package Advanced

import (
	"fmt"
	"regexp"
)

func Namedmatches() {
	/*
		匹配项只是按顺序存储在数组中，这有点尴尬。出现了两种不同的问题。
		首先，当您在正则表达式的某处插入一个新组时，以下匹配项中的所有数组索引都必须递增。这很麻烦。
		其次，字符串可能是在运行时构建的，并且可能包含许多我们无法控制的括号。
		这意味着我们不知道我们精心构造的括号在哪个索引处匹配。
		为了解决这个问题，引入了命名匹配。它们允许为可用于查找结果的匹配项提供符号名称。
	*/
	re := regexp.MustCompile("(?P<first_char>.)(?P<middle_part>.*)(?P<last_char>.)")

	//[ first_char middle_part last_char]
	name_set := re.SubexpNames()

	//[Super S upe r]
	results_set := re.FindAllStringSubmatch("Super", -1)[0]

	md := map[string]string{}

	for i, value := range results_set {
		fmt.Printf("%d. match='%s'\tname='%s'\n", i, value, name_set[i])
		md[name_set[i]] = value
	}
	fmt.Printf("The names are  : %v\n", name_set)
	fmt.Printf("The matches are: %v\n", results_set)
	fmt.Printf("The first character is %s\n", md["first_char"])
	fmt.Printf("The last  character is %s\n", md["last_char"])

	/*
			在此示例中，字符串“Super”与包含三个部分的正则表达式匹配：
			first_char : 单个字符 (.)
			middle_part : 由一系列字符组成的中间部分
			last_char : 最后一个字符 (.)
			为了简化结果的使用，我们将所有名称存储在 name_set 中，
		并将它们与匹配结果 results_set 一起压缩到一个新映射中，我们将结果作为命名变量的值存储在名为 md 。
	*/
}
