# ミドルウェア について
ミドルウェアとは リクエストを事前に解釈し エラー処理をまとめて行ったり 認証チェックをしたり する機構のこと  
標準ライブラリでは [http.TimeoutHandler](https://pkg.go.dev/net/http#TimeoutHandler) が該当するらしい  

```go
func TimeoutHandler(h Handler, dt time.Duration, msg string) Handler {
	return &timeoutHandler{
		handler: h,
		body:    msg,
		dt:      dt,
	}
}

type timeoutHandler struct {
	handler Handler
	body    string
	dt      time.Duration

	// When set, no context will be created and this context will
	// be used instead.
	testContext context.Context
}
```
>TimeoutHandler は、与えられた制限時間内に h を実行する Handler を返します。  
>新しいハンドラは各リクエストを処理するために h.ServeHTTP を呼び出しますが、もし呼び出しがその制限時間を超えて実行されると、ハンドラは 503 Service Unavailable エラーとそのボディに与えられたメッセージを持って応答します。(msg が空の場合、適切なデフォルトメッセージが送信されます)  
>このようなタイムアウトの後、 h によるその ResponseWriter への書き込みは、 ErrHandlerTimeout を返します。  

timeoutHandler は 2つのメソッドを持っているっぽい  
```go
func (h *timeoutHandler) errorBody() string
func (h *timeoutHandler) ServeHTTP(w ResponseWriter, r *Request)
```

field はすべて private だし 戻り値は Handler だから 外部からは `(h *timeoutHandler) ServeHTTP` しか見えないと思われる  

使い方は クライアントとしてレスポンスが無いときのタイムアウトと サーバとしての処理が長すぎて返すタイムアウトがあるっぽい  

### http.TimeoutHandler の使い方
[Overview](https://pkg.go.dev/net/http#pkg-overview) より  
クライアントとしてのタイムアウトのサンプル  
>プロキシ、TLS設定、キープアライブ、圧縮、その他の設定を制御するには、Transportを作成します。

```go
tr := &http.Transport{
	MaxIdleConns:       10,
	IdleConnTimeout:    30 * time.Second,
	DisableCompression: true,
}
client := &http.Client{Transport: tr}
resp, err := client.Get("https://example.com")
```

サーバとしてのタイムアウトのサンプル  
>カスタムサーバーを作成することで、サーバーの動作をより細かく制御することができます。

```go
s := &http.Server{
	Addr:           ":8080",
	Handler:        myHandler,
	ReadTimeout:    10 * time.Second,
	WriteTimeout:   10 * time.Second,
	MaxHeaderBytes: 1 << 20,
}
log.Fatal(s.ListenAndServe())
```

http.Server の ReadTimeout フィールドの説明を読む  
>ReadTimeout is the maximum duration for reading the entire request, including the body.  
>A zero or negative value means there will be no timeout.  
>Because ReadTimeout does not let Handlers make per-request decisions on each request body's acceptable deadline or upload rate, most users will prefer to use ReadHeaderTimeout.  
>It is valid to use them both.  
>>ReadTimeout は、ボディを含むリクエスト全体を読み取るための最大時間です。  
>>0または負の値は、タイムアウトがないことを意味する。ReadTimeoutは、各リクエストボディの許容期限やアップロードレートについて、 ハンドラにリクエストごとの決定をさせないので、ほとんどのユーザーは ReadHeaderTimeout を使うことを好むだろう。  
>>両方を使用することは有効です。  

http.Server の ReadHeaderTimeout フィールドの説明を読む  
>ReadHeaderTimeout is the amount of time allowed to read request headers.  
>The connection's read deadline is reset after reading the headers and the Handler can decide what is considered too slow for the body.  
>If ReadHeaderTimeout is zero, the value of ReadTimeout is used.  
>If both are zero, there is no timeout.  
>>ReadHeaderTimeout は、リクエストヘッダを読み取るために許容される時間です。  
>>接続の読み取り期限はヘッダを読み取った後にリセットされ、Handlerはボディに対して遅すぎると考えられるものを決定することができます。  
>>ReadHeaderTimeout がゼロの場合、 ReadTimeout の値が使用されます。  
>>両方がゼロの場合、タイムアウトは発生しない。  

http.Server の WriteTimeout フィールドの説明を読む  
>WriteTimeout is the maximum duration before timing out writes of the response.  
>It is reset whenever a new request's header is read.  
>Like ReadTimeout, it does not let Handlers make decisions on a per-request basis.  
>A zero or negative value means there will be no timeout.  
>>WriteTimeout は、レスポンスの書き込みがタイムアウトするまでの最大時間です。  
>>新しいリクエストのヘッダーが読み込まれるたびにリセットされます。  
>>ReadTimeout と同様に、 Handler にリクエストごとの判断をさせない。  
>>ゼロまたは負の値は、タイムアウトがないことを意味する。  

### 実際にどこでタイムアウトの判定をしているのか
`(c *conn) serve(ctx context.Context)` の中でも呼ばれている `(c *conn) readRequest(ctx context.Context) (w *response, err error)` の中で 判定しているっぽい  
```go
// 一部抜粋
func (c *conn) readRequest(ctx context.Context) (w *response, err error) {
	t0 := time.Now()
	if d := c.server.readHeaderTimeout(); d > 0 {
		hdrDeadline = t0.Add(d)
	}
	if d := c.server.ReadTimeout; d > 0 {
		wholeReqDeadline = t0.Add(d)
	}

	// Adjust the read deadline if necessary.
	if !hdrDeadline.Equal(wholeReqDeadline) {
		c.rwc.SetReadDeadline(wholeReqDeadline)
	}
}

func (s *Server) readHeaderTimeout() time.Duration {
	if s.ReadHeaderTimeout != 0 {
		return s.ReadHeaderTimeout
	}
	return s.ReadTimeout
}

// 一部抜粋
type conn struct {
	rwc net.Conn
}

type Conn interface {
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
	Close() error
	LocalAddr() Addr
	RemoteAddr() Addr
	SetDeadline(t time.Time) error
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
}
```
`(s *Server) readHeaderTimeout()` は ReadHeaderTimeout の設定があれば ReadHeaderTimeout なければ ReadTimeout を返すだけ  

Conn interface の SetReadDeadline メソッドの[説明](https://pkg.go.dev/net#Conn.SetReadDeadline)を読む  
>SetReadDeadline sets the deadline for future Read calls and any currently-blocked Read call.  
>A zero value for t means Read will not time out.  
>>SetReadDeadline は、将来の Read コールおよび現在ブロックされている Read コールのデッドラインを設定します。  
>>t に0を指定すると、 Read がタイムアウトしないことを意味します。  

実装は net/net.go  
```go
// SetReadDeadline implements the Conn SetReadDeadline method.
// SetReadDeadline は Conn の SetReadDeadline メソッドを実装したものです。
func (c *conn) SetReadDeadline(t time.Time) error {
	if !c.ok() {
		return syscall.EINVAL
	}
	if err := c.fd.SetReadDeadline(t); err != nil {
		return &OpError{Op: "set", Net: c.fd.net, Source: nil, Addr: c.fd.laddr, Err: err}
	}
	return nil
}
```

実際の タイムアウト判定は `c.rwc.SetReadDeadline()` でエラーがあれば エラーを返すっぽい  

### ミドルウェア を自作した時の使い方(全部引用)
参考: [HTTP Middleware の作り方と使い方](https://tutuz-tech.hatenablog.com/entry/2020/03/23/220326)  
原典: [Making and Using HTTP Middleware](https://www.alexedwards.net/blog/making-and-using-middleware)  

Go では、HTTP リクエストの制御の流れが以下のようになるように、 ServeMux とアプリケーションハンドラの間でミドルウェアを使用するのが一般的です。  
ServeMux => Middleware Handler => Application Handler  
Go でのミドルウェアの作成と使用は基本的にシンプルです。以下のようにしたいです。  
- http.Handler インターフェースを満たすようにミドルウェアを実装  
- ミドルウェアハンドラと通常のアプリケーションハンドラの両方を含むハンドラのチェーンを構築  

```go
func messageHandler(message string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(message))
	})
}
```
このハンドラでは、匿名関数の中にロジック (シンプルな w.Write という書き込みです) を配置し、メッセージ変数の上でクロージャを形成しています。  
そして、このクロージャを http.HandlerFunc アダプタを使用してハンドラに変換し、ハンドラを返します。  

同じアプローチを使用して、ハンドラのチェーンを作ることもできます。  
(上記のように) クロージャに文字列を渡すのではなく、 チェーンの中で次のハンドラを変数として渡し、その次のハンドラの ServeHTTP() メソッドを呼び出することで制御を次のハンドラに移すことができます。  
これにより、ミドルウェアを構築するための完全なパターンが得られます。  
```go
func exampleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Our middleware logic goes here...
		next.ServeHTTP(w, r)
	})
}
```
このミドルウェアの関数は func(http.Handler) http.Handler シグネチャを持っていることに気づくでしょう。  
これはハンドラをパラメータとして受け取り、ハンドラを返します。これは二つの理由で便利です。  
- これは http.Handler を返すので、ミドルウェアの関数を net/http パッケージで提供されている標準の ServeMux に直接登録することができます。  
- ミドルウェア関数を互いに入れ子にすることで、任意の長いハンドラチェーンを作ることができます。例えば、以下のようになります。  
```go
http.Handle("/", middlewareOne(middlewareTwo(finalHandler)))
```
ログメッセージを標準出力に書き込むだけのミドルウェアの例を見てみましょう。  
```go
func middlewareOne(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing middlewareOne")
		next.ServeHTTP(w, r)
		log.Println("Executing middlewareOne again")
	})
}

func middlewareTwo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing middlewareTwo")
		if r.URL.Path != "/" {
			return
		}
		next.ServeHTTP(w, r)
		log.Println("Executing middlewareTwo again")
	})
}

func final(w http.ResponseWriter, r *http.Request) {
	log.Println("Executing finalHandler")
	w.Write([]byte("OK"))
}

func main() {
	finalHandler := http.HandlerFunc(final)

	http.Handle("/", middlewareOne(middlewareTwo(finalHandler)))
	http.ListenAndServe(":3000", nil)
}
```
ミドルウェアハンドラから return を発行することで、チェーンを介して伝播する制御を任意の時点で停止することができます。  
