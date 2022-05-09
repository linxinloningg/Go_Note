package Basic_knowledge

import (
	"fmt"
	"regexp"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
func Character(letter byte) {
	switch letter {
	//字符类 '\w' 表示 [A-Za-z0-9_] 类中的任何字符，助记符：'word'。
	case 'w':
		{
			r, err := regexp.Compile(`H\wllo`)
			// Will print 'true'.
			checkErr(err)
			fmt.Printf("\n%v", r.MatchString("Hello Regular Expression."))
		}
	// 	字符类 '\d' 表示任何数字。
	case 'd':
		{
			r, err := regexp.Compile(`\d`)
			// Will print 'true':
			checkErr(err)
			fmt.Printf("\n%v", r.MatchString("Seven times seven is 49."))
			// Will print 'false':
			fmt.Printf("\n%v", r.MatchString("Seven times seven is forty-nine."))
		}
	// 字符类 '\s' 表示以下任何空格：TAB、SPACE、CR、LF。或者更准确地说 [\t\n\f\r ]。
	case 's':
		{
			r, err := regexp.Compile(`\s`)
			// Will print 'true':
			checkErr(err)
			fmt.Printf("\n%v", r.MatchString("/home/bill/My Documents"))
		}
	// 可以使用大写的“\D”、“\S”、“\W”来否定字符类。因此，'\D' 是任何不是'\d' 的字符。
	case 'S':
		{
			r, err := regexp.Compile(`\S`) // Not a whitespace
			//将打印'true'，显然这里有非空格：
			checkErr(err)
			fmt.Printf("\n%v", r.MatchString("/home/bill/My Documents"))
		}
	// 检查字符串是否有任何不是单词字符的内容。
	case 'W':
		{
			r, err := regexp.Compile(`\W`) // Not a \w character.
			checkErr(err)
			fmt.Printf("\n%v", r.MatchString("555-shoe")) // true: has a non-word char: The hyphen
			fmt.Printf("\n%v", r.MatchString("555shoe"))  // false: has no non-word char.
		}
	default:

	}
}
func Class() {
	/*
		您可以在任何位置要求一组（或类）字符，而不是文字字符。
		在这个例子中，[uio] 是一个“字符类”。方括号中的任何字符都将满足正则表达式。
		因此，这个正则表达式将匹配 'Hullo'、'Hillo' 和 'Hollo'
	*/

	// Will print 'Hullo'.
	fmt.Printf(regexp.MustCompile(`H[uio]llo`).FindString("Hello Regular Expression. Hullo again."))

	/*
		否定字符类反转该类的匹配。在这种情况下，它将匹配所有字符串“H.llo”，其中点不是“o”、“i”或“u”。
		它不会匹配“Hullo”、“Hillo”、“Hollo”，但会匹配“Hallo”甚至“H9llo”。
	*/
	// ^ 非
	r := regexp.MustCompile(`H[^uio]llo`)
	fmt.Printf("\n%v ", r.MatchString("Hillo")) // false
	fmt.Printf("\n%v ", r.MatchString("Hallo")) // true
	fmt.Printf("\n%v ", r.MatchString("H9llo")) // true
}
