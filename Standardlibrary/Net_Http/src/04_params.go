package main

import "net/url"

func main() {
	// 初始化
	data := url.Values{}

	//设定参数
	data.Set("key", "value")

	//将参数添加到url
	// url.RawQuery = data.Encode()
}

/*
func HiHttp(){
	apiUrl := "http://gitlabcto.xxx.com.cn/api/v4/projects"
	// URL param
	data := url.Values{}
	data.Set("private_token", "-raXU2B8rVbRAFxxxxxx")
	data.Set("per_page", "100")
	data.Set("page", "30")
	//把string转换为url
	u, err := url.ParseRequestURI(apiUrl)
	if err != nil {
		fmt.Printf("parse url requestUrl failed,err:%v\n", err)
	}
	//将参数添加到请求url
	u.RawQuery = data.Encode() // URL encode
	fmt.Println(u.String())
	resp, err := http.Get(u.String())
	if err != nil {
		fmt.Println("post failed, err:%v\n", err)
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("get resp failed,err:%v\n", err)
		return
	}
	fmt.Println(string(b))
}

*/
