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

## テンプレート構文
### [Overview](https://pkg.go.dev/text/template#pkg-overview) より  

>テンプレートの入力テキストは，UTF-8でエンコードされた任意の形式のテキストである．  
>アクション（データ評価または制御構造）は、"{{"と"}}"で区切られ、アクション以外のすべてのテキストは変更されずに出力にコピーされます。  
>生の文字列を除いて、アクションは改行をまたぐことはできませんが、コメントは可能です。  

### [Text and spaces](https://pkg.go.dev/text/template#hdr-Text_and_spaces) より  

>テンプレートのソースコードを整形するために、アクションの左のデリミタ(デフォルトでは"{{")の直後に マイナス記号と空白が続く場合、直前のテキストからすべての後続の空白が切り取られます。  
>同様に、右のデリミタ("}}")の前に空白とマイナス記号がある場合、直後のテキストからすべての先行する空白が切り取られる。  
>これらのトリムマーカーでは、ホワイトスペースが存在しなければならない。  
>"{{- 3}}" は "{{3}}" のようですが、直前のテキストを切り捨てます。  
>一方 "{{-3}}" は数字の -3 を含むアクションとして解析されます。  

delimiter(デリミタ)は 区切り文字のこと  
`"{{23 -}} < {{- 45}}"` は `"23<45"` になる  

>このトリミングでは、空白文字の定義は Go と同じで、スペース、水平タブ、キャリッジリターン、ニューラインである。  

### [Actions(アクション)](https://pkg.go.dev/text/template#hdr-Actions) より  
>以下はアクションの一覧である。"引数 "と "パイプライン "はデータの評価で、詳細はこの後の対応するセクションで定義されます。

```
コメントの書き方, コメント中には改行を使うことができる
{{/* a comment */}}

前後の空白を削除したコメント
{{- /* a comment with white space trimmed from preceding and following text */ -}}

デフォルトのテキスト, パイプラインのデフォルトテキストで使える
{{pipeline}}

if 文
パイプラインの値が空(false, 0, nil など)なら if の内側は実行されない
{{if pipeline}} T1 {{end}}
{{if pipeline}} T1 {{else}} T0 {{end}}
{{if pipeline}} T1 {{else if pipeline}} T0 {{end}}

for 文
パイプラインの値は 配列, スライス, マップ, チャネル でないといけない
値が0なら else にいく
{{range pipeline}} T1 {{end}}
{{range pipeline}} T1 {{else}} T0 {{end}}

{{range pipeline}} を制御する構文
{{break}}
{{continue}}

指定された名前のテンプレートが、nil データで実行される
{{template "name"}}

指定された名前のテンプレートを実行します
指定された名前のテンプレートが、パイプラインの値にドットセット をパイプラインの値に設定して実行されます
{{template "name" pipeline}}

ブロックは、テンプレートを定義するための略記法
{{define "name"}} T1 {{end}} を定義し、その場 {{template "name" pipeline}} を実行します
{{block "name" pipeline}} T1 {{end}}

if 文と似ている? 使い分けがちょっと分からない
{{with pipeline}} T1 {{end}}
{{with pipeline}} T1 {{else}} T0 {{end}}
```

template の中に if 文なども書けるみたい  
php や jsp に近いため 乱用するとカオスコードになりそうだから乱用しすぎない方がいいかも  

### [Arguments(引数)](https://pkg.go.dev/text/template#hdr-Arguments) より  

>引数は単純な値であり、以下のいずれかで示される。  

真偽値, 文字列, 文字, 整数, 浮動小数, 虚数, 複素数  
. は その時(range, with, 渡された struct そのまま, map そのまま)の値
変数は "$" から始まる英数字の文字  
struct なら ".field"  
マップ なら ".key"  
".method" もできる  

### [Pipelines(パイプライン)](https://pkg.go.dev/text/template#hdr-Pipelines) より  

>パイプラインは、連鎖する可能性のある一連の「コマンド」です。  
>コマンドは単純な値（引数）、または関数やメソッドの呼び出しであり、複数の引数を持つこともあります。  

>パイプラインは、一連のコマンドをパイプライン文字'|'で区切ることによって「連鎖」させることができます。  
>連鎖したパイプラインでは、各コマンドの結果は、次のコマンドの最後の引数として渡されます。  

### [Variables(変数)](https://pkg.go.dev/text/template#hdr-Variables) より  

>$variable := pipeline

Examples  
すべて output と出力される  
```
{{"\"output\""}}
	文字列の定数
{{`"output"`}}
	生の文字列の定数
{{printf "%q" "output"}}
	関数の呼び出し
{{"output" | printf "%q"}}
	最終引数が前のコマンドから来た関数呼び出し。コマンドから来る関数呼び出し。
{{printf "%q" (print "out" "put")}}
	括弧でくくられた引数。
{{"put" | printf "%s%s" "out" | printf "%q"}}
	もっと凝った呼び方。
{{"output" | printf "%s" | printf "%q"}}
	もっと長いチェーン。
{{with "output"}}{{printf "%q" .}}{{end}}
	ドットを使ったwithアクション。
{{with $x := "output" | printf "%q"}}{{$x}}{{end}}
	変数を作成・使用するwithアクション。
{{with $x := "output"}}{{printf "%q" $x}}{{end}}
	変数を他のアクションで使用するwithアクション。
{{with $x := "output"}}{{$x | printf "%q"}}{{end}}
	同じですが、パイプライン化されています。
```

