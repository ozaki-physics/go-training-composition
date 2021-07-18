err 処理に色々な書き方があるっぽいけど どんな書き方がいいのだろう

```go
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return
	}

	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		panic(err)
	}
```

`panic()` はめったに使わなさそう
エラーハンドリングもできないような どうしようもないエラーで使うっぽい

log パッケージ と fmt パッケージ は何が違うのか
→ 使用する目的が違う
ただ コンソールに出力したいなら fmt でもいい
しかし ログとして使いたいなら log が良い
デバッグで使いたい程度なら どっちでもいい

fmt より log の方がコンソールに多くの情報を出力してくれる
感覚的には fmt + os = log な感じ
`log.Fatal()`, `log.Fatalln()`, `log.Fatalf()` は
コンソール出力後に os.Exit(1) を呼んでくれるらしい

```go
	// Hello fmt
	fmt.Println("Hello fmt")
	// 2021/07/18 02:21:04 Hello log
	log.Println("Hello log")
```

>Fatal 関連のメソッドはプログラムを終了させるので基本的には main 関数の中で使い
>それ以外の場合はエラーを返して伝搬させていくべき

>Panic 関連の関数は対応する Print 関連の関数でエラーメッセージを表示した後、パニックを起こします。
>パニックは呼び出し元にエラーを伝搬させずにプログラムを終了させるので
>バグに分類される分岐（switch 文など）に到達した場合に使用

__近いうちに log パッケージで遊ぶ必要もありそうだな__
今後 表示として使う fmt と デバッグとして使う log を積極的に使っていこうかな

```go
// init 関数は main 関数の前に実行される初期化関数
func init(){
  log.SetPrefix("[TEST]")  // 接頭辞の設定
}

func main(){
  // [TEST]
  fmt.Println(log.Prefix())  // 接頭辞の取得
  // [TEST]2017/09/28 23:33:59 Hello, world!
  log.Print("Hello, world!")
}
```

参考文献
[Go 言語の log パッケージを使ってみる](https://waman.hatenablog.com/entry/2017/09/29/011614)
