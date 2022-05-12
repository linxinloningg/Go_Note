package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

/*
定义Query
//func (u *URL) Query() Values
req, err := http.NewRequest("GET", ""https://xxxx.xxxx.com.cn", nil)
q := req.URL.Query()
*/

/*
添加 Query
//func (v Values) Add(key string, value string)
	q.Add("per_page", "3")
	q.Add("page", "1")

*/
func GetHttp() {
	client := &http.Client{}
	apiURL := "http://gitlabcto.xxx.com.cn/api/v4/projects"

	req, err := http.NewRequest("GET", apiURL, nil)
	//添加查询参数
	q := req.URL.Query()
	q.Add("private_token", "-raXU2B8rVbRAFdYEqEg")
	q.Add("per_page", "3")
	q.Add("page", "1")
	req.URL.RawQuery = q.Encode()
	fmt.Println(req.URL.String())

	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		fmt.Printf("post failed, err:%v\n\n", err)
		return
	}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("get resp failed,err:%v\n\n", err)
		return
	}
	fmt.Println(string(b))
}
