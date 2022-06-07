package main

import (
	"fmt"
	"htmlquery/src/RealDataSpider"
)

func main() {
	var Spider RealDataSpider.Spider = &RealDataSpider.RealDataSpider{}

	data := Spider.Run([]string{"601238", "600703"}, "default")

	for key, datadict := range data {
		for name, value := range datadict {
			fmt.Printf("%s: %s : %s", key, name, value)
		}
	}
}
