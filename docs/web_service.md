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
