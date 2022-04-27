package webServer

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// greet helloweb で補完が動いてコードが生成されたサーバ
func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! %s", time.Now())
}

// Main helloweb で補完が動いてコードが生成されたサーバ
func Main() {
	// http.HandleFunc(a, b) で a は path, b は path のリクエストがきたときに実行する関数
	http.HandleFunc("/", greet)
	// Web サーバを開始
	// ListenAndServe() の第2引数が nil なら DefaultServeMux が Handler として指定される
	http.ListenAndServe(":8080", nil)
}

// MainHtml 直接 html を固定値で配信する
func MainHtml() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
			<html>
				<head>
					<title>go サーバー</title>
				</head>
				<body>
					<h1>直接出力する</h1>
				</body>
			</html>
		`))
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServer:", err)
	}
}

// MainUrl path として "/" が特別なことを確認
// "/" は URL が存在しない path でも とりあえず "/" に飛ばす
func MainUrl() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, World!")
	})
	http.HandleFunc("/img", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "image")
	})
	http.ListenAndServe(":8080", nil)
}
