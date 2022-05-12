package main

import (
	"fmt"
	"net/http"
	"strings"
)

func PostHttpWithJson() {
	client := &http.Client{}
	url := "https://harbocto.xxx.com.cn/api/v2.0/projects"
	data := `{"project_name": "liubei","metadata": {"public": "true"}}`
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic 5YiY5aSHOuaIkeS4jeS8muWRiuivieS9oOWvhueggQ==")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Printf("执行状态为：%d", resp.StatusCode)

}
