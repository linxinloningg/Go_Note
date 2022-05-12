#### Client结构体

```go
type Client struct {
    Transport     RoundTripper
    CheckRedirect func(req *Request, via []*Request) error
    Jar           CookieJar
    Timeout       time.Duration
}
```

Client结构体有以下方法：

* ```
  send(req *http.Request, deadline time.Time) (resp *http.Response, didTimeout func() bool, err error)
  ```

* ```
  deadline() time.Time
  ```

* ```
  transport() http.RoundTripper
  ```

* ```
  Get(url string) (resp *http.Response, err error)
  ```

* ```
  checkRedirect(req *http.Request, via []*http.Request) error
  ```

* ```
  Do(req *http.Request) (*http.Response, error)
  ```

* ```
  do(req *http.Request) (retres *http.Response, reterr error)
  ```

* ```
  makeHeadersCopier(ireq *http.Request) func(*http.Request)
  ```

* ```
  Post(url string, contentType string, body io.Reader) (resp *http.Response, err error)
  ```

* ```
  PostForm(url string, data url.Values) (resp *http.Response, err error)
  ```

* ```
  Head(url string) (resp *http.Response, err error)
  ```

* ```
  CloseIdleConnections()
  ```

#### 初始化结构体

```go
var DefaultClient = &http.Client{}
```

####  Get() 方法

[示例](https://github.com/linxinloningg/Go_Note/blob/main/Standardlibrary/Net_Http/src/02_get.go)

```go
// func Get(url string) (resp *Response, err error)
resp, err := http.Get("http://xxxx.xxx.com")
```

#### url

[示例](https://github.com/linxinloningg/Go_Note/blob/main/Standardlibrary/Net_Http/src/03_url.go)

```go
// func ParseRequestURI(rawURL string) (*URL, error)
apiUrl := "http://gitlabcto.xxx.com.cn/api/v4/projects"
u,err := url.ParseRequestURI(apiUrl)
```

#### params（参数）

作用：在url后拼接参数 [示例](https://github.com/linxinloningg/Go_Note/blob/main/Standardlibrary/Net_Http/src/04_params.go)

初始化

```go
data := url.Values{}
```

设定参数

```go
data.Set("key","value")
```

将参数添加到url

```go
url.RawQuery = data.Encode()
```

#### Post()方法

[示例](https://github.com/linxinloningg/Go_Note/blob/main/Standardlibrary/Net_Http/src/05_post.go)

```go
//func Post(url string, contentType string, body io.Reader) (resp *Response, err error)
http.Post(url, contentType, strings.NewReader(data))
```

#### Do() 方法

[示例](https://github.com/linxinloningg/Go_Note/blob/main/Standardlibrary/Net_Http/src/06_Do.go)

```go
//func (c *Client) Do(req *Request) (*Response, error)
resp, _ := client.Do(req)
```

>#### NewRequest
>
>语法
>
>```go
>func NewRequest(method string, url string, body io.Reader) (*Request, error)
>```
>
>示例：Get 方法，没有body的情况
>
>```go
>req, err := http.NewRequest("GET", ""https://xxxx.xxxx.com.cn", nil)
>```
>
>示例：Post方法，body是json字串
>
>```go
>req, err := http.NewRequest("POST", url, strings.NewReader(`{"project_name": "liubei","metadata": {"public": "true"}}`)
>```
>
>示例：Post方法，body是[]byte
>
>```go
>var js []byte = xxxxxx
>req, err := http.NewRequest("POST", url, bytes.NewBuffer(js))
>```

#### Header 添加头信息

[示例](https://github.com/linxinloningg/Go_Note/blob/main/Standardlibrary/Net_Http/src/07_with_header.go)

>Set 函数(相同key的值会被替换为最新的)
>
>```go
>//func (h Header) Set(key string, value string)
>req.Header.Set("aaa", "111")
>req.Header.Set("aaa", "222")
>req.Header.Set("aaa", "333")
>fmt.Printf("%+v\n",req.Header)
>```
>
>Add 函数(相同key的值会被追加)
>
>```go
>//func (h Header) Add(key string, value string)
>req.Header.Add("aaa", "111")
>req.Header.Add("aaa", "222")
>req.Header.Add("aaa", "333")
>fmt.Printf("%+v\n",req.Header)
>```

#### Query

[示例](https://github.com/linxinloningg/Go_Note/blob/main/Standardlibrary/Net_Http/src/08_with_query.go)

```go
//定义Query
// func (u *URL) Query() Values
req, err := http.NewRequest("GET", ""https://xxxx.xxxx.com.cn", nil)
q := req.URL.Query()
//添加 Query
//func (v Values) Add(key string, value string)
```