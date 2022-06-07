//
package main

import (
	"Spider/src/core/common/page"
	"Spider/src/core/common/request"
	"Spider/src/core/spider"
	"fmt"
	"strings"
)

type MyPageProcesser struct {
}

func NewMyPageProcesser() *MyPageProcesser {
	return &MyPageProcesser{}
}

//网页解析器
// Parse html dom here and record the parse result that we want to crawl.
// Package goquery (http://godoc.org/github.com/PuerkitoBio/goquery) is used to parse html.
func (this *MyPageProcesser) Process(p *page.Page) {
	if !p.IsSucc() {
		println(p.Errormsg())
		return
	}

	//返回绑定到目标爬取结果的 goquery 对象。
	query := p.GetHtmlParser()

	name := query.Find(".lemmaTitleH1").Text()
	name = strings.Trim(name, " \t\n")

	summary := query.Find(".card-summary-content .para").Text()
	summary = strings.Trim(summary, " \t\n")

	// 希望通过管道保存的实体
	p.AddField("name", name)
	p.AddField("summary", summary)

}

func (this *MyPageProcesser) Finish() {
	fmt.Printf("TODO:before end spider \r\n")
}

func main() {
	req := request.NewRequest("http://baike.baidu.com/view/1628025.htm?fromtitle=http&fromid=243074&type=syn",
		"html", "baidu_baike", "GET", "", nil, nil, nil, nil)

	// spider input:
	//  PageProcesser ;
	//  task name used in Pipeline for record;
	gospider := spider.NewSpider(NewMyPageProcesser(), "BaiduSpider").
		AddRequest(req)

	//GetWithParams参数：
	//  1. 网址。
	//  2. 响应类型为“html”或“json”或“jsonp”或“文本”。
	//  3. urltag是标记url的名称，用于在PageProcesser和Pipeline中区分不同的url。
	//  4. 方法是POST或GET。
	//  5. postdata是发送到服务器的主体字符串。
	//  6. 标头是http请求的标头。
	//  7. Cookies

	//处理一个 url 并返回具有其他设置的 PageItems。
	pageItems := gospider.GetByRequest(req)
	//pageItems := gospider.Get("http://baike.baidu.com/view/1628025.htm?fromtitle=http&fromid=243074&type=syn", "html")

	//GetRequest返回PageItems的请求。
	//获取Url
	url := pageItems.GetRequest().GetUrl()
	println("-----------------------------------spider.Get---------------------------------")
	println("url\t:\t" + url)

	//GetAll返回所有键结果。
	for name, value := range pageItems.GetAll() {
		println(name + "\t:\t" + value)
	}

	/*
		println("\n--------------------------------spider.GetAll---------------------------------")
		urls := []string{
			"http://baike.baidu.com/view/1628025.htm?fromtitle=http&fromid=243074&type=syn",
			"http://baike.baidu.com/view/383720.htm?fromtitle=html&fromid=97049&type=syn",
		}
		var reqs []*request.Request
		for _, url := range urls {
			req := request.NewRequest(url, "html", "", "GET", "", nil, nil, nil, nil)
			reqs = append(reqs, req)
		}
		pageItemsArr := gospider.SetThreadnum(2).GetAllByRequest(reqs)
		//pageItemsArr := gospider.SetThreadnum(2).GetAll(urls, "html")
		for _, item := range pageItemsArr {
			url = item.GetRequest().GetUrl()
			println("url\t:\t" + url)
			fmt.Printf("item\t:\t%s\n", item.GetAll())
		}*/
}
