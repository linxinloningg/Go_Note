// example for request meta

package main

import (
	"Spider/src/core/common/page"
	"Spider/src/core/common/request"
	"Spider/src/core/spider"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strings"
)

const (
	wkSohuUrl       = "http://yule.sohu.com/gossip/index.shtml"
	wkSohuYule      = `http://changyan.sohu.com/node/html?appid=cyqemw6s1&client_id=cyqemw6s1&topicsid=%s&spSize=5`
	wkSohuPic       = `http://changyan.sohu.com/node/html?appid=cyqemw6s1&client_id=cyqemw6s1&topicsid=9000%s&spSize=5`
	maxWKSouhuLayer = 3 // max grab page
)

var rxYule = regexp.MustCompile(`^http://yule\.sohu\.com/.*?/n(.*?).shtml`)      // gossip section
var rxPic = regexp.MustCompile(`^http://pic\.yule\.sohu\.com/group-(.*?).shtml`) // picture section

type MyPageProcesser struct {
}

type ChangyanListDataJson struct {
	OuterCmtSum      int `json:"outer_cmt_sum"`
	ParticipationSum int `json:"participation_sum"`
}

type ChangyanJson struct {
	ListData ChangyanListDataJson `json:"listData"`
}

func NewMyPageProcesser() *MyPageProcesser {
	return &MyPageProcesser{}
}

func addRequest(p *page.Page, tag, url, content string) {
	req := request.NewRequest(url, "json", tag, "GET", "", nil, nil, nil, content)
	p.AddTargetRequestWithParams(req)
}

func (this MyPageProcesser) Process(p *page.Page) {

	query := p.GetHtmlParser()

	function := func(i int, s *goquery.Selection) {
		url, isExsit := s.Attr("href")
		if isExsit {
			reg := regexp.MustCompile(`do not know what is this`)
			var fmtStr string

            //MatchString报告字符串是否为空
            //包含正则表达式re的任何匹配项。
			if rxYule.MatchString(url) {
				reg = rxYule
				fmtStr = wkSohuYule
			}

			if rxPic.MatchString(url) {
				reg = rxPic
				fmtStr = wkSohuPic
			}

			regxpArrag := reg.FindStringSubmatch(url)
			if len(regxpArrag) == 2 {
				addRequest(p, "changyan", fmt.Sprintf(fmtStr, regxpArrag[1]), s.Text())
			}
		}
	}

	if p.GetUrlTag() == "index" {
		query.Find(`div[class="main area"] div[class="lc"] ul li a`).Each(function)
	}

	if p.GetUrlTag() == "changyan" {
		jsonMap := ChangyanJson{}
		err := json.NewDecoder(strings.NewReader(p.GetBodyStr())).Decode(&jsonMap)
		if err == nil {
			content, ok := p.GetRequest().GetMeta().(string)
			if ok {
				fmt.Println("Title:", content, " CommentCount:", jsonMap.ListData.OuterCmtSum, " ParticipationCount:", jsonMap.ListData.ParticipationSum)
			}
		}
	}
}

func (this *MyPageProcesser) Finish() {
	fmt.Printf("TODO:before end spider \r\n")
}

func main() {
	req := request.NewRequest(wkSohuUrl, "html", "index", "GET", "", nil, nil, nil, nil)

	sohuSpider := spider.NewSpider(NewMyPageProcesser(), "Sohu").
		AddRequest(req).
		SetSleepTime("rand", 500, 1000).
		SetThreadnum(2)

	for i := 1; i < maxWKSouhuLayer; i++ {
		url := fmt.Sprintf("http://yule.sohu.com/gossip/index_%d.shtml", 5301-i) // magic num
		req := request.NewRequest(url, "html", "index", "GET", "", nil, nil, nil, nil)
		sohuSpider.AddRequest(req)
	}

	sohuSpider.Run()
}
