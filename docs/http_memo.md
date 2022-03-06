# http について
## GET
### 1. http.Get 関数を実行する
一番単純だけど URL パラメータぐらいしか付与できなさそう  
[func Get](https://pkg.go.dev/net/http@go1.17.8#Get)  
>Get is a wrapper around DefaultClient.Get.  
>>Getは、DefaultClient.Getのラッパーです。  

URL パラメータを付けるには net/url の `values := url.Values{}` を使う  
`values.Add("key", "value")` して `values.Encode()` で URL パラメータ(`key=value&key=value` の形)ができるから  
`url?` の後ろに加える  
### 2. Client 型の Get(url) メソッドを実行する
`type Request` 自体は Go に生成してもらうので `type Client` の調整しかできない  
`type Client` には トランスポート, リダイレクト, クッキー, タイムアウト を設定できる  
[type Client](https://pkg.go.dev/net/http@go1.17.8#Client)  
### 3. Client 型の Do(request) メソッドを実行する
`type Request` から自分で作るため リクエストヘッダのフィールド や form データのフィールドなど 設定できる  

### 調べる
[func (*Client) Do](https://pkg.go.dev/net/http@go1.17.8#Client.Do) である `func (c *Client) Do(req *Request) (*Response, error)` を調べる  
>Generally Get, Post, or PostForm will be used instead of Do.  
>>一般的にはDoの代わりにGet、Post、PostFormが使用される。  

と書かれていたから 次は [func (*Client) Get](https://pkg.go.dev/net/http@go1.17.8#Client.Get) である `func (c *Client) Get(url string) (resp *Response, err error)` を調べる  
中を見たら url を使って `NewRequest("GET", url, nil)` で `type Request` を生成して `c.Do(req)` を返していた  
リクエストヘッダのフィールドは `type Request` に存在するため Get() の中で request を生成してはダメ  
よって リクエストヘッダを扱いたいときは Do() を使うことになる  

リクエストヘッダを設定する場合 Request 内の Header フィールドに設定する  
Set() や Add() を使う  
`req.Header.Add("User-Agent", "foobar")`

## POST
### 1. http.Post 関数を実行する
[func Post](https://pkg.go.dev/net/http@go1.17.8#Post)  
>Post is a wrapper around DefaultClient.Post.  
>To set custom headers, use NewRequest and DefaultClient.Do.  
>>Postは、DefaultClient.Postのラッパーです。  
>>カスタムヘッダを設定するには、NewRequestとDefaultClient.Doを使用します。  

`func Post(url, contentType string, body io.Reader) (resp *Response, err error)`  
### 2. http.PostForm(url, data) 関数を実行する
[func PostForm](https://pkg.go.dev/net/http@go1.17.8#PostForm)  
`func PostForm(url string, data url.Values) (resp *Response, err error)`  
### 3. Client 型の Post(url, bodyType, body) メソッドを実行する
[func (*Client) Post](https://pkg.go.dev/net/http@go1.17.8#Client.Post)  
`func (c *Client) Post(url, contentType string, body io.Reader) (resp *Response, err error)`  
>To set custom headers, use NewRequest and DefaultClient.Do.  
>>カスタムヘッダを設定するには、NewRequestとDefaultClient.Doを使用します。  

中身を見たら  
`req, err := NewRequest("POST", url, body)` してたり  
`req.Header.Set("Content-Type", contentType)` してたり  
だった
### 4. Client 型の PostForm(url, data) メソッドを実行する
[func (*Client) PostForm](https://pkg.go.dev/net/http@go1.17.8#Client.PostForm)  
`func (c *Client) PostForm(url string, data url.Values) (resp *Response, err error)`  

>PostForm issues a POST to the specified URL, with data's keys and values URL-encoded as the request body.  
>The Content-Type header is set to application/x-www-form-urlencoded.  
>To set other headers, use NewRequest and Client.Do.  
>>PostForm は指定された URL に POST を発行し、データのキーと値はリクエストボディとして URL エンコードされます。  
>>Content-Typeヘッダはapplication/x-www-form-urlencodedに設定されています。  
>>その他のヘッダを設定するには、NewRequestとClient.Doを使用します。  

中身を見たら  
`c.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))` を返しているだけ  
### 5. Client 型の Do(request) メソッドを実行する
一番細かく設定できる  

## まとめ
GET も POST も `func (*Client) Do` が根本にある  

## その他
```Go
func (r *Request) SetBasicAuth(username, password string) {
	r.Header.Set("Authorization", "Basic "+basicAuth(username, password))
}
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
```

curl の -u は BASIC 認証に該当する  
`-u user:pass`  
だから Go 言語では SetBasicAuth() を使えばよい  

## 参考
[Go net/httpパッケージの概要とHTTPクライアント実装例](https://qiita.com/jpshadowapps/items/463b2623209479adcd88)  
