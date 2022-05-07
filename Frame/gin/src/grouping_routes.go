package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
如果有一组路由，前缀都是/api/v1开头，是否每个路由都需要加上/api/v1这个前缀呢？
答案是不需要，分组路由可以解决这个问题。利用分组路由还可以更好地实现权限控制
例如将需要登录鉴权的路由放到同一分组中去，简化权限控制。
*/
func main() {
	// 使用了gin.Default()生成了一个实例,即 WSGI 应用程序
	r := gin.Default()

	/*
		// r.Get("/", ...)声明了一个路由，告诉 Gin 什么样的URL 能触发传入的函数，
		这个函数返回我们想要显示在用户浏览器中的信息。
			r.GET("/", func(c *gin.Context) {
				c.String(200, "Hello, Geektutu")
			})
	*/
	// group routes 分组路由
	defaultHandler := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"path": c.FullPath(),
		})
	}
	// group: v1
	v1 := r.Group("/v1")
	{
		v1.GET("/posts", defaultHandler)
		v1.GET("/series", defaultHandler)
	}
	// group: v2
	v2 := r.Group("/v2")
	{
		v2.GET("/posts", defaultHandler)
		v2.GET("/series", defaultHandler)
	}

	//r.Run()函数来让应用运行在本地服务器上，默认监听端口是 _8080_，
	//可以传入参数设置端口，例如r.Run(":9999")即运行在 _9999_端口
	err := r.Run() // listen and serve on 0.0.0.0:8080
	print(err)

	//curl http://localhost:8080/v1/posts
	//curl http://localhost:8080/v2/posts
}

/*
路由(Route)
路由方法有 GET, POST, PUT, PATCH, DELETE 和 OPTIONS，还有Any，可匹配以上任意类型的请求。
*/
