package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	// 使用了gin.Default()生成了一个实例,即 WSGI 应用程序
	r := gin.Default()

	/*
		r.GET("/get", func(c *gin.Context) {
				ids := c.QueryMap("ids")

				c.JSON(http.StatusOK, gin.H{
					"ids":   ids,
				})
			})
	*/
	//Map参数(字典参数)

	r.POST("/form", func(c *gin.Context) {
		names := c.PostFormMap("names")

		c.JSON(http.StatusOK, gin.H{
			"names": names,
		})
	})

	//r.Run()函数来让应用运行在本地服务器上，默认监听端口是 _8080_，
	//可以传入参数设置端口，例如r.Run(":9999")即运行在 _9999_端口
	err := r.Run() // listen and serve on 0.0.0.0:8080
	print(err)
	//curl -g "http://localhost:8080/get?ids[Jack]=001&ids[Tom]=002"

	// curl http://localhost:8080/form -X POST -d 'names[a]=Sam'
}

/*
路由(Route)
路由方法有 GET, POST, PUT, PATCH, DELETE 和 OPTIONS，还有Any，可匹配以上任意类型的请求。
*/
