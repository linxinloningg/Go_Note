package Basic_knowledge

import (
	"fmt"
	"regexp"
)

func AnchorandBoundaries() {
	// 插入符号 ^ 表示“行首”。
	s1 := "Never say never."
	// Do we have an 'N' at the beginning?
	fmt.Printf("\n%v ", regexp.MustCompile(`^N`).MatchString(s1)) // true

	// Do we have an 'n' at the beginning?
	fmt.Printf("\n%v ", regexp.MustCompile(`^n`).MatchString(s1)) // false

	// 美元符号 $ 表示“行尾”。
	s2 := "All is well that ends well"
	fmt.Printf("\n%v ", regexp.MustCompile(`well$`).MatchString(s2)) // true

	// true, but matches with first
	fmt.Printf("\n%v ", regexp.MustCompile(`well`).MatchString(s2))
	// occurrence of 'well'

	// false, not at end of line.
	fmt.Printf("\n%v ", regexp.MustCompile(`ends$`).MatchString(s2))

	/*
		我们看到“well”匹配。为了弄清楚正则表达式在哪里匹配，让我们看一下索引。
		FindStringIndex 函数返回一个包含两个条目的数组。
		第一个条目是正则表达式匹配的索引（当然从 0 开始）。
		第二个是正则表达式结束的索引。
	*/

	s3 := "All is well that ends well"
	//    012345678901234567890123456
	//              1         2
	fmt.Printf("\n%v", regexp.MustCompile(`well$`).FindStringIndex(s3)) // Prints [22 26]

	fmt.Printf("\n%v ", regexp.MustCompile(`well`).MatchString(s3)) // true, but matches with first
	// occurrence of 'well'
	fmt.Printf("\n%v", regexp.MustCompile(`well`).FindStringIndex(s3)) // Prints [7 11], the match starts at 7 and end before 11.

	fmt.Printf("\n%v ", regexp.MustCompile(`ends$`).MatchString(s3)) // false, not at end of line.

	//您可以使用 '\b' 找到单词边界。 FindAllStringIndex 函数捕获容器数组中正则表达式的所有命中。
	s4 := "How much wood would a woodchuck chuck in Hollywood?"
	//    012345678901234567890123456789012345678901234567890
	//              10        20        30        40        50
	//             -1--         -2--                    -3--
	// Find words that *start* with wood
	//    1      2
	fmt.Printf("%v", regexp.MustCompile(`\bwood`).FindAllStringIndex(s4, -1))
	// [[9 13] [22 26]]

	// Find words that *end* with wood
	//   1      3
	fmt.Printf("%v", regexp.MustCompile(`wood\b`).FindAllStringIndex(s4, -1))
	// [[9 13] [46 50]]

	// Find words that *start* and *end* with wood
	//   1
	fmt.Printf("%v", regexp.MustCompile(`\bwood\b`).FindAllStringIndex(s4, -1))
	// [[9 13]]

}
