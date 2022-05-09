package Basic_knowledge

import (
	"fmt"
	"regexp"
)

func Alternatives() {
	fmt.Printf("\n%v", regexp.MustCompile(`Jim|Tim`).MatchString("Dickie, Tom and Tim")) // true
	fmt.Printf("\n%v", regexp.MustCompile(`Jim|Tim`).MatchString("Jimmy, John and Jim")) // true

	s := "Clara was from Santa Barbara and Barbara was from Santa Clara"
	//                   -------------                      -----------
	fmt.Printf("\n%v", regexp.MustCompile(`Santa Clara|Santa Barbara`).FindAllStringIndex(s, -1))
	// [[15 28] [50 61]]

	// Equivalent
	v := "Clara was from Santa Barbara and Barbara was from Santa Clara"
	//                   -------------                      -----------
	fmt.Printf("\n%v", regexp.MustCompile(`Santa (Clara|Barbara)`).FindAllStringIndex(v, -1))
	// [[15 28] [50 61]]
}
