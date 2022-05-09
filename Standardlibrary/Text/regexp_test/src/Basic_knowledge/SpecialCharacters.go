package Basic_knowledge

import (
	"fmt"
	"regexp"
)

func SpecialCharacters() {
	// Will print 'cat'.
	r1, err := regexp.Compile(`.at`)
	if err != nil {
		panic(err)
	}
	fmt.Printf(r1.FindString("The cat sat on the mat."))

	// more dot.
	s := "Nobody expects the Spanish inquisition."

	r2, err := regexp.Compile(`e.`)
	res := r2.FindAllString(s, -1) // negative: 所有匹配
	// Prints [ex ec e ]. 最后一项是'e'和一个空格.
	fmt.Printf("%v", res)
	res = r2.FindAllString(s, 2) // 找到 2 个或更少的匹配项
	// Prints [ex ec].
	fmt.Printf("%v", res)
}
