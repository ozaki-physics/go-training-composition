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

// middlewareOne ミドルウェア1層目
func middlewareOne(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("--- URL: " + r.URL.Path)
		log.Println("ミドルウェア1層目")
		next.ServeHTTP(w, r)
		log.Println("ミドルウェア1層目 再度")
	})
}

// middlewareTwo ミドルウェア2層目
func middlewareTwo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("ミドルウェア2層目")
		// ミドルウェアの途中で処理を止める
		// この書き方だと "/" 以外("/img"など)だと 次の層に進まないx
		if r.URL.Path != "/" {
			return
		}
		next.ServeHTTP(w, r)
		log.Println("ミドルウェア2層目 再度")
	})
}

// final ミドルウェア最後の層
func final(w http.ResponseWriter, r *http.Request) {
	log.Println("最後の Handler を実行")
	w.Write([]byte("OK"))
}

// MiddlewareRoot ミドルウェアを実行する関数
// ブラウザでアクセスすると /favicon.ico にもアクセスするから二重になる
func MiddlewareRoot() {
	// 型変換
	finalHandler := http.HandlerFunc(final)
	makeHandle := middlewareOne(middlewareTwo(finalHandler))
	// ここで URL に対応する http.Handler を DefaultServeMux に登録
	http.Handle("/", makeHandle)
	http.ListenAndServe(":8080", nil)
}
