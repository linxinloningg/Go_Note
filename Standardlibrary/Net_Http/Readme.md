# Client结构体

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

## 初始化结构体

```
var DefaultClient = &http.Client{}
```

#  Get() 方法

```
// func Get(url string) (resp *Response, err error)
resp, err := http.Get("http://xxxx.xxx.com")
```

### url

```
// func ParseRequestURI(rawURL string) (*URL, error)
apiUrl := "http://gitlabcto.xxx.com.cn/api/v4/projects"
u,err := url.ParseRequestURI(apiUrl)
```

### params（参数）

作用：在url后拼接参数

#### 初始化

```go
data := url.Values{}
```

#### 设定参数

```go
data.Set("key","value")
```

#### 将参数添加到url

```go
url.RawQuery = data.Encode()
```

