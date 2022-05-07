### Go 简明教程

* [第一个Gin程序](https://github.com/linxinloningg/Go_Note/blob/main/Frame/gin/src/test.go)
  * 首先，我们使用了`gin.Default()`生成了一个实例，这个实例即 WSGI 应用程序。
  * 接下来，我们使用`r.Get("/", ...)`声明了一个路由，告诉 Gin 什么样的URL 能触发传入的函数，这个函数返回我们想要显示在用户浏览器中的信息。
  * 最后用 `r.Run()`函数来让应用运行在本地服务器上，默认监听端口是 _8080_，可以传入参数设置端口，例如`r.Run(":9999")`即运行在 _9999_端口。

* 路由(Route)

  * [静态路由](https://github.com/linxinloningg/Go_Note/blob/main/Frame/gin/src/static_routing.go)
  * [动态路由](https://github.com/linxinloningg/Go_Note/blob/main/Frame/gin/src/dynamic_routing.go)
  * 获取Query参数
    * [GET](https://github.com/linxinloningg/Go_Note/blob/main/Frame/gin/src/query.go)
    * [POST](https://github.com/linxinloningg/Go_Note/blob/main/Frame/gin/src/post.go)
    * [POST和GET混合](https://github.com/linxinloningg/Go_Note/blob/main/Frame/gin/src/post_query.go)
    * [Map参数(字典参数)](https://github.com/linxinloningg/Go_Note/blob/main/Frame/gin/src/map.go)
    * [重定向(Redirect)](https://github.com/linxinloningg/Go_Note/blob/main/Frame/gin/src/redirect.go)

  路由方法有 **GET, POST, PUT, PATCH, DELETE** 和 **OPTIONS**，还有**Any**，可匹配以上任意类型的请求。

* [分组路由(Grouping Routes)](https://github.com/linxinloningg/Go_Note/blob/main/Frame/gin/src/grouping_routes.go)

  如果有一组路由，前缀都是`/api/v1`开头，是否每个路由都需要加上`/api/v1`这个前缀呢？答案是不需要，分组路由可以解决这个问题。

* [上传文件](https://github.com/linxinloningg/Go_Note/blob/main/Frame/gin/src/upload_file.go)

* [HTML模板(Template)](https://github.com/linxinloningg/Go_Note/blob/main/Frame/gin/src/template.go)

* [中间件(Middleware)](https://github.com/linxinloningg/Go_Note/blob/main/Frame/gin/src/middleware.go)

* [热加载调试 Hot Reload](https://github.com/linxinloningg/Go_Note/blob/main/Frame/gin/src/%E7%83%AD%E5%8A%A0%E8%BD%BD%E8%B0%83%E8%AF%95%20Hot%20Reload)



