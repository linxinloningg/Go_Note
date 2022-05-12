package main

import "net/http"

func main() {
	var DefaultClient = &http.Client{}
	print(DefaultClient)
}
