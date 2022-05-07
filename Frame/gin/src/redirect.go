package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	// 使用了gin.Default()生成了一个实例,即 WSGI 应用程序
	r := gin.Default()
	r.GET("/index", func(c *gin.Context) {
		c.String(200, "redirect")
	})

	r.GET("/redirect", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/index")
	})

	r.GET("/goindex", func(c *gin.Context) {
		c.Request.URL.Path = "/index"
		r.HandleContext(c)
	})

	//r.Run()函数来让应用运行在本地服务器上，默认监听端口是 _8080_，
	//可以传入参数设置端口，例如r.Run(":9999")即运行在 _9999_端口
	err := r.Run() // listen and serve on 0.0.0.0:8080
	print(err)
}

/*
路由(Route)
路由方法有 GET, POST, PUT, PATCH, DELETE 和 OPTIONS，还有Any，可匹配以上任意类型的请求。
*/
