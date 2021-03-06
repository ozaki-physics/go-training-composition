# go で server を立てる
server を立てるときに使うパッケージ `"net/http"`  
そのとき よく使う2個のメソッドから調べていく  
- `http.ListenAndServe()`  
- `http.HandleFunc()`  

引用は 一部抜き出して翻訳している  

## `http.ListenAndServe()` を調べる
server を立てるメソッド  
使うときは `http.ListenAndServe(":8080", nil)` とすることが多い  

[func ListenAndServe](https://pkg.go.dev/net/http#ListenAndServe) より  
```go
func ListenAndServe(addr string, handler Handler) error {
	server := &Server{Addr: addr, Handler: handler}
	return server.ListenAndServe()
}
```
>ListenAndServe は、TCPネットワークアドレス addr でリッスンし、受信接続の要求を処理するハンドラで Serve を呼び出します。  
>ハンドラは通常 nil であり、その場合 DefaultServeMux が使用される。  
>ListenAndServe は常に nil でないエラーを返します。  

[type Server](https://pkg.go.dev/net/http#Server) は 内部では Addr, Handler の field しか設定してないけど タイムアウトとかもっと設定することもできる  

使うときに 第2引数に nil を渡すことが多い つまり DefaultServeMux が使われていると思われる  

[type Handler](https://pkg.go.dev/net/http#Handler) より  
```go
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
```
>ハンドラは、 HTTP リクエストに応答します。  
>ServeHTTP は、応答ヘッダとデータを ResponseWriter に書き込んでから、リターンする必要があります。  
>HTTP クライアントソフトウェアや HTTP プロトコルのバージョン、クライアントとGoサーバーの間に存在する仲介エンティティによっては、 ResponseWriter に書き込んだ後に Request.Body から読み取ることができない場合があります。  

### `server.ListenAndServe()` をもう少し深く調べていく
[func (*Server) ListenAndServe](https://pkg.go.dev/net/http#Server.ListenAndServe) より  
```go
func (srv *Server) ListenAndServe() error
```
>ListenAndServe は、 TCP ネットワークアドレス srv.Addr でリッスンし、Serveを呼び出して着信接続の要求を処理します。  
>受け入れられた接続は、 TCP キープアライブを有効にするように構成されます。  

```go
func (srv *Server) ListenAndServe() error {
	if srv.shuttingDown() {
		return ErrServerClosed
	}
	addr := srv.Addr
	if addr == "" {
		addr = ":http"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return srv.Serve(ln)
}
```
大事そうなのは `ln, err := net.Listen("tcp", addr)` して `return srv.Serve(ln)` だから `srv.Serve()` を調べる  
[func (*Server) Serve](https://pkg.go.dev/net/http#Server.Serve) より  
```go
// 一部抜粋
func (srv *Server) Serve(l net.Listener) error {
	ctx := context.WithValue(baseCtx, ServerContextKey, srv)
	for {
		rw, err := l.Accept()

		connCtx := ctx

		c := srv.newConn(rw)
		c.setState(c.rwc, StateNew, runHooks) // before Serve can return
		go c.serve(connCtx)
	}
}
```
>ServeはListener の l で着信接続を受け付け、それぞれに新しいサービスゴルーチンを生成します。  
>サービスゴルーチンはリクエストを読み、 srv.Handler を呼び出して返信します。  

大事そうなのは `go c.serve(connCtx)` だから 調べる  

```go
// 一部抜粋
// Serve a new connection.
func (c *conn) serve(ctx context.Context) {
	ctx = context.WithValue(ctx, LocalAddrContextKey, c.rwc.LocalAddr())
	if tlsConn, ok := c.rwc.(*tls.Conn); ok {
	}

	// HTTP/1.x from here on.
	ctx, cancelCtx := context.WithCancel(ctx)
	for {
		w, err := c.readRequest(ctx)
		if err != nil {
			switch {
			case err == errTooLarge:
			case isUnsupportedTEError(err):
			case isCommonNetReadError(err):
			default:
				if v, ok := err.(statusError); ok {
					fmt.Fprintf(c.rwc, "HTTP/1.1 %d %s: %s%s%d %s: %s", v.code, StatusText(v.code), v.text, errorHeaders, v.code, StatusText(v.code), v.text)
					return
				}
				publicErr := "400 Bad Request"
				fmt.Fprintf(c.rwc, "HTTP/1.1 "+publicErr+errorHeaders+publicErr)
				return
			}
		}

		serverHandler{c.server}.ServeHTTP(w, w.req)
  }
}
```
大事そうなのは `serverHandler{c.server}.ServeHTTP(w, w.req)` だから調べる  
`c.rwc.(*tls.Conn)` で https の処理をしていそう  
ちなみに `c.readRequest(ctx)` の中で タイムアウトの処理をしている  
```go
// serverHandler delegates to either the server's Handler or
// DefaultServeMux and also handles "OPTIONS *" requests.
// serverHandlerはサーバーのHandlerまたはDefaultServeMuxに委ね、"OPTIONS *" リクエストも扱います。
type serverHandler struct {
	srv *Server
}

func (sh serverHandler) ServeHTTP(rw ResponseWriter, req *Request) {
	handler := sh.srv.Handler
	if handler == nil {
		handler = DefaultServeMux
	}
	if req.RequestURI == "*" && req.Method == "OPTIONS" {
		handler = globalOptionsHandler{}
	}

	if req.URL != nil && strings.Contains(req.URL.RawQuery, ";") {
		var allowQuerySemicolonsInUse int32
		req = req.WithContext(context.WithValue(req.Context(), silenceSemWarnContextKey, func() {
			atomic.StoreInt32(&allowQuerySemicolonsInUse, 1)
		}))
		defer func() {
			if atomic.LoadInt32(&allowQuerySemicolonsInUse) == 0 {
				sh.srv.logf("http: URL query contains semicolon, which is no longer a supported separator; parts of the query may be stripped when parsed; see golang.org/issue/25192")
			}
		}()
	}

	handler.ServeHTTP(rw, req)
}
```

[type Server](https://pkg.go.dev/net/http?utm_source=gopls#Server) より  
>Server は、 HTTP サーバーを動作させるためのパラメータを定義します。  

ここでやっと  
>srv.Handler を呼び出して返信します  

の `srv.Handler` が出てきたし ServeHTTP メソッド つまり Handler interface の実装が使われる部分が出てきた  
つまり 自作した ServeHTTP(ResponseWriter, *Request) は ここで呼ばれて動作すると思われる  

## `http.HandleFunc()` を調べる
[func HandleFunc](https://pkg.go.dev/net/http?utm_source=gopls#HandleFunc) より  
```go
func HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	DefaultServeMux.HandleFunc(pattern, handler)
}
```
>HandleFunc は、与えられたパターンに対応するハンドラ関数を DefaultServeMux に登録する。  
>ServeMux のドキュメントでは、パターンがどのようにマッチングされるかを説明しています。  

### `DefaultServeMux` を もう少し深く調べていく  
[DefaultServeMux](https://pkg.go.dev/net/http#DefaultServeMux)  
```go
// DefaultServeMux is the default ServeMux used by Serve.
var DefaultServeMux = &defaultServeMux

var defaultServeMux ServeMux

type ServeMux struct {
	mu    sync.RWMutex
	m     map[string]muxEntry
	es    []muxEntry // slice of entries sorted from longest to shortest.
	hosts bool       // whether any patterns contain hostnames
}

type muxEntry struct {
	h       Handler
	pattern string
}
```
[type ServeMux](https://pkg.go.dev/net/http#ServeMux)  
ServeMux の field は private なことに注意  
>ServeMux は、 HTTP リクエストのマルチプレクサーです。
>受信した各リクエストの URL を登録されたパターンのリストと照合し、その URL に最も近いパターンのハンドラを呼び出します。  
>パターンは、"/favicon.ico"のような固定ルートパス、または"/images/"(末尾のスラッシュに注意)のようなルートサブツリーの名前を指定します。  
>長いパターンは短いパターンよりも優先されます。  
>たとえば、"/images/" と "/images/thumbnails/" の両方についてハンドラが登録されている場合、"/images/thumbnails/" で始まるパスについては後者のハンドラが呼ばれ、 "/images/" サブツリー内の他のパスについては前者がリクエストを受け取ります。  
>スラッシュで終わるパターンはルート化されたサブツリーを指定するため、パターン"/"は Path == "/" のURLだけでなく、他の登録パターンにマッチしないすべてのパスにマッチすることに注意してください。  
>サブツリーが登録されている場合に、末尾のスラッシュなしでサブツリー ルートを指定するリクエストを受信すると、 ServeMux はそのリクエストをサブツリー ルートにリダイレクトします(末尾のスラッシュを追加します)。  
>この動作は、末尾のスラッシュを除いたパスに対する別の登録で上書きすることができます。たとえば、"/images/"を登録すると、"/images"が別途登録されていない限り、 ServeMux は"/images"への要求を"/images/"にリダイレクトします。  
>パターンはオプションでホスト名で始めることができ、そのホスト上の URL のみにマッチを制限することができます。  
>ホスト固有のパターンは一般的なパターンよりも優先されます。  
>つまり、ハンドラは "/codesearch" と "codesearch.google.com/" という二つのパターンを登録することで、 "http://www.google.com/" に対するリクエストも引き受けることができるようになります。  
>ServeMux は、 URL リクエストパスと Host ヘッダのサニタイズも行い、ポート番号を削除し、"."や"..要素"、繰り返されるスラッシュを含むリクエストは、同等のクリーンなURLにリダイレクトされます。  

要約すると  
pattern(url のこと) には rooted paths(ex. /favicon.ico) と rooted subtrees(ex. /images/) が登録できる  
pattern は長いが優先される  
"/images/" と "/images/thumbnails/" が登録されていて リクエストが "/images/thumbnails/" で始まるなら後者, "/images/なんとか" は前者が呼ばれる  
pattern に "/" だけだと リクエストが "/" だけのものではなく 登録のないリクエストすべて が "/" を呼び出す  
つまり "/images/" の登録がないのに リクエストが "/images/" なら "/" にリダイレクトされる感じ  
また "/images/" の登録があるのに リクエストが "/images" なら "/images/" にリダイレクトする  

ServeMux struct は4個のメソッドを持っている  
- `func (mux *ServeMux) Handle(pattern string, handler Handler)`  
  `DefaultServeMux.HandleFunc()` と一緒に調べる  
- `func (mux *ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, *Request))`  
  `DefaultServeMux.HandleFunc()` と一緒に調べる  
- `func (mux *ServeMux) Handler(r *Request) (h Handler, pattern string)`  
  中で ServeMux の private メソッド `(mux *ServeMux) handler(host, path string) (h Handler, pattern string)` が呼ばれている  
  `(mux *ServeMux) handler` で ServeMux の中に 存在する pattern を探して Handler を返す  
- `func (mux *ServeMux) ServeHTTP(w ResponseWriter, r *Request)`  
  中で `(mux *ServeMux) Handler()` を呼んでいる  


ServeMux を自分で生成するメソッド [func NewServeMux() *ServeMux](https://pkg.go.dev/net/http#NewServeMux) もある  

### `DefaultServeMux.HandleFunc()` を もう少し深く調べていく  
[func (*ServeMux) HandleFunc](https://pkg.go.dev/net/http#ServeMux.HandleFunc) より  
>HandleFunc は、与えられたパターンに対するハンドラ関数を登録する。  

```go
func (mux *ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	if handler == nil {
		panic("http: nil handler")
	}
	mux.Handle(pattern, HandlerFunc(handler))
}
```
HandlerFunc(handler) は ただ 型変換してる感じ  

[func (*ServeMux) Handle](https://pkg.go.dev/net/http#ServeMux.Handle)
>ハンドルは与えられたパターンに対応するハンドラを登録します。もし、pattern に対応するハンドラが既に存在する場合、Handle はパニックに陥ります。  

```go
func (mux *ServeMux) Handle(pattern string, handler Handler)
```
ServeMux の field m (型 map[string]muxEntry) に pattern と 呼び出しを登録するみたい  

公式ドキュメントのサンプルには 以下が書かれていた  
```go
// apiHandler は ServeHTTP() を持つため Handler interface を実装したことになる
type apiHandler struct{}
func (apiHandler) ServeHTTP(http.ResponseWriter, *http.Request) {}

func main() {
	mux := http.NewServeMux()
  // 以下2個の書き方は実質同じで struct を用意するか否かの違いだと思う
	mux.Handle("/api/", apiHandler{})
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		// The "/" pattern matches everything, so we need to check
		// that we're at the root here.
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
		fmt.Fprintf(w, "Welcome to the home page!")
	})
}
```

## 軽く認識をまとめる
go の server は ServeMux struct に登録された pattern と Handler interface によって リクエストが処理されている  
そして 毎回 ServeMux を用意しなくていいように DefaultServeMux が存在する  
`http.HandleFunc()` は DefaultServeMux に ハンドラを追加するメソッドであり  
`http.ListenAndServe()` の第2引数を nil にすると DefaultServeMux が使われる  
これで ルーティング の登録についてはなんとなく分かった  
ただ 実際に 呼び出される処理はどうなっているんだろう  

## `http.FileServer()`を調べる(昔 調べたメモに書いてあったから統合する)
[func FileServer](https://pkg.go.dev/net/http#FileServer) より  
net/http/fs.go に書いてある  
```go
func FileServer(root FileSystem) Handler {
	return &fileHandler{root}
}
```
`Handler 型`を返す関数 つまり`ServeHTTP()`関数だけを持ち HTTP リクエストを受けてレスポンスを返すことができる  
静的な web server として動作しそう  

>FileServer は、root にルートされたファイルシステムの内容を HTTP リクエストに提供するハンドラを返す。  
>特殊なケースとして、返されたファイルサーバーは "/index.html" で終わるいかなるリクエストも、最後の "index.html" を除いた同じパスにリダイレクトする。  
>オペレーティングシステムのファイルシステム実装を使用するには、 http.Dir を使用します。  

公式ドキュメントにあったサンプル
```go
func main() {
	// Simple static webserver:
	log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir("/usr/share/doc"))))
}
```

### `fileHandler struct` をもう少し深く調べていく
```go
type fileHandler struct {
	root FileSystem
}
type FileSystem interface {
	Open(name string) (File, error)
}
type File interface {
	io.Closer
	io.Reader
	io.Seeker
	Readdir(count int) ([]fs.FileInfo, error)
	Stat() (fs.FileInfo, error)
}
type Dir string
```

[type FileSystem](https://pkg.go.dev/net/http#FileSystem) より  
>FileSystemは、名前付きファイルのコレクションへのアクセスを実装しています。ファイルパスの要素は、ホストのオペレーティングシステムの慣習に関係なく、スラッシュ ('/', U+002F) 文字で区切られます。  
>FileSystem を Handler に変換するには、 FileServer 関数を参照してください。  
>このインターフェースは fs.FS インターフェースよりも古いものであり、代わりに FS アダプタ関数が fs.FS を FileSystem に変換することができます。  

[type File](https://pkg.go.dev/net/http#File) より  
>FileSystem の Open メソッドから返され、 FileServer の実装によって提供される。  
>メソッドは、 *os.File のメソッドと同じように動作する必要があります。  

[type Dir](https://pkg.go.dev/net/http#Dir) より  
>Dir は、特定のディレクトリツリーに限定されたネイティブファイルシステムを使用して FileSystem を実装しています。  
> FileSystem.Open メソッドは '/' で区切られたパスを受け取りますが、Dir の文字列値は URL ではなくネイティブファイルシステム上のファイル名なので filepath.Separator で区切られており、必ずしも '/' とは限りません。  
>Dir は機密性の高いファイルやディレクトリを公開する可能性があることに注意してください。  
>Dir はディレクトリツリーの外を指すシンボリックリンクをたどります。  
>これは、ユーザが任意のシンボリックリンクを作成できるようなディレクトリから 提供される場合、特に危険です。  
>Dir はまた、ピリオドで始まるファイルやディレクトリへのアクセスを許可します。  
>これは、.git のような機密ディレクトリや、.htpasswd のような機密ファイルを公開する可能性があります。  
>ピリオドで始まるファイルを除外するには、そのファイル/ディレクトリをサーバーから削除するか、カスタム FileSystem の実装を作成してください。  
>空の Dir は"."として扱われます。  

つまり `FileServer()` のための基盤って感じがする  
type Dir の適切な使い方は 今度勉強しよう  

## http メソッドによる 処理の振り分け
[http の Constants](https://pkg.go.dev/net/http#pkg-constants) より  
http.MethodGet などが存在する(net/http/method.go)  
`http.Request.Method` を取り出して switch などで 定数(http.MethodGet など) 振り分ける  

## ミドルウェア について
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

## 静的ファイルの配信(テンプレートなし)
```go
	makeHandle := http.FileServer(http.Dir("web"))
	http.Handle("/", makeHandle)
```
`http.Handle("/", http.FileServer(http.Dir("web")))` で web ディレクトリ内の静的ファイルを配信できる  
URL は web ディレクトリが ルート URL になる  
http://localhost:8088/hello.html アクセスできる  
http://localhost:8088/sub/hello.html アクセスできない  
http://localhost:8088/sub/sub01.html アクセスできる  
http://localhost:8088/sub/sub/sub01.html アクセスできない  

プレフィックスをつけてみる  
```go
	makeHandle01 := http.StripPrefix("/sub/", makeHandle)
	http.Handle("/", makeHandle01)
```
http://localhost:8088/hello.html アクセスできない  
http://localhost:8088/sub/hello.html アクセスできる  
http://localhost:8088/sub/sub01.html アクセスできない  
http://localhost:8088/sub/sub/sub01.html アクセスできる  
`http.StripPrefix("/sub/", makeHandle)` をすると ルート URL に 引数01番目 を追加した URL が web ディレクトリと つながるっぽい  
http.StripPrefix で プレフィックスって言ってるぐらいだから URL に 接頭辞が付く  

受け付けるパスを変えてみる  
```go
	makeHandle01 := http.StripPrefix("/sub/", makeHandle)
	http.Handle("/aaa/", makeHandle01)
```
http://localhost:8088/hello.html アクセスできない  
http://localhost:8088/sub/hello.html アクセスできない  
http://localhost:8088/sub/sub01.html アクセスできない  
http://localhost:8088/sub/sub/sub01.html アクセスできない  

http://localhost:8088/aaa/hello.html アクセスできない  
http://localhost:8088/aaa/sub/hello.html アクセスできない  
http://localhost:8088/sub/aaa/hello.html アクセスできない  

http://localhost:8088/aaa/sub01.html アクセスできない  
http://localhost:8088/aaa/sub/sub01.html アクセスできない  
http://localhost:8088/aaa/sub/sub/sub01.html アクセスできない  
まず サーバが /aaa をついてないとリッスンしてない  
また /aaa をつけても /aaa ディレクトリなんて存在しないから すべてアクセスできないと思われる  

プレフィックスを変えてみる  
```go
	makeHandle01 := http.StripPrefix("/aaa/", makeHandle)
	http.Handle("/aaa/", makeHandle01)
```
http://localhost:8088/hello.html アクセスできない  
http://localhost:8088/sub/hello.html アクセスできない  
http://localhost:8088/sub/sub01.html アクセスできない  
http://localhost:8088/sub/sub/sub01.html アクセスできない  

http://localhost:8088/aaa/hello.html アクセスできる  
http://localhost:8088/aaa/sub/hello.html アクセスできない  
http://localhost:8088/sub/aaa/hello.html アクセスできない  

http://localhost:8088/aaa/sub01.html アクセスできない  
http://localhost:8088/aaa/sub/sub01.html アクセスできる  
http://localhost:8088/aaa/sub/sub/sub01.html アクセスできない  
プレフィックスに /aaa をつけたので URL に /aaa があったら web ディレクトリにつながるようになったと思われる  

## go のサーバを https にする方法(TLS)
### https に関する基礎知識
https は HTTP over SSL/TLS のこと
SSL(Secure Sockets Layer) とは インターネット上で データを暗号化して送受信するプロトコル  
TLS(Transport Layer Security) とは SSL の後継規格  
今は 全部 TLS になっており SSL は使ってはならないことになっている  
2021/03/xx には TLS 1.1 も非推奨になっている  
2015年あたりまでは SSL だったため 俗に SSL/TLS という両方書きが普及しているらしい  
つまり https にすると インターネット上の通信が暗号化される  

TLS のシーケンス  
1. サーバーは公開鍵の正当性を証明するために、認証局に登録申請を行う  
2. 認証局がサーバーに対して、デジタル証明書を発行する  
3. クライアントがサーバーに接続すると、サーバは認証局の署名付き証明書を返信する  
4. クライアントは証明書が本物であることを確認する  
5. クライアントは乱数を発生させ、サーバの公開鍵で暗号化してサーバーに送る  
  この乱数は実際にHTTP通信を行う際の鍵を作ることに利用される。  
6. サーバーは秘密鍵で、暗号化された乱数を復号する。  
7. サーバとクライアントで共通鍵を使って暗号化通信を行う。  
  サーバとクライアントの双方で、マスターシークレット(MS)と呼ばれる鍵を生成し、この MS から生成した共通鍵を使って暗号化通信を行う。  
  サーバーとクライアントが実際のデータの送受信を行う際には、共通鍵暗号を利用する。  
  このとき、この乱数から2種類の鍵を生成し、"クライアント -> サーバ"と"サーバ -> クライアント"で異なる鍵を用いる。  

https 通信であっても盗み見される情報  
- CONNECT メソッドのホスト情報  
  https がプロキシ経由で通信する際、例えば https://www.yahoo.co.jp と通信するクライアントは、プロキシに対して "CONNECT www.yahoo.co.jp" というメソッドを平文で発行しますので、これを盗み見すれば宛先情報が分かります。  
- デジタル証明書のサブジェクト代替名/コモンネーム  
  https のセッション開始時 (暗号化通信が始まる前) の "Certificate" メッセージではサーバがデジタル証明書を提示しますが、そのデジタル証明書には サブジェクト代替名 (SANs) という URL のホスト情報が含まれています。  
  これを盗み見することでやはり宛先情報が分かります。  
- Client Hello の ServerName Extension  
  Extension には通信先サーバの URL 情報が載せられています。  
  元々は https を使うバーチャルホストへの対応のため (HTTP GET の前に URL が分かっていないと、対応した証明書を提示できず TLS コネクションが開始できないため) に考えられましたが、中間の NW 機器 (主に UTM) 等が宛先を認識して制御する目的でもよく使われます。  

参考:  
[SSL/TLSについて調べたことまとめ](https://qiita.com/Takatoshi_Hiki/items/d98a2d7f52708eac5324)  
[【図解】https(SSL/TLS)の仕組みとシーケンス,パケット構造 〜暗号化の範囲, Encrypted Alert, ヘッダやレイヤについて～](https://milestone-of-se.nesuke.com/nw-basic/tls/https-structure/)  

### 実際に go のコードを見る
[func ListenAndServeTLS](https://pkg.go.dev/net/http#ListenAndServeTLS) より  
```go
func ListenAndServeTLS(addr, certFile, keyFile string, handler Handler) error {
	server := &Server{Addr: addr, Handler: handler}
	return server.ListenAndServeTLS(certFile, keyFile)
}
```
>ListenAndServeTLS は、 HTTPS 接続を必要とする以外は、 ListenAndServe と同じように動作します。  
>さらに、サーバーの証明書と一致する秘密鍵を含むファイルを提供する必要があります。  
>証明書が認証局によって署名されている場合、 certFile はサーバーの証明書、任意の中間体、およびCAの証明書を連結したものである必要があります。  

サンプル  
```go
func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hello, TLS!\n")
	})

	// One can use generate_cert.go in crypto/tls to generate cert.pem and key.pem.
  // crypto/tls の generate_cert.go を使って cert.pem と key.pem を生成することができます。
	log.Printf("About to listen on 8443. Go to https://127.0.0.1:8443/")
	err := http.ListenAndServeTLS(":8443", "cert.pem", "key.pem", nil)
	log.Fatal(err)
}
```
サンプル見ても分かるように公開鍵と秘密鍵をサーバ(Mux)に設定する必要があるっぽい  
中身としては DefaultServeMux の代わりとなる Mux を作って `(srv *Server) ListenAndServeTLS()` を呼んでいる  

コメントに書いてある [crypto/tls の generate_cert.go](https://github.com/golang/go/blob/master/src/crypto/tls/generate_cert.go) は main パッケージ かつ build ignore している  
つまり `go run $(go env GOROOT)/src/crypto/tls/generate_cert.go --host=localhost` と使う  
`curl --cacert webServer/sample_cert.pem https://localhost:8088` でアクセスできた  
http では `Client sent an HTTP request to an HTTPS server` と言われ ちゃんと https サーバだからダメだよって言われた

[crypto/tls](https://pkg.go.dev/crypto/tls) に [Example (HttpServer)](https://pkg.go.dev/crypto/tls#example-X509KeyPair-HttpServer) があった  
```go
func main() {
	certPem := []byte(`-----BEGIN CERTIFICATE-----
MIIBhTCCASugAwIBAgIQIRi6zePL6mKjOipn+dNuaTAKBggqhkjOPQQDAjASMRAw
DgYDVQQKEwdBY21lIENvMB4XDTE3MTAyMDE5NDMwNloXDTE4MTAyMDE5NDMwNlow
EjEQMA4GA1UEChMHQWNtZSBDbzBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABD0d
7VNhbWvZLWPuj/RtHFjvtJBEwOkhbN/BnnE8rnZR8+sbwnc/KhCk3FhnpHZnQz7B
5aETbbIgmuvewdjvSBSjYzBhMA4GA1UdDwEB/wQEAwICpDATBgNVHSUEDDAKBggr
BgEFBQcDATAPBgNVHRMBAf8EBTADAQH/MCkGA1UdEQQiMCCCDmxvY2FsaG9zdDo1
NDUzgg4xMjcuMC4wLjE6NTQ1MzAKBggqhkjOPQQDAgNIADBFAiEA2zpJEPQyz6/l
Wf86aX6PepsntZv2GYlA5UpabfT2EZICICpJ5h/iI+i341gBmLiAFQOyTDT+/wQc
6MF9+Yw1Yy0t
-----END CERTIFICATE-----`)
	keyPem := []byte(`-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIIrYSSNQFaA2Hwf1duRSxKtLYX5CB04fSeQ6tF1aY/PuoAoGCCqGSM49
AwEHoUQDQgAEPR3tU2Fta9ktY+6P9G0cWO+0kETA6SFs38GecTyudlHz6xvCdz8q
EKTcWGekdmdDPsHloRNtsiCa697B2O9IFA==
-----END EC PRIVATE KEY-----`)
	cert, err := tls.X509KeyPair(certPem, keyPem)
	if err != nil {
		log.Fatal(err)
	}
	cfg := &tls.Config{Certificates: []tls.Certificate{cert}}
	srv := &http.Server{
		TLSConfig:    cfg,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}
	log.Fatal(srv.ListenAndServeTLS("", ""))
}
```

`tls.X509KeyPair()` の戻り値が tls.Certificate 型 になっている  
`http.Server.TLSConfig` に tls.Config 型 で 渡す必要があるから  
tls.Config の中に 戻り値(tls.Certificate 型)を入れ込む  
`var b []int = []int{1,2,3}` と定義するのと同様に 要素数1で tls.Config.Certificates に入れる  

### `(srv *Server) ListenAndServeTLS()` をもう少し深く調べていく
[func (*Server) ListenAndServeTLS](https://pkg.go.dev/net/http#Server.ListenAndServeTLS) より  
```go
func (srv *Server) ListenAndServeTLS(certFile, keyFile string) error {
	if srv.shuttingDown() {
		return ErrServerClosed
	}
	addr := srv.Addr
	if addr == "" {
		addr = ":https"
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	defer ln.Close()

	return srv.ServeTLS(ln, certFile, keyFile)
}
```

>ListenAndServeTLS は、 TCP ネットワークアドレス srv.Addr でリッスンし、 ServeTLS を呼び出して受信する TLS 接続のリクエストを処理します。  
>受け入れられた接続は、 TCP キープアライブを有効にするように構成されます。  
>サーバーの  TLSConfig.Certificates または TLSConfig.GetCertificate のいずれも入力されていない場合、サーバーの証明書と一致する秘密鍵を含むファイル名を提供する必要があります。  
>証明書が認証局によって署名されている場合、 certFile は、サーバーの証明書、任意の中間体、および CA の証明書を連結したものでなければなりません。  

このメソッドは `http.ListenAndServeTLS()` の中から呼ばれてるし crypto/tls の Example (HttpServer) でも呼ばれてる  
ただ pem キーの渡すタイミングが違う  
crypto/tls の Example (HttpServer) では Mux に入れてる  
http.ListenAndServeTLS() では 引数で渡している  
結局 どこに pem キーが存在すればいいの?  

### `(srv *Server) ServeTLS()` をもう少し深く調べていく
[func (*Server) ServeTLS](https://pkg.go.dev/net/http#Server.ServeTLS) より  
```go
func (srv *Server) ServeTLS(l net.Listener, certFile, keyFile string) error {
	// 一部抜粋
	config := cloneTLSConfig(srv.TLSConfig)

	if !configHasCert || certFile != "" || keyFile != "" {
		var err error
		config.Certificates = make([]tls.Certificate, 1)
		config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	}

	tlsListener := tls.NewListener(l, config)
	return srv.Serve(tlsListener)
}
```

>ServeTLS は、 Listener lで着信接続を受け付け、それぞれに新しいサービスゴルーチンを生成します。  
>サービスゴルーチンは TLS のセットアップを行い、リクエストを読み、 srv.Handler を呼び出して返信します。  
>サーバーの TLSConfig.Certificates と TLSConfig.GetCertificate のいずれにも値が入力されていない場合、サーバーの証明書と一致する秘密鍵の入ったファイルを提供しなければならない。  
>証明書が認証局によって署名されている場合、 certFile は、サーバーの証明書、任意の中間体、およびCAの証明書を連結したものである必要があります。  

結局 keyFile は `tls.LoadX509KeyPair()` に渡されている  
そして `tls.LoadX509KeyPair()` から `tls.X509KeyPair()` に渡されているし `cloneTLSConfig()` は tls.Config 型を返すので  
結局 tls.Config の中に tls.Certificate 型 を入れ込む ことになり crypto/tls の Example (HttpServer) と同じ流れになった  
よって pem キーは Mux の中に存在しないとダメと思われる  

また `srv.Serve(tlsListener)` と http と同じ処理(`(srv *Server) Serve(l net.Listener)` つまり `server.ListenAndServe()`)にもなっている  

参考:  
[Go言語と暗号技術（AESからTLS）](https://deeeet.com/writing/2015/11/10/go-crypto/)  
- [TLSを使ったサーバとクライアントの実装例](https://github.com/tcnksm/go-crypto/tree/master/tls)  
- [httpsのサーバ実装例](https://github.com/tcnksm/go-crypto/blob/master/https/server.go)  

また クライアントから https 通信するときも何かしら設定がいる?  
- [httpsのクライアント実装例](https://github.com/tcnksm/go-crypto/blob/master/https/client.go)  

## 軽く気になったこと
いつドメインの設定をするんだろう  
-> DNS 解決だからサーバには無理  
-> サーバはただのIPアドレスしか無く IP アドレスは好きに設定できるわけでないから  
-> ただ go で DNS サーバも作ることはできるっぽい  

フツーに使ってたら 毎回 localhost:8080 になる  
-> Server struct の Addr フィールド に 8080 を渡しているから  
