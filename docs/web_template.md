# 静的ファイルの配信(テンプレートあり)
## 概要
テンプレートを使うと 汎用的なテキストの中に 固有のテキストを埋め込むことができる  
html に go から変数を渡して レスポンスすることができる  
変数は map でも struct でもよい  
そのためテンプレートを利用するのが一般的  
テンプレートを使う場合はテンプレートのコンパイルが必要  
テンプレートのコンパイルとは テンプレートを解釈してデータを埋め込める状態にすること  

Go 標準パッケージに テキスト向けの `text/template` と html 向けの `html/template` がある  
`html/template` だと挿入するコンテキストに 不正なスクリプトの確認, URLで使用できない文字をエンコードすること が可能  
[html/template](https://pkg.go.dev/html/template) より  

>パッケージ template (html/template) は、コードインジェクションに対して安全な HTML 出力を生成するためのデータ駆動型テンプレートを実装しています。  
>これはパッケージ text/template と同じインタフェースを提供し、出力が HTML であるときは常に text/template の代わりに使用されるべきです。  
>このドキュメントは、パッケージのセキュリティ機能に焦点を当てています。  
>テンプレート自体のプログラム方法については、 text/template のドキュメントを参照してください。  

テンプレートの `{{` と `}}` で囲まれた部分はアノテーション(注釈, 目印)を表していて データを埋め込む場所と値を示す  
`{{.Host}}` なら `Template.Execute(w, data)` より `data.Host` の値で置換される  

主に使われている4種を調べる  
- `template.Must()`  
- `template.New()`  
- `func (*Template) Parse()`  
- `func (*Template) Execute()`  

実際は テンプレートの作成(`Template struct`) と 値の埋め込み(`Execute()`) が大事そう  
```go
	// テンプレート作成 と 格納
	t, err := h_template.New("webpage").Parse(tpl)
	// テンプレートに値を埋め込む
	err = t.Execute(os.Stdout, data)
```

### `template.Must()` を もう少し深く調べていく  
[func Must(t *Template, err error) *Template](https://pkg.go.dev/text/template#Must) より  
```go
func Must(t *Template, err error) *Template {
	if err != nil {
		panic(err)
	}
	return t
}
```
>Must は、(*Template, error) を返す関数への呼び出しをラップし、エラーが non-nil である場合にパニック状態にするヘルパーです。  
>これは、次のような変数の初期化で使うことを意図しています。  
>`var t = template.Must(template.New("name").Parse("text"))`  

### `template.New()` を もう少し深く調べていく  
[func New(name string) *Template](https://pkg.go.dev/text/template#New) より  
```go
func New(name string) *Template {
	t := &Template{
		name: name,
	}
	t.init()
	return t
}
```

>New は与えられた名前の新しい未定義のテンプレートを割り当てる。  

struct Template が重要そう  

### `type Template struct` を もう少し深く調べていく  
[type Template](https://pkg.go.dev/text/template#Template) より  
```go
type Template struct {
	name string
	*parse.Tree
	*common
	leftDelim  string
	rightDelim string
}
```

>Template は、解析されたテンプレートの表現である。  
>parse.Tree フィールドは html/template が使用するためにのみエクスポートされ、他のすべてのクライアントは unexported として扱う必要があります。  

解析されたテンプレートごとに名前を設定できるっぽい  

### `func (*Template) Parse()` を もう少し深く調べていく  
[func (t *Template) Parse(text string) (*Template, error)](https://pkg.go.dev/text/template#Template.Parse) より  
```go
func (t *Template) Parse(text string) (*Template, error) {
	t.init()
	t.muFuncs.RLock()
	trees, err := parse.Parse(t.name, text, t.leftDelim, t.rightDelim, t.parseFuncs, builtins())
	t.muFuncs.RUnlock()
	if err != nil {
		return nil, err
	}
	// Add the newly parsed trees, including the one for t, into our common structure.
	for name, tree := range trees {
		if _, err := t.AddParseTree(name, tree); err != nil {
			return nil, err
		}
	}
	return t, nil
}
```

>テキスト中の名前付きテンプレート定義 ({{define ...}} または {{block ...}}文) は t に関連する追加のテンプレートを定義し、 t の定義自体から削除されます。  
>テンプレートは、 Parse の連続した呼び出しで再定義することができます。  
>空白とコメントだけの本文を持つテンプレート定義は空とみなされ、既存のテンプレートの本文を置き換えることはありません。  
>このため、 Parse を使用すると、メインのテンプレート本体を上書きすることなく、新しい名前付きテンプレート定義を追加することができます。  

似ているメソッドに  
- [func (*Template) ParseFS](https://pkg.go.dev/text/template#Template.ParseFS)  
- [func (*Template) ParseFiles](https://pkg.go.dev/text/template#Template.ParseFiles)  
`template.ParseFiles()` は ファイルをテンプレート化するのに使うが 結局は 中で `Parse()` を使っている  
- [func (*Template) ParseGlob](https://pkg.go.dev/text/template#Template.ParseGlob)  




### `func (*Template) Execute()` を もう少し深く調べていく  
[func (*Template) Execute](https://pkg.go.dev/text/template#Template.Execute) より  
```go
func (t *Template) Execute(wr io.Writer, data interface{}) error {
	return t.execute(wr, data)
}

func (t *Template) execute(wr io.Writer, data interface{}) (err error) {
	defer errRecover(&err)
	value, ok := data.(reflect.Value)
	if !ok {
		value = reflect.ValueOf(data)
	}
	state := &state{
		tmpl: t,
		wr:   wr,
		vars: []variable{{"$", value}},
	}
	if t.Tree == nil || t.Root == nil {
		state.errorf("%q is an incomplete or empty template", t.Name())
	}
	state.walk(value, t.Root)
	return
}
```

>実行は、指定されたデータオブジェクトに解析されたテンプレートを適用し、出力を wr に書き込みます。  
>テンプレートの実行や出力の書き込みでエラーが発生した場合，実行は停止しますが，部分的な結果はすでに出力ライタに書き込まれている可能性があります．  
>テンプレートは安全に並列実行できますが，並列実行がライタを共有する場合，出力はインターリーブされる可能性があります．
>data が reflect.Value の場合，テンプレートは fmt.Print のように reflect.Value が保持する具象値に対して適用されます．

実際にテンプレートに値を埋め込むメソッドと思われる  

