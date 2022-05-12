package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func HiHttp() {
	resp, err := http.Get("http://www.baidu.com/")
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	//body是一个二进制组，用string(body[:])转换为字串
	fmt.Printf("%s", string(body[:]))
}

func main() {
	HiHttp()
}
