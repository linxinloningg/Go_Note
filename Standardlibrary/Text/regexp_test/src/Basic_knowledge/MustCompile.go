package Basic_knowledge

import (
	"fmt"
	"regexp"
)

func MustCompile() {
	// r := regexp.MustCompile(`Hello`)
	var r = regexp.MustCompile(`Hello`)

	if r.MatchString("Hello Regular Expression.") == true {
		fmt.Printf("Match ") // Will print 'Match' again
	} else {
		fmt.Printf("No match ")
	}
}
