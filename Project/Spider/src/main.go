//
package main

/*
Packages must be imported:
    "Spider/src/core/common/page"
    "Spider/src/core/spider"
Pckages may be imported:
    "Spider/src/core/pipeline": scawler result persistent;
    "github.com/PuerkitoBio/goquery": html dom parser.
*/
import (
	"Spider/src/core/common/page"
	"Spider/src/core/pipeline"
	"Spider/src/core/spider"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type MyPageProcesser struct {
}

func NewMyPageProcesser() *MyPageProcesser {
	return &MyPageProcesser{}
}

// 在这里解析HTMLDOM，并记录我们想要分页的解析结果。
// Package goquery (http://godoc.org/github.com/PuerkitoBio/goquery) 用于解析html.
func (this *MyPageProcesser) Process(p *page.Page) {
	if !p.IsSucc() {
		println(p.Errormsg())
		return
	}

	query := p.GetHtmlParser()
	var urls []string
	query.Find("h3[class='wb-break-all'] a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		urls = append(urls, "http://github.com/"+href)
	})
	// 这些URL将由其他协同程序保存和控制。
	p.AddTargetRequests(urls, "html")

	name := query.Find(".author a").Text()
	name = strings.Trim(name, " \t\n")
	repository := query.Find("strong[itemprop='name'] a").Text()
	repository = strings.Trim(repository, " \t\n")
	//readme, _ := query.Find("#readme").Html()
	if name == "" {
		p.SetSkip(true)
	}

	// 我们希望通过管道保存的实体
	p.AddField("author", name)
	p.AddField("project", repository)
	//p.AddField("readme", readme)
}

func (this *MyPageProcesser) Finish() {
}

func main() {
	// spider input:
	//  PageProcesser ;
	//  task name used in Pipeline for record;
	spider.NewSpider(NewMyPageProcesser(), "TaskName").

		// 开始url，html是响应类型（“html”或“json”）
		AddUrl("https://github.com/hu17889?tab=repositories", "html").

		AddPipeline(pipeline.NewPipelineConsole()). // print result on screen

		SetThreadnum(3). // crawl request by three Coroutines

		Run()

}
