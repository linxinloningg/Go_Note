package Basic_knowledge

import (
	"fmt"
	"regexp"
)

/*
CompilePOSIX 和 MustCompilePOSIX 运行的引擎略有不同。  规则按照 POSIX ERE（扩展正则表达式）实现；
从 Go 的角度来看，这意味着一组受限制的规则，即 egrep 支持的规则。
因此，在 POSIX 版本中找不到 Go 的标准 re2 引擎支持的一些细节，例如 \A。
*/

func CompilePOSIX() {
	/*
		s := "ABCDEEEEE"
		rr := regexp.MustCompile(`\AABCDE{2}|ABCDE{4}`)
		rp := regexp.MustCompilePOSIX(`\AABCDE{2}|ABCDE{4}`)
		fmt.Println(rr.FindAllString(s, 2))
		fmt.Println(rp.FindAllString(s, 2))
	*/
	s := "ABCDEEEEE"
	rr := regexp.MustCompile(`ABCDE{2}|ABCDE{4}`)
	rp := regexp.MustCompilePOSIX(`ABCDE{2}|ABCDE{4}`)
	fmt.Println(rr.FindAllString(s, 2))
	fmt.Println(rp.FindAllString(s, 2))

	/*
		[ABCDEE]    <- first acceptable match
		[ABCDEEEE]  <- But POSIX wants the longer match
	*/
}
