package Advanced

import (
	"fmt"
	"regexp"
)

/*
FindAllStringSubmatch
Find返回一个保管正则表达式re在b中的最左侧的一个匹配结果以及（可能有的）分组匹配的结果的起止位置的切片。
匹配结果和分组匹配结果可以通过起止位置对b做切片操作得到：b[loc[2*n]:loc[2*n+1]]。如果没有匹配到，会返回nil。
*/
func Groups() {

	// 有时您想匹配一个字符串，但只想查看特定的切片。在上一章中，我们总是查看整个匹配字符串。

	//[[cat] [sat] [mat]]
	fmt.Printf("%v\n", regexp.MustCompile(`.at`).FindAllStringSubmatch("The cat sat on the mat.", -1))

	//括号允许捕获您真正感兴趣的那段字符串，而不是整个正则表达式。
	//[[cat c] [sat s] [mat m]]
	// want to know what is in front of 'at'
	fmt.Printf("%v\n", regexp.MustCompile(`(.)at`).FindAllStringSubmatch("The cat sat on the mat.", -1))

	//您可以拥有多个组。
	// Prints [[ex e x] [ec e c] [e  e  ]]
	// Prepare our regex
	fmt.Printf("%v\n", regexp.MustCompile(`(e)(.)`).FindAllStringSubmatch("Nobody expects the Spanish inquisition.", -1))

	/*
		FindAllStringSubmatch 函数将为每个匹配返回一个数组，其中第一个字段中包含整个匹配项，其余字段中包含组的内容。然后将所有匹配的数组捕获到容器数组中。

		如果您有一个未出现在字符串中的可选组，则结果数组的单元格中将有一个空字符串，换句话说，结果数组中的字段数始终与组数加一匹配。
	*/

	result1 := regexp.MustCompile(`(Mr)(s)?\. (\w+) (\w+)`).FindStringSubmatch("Mr. Leonard Spock")

	for k, v := range result1 {
		fmt.Printf("%d. %s\n", k, v)
	}
	// Prints
	// 0. Mr. Leonard Spock
	// 1. Mr
	// 2.
	// 3. Leonard
	// 4. Spock

	// 您不能有部分重叠的组。如果我们希望第一个正则表达式匹配 'expects the' 而另一个匹配 'the Spanish'，括号的解释会有所不同。
	//座右铭是：opened for 'the' is closed after 'the'.。

	// Wanted regex1          --------------
	// Wanted regex2                   --------------
	result2 := regexp.MustCompile(`(expects (...) Spanish)`).FindStringSubmatch("Nobody expects the Spanish inquisition.")

	for k, v := range result2 {
		fmt.Printf("%d. %s\n", k, v)
	}
	// 0. expects the Spanish
	// 1. expects the Spanish
	// 2. the

}
