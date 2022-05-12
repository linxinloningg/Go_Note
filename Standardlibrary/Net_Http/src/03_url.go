package main

import "net/url"

func main() {
	// func ParseRequestURI(rawURL string) (*URL, error)
	apiUrl := "http://gitlabcto.xxx.com.cn/api/v4/projects"
	url, err := url.ParseRequestURI(apiUrl)
	if err != nil {
		panic(err)
	}
	print(url)
}
