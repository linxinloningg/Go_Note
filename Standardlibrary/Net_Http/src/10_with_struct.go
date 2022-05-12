package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Body struct {
	ProjectName string `json:"project_name"`
	Metadata    struct {
		Public string `json:"public"`
	} `json:"metadata"`
}

func PostHttpWithStruct() {
	client := &http.Client{}
	url := "https://harbocto.xxx.com.cn/api/v2.0/projects"

	var s Body
	s.ProjectName = "liubei-01"
	s.Metadata.Public = "true"

	js, err := json.MarshalIndent(&s, "", "\t")
	fmt.Println(string(js))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(js))
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	req.Header.Set("Authorization", "Basic 5YiY5aSHOuaIkeS4jeS8muWRiuivieS9oOWvhueggQ==")
	req.Header.Set("Content-Type", "application/json")

	resp, _ := client.Do(req)
	defer resp.Body.Close()

	fmt.Printf("执行状态为：%d", resp.StatusCode)

}
