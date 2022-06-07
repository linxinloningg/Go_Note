### 函数

* ```go
  func NewSpider(pageinst page_processer.PageProcesser, taskname string) *Spider
  ```

  Spider 是所有其他模块的调度器模块，如下载器、管道、调度器等。 taskname也可以是空字符串，也可以在Pipeline中用来记录哪个task爬取的结果； 

* ```go
  func (this * Spider ) AddPipeline(p pipeline . Pipeline ) * Spider 
  ```

* ```go
  func (this * Spider ) AddRequest(req * request . Request ) * Spider 
  ```

  将请求添加到计划 

* ```go
  func (this * Spider ) AddRequests(reqs []* request . Request ) * Spider 
  ```

* ```go
  func (this * Spider ) AddUrl(url string , respType string ) * Spider
  ```

* ```go
  func (this * Spider ) AddUrlEx(url string , respType string , headerFile string , proxyHost string ) * Spider 
  ```

* ```go
  func (this * Spider ) AddUrlWithHeaderFile(url string , respType string , headerFile string ) * Spider 
  ```

* ```go
  func (this * Spider ) AddUrls(urls [] string , respType string ) * Spider 
  ```

* ```go
  func (this * Spider ) AddUrlsEx(urls [] string , respType string , headerFile string , proxyHost string ) * Spider 
  ```

* ```go
  func (this * Spider ) AddUrlsWithHeaderFile(urls [] string , respType string , headerFile string ) * Spider 
  ```

* ```go
  func (this * Spider ) CloseFileLog() * Spider 
  ```

  CloseFileLog 关闭文件日志。 

* ```go
  func (this * Spider ) CloseStrace() * Spider 
  ```

  CloseStrace 关闭 strace。 

* ```go
  func (this *Spider) Get(url string, respType string) *page_items.PageItems
  ```

  处理一个 url 并返回 PageItems。 

* ```go
  func (this *Spider) GetAll(urls []string, respType string) []*page_items.PageItems
  ```

  处理几个 url 并返回 PageItems 切片。 

* ```go
  func (this *Spider) GetAllByRequest(urls []string, respType string) []*page_items.PageItems
  ```

  处理几个 url 并返回 PageItems 切片 

* ```go
  func (this *Spider) GetByRequest(reqs []*request.Request) []*page_items.PageItems
  ```

  处理一个 url 并返回具有其他设置的 PageItems。

* ```go
  func (this *Spider) GetDownloader() downloader.Downloader
  ```

* ```go
  func (this *Spider) GetExitWhenComplete() bool
  ```

* ```go
  func (this *Spider) GetScheduler() scheduler.Scheduler
  ```

* ```go
  func (this *Spider) GetThreadnum() uint
  ```

* ```go
  func (this *Spider) OpenFileLog(filePath string) *Spider
  ```

  OpenFileLog 初始化日志路径并打开日志。  如果打开日志，spider 中的错误信息或其他有用信息将记录在文件路径的文件中。   日志命令是 mlog.LogInst().LogError("info") 或 mlog.LogInst().LogInfo("info")。  Spider 的默认日志是关闭的。  文件路径是绝对路径。

* ```go
  func (this *Spider) OpenFileLogDefault() *Spider
  ```

  OpenFileLogDefault 使用默认文件路径打开文件日志，例如“WD/log/log.2014-9-1”。

* ```go
  func (this *Spider) OpenStrace() *Spider
  ```

  OpenStrace 打开 strace，在屏幕上输出进度信息。  Spider 的默认 strace 已打开。

* ```go
  func (this *Spider) Run()
  ```

* ```go
  func (this *Spider) SetDownloader(d downloader.Downloader) *Spider
  ```

* ```go
  func (this *Spider) SetExitWhenComplete(e bool) *Spider
  ```

  如果在每个爬网任务完成后退出。  如果你想让Spider一直在内存中并从外部添加 url，你可以将其设置为 true。

* ```go
  func (this *Spider) SetScheduler(s scheduler.Scheduler) *Spider
  ```

* ```go
  func (this *Spider) SetSleepTime(sleeptype string, s uint, e uint) *Spider
  ```

  SetSleepTime 设置每次爬网任务后的睡眠时间。  单位是毫秒。  如果 sleeptype 是“固定的”，则 s 是睡眠时间，而 e 是useless。  如果 sleeptype 为“rand”，则睡眠时间为 s 和 e 之间的。

* ```go
  func (this *Spider) SetThreadnum(i uint) *Spider
  ```

* ```go
  func (this *Spider) Taskname() string
  ```

