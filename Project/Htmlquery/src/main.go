package main

import "htmlquery/src/spider"

func main() {
	tempUrl := "https://www.kuaidaili.com/free/inha/%d/"
	var testSpider spider.Spider = &spider.KuaidailiSpider{}
	data := testSpider.Run(tempUrl)
	print(data)
}
