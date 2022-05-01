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

func httpGet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "GET メソッドのレスポンスだよ \n")
}
func httpPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "POST メソッドのレスポンスだよ \n")
}

func httpMethod(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		httpGet(w, r)
	case http.MethodPost:
		httpPost(w, r)
	default:
		fmt.Fprintln(w, "許可されたメソッドじゃないよ")
	}
}

// MainHttpMethod http メソッドに応じて処理を分ける
func MainHttpMethod() {
	http.HandleFunc("/", httpMethod)
	http.ListenAndServe(":8080", nil)
}

// MainFileServer 静的ファイルを配信
// http://localhost:8080/ だと web/index.html
// http://localhost:8080/hello.html だと web/hello.html
func MainFileServer() {
	fileHandler := http.FileServer(http.Dir("web"))
	makeHandle := middlewareOne(fileHandler)
	http.Handle("/", makeHandle)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServer:", err)
	}
}

// MainFileServer02 フォルダ階層を省略して URL と対応させる
func MainFileServer02() {
	makeHandle := http.FileServer(http.Dir("web"))
	// makeHandle := http.FileServer(http.Dir("web/sub"))
	// makeHandle = http.StripPrefix("/sub/", makeHandle)
	makeHandle = http.StripPrefix("/aaa/", makeHandle)
	// http.Handle("/", makeHandle)
	http.Handle("/aaa/", makeHandle)
	http.ListenAndServe(":8080", nil)
}
