package spider

import (
	"fmt"
	"htmlquery/src/core"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type Spider interface {
	GetUrls(tempUrl string) []string
	ParseUrl(url string) string
	GetContent(html string) []string
	Run(tempUrl string) []string
}
type KuaidailiSpider struct {
}

func (this *KuaidailiSpider) GetUrls(tempUrl string) []string {
	urls := make([]string, 0)
	for i := 1; i < 4; i++ {
		urls = append(urls, fmt.Sprintf(tempUrl, i))
	}
	return urls
}

func (this *KuaidailiSpider) ParseUrl(url string) string {
	client := &http.Client{Timeout: time.Second * 5}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3776.0 Safari/537.36")

	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	response, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil && data == nil {
		log.Fatalln(err)
	}
	return fmt.Sprintf("%s", data)
}

func (this *KuaidailiSpider) GetContent(html string) []string {
	var proxies []string
	root, _ := core.Parse(strings.NewReader(html))
	tr := core.Find(root, "//*[@id='list']/table/tbody/tr")
	for _, row := range tr {
		item := core.Find(row, "./td")
		ip := core.InnerText(item[0])
		port := core.InnerText(item[1])
		//type_ := core.InnerText(item[3])
		p := ip + ":" + port
		proxies = append(proxies, p)
	}
	return proxies
}

func (this *KuaidailiSpider) Run(tempUrl string) []string {
	data := make([]string, 0)
	urls := this.GetUrls(tempUrl)
	for _, url := range urls {
		html := this.ParseUrl(url)
		content := this.GetContent(html)
		data = append(data, content...)
	}
	return data
}
