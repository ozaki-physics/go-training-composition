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
