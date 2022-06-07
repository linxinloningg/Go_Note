package DailymarketSpider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Spider interface {
	Init()
	GetUrls() []string
	ParseUrl(url string) string
	GetStockdata(content string) map[string]map[string]interface{}
	Run() map[string]map[string]interface{}
}

type DailymarketSpider struct {
	GsApi   string
	headers map[string]string
}

func (this *DailymarketSpider) Init() {
	this.GsApi = "http://76.push2.eastmoney.com/api/qt/clist/get?cb=jQuery112408744624686429123_1578798932591&pn=%d&pz=20&po=1&np=1&ut=bd1d9ddb04089700cf9c27f6f7426281&fltt=2&invt=2&fid=f3&fs=m:0+t:6,m:0+t:13,m:0+t:80,m:1+t:2,m:1+t:23&fields=f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f12,f13,f14,f15,f16,f17,f18,f20,f21,f23,f24,f25,f22,f11,f62,f128,f136,f115,f152&_=1586266306109"

	this.headers = map[string]string{"User-Agent": "Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.82 Safari/537.36"}
}

func (this *DailymarketSpider) GetUrls() []string {
	urllist := make([]string, 0)

	content := this.ParseUrl(fmt.Sprintf(this.GsApi, 1))

	pattern := regexp.MustCompile(`"total":(?P<Total>.+?),.*?`)

	result := pattern.FindStringSubmatch(content)

	if len(result) == 2 {
		page_number, err := strconv.Atoi(result[1])
		if err != nil {
			panic(err)
		}
		page_number = page_number / 20
		for i := 1; i < page_number; i++ {
			urllist = append(urllist, fmt.Sprintf(this.GsApi, i))
		}
	}
	return urllist
}

func (this *DailymarketSpider) ParseUrl(url string) string {
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
	fmt.Printf("爬取：%s\n", url)
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

func (this *DailymarketSpider) GetStockdata(content string, ) map[string]map[string]interface{} {

	content = strings.Replace(content, "jQuery112408744624686429123_1578798932591", "", -1)
	content = strings.Replace(content, "(", "", -1)
	content = strings.Replace(content, ")", "", -1)
	content = strings.Replace(content, ";", "", -1)

	var mapData map[string]interface{}
	err := json.Unmarshal([]byte(content), &mapData)
	if err != nil {
		panic(err)
	}

	data := mapData["data"].(map[string]interface{})["diff"].([]interface{})

	pagedata := make(map[string]map[string]interface{})
	for _, value := range data {
		name := value.(map[string]interface{})["f14"]

		delete(value.(map[string]interface{}), "f14")

		pagedata[fmt.Sprintf("%s", name)] = value.(map[string]interface{})
	}
	return pagedata
}

func (this *DailymarketSpider) Run() map[string]map[string]interface{} {
	this.Init()

	data := make(map[string]map[string]interface{})
	urllist := this.GetUrls()
	for _, url := range urllist {
		content := this.ParseUrl(url)
		pagedata := this.GetStockdata(content)
		for key, value := range pagedata {
			data[key] = value
		}
	}
	return data
}
