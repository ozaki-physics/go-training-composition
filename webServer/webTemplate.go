package webServer

import (
	h_template "html/template"
	"log"
	"net/http"
	"os"
	t_template "text/template"
)

// SampleFirstTextTemplate 一番簡単なテンプレートの使い方
func SampleFirstTextTemplate() {
	type Inventory struct {
		Material string
		Count    uint
	}
	sweaters := Inventory{"wool", 17}
	tmpl, err := t_template.New("test").Parse("{{.Count}} items are made of {{.Material}}\n")
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, sweaters)
	if err != nil {
		panic(err)
	}
}

// SampleTextTemplate 公式リファレンスのサンプルのちょっと改造
// See: [text/template にあったサンプル](https://pkg.go.dev/text/template#example-Template) より
func SampleTextTemplate() {
	// Define a template.
	const letter = `
親愛なる {{.Name}},

{{if .Attended}}
結婚式でお会いできて嬉しかったです。
{{- else}}
結婚式に来れなくて残念です。
{{- end}}

{{with .Gift -}}
素敵な {{.}} をありがとうございました。
{{end}}

よろしくお願いします。
ジョシー
`

	// Prepare some data to insert into the template.
	// テンプレートに挿入するデータを用意する。
	type Recipient struct {
		Name, Gift string
		Attended   bool // 出席
	}
	// Recipient 受取人
	var recipients = []Recipient{
		{"ミルドレッドおばさん", "bone china tea set", true},
		{"ジョン叔父さん", "moleskin pants", false},
		{"いとこのロドニー", "", false},
	}

	// Create a new template and parse the letter into it.
	// 新しいテンプレートを作成し、そのテンプレートに手紙を解析します。
	t := t_template.Must(t_template.New("letter").Parse(letter))

	// Execute the template for each recipient.
	// 宛先ごとにテンプレートを実行します。
	for _, r := range recipients {
		// とりあえずコンソールに出力する
		err := t.Execute(os.Stdout, r)
		if err != nil {
			log.Println("executing template:", err)
		}
	}
}

// SampleHtmlTemplate 公式リファレンスのサンプルのちょっと改造
// See: [html/template にあったサンプル](https://pkg.go.dev/html/template#example-package) より
func SampleHtmlTemplate() {
	const tpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.Title}}</title>
	</head>
	<body>
		{{- range .Items}}
			<div>{{ . }}</div>
		{{- else}}
			<div>行なし</div>
		{{- end}}
	</body>
</html>
`

	check := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}
	t, err := h_template.New("webpage").Parse(tpl)
	check(err)

	type base struct {
		Title string
		Items []string
	}

	// 1ページ目
	data := base{
		Title: "1ページ目",
		Items: []string{
			"写真",
			"ブログ",
		},
	}

	// とりあえずコンソールに出力する
	err = t.Execute(os.Stdout, data)
	check(err)

	// 2ページ目
	noItems := base{
		Title: "2ページ目",
		Items: []string{},
	}

	// とりあえずコンソールに出力する
	err = t.Execute(os.Stdout, noItems)
	check(err)
}

// textTemplate ServeHTTP() メソッドを持つための struct
type textTemplate struct {
	fileName string               // ファイル名の格納
	data     interface{}          // html に埋め込む値の構造体
	templ    *t_template.Template // templ コンパイルされたテンプレートの参照を保持
}

// ServeHTTP http.Handle() の引数に渡すため
func (t *textTemplate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.templ = t_template.Must(t_template.ParseFiles(t.fileName))
	// テンプレートをコンパイル(値の埋め込み など)
	t.templ.Execute(w, t.data)
}

// MainTextTemplateServer text/template を使った方法
func MainTextTemplateServer() {
	// html に渡す値たち
	// field は public じゃないと html に埋め込めない
	type data struct {
		Name string
		Age  int
	}

	d := data{Name: "ozaki", Age: 25}
	t := &textTemplate{fileName: "web/template01.html", data: d}
	http.Handle("/template", t)

	// サーバを立てる
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServer:", err)
	}
}

// htmlTemplate ServeHTTP() メソッドを持つための struct
type htmlTemplate struct {
	fileName string               // ファイル名の格納
	data     interface{}          // html に埋め込む値の構造体
	templ    *h_template.Template // templ コンパイルされたテンプレートの参照を保持
}

// ServeHTTP http.Handle() の引数に渡すため
func (t *htmlTemplate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.templ = h_template.Must(h_template.ParseFiles(t.fileName))
	// テンプレートをコンパイル(値の埋め込み など)
	t.templ.Execute(w, t.data)
}

// MainHtmlTemplateServer html/template を使った方法
func MainHtmlTemplateServer() {
	// html に渡す値たち
	type data struct {
		Name string
		Age  int
	}

	d := data{Name: "ozaki", Age: 25}
	t := &htmlTemplate{fileName: "web/template02.html", data: d}
	http.Handle("/template", t)

	// サーバを立てる
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServer:", err)
	}
}

// template02 code injection を試す
type template02 struct {
	fileName string               // ファイル名の格納
	data     interface{}          // html に埋め込む値の構造体
	t_templ  *t_template.Template // templ コンパイルされたテンプレートの参照を保持
	h_templ  *h_template.Template // templ コンパイルされたテンプレートの参照を保持
}

// MainTemplateInjection code injection をしてみる
func MainTemplateInjection() {
	type data struct {
		Name string
		Age  int
	}

	// code injection をしてみる
	d := data{
		Name: "<script>alert(\"インジェクション\")</script>",
		Age:  26,
	}

	t02 := &template02{fileName: "web/template01.html", data: d}

	// コードインジェクションができる
	http.HandleFunc("/template-injection", func(w http.ResponseWriter, r *http.Request) {
		t02.t_templ = t_template.Must(t_template.ParseFiles(t02.fileName))
		// テンプレートをコンパイル(値の埋め込み など)
		t02.t_templ.Execute(w, t02.data)
	})

	// コードインジェクションができない
	http.HandleFunc("/template", func(w http.ResponseWriter, r *http.Request) {
		t02.h_templ = h_template.Must(h_template.ParseFiles(t02.fileName))
		// テンプレートをコンパイル(値の埋め込み など)
		t02.h_templ.Execute(w, t02.data)
	})

	// サーバを立てる
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServer:", err)
	}
}
