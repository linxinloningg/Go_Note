package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	// 使用了gin.Default()生成了一个实例,即 WSGI 应用程序
	r := gin.Default()

	// 单个文件
	r.POST("/upload", func(c *gin.Context) {
		file, _ := c.FormFile("file")
		// c.SaveUploadedFile(file, dst)
		c.String(http.StatusOK, "%s uploaded!", file.Filename)
	})

	// 多个文件
	r.POST("/multiupload", func(c *gin.Context) {
		// Multipart form
		form, _ := c.MultipartForm()
		files := form.File["upload[]"]

		for _, file := range files {
			log.Println(file.Filename)
			// c.SaveUploadedFile(file, dst)
		}
		c.String(http.StatusOK, "%d files uploaded!", len(files))
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
