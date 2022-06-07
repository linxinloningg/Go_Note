package main

import (
	"fmt"
	"htmlquery/src/IndexDataSpider"
)

func main() {
	var Spider IndexDataSpider.Spider = &IndexDataSpider.IndexDataSpider{}
	data := Spider.Run("default")
	fmt.Printf("%v", data)
}
