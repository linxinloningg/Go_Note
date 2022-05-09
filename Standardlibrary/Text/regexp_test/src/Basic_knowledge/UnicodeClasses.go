package Basic_knowledge

import (
	"fmt"
	"regexp"
)

func Unicode() {
	/*
		Unicode 按块组织，通常按主题或语言分组。在本章中，我给出了一些示例，因为几乎不可能涵盖所有示例（而且它并没有真正的帮助）。
		请参阅 re2 引擎的完整 unicode 列表。
	*/

	//希腊语(Greek):
	r := regexp.MustCompile(`\p{Greek}`)

	if r.MatchString("This is all Γςεεκ to me.") == true {
		fmt.Println("Match ") // Will print 'Match'
	} else {
		fmt.Println("No match ")
	}

	/*
		在 Windows-1252 代码页上有一个 µ，但它不符合条件，
		因为 \p{Greek} 仅涵盖 http:en.wikipedia.orgwikiGreek_and_Coptic 范围 U+0370..U+03FF。
	*/

	if r.MatchString("the µ is right before ¶") == true {
		fmt.Println("Match ")
	} else {
		fmt.Println("No match ") // Will print 'No match'
	}

	//来自希腊语和科普特语代码页的一些非常酷的字母虽然可能是科普特语，但它们符合“希腊语”的条件，所以要小心。
	if r.MatchString("ϵ϶ϓϔϕϖϗϘϙϚϛϜ") == true {
		fmt.Println("Match ") // Will print 'Match'
	} else {
		fmt.Println("No match ")
	}

	// 盲文(Braille)
	if regexp.MustCompile(`\p{Braille}`).MatchString("This is all ⢓⢔⢕⢖⢗⢘⢙⢚⢛ to me.") == true {
		fmt.Println("Match ") // Will print 'Match'
	} else {
		fmt.Println("No match ")
	}
}
