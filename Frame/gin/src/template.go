package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {

	// 使用了gin.Default()生成了一个实例,即 WSGI 应用程序
	r := gin.Default()

	type student struct {
		Name string
		Age  int8
	}

	r.LoadHTMLGlob("templates/test.html")

	stu1 := &student{Name: "Geektutu", Age: 20}
	stu2 := &student{Name: "Jack", Age: 22}

	r.GET("/test", func(c *gin.Context) {
		c.HTML(http.StatusOK, "test.tmpl", gin.H{
			"title":  "Gin",
			"stuArr": [2]*student{stu1, stu2},
		})
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
