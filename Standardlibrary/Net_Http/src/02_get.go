package main

import "net/http"

func main() {
	var Client = &http.Client{}
	response, err := Client.Get("http://gitlabcto.xxx.com.cn/api/v4/projects")
	if err != nil {
		panic(err)
	}
	print(response)
}
