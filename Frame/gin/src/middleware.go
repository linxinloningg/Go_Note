package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		// 给Context实例设置一个值
		c.Set("geektutu", "1111")
		// 请求前
		c.Next()
		// 请求后
		latency := time.Since(t)
		log.Print(latency)
	}
}
func main() {
	// 使用了gin.Default()生成了一个实例,即 WSGI 应用程序
	r := gin.Default()

	// 作用于全局
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 作用于单个路由
	r.GET("/benchmark", Logger(), benchEndpoint)

	/*
		// 作用于某个组
		authorized := r.Group("/")
		authorized.Use(AuthRequired())
		{
			authorized.POST("/login", loginEndpoint)
			authorized.POST("/submit", submitEndpoint)
		}

	*/

	//r.Run()函数来让应用运行在本地服务器上，默认监听端口是 _8080_，
	//可以传入参数设置端口，例如r.Run(":9999")即运行在 _9999_端口
	err := r.Run() // listen and serve on 0.0.0.0:8080
	print(err)
}

/*
路由(Route)
路由方法有 GET, POST, PUT, PATCH, DELETE 和 OPTIONS，还有Any，可匹配以上任意类型的请求。
*/