### [Functions(関数)](https://pkg.go.dev/text/template#hdr-Functions) より  

>実行中の関数は、2つの関数マップに格納されます。  
>デフォルトでは、テンプレートに関数は定義されていませんが、 Funcs メソッドを使用して関数を追加することができます。  
>定義済みのグローバル関数の名前は以下の通りです  

```
and
call
html: 引数のテキスト表現をエスケープしたHTMLを返す。 html/template では基本使えない
index
slice
js: 引数のテキスト表現をエスケープしたJavaScriptを返します。
len
not
or
print: fmt.Sprint の alias
printf: fmt.Sprintf の alias
println: fmt.Sprintln の alias
urlquery: URL クエリを埋め込むのに適した形 html/template では基本使えない
eq: ==
ne: !=
lt: <
le: <=
gt: >
ge: >=
```

### [Nested template definitions(ネストしたテンプレート定義)](https://pkg.go.dev/text/template#hdr-Nested_template_definitions) より  

>あるテンプレートを解析するとき、別のテンプレートが定義され、解析中のテンプレートに関連付けられることがあります。  
>テンプレートの定義は、Goプログラムにおけるグローバル変数のように、テンプレートのトップレベルに表示されなければなりません。  
> このような定義の構文は、各テンプレート宣言を "define "と "end "アクションで囲むことです。  

```
const text =`
{{define "T1"}}ONE{{end}}
{{define "T2"}}TWO{{end}}
{{define "T3"}}{{template "T1"}} {{template "T2"}}{{end}}
{{template "T3"}}
`

出力: ONE TWO
```

## サンプル
[text/template にあったサンプル](https://pkg.go.dev/text/template#example-Template)  

[html/template にあったサンプル](https://pkg.go.dev/html/template#example-package)  

### テンプレートの一部共通化
```
{{/* テンプレートを入れ子にする */ -}}
{{define "T1"}}ONE{{end}}
{{define "T2"}}TWO{{end}}
{{define "T3"}}{{template "T1"}} {{template "T2"}}{{end}}
{{template "T3"}}
```
を使う  

3ファイル用意して読み込む  
```go
type d struct {
	Header struct {
		Title    string
		UserName string
	}
	Message string
}

t := h_template.Must(h_template.ParseFiles(
	"web/template03.html",
	"web/template03_header.html",
	"web/template03_footer.html",
))
t.Execute(w, d)
```
読み込む順番は 一番上が大元になる テンプレートじゃないとダメ  

web/template03.html  
```html
{{template "header" .}}
  <p>{{.Message}}</p>
{{template "footer"}}
```
web/template03_header.html  
```html
{{define "header"}}
<!DOCTYPE html>
<html lang="ja">

<head>
  <meta charset="UTF-8">
  <title>{{.Header.Title}}</title>
</head>

<body>
  <h1>{{.Header.Title}}</h1>
  <p>ようこそ {{.Header.UserName}} さん!</p>
{{end}}
```
web/template03_footer.html  
```html
{{define "footer"}}
<div>フッター</div>
</body>

</html>
{{end}}
```

web/template03.html での `{{ . }}` は `d struct` そのものだから  
web/template03_header.html で Title などの値を取り出すときは `{{ .Header.Title }}` と書く必要がある  
冗長でイヤなら `{{template "header" .}}` で `d struct` ごとではなく `d.Header` で渡す  
つまり `{{template "header" .Header}}` にすれば `{{ .Title }}` で値が取り出せるようになる  

### 参考文献
[Go言語(golang) テンプレートの使い方](https://golang.hateblo.jp/entry/golang-text-html-template)  
[Goのテンプレートをちゃんと使ってみる](https://qiita.com/rock619/items/925575c2878b131a16b5)  
[Go の html/template でヘッダーやフッター等の共通化を実現する方法](https://mikan.github.io/2019/12/08/implementing-header-and-footer-with-golang-html-template/)

## テンプレートのコンパイルを1回にする
リクエストごとに処理を行うメソッドで テンプレートファイルを読み込むのは効率が良くない気がする  
よって テンプレートのコンパイルは 1回だけ実行する方が効率的と思われる  
`sync.Once.Do()` を使うことで 複数の goroutine から `ServeHTTP()` メソッドを呼び出されてもコンパイルは1回しか実行しないことを保証する  

また `ServeHTTP()` メソッドの中でテンプレートをコンパイルすると 必要になるまで処理を後回しにできる  
このことを 遅延初期化(lazy initialization)という  
めったに呼ばれない処理の中で遅延初期化が使われているとエラーに気づかないという問題もある  

良い使い方は分からないが 一例として  
```go
type onceTemplate struct {
	fileName string               // ファイル名の格納
	once     sync.Once            // コンパイルするために使う
	templ    *h_template.Template // templ コンパイルされたテンプレートの参照を保持
}

func (t *onceTemplate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(t.fileName))
	})
	t.templ.Execute(w, t.data)
}
```
