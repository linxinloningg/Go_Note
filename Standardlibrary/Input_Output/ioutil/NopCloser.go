package main

import (
	"io"
	"io/ioutil"
)

/*
NopCloser用一个无操作的Close方法包装r返回一个ReadCloser接口。
*/
func main() {
	rc, ok := body.(io.ReadCloser)
	if !ok && body != nil {
		rc = ioutil.NopCloser(body)
	}
}
