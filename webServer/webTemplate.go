package webServer

import (
	h_template "html/template"
	"log"
	"net/http"
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
