// Package page contains result catched by Downloader.
// And it alse has result parsed by PageProcesser.
package page

import (
	"Spider/src/core/common/mlog"
	"Spider/src/core/common/page_items"
	"strings"

	//"fmt"
	"Spider/src/core/common/request"
	"github.com/PuerkitoBio/goquery"
	"github.com/bitly/go-simplejson"
	"net/http"
)

// 页面:要爬网的实体。
type Page struct {
	//当爬网过程失败时，isfail为真，errormsg为失败重置。
	isfail   bool
	errormsg string

	//请求由包含url和相关信息的爬行器进行爬网。
	req *request.Request

	// body:抓取结果的纯文本。
	body string

	header  http.Header
	cookies []*http.Cookie

	// docParser:包含 html 结果的 goquery 对象的指针。
	docParser *goquery.Document

	// jsonMap:json 结果。
	jsonMap *simplejson.Json

	// pItems：用于在 PageProcesser 中保存 Key-Values 的对象。
	// And pItems is output in Pipline.
	pItems *page_items.PageItems

	// targetRequests:要放入调度程序的请求。
	targetRequests []*request.Request
}

// NewPage 返回初始化的 Page 对象。
func NewPage(req *request.Request) *Page {
	return &Page{pItems: page_items.NewPageItems(req), req: req}
}

// SetHeader 设置http响应的header
func (this *Page) SetHeader(header http.Header) {
	this.header = header
}

// GetHeader 返回http响应的头部
func (this *Page) GetHeader() http.Header {
	return this.header
}

// SetCookies 设置http 响应的 cookie
func (this *Page) SetCookies(cookies []*http.Cookie) {
	this.cookies = cookies
}

// GetCookies 返回http响应的cookies
func (this *Page) GetCookies() []*http.Cookie {
	return this.cookies
}

// IsSucc 测试下载过程是否成功。
func (this *Page) IsSucc() bool {
	return !this.isfail
}

// Errormsg 显示下载错误信息。
func (this *Page) Errormsg() string {
	return this.errormsg
}

// SetStatus 保存有关下载过程的状态信息。
func (this *Page) SetStatus(isfail bool, errormsg string) {
	this.isfail = isfail
	this.errormsg = errormsg
}

// AddField 将 键-值 字符串对保存到 PageItems 为 Pipeline 做准备
func (this *Page) AddField(key string, value string) {
	this.pItems.AddItem(key, value)
}

// GetPageItems 返回记录在 PageProcesser 中解析的 键-值 对的 PageItems 对象。
func (this *Page) GetPageItems() *page_items.PageItems {
	return this.pItems
}

// SetSkip 设置 PageItems 的标签“跳过”。
// PageItems 将不会保存在跳过设置为 true 的管道中
func (this *Page) SetSkip(skip bool) {
	this.pItems.SetSkip(skip)
}

// GetSkip returns skip label of PageItems.
func (this *Page) GetSkip() bool {
	return this.pItems.GetSkip()
}

// SetRequest 设置此页面的请求对象。
func (this *Page) SetRequest(r *request.Request) *Page {
	this.req = r
	return this
}

// GetRequest 返回此页面的请求对象。
func (this *Page) GetRequest() *request.Request {
	return this.req
}

// GetUrlTag 返回网址的Tag。
func (this *Page) GetUrlTag() string {
	return this.req.GetUrlTag()
}

// AddTargetRequest 添加一个新的请求等待抓取。
func (this *Page) AddTargetRequest(url string, respType string) *Page {
	this.targetRequests = append(this.targetRequests, request.NewRequest(url, respType, "", "GET", "", nil, nil, nil, nil))
	return this
}

// AddTargetRequests 添加多个新的等待抓取的请求。
func (this *Page) AddTargetRequests(urls []string, respType string) *Page {
	for _, url := range urls {
		this.AddTargetRequest(url, respType)
	}
	return this
}

// AddTargetRequestWithProxy 添加一个新的带有代理的请求等待抓取。
func (this *Page) AddTargetRequestWithProxy(url string, respType string, proxyHost string) *Page {

	this.targetRequests = append(this.targetRequests, request.NewRequestWithProxy(url, respType, "", "GET", "", nil, nil, proxyHost, nil, nil))
	return this
}

// AddTargetRequestsWithProxy 添加多个新的带有代理的请求等待抓取。
func (this *Page) AddTargetRequestsWithProxy(urls []string, respType string, proxyHost string) *Page {
	for _, url := range urls {
		this.AddTargetRequestWithProxy(url, respType, proxyHost)
	}
	return this
}

// AddTargetRequest 添加一个带有Header的新请求等待抓取。
func (this *Page) AddTargetRequestWithHeaderFile(url string, respType string, headerFile string) *Page {
	this.targetRequests = append(this.targetRequests, request.NewRequestWithHeaderFile(url, respType, headerFile))
	return this
}

// AddTargetRequest 添加一个带有设置的新的等待抓取的请求。
// The respType is "html" or "json" or "jsonp" or "text".
// The urltag is name for marking url and distinguish different urls in PageProcesser and Pipeline.
// The method is POST or GET.
// The postdata is http body string.
// The header is http header.
// The cookies is http cookies.
func (this *Page) AddTargetRequestWithParams(req *request.Request) *Page {
	this.targetRequests = append(this.targetRequests, req)
	return this
}

// AddTargetRequests 添加多个带有设置的新的等待抓取的请求。
func (this *Page) AddTargetRequestsWithParams(reqs []*request.Request) *Page {
	for _, req := range reqs {
		this.AddTargetRequestWithParams(req)
	}
	return this
}

// GetTargetRequests 返回将放入调度程序的目标请求
func (this *Page) GetTargetRequests() []*request.Request {
	return this.targetRequests
}

// SetBodyStr 保存在 Page 中抓取的纯字符串。
func (this *Page) SetBodyStr(body string) *Page {
	this.body = body
	return this
}

// GetBodyStr 返回抓取的纯字符串。
func (this *Page) GetBodyStr() string {
	return this.body
}

// SetHtmlParser 设置绑定到目标爬取结果的 goquery 对象。
func (this *Page) SetHtmlParser(doc *goquery.Document) *Page {
	this.docParser = doc
	return this
}

// GetHtmlParser 返回绑定到目标爬取结果的 goquery 对象。
func (this *Page) GetHtmlParser() *goquery.Document {
	return this.docParser
}

// GetHtmlParser 重置绑定到目标爬取结果的 goquery 对象。
func (this *Page) ResetHtmlParser() *goquery.Document {
	r := strings.NewReader(this.body)
	var err error
	this.docParser, err = goquery.NewDocumentFromReader(r)
	if err != nil {
		mlog.LogInst().LogError(err.Error())
		panic(err.Error())
	}
	return this.docParser
}

// SetJson 保存 json 结果。
func (this *Page) SetJson(js *simplejson.Json) *Page {
	this.jsonMap = js
	return this
}

// SetJson 返回 json 结果。
func (this *Page) GetJson() *simplejson.Json {
	return this.jsonMap
}
