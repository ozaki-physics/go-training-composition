# http サーバの自動テストを勉強
テストの書き方は[こっち](./test_memo.md)を見る  

go のテストする仕様なのか 以下の条件を満たさないと エディタに怒られた  

ブラックボックス化してテストしたいなら ファイル名を `なんとか_test` にして  
パッケージ名を `テスト対象のパッケージ名_test` にする  
すると パッケージ名が異なる扱いになり import が必要になり public メソッドしか見れなくなる  

サーバを立てなくても HandlerFunc 型 つまり `ServeHTTP(w ResponseWriter, r *Request)` メソッドを持っていればテストできるのが便利そう

[httptest](https://pkg.go.dev/net/http/httptest) のサンプルを見ると  
`func(w http.ResponseWriter, r *http.Request){}` メソッドを作って  
`http.HandlerFunc()` で型変換して  
それぞれ  
- `httptest.NewServer()`  
- `httptest.NewUnstartedServer()`  
- `httptest.NewTLSServer()`  

に渡している  
実際 それぞれの中を見る  

`httptest.NewServer()`, `httptest.NewTLSServer()` は サーバを返す かつ 起動してくれる  
サーバを閉じるときは Close する必要がある  
中では `httptest.NewUnstartedServer()` を呼んでいる  

`httptest.NewUnstartedServer()` はサーバを返すけど 起動はしない
返された方で Start か StartTLS を実行する必要がある  
閉じるときは Close する  

立てたサーバに対して アクセスしている  
その立てたサーバに  
`r := httptest.NewRequest("GET", "http://example.com/foo", nil)`  
`w := httptest.NewRecorder()`  
を渡してるっぽい  

`httptest.NewRequest()` は  

>NewRequest は、テストのために http.Handler に渡すのに適した、新しい着信サーバー リクエストを返します。  


`httptest.NewRecorder()` は [type ResponseRecorder](https://pkg.go.dev/net/http/httptest#ResponseRecorder) を返す  

>ResponseRecorder は http.ResponseWriter の実装で、その変異を記録し、後でテストで検査できるようにするものです。  

これらを `テスト対象のパッケージ名_test` の方に書く  
僕は `httptest.NewRequest()`, `httptest.NewRecorder()` を使う方が好みかなぁ  
まだちゃんと `httptest.NewUnstartedServer()` が理解できてないのもある笑  

`w := httptest.NewRecorder()` で w から直接ステータスコードや Body を取り出すパターン と  
`resp := w.Result()` して resp から ステータスコードや Body を取り出すパターン は  
何が違うのだろう  
