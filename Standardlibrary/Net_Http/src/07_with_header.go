package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

/*
Set 函数
//func (h Header) Set(key string, value string)
	req.Header.Set("aaa", "111")
	req.Header.Set("aaa", "222")
	req.Header.Set("aaa", "333")
	fmt.Printf("%+v\n",req.Header)

*/

/*
Add 函数
//func (h Header) Add(key string, value string)
	req.Header.Add("aaa", "111")
	req.Header.Add("aaa", "222")
	req.Header.Add("aaa", "333")
	fmt.Printf("%+v\n",req.Header)

*/
func GetHttpWithHeader() {
	client := &http.Client{}
	apiURL := "https://harbocto.xxxx.com.cn/api/v2.0/projects/9"

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	req.Header.Add("Authorization", "Basic 5YiY5aSHOuaIkeS4jeS8muWRiuivieS9oOWvhueggQ==")
	req.Header.Add("Content-Type", "application/json")

	resp, _ := client.Do(req)
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("get resp failed,err:%v\n\n", err)
		return
	}
	fmt.Println(string(b))
}


