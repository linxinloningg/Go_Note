package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

/*
此处是http包提供的一个简单的post方法，不能添加header，如果要添加header需要使用Do()方法。
因为http包的Post() 方法没有header，因此harbor的认证是过不去的。
建议使用Do()方法
*/

func PostHttp() {
	url := "https://harbocto.xxx.com.cn/api/v2.0/projects"
	contentType := "application/json"
	data := `{"project_name": "liuBei","metadata": {"public": "true"}}`
	resp, err := http.Post(url, contentType, strings.NewReader(data))
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Println(string(b))
}

func main() {
	PostHttp()
}
