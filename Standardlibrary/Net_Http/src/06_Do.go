package main

import (
	"log"
	"net/http"
)

/*
NewRequest
// func NewRequest(method string, url string, body io.Reader) (*Request, error)
	url := "https://xxxx.xxxx.com.cn"

	//Get 方法，没有body的情况
	getReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	print(getReq)

	//Post方法，body是json字串
	postReq, err := http.NewRequest("POST", url, strings.NewReader(`{"project_name": "liubei","metadata": {"public": "true"}}`))
	if err != nil {
		panic(err)
	}
	print(postReq)

	//Post方法，body是[]byte
	js := []byte("xxx")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(js))
	if err != nil {
		panic(err)
	}
	print(req)
*/
func main() {
	// func (c *Client) Do(req *Request) (*Response, error)
	var Client = &http.Client{}

	url := "http://gitlabcto.xxx.com.cn/api/v4/projects"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3776.0 Safari/537.36")

	response, _ := Client.Do(req)
	print(response)
}
