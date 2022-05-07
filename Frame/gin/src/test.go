package main

import "github.com/gin-gonic/gin"

func main() {
	app := gin.Default()

	fun := func(c *gin.Context) { c.String(200, "Hello, World") }

	app.GET("/", fun)

	err := app.Run()
	print(err)
}
