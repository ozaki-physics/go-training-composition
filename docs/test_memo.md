## テストコードはどこに書くのが良いか

テストで1つのパッケージにするのか?  
[Standard Go Project Layout](https://github.com/golang-standards/project-layout/blob/master/README_ja.md) によると  
/test ディレクトリ は  

>追加の外部テストアプリとテストデータ

らしいからテストだけをまとめたパッケージを作るのは違う?  

[標準パッケージの fmt](https://github.com/golang/go/tree/master/src/fmt) では  
わざわざ test をまとめたパッケージは見当たらない  
同じパッケージ内に *_test.go を作っている  
ただ `package fmt_test` という名前がついている  

[他言語プログラマが最低限、気にすべきGoのネーミングルール](https://zenn.dev/keitakn/articles/go-naming-rules)  
標準パッケージでテストコードには `パッケージ名_test` と書かれていた  
コンパイル するときには 別パッケージ扱いになる  
つまり パッケージレベルで非公開のメソッドとかは テストパッケージからは呼び出せなくなる(ブラックボックステスト化)  

ファイル名が アンダースコア や ドット から始まると テストコマンドでも実行されない  

[Goのテストに入門してみよう！](https://future-architect.github.io/articles/20200601/)  

## コマンド
`go test` や `go test -v` だけだと カレントディレクトリ以下のテストが実行される  
`go test パッケージ名` でそのパッケージ名や path だけテストされる  
例: `go test -v ./pkg03/`  

-v は詳細が出てくる  
-run を付けても付けなくてもいい  
-cover でカバレッジが分かる  
-count=1 でテストのキャッシュがされなくなる  

## テストの書く種類
基本的な書き方  
`go test` で実行して欲しいメソッドには `func TestXxx(*testing.T)` としないと認識しない  

```go
func TestAbs(t *testing.T) {
  got := Abs(-1)
  if got != 1 {
      t.Errorf("Abs(-1) = %d; want 1", got)
  }
}
```

[サンプル](pkg03/example03_test.go) を `go test -v ./pkg03/example03_test.go` 実行したときの出力は  

```console
=== RUN   TestAbs
--- PASS: TestAbs (0.00s)
=== RUN   TestReverseAbss
--- PASS: TestReverseAbss (0.00s)
PASS
ok      command-line-arguments  0.005s
```

サブテストのある書き方  
サブテストにすると 特定のサブテストのみ実行や テストの並列化ができる  
`go test ./pkg03/example03_test.go -v -run TestAbs/all_ok`  
`-run` は テストしたいファイルのパスの後ろに書かなければいけない  
テストで半角スペース開けると アンダースコアに変換されてコンソールに出力される  

```go
func TestAbs(t *testing.T) {
	t.Run("いけた?", func(t *testing.T) {
		got := pkg03.Abs(-1)
		if got != 1 {
			t.Errorf("Abs(-1) = %d; want 1", got)
		}
	})
}
```

```console
=== RUN   TestAbs
=== RUN   TestAbs/いけた?
--- PASS: TestAbs (0.00s)
    --- PASS: TestAbs/いけた? (0.00s)
=== RUN   TestReverseAbss
--- PASS: TestReverseAbss (0.00s)
PASS
ok      command-line-arguments  0.009s
```

テスト全体に対する前処理や後処理がある書き方  
テストの書き方は サブテストでも問題ない  
ただし 1つのパッケージの中に `TestMain` は1個しか書けない  

```go
func TestMain(m *testing.M) {
	fmt.Println("test の前処理")
	code := m.Run()
	fmt.Println("test の後処理")
	// Exit を書くのは慣例
	os.Exit(code)
}
```

```console
test の前処理
=== RUN   TestAbs
=== RUN   TestAbs/いけた?
--- PASS: TestAbs (0.00s)
    --- PASS: TestAbs/いけた? (0.00s)
=== RUN   TestReverseAbss
--- PASS: TestReverseAbss (0.00s)
PASS
test の後処理
ok      command-line-arguments  0.004s
```

Go 言語は テストデータのスライスを作って for 文で回すやり方が一般的っぽい [Testing](https://go.dev/doc/code#Testing)  

```go
func TestReverseRunes(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"Hello, world", "dlrow ,olleH"},
		{"Hello, 世界", "界世 ,olleH"},
		{"", ""},
	}
	for _, c := range cases {
		got := ReverseRunes(c.in)
		if got != c.want {
			t.Errorf("ReverseRunes(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}
```

僕は テストしたいメソッドが簡単なものなら サブテスト `t.Run()` しなくて良いと思うが  
1個のメソッドの中でとりあえずどこまで進んだか確かめるためにも 複数アサーションを書くなら 絶対サブテストの方が良いと思う  
つまり 単純なメソッドでもとりあえずサブテストを書いたほうが良いかも  
また expect や actual より want と got で書く方が Go 言語っぽいらしい(確かに前者は長いし)  
