package training_middleware

import (
	"log"
	"net/http"
)

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
