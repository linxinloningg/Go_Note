package RealDataSpider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type Spider interface {
	Init()
	GetUrls(codes []string, parameters string) map[string]map[string]string
	ParseUrl(params map[string]string) string
	GetStockdata(code string, content string) map[string]map[string]string
	Run(codes []string, parameters string) map[string]map[string]string
}

type RealDataSpider struct {
	GsApi      string
	headers    map[string]string
	datatype   map[string]string
	parameters string
}

func (this *RealDataSpider) Init() {
	this.GsApi = "http://push2.eastmoney.com/api/qt/stock/get"

	this.headers = map[string]string{"User-Agent": "Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.82 Safari/537.36"}

	this.datatype = map[string]string{"f43": "最新价", "f44": "最高", "f45": "最低", "f46": "今开", "f47": "成交量", "f48": "成交额", "f50": "量比", "f51": "涨停", "f52": "跌停", "f55": "收益", "f57": "代码", "f58": "名称", "f60": "昨收", "f71": "均价", "f92": "美股净资产", "f105": "净利润", "f116": "总市值", "f117": "流通市值", "f162": "市盈动", "f167": "市净", "f168": "换手", "f169": "涨跌", "f170": "涨幅", "f173": "ROE", "f183": "总营收", "f186": "毛利率", "f187": "净利率", "f188": "负债率", "f191": "委比", "f192": "委差", "f193": "主力净比", "f194": "超大单净比", "f195": "大单净比", "f196": "中单净比", "f197": "小单净比"}
	// this.parameters = append(this.parameters, "f43", "f44", "f45", "f46", "f47", "f48", "f50", "f51", "f52", "f55", "f57", "f58", "f60", "f71", "f92", "f105", "f116", "f117", "f162", "f167", "f168", "f169", "f170", "f173", "f183", "f186", "f187", "f188", "f191", "f192", "f193", "f194", "f195", "f196", "f197")
	//this.parameters = []string{"f43", "f44", "f45", "f46", "f47", "f48", "f50", "f51", "f52", "f55", "f57", "f58", "f60", "f71", "f92", "f105", "f116", "f117", "f162", "f167", "f168", "f169", "f170", "f173", "f183", "f186", "f187", "f188", "f191", "f192", "f193", "f194", "f195", "f196", "f197"}
	this.parameters = "f43,f44,f45,f46,f47,f48,f50,f51,f52,f55,f57,f58,f60,f71,f92,f105,f116,f117,f162,f167,f168,f169,f170,f173,f183,f186,f187,f188,f191,f192,f193,f194,f195,f196,f197"
}

func (this *RealDataSpider) GetUrls(codes []string, parameters string) map[string]map[string]string {

	if parameters != "default" {
		this.parameters = parameters
	}
	params_dict := make(map[string]map[string]string)
	for _, code := range codes {
		params := make(map[string]string)
		if code[:1] == "6" {
			/*
				{
					"ut": "fa5fd1943c7b386f172d6893dbfba10b",
					"invt": "2",
					"fltt": "2",
					"fields": self.parameters,
					"secid": "%s.%s" % (1, code),
					"cb": "jQuery112405831440079032297_1578892365285",
					"_": "1578892365407"
				}
			*/
			params["ut"] = "fa5fd1943c7b386f172d6893dbfba10b"
			params["invt"] = "2"
			params["fltt"] = "2"
			params["fields"] = this.parameters
			params["secid"] = fmt.Sprintf("%d.%s", 1, code)
			params["cb"] = "jQuery112405831440079032297_1578892365285"
			params["_"] = "1578892365407"
			params_dict[code] = params
		}
		if code[:1] == "3" || code[:1] == "0" {
			/*
				{
					"ut": "fa5fd1943c7b386f172d6893dbfba10b",
					"invt": "2",
					"fltt": "2",
					"fields": self.parameters,
					"secid": "%s.%s" % (0, code),
					"cb": "jQuery112405831440079032297_1578892365285",
					"_": "1578892365407"
				}
			*/
			params["ut"] = "fa5fd1943c7b386f172d6893dbfba10b"
			params["invt"] = "2"
			params["fltt"] = "2"
			params["fields"] = this.parameters
			params["secid"] = fmt.Sprintf("%d.%s", 0, code)
			params["cb"] = "jQuery112405831440079032297_1578892365285"
			params["_"] = "1578892365407"
			params_dict[code] = params
		}
	}
	return params_dict
}

func (this *RealDataSpider) ParseUrl(params map[string]string) string {
	client := &http.Client{Timeout: time.Second * 5}

	req, err := http.NewRequest("GET", this.GsApi, nil)

	param := req.URL.Query()

	if err != nil {
		log.Fatalln(err)
	}
	for key, value := range params {
		param.Set(key, value)
	}

	//将参数添加到请求url
	req.URL.RawQuery = param.Encode() // URL encode

	req.Header.Add("User-Agent", this.headers["User-Agent"])

	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	response, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil && param == nil {
		log.Fatalln(err)
	}
	return fmt.Sprintf("%s", data)
}

func (this *RealDataSpider) GetStockdata(code string, content string) map[string]map[string]string {
	stockdata := make(map[string]map[string]string)
	content = strings.Replace(content, "jQuery112405831440079032297_1578892365285", "", -1)
	content = strings.Replace(content, "(", "", -1)
	content = strings.Replace(content, ")", "", -1)
	content = strings.Replace(content, ";", "", -1)

	var mapData map[string]interface{}
	err := json.Unmarshal([]byte(content), &mapData)
	if err != nil {
		panic(err)
	}

	getdata := make(map[string]string)
	for _, parameter := range strings.Split(this.parameters, ",") {
		data := mapData["data"].(map[string]interface{})[parameter]

		switch data.(type) {
		case float64:
			{
				data := data.(float64)
				getdata[this.datatype[parameter]] = fmt.Sprintf("%.2f", data)
			}
		case string:
			{
				data := data.(string)
				getdata[this.datatype[parameter]] = data
			}
		}
		stockdata[code] = getdata

	}

	return stockdata
}

/*func (this *RealDataSpider) SaveStockdata(stockdata map[string]string) map[string]string {
	for code, datadict := range stockdata {
		for parameter, data := range datadict {
			for datatype, zhname := range this.datatype {
				if parameter == datatype {
					parameter = zhname
				}
			}
			fmt.Sprintf("[%s] [%s] [%s]: %s。", code, "datetime.now()", parameter, data)
		}
	}
}*/
func (this *RealDataSpider) Run(codes []string, parameters string) map[string]map[string]string {
	this.Init()
	data := make(map[string]map[string]string)
	gpurl_dict := this.GetUrls(codes, parameters)
	for code, url := range gpurl_dict {
		content := this.ParseUrl(url)
		stockdata := this.GetStockdata(code, content)
		data[code] = stockdata[code]

	}
	return data
}
