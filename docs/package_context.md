# context パッケージ の勉強
## 参考書籍
『詳解Go言語Webアプリケーション開発』
ISBNコード: 978-4-86354-372-0
著者: 清水 陽一郎
初版: 2022/08/01
## context の役割
- キャンセル や デッドライン の伝播  
- リクエスト, トランザクションスコープ の メタデータを関数やゴルーチン間で伝播  

`xxx` 型の値 = `xxx` 型のインスタンス, オブジェクト と定義する  
たぶん Go が関数も値として扱えることに起因した表現と思われる  

[公式 context](https://pkg.go.dev/context)  
>Package context defines the Context type, which carries deadlines, cancellation signals, and other request-scoped values across API boundaries and between processes.  
>Package context は Context 型を定義し デッドライン, キャンセルシグナル, その他のリクエスト に対応した値を API 境界やプロセス間で伝達します。  

>Do not store Contexts inside a struct type;  
>instead, pass a Context explicitly to each function that needs it.  
>The Context should be the first parameter, typically named ctx:  
>その代わり Context を必要とする各関数には明示的に Context を渡してください  
>Context は最初のパラメータであるべきで 典型的な名前は ctx です。  

__HTTP サーバの開発で context は必須__  
なぜなら エンドポイント実装内部では クライアントからの通信切断やタイムアウトは `context.Context` 型の値 からしか検知できない  
キャンセルやタイムアウトしたのに 処理を続けて不整合になるのは問題だから  

`net/http` パッケージ の `*http.Request` 型の値 の `Context()` メソッドで `context.Context` 型の値を取得できる  
[func (*Request) Context](https://pkg.go.dev/net/http#Request.Context)  

各メソッドで `context.Context` を引数に渡すことで `context.WithValue()` から __メタデータをまとめて渡せる__  
トレースやメトリックを計測するツールなどで 分散処理するための トレース ID や リクエスト ID など伝播できる  

自作ハンドラでは context を使わなくても 呼び出し先のメソッドから要求される可能性がある(`database/sql` など)  
例えば トランザクションを開始するメソッド [`func (c *Conn) BeginTx(ctx context.Context, opts *TxOptions) (*Tx, error)`](https://pkg.go.dev/database/sql#Conn.BeginTx)  

ルートとなる context は [`func Background() Context`](https://pkg.go.dev/context#Background)  
>Background returns a non-nil, empty Context.  
>It is never canceled, has no values, and has no deadline.  
>It is typically used by the main function, initialization, and tests, and as the top-level Context for incoming requests.  
>Background は、Nilでない、空のContextを返す。  
>これは決してキャンセルされず、何の値も持たず、デッドラインもない。  
>これは通常、メイン関数、初期化、テスト、およびリクエストを受信するためのトップレベルの Context として使用されます。  

任意のタイミングでキャンセルや完了通知をするためのメソッド  
[`func WithCancel(parent Context) (ctx Context, cancel CancelFunc)`](https://pkg.go.dev/context#WithCancel)  

時間制限を設定する方法は2種  
- 指定した時刻を経過したらキャンセルする [`func WithDeadline(parent Context, d time.Time) (Context, CancelFunc)`](https://pkg.go.dev/context?utm_source=gopls#WithDeadline)  
- 指定した時間を経過したらキャンセルする [`func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)`](https://pkg.go.dev/context#WithTimeout)  

重い処理の前にキャンセル済みか把握するには `context.Context.Err()`  
```go
type Context interface {
  // 中略
	// If Done is not yet closed, Err returns nil.
	// Doneがまだ閉じられていない場合、Errはnilを返す。
	// If Done is closed, Err returns a non-nil error explaining why:
	// Doneが閉じられている場合、Errはその理由を説明する非Nilのエラーを返す。
	// Canceled if the context was canceled or DeadlineExceeded if the context's deadline passed.
	// コンテキストがキャンセルされた場合、Canceled または コンテキストのデッドラインが過ぎた場合は、DeadlineExceeded。
	// After Err returns a non-nil error, successive calls to Err return the same error.
	// Errがnilでないエラーを返した後、連続してErrを呼び出すと同じエラーが返される。
	Err() error
}
```
キャンセルされるまで待つには `context.Context.Done()`  
```go
type Context interface {
  // 中略
	// Done returns a channel that's closed when work done on behalf of this context should be canceled.
	// Done は、このコンテキストのために行われた作業がキャンセルされるべきときに閉じられるチャネルを返します。
  // Done may return nil if this context can never be canceled.
  // Done は、このコンテキストが決してキャンセルされない場合は nil を返すことがあります。
  // Successive calls to Done return the same value.
  // Done を連続して呼び出すと、同じ値が返される。
	// The close of the Done channel may happen asynchronously, after the cancel function returns.
	// Done チャンネルのクローズは、cancel 関数が戻った後に非同期的に行われることがある。
	//
	// WithCancel arranges for Done to be closed when cancel is called;
	// WithDeadline arranges for Done to be closed when the deadline
	// expires; WithTimeout arranges for Done to be closed when the timeout
	// elapses.
	//
	// Done is provided for use in select statements:
	//
	//  // Stream generates values with DoSomething and sends them to out
	//  // until DoSomething returns an error or ctx.Done is closed.
	//  func Stream(ctx context.Context, out chan<- Value) error {
	//  	for {
	//  		v, err := DoSomething(ctx)
	//  		if err != nil {
	//  			return err
	//  		}
	//  		select {
	//  		case <-ctx.Done():
	//  			return ctx.Err()
	//  		case out <- v:
	//  		}
	//  	}
	//  }
	//
	// See https://blog.golang.org/pipelines for more examples of how to use
	// a Done channel for cancellation.
	Done() <-chan struct{}
}
```

context に データを保持させるには [`func WithValue(parent Context, key, val any) Context`](https://pkg.go.dev/context#WithValue) を使う  
context から データの取得は `Context.Value()`  
```go
type Context interface {
  // 中略
	// Value returns the value associated with this context for key, or nil if no value is associated with key.
	// Value は、このコンテキストの key に関連付けられた値を返すが、 key に値が関連付けられていない場合は nil を返す。
  // Successive calls to Value with the same key returns the same result.
  // 同じキーで連続して Value を呼び出すと、同じ結果が返される。
	//
	// Use context values only for request-scoped data that transits processes and API boundaries, not for passing optional parameters to functions.
	// コンテキスト値は、プロセスやAPI境界を通過するリクエストスコープのデータにのみ使用し、関数にオプションのパラメータを渡す際には使用しない。
	//
	// A key identifies a specific value in a Context.
	// キーとは、Context 内の特定の値を識別するものである。
  // Functions that wish to store values in Context typically allocate a key in a global variable then use that key as the argument to context.WithValue and Context.Value.
  // Contextに値を格納したい関数は、通常グローバル変数にキーを割り当て、そのキーをcontext.WithValueおよびContext.Valueの引数として使用する。
  // A key can be any type that supports equality;
  // キーは、等号をサポートする任意の型にすることができる。
	// packages should define keys as an unexported type to avoid collisions.
	// パッケージは、衝突を避けるためにキーを unexported 型として定義する必要があります。
  // 中略
  Value(key any) any
}
```
key には __空構造体 を使って 独自型__ にするのが一般的  
なぜなら プリミティブな値だと 他のパッケージと key が衝突する可能性があるから  

`context.Context` 型の値 を 構造体のフィールドに格納するのは アンチパターン  
`context.Context` 型の値 に含めてよい情報は API リクエストスコープ(トランザクションスコープ)  
- リクエスト元 IP アドレス  
- ユーザーエージェント  
- リファラ  
- ロードバランサ や SaaS によって割り振られた ID  
- リクエスト受信時刻  

Logic に関わる値を含めてはダメ  
ただ 認証, 認可 に関する情報は アーキテクチャ全体をシンプルにするため妥協してよいこともある  

マイクロサービスアーキテクチャ的に context をサーバ間で共有したいときは context のデータを HTTP ヘッダに詰め替えるなどの処理が必要  

既存のコードが `context.Context` 型の値 を受け取ってないが 今後導入していくときは 空な context である `context.TODO()` メソッドが便利  
公式がそんなことまでフォローしてるの優しすぎで良い言語だな笑  
