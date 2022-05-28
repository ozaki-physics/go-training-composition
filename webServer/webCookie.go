package webServer

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

// cookie を格納する
func setCookies(w http.ResponseWriter, r *http.Request) {
	cookie01 := &http.Cookie{
		Name:   "hoge",
		Value:  "bar",
		MaxAge: 30,
	}
	http.SetCookie(w, cookie01)

	cookie02 := &http.Cookie{
		Name:   "hello",
		Value:  "world",
		MaxAge: 30,
	}
	http.SetCookie(w, cookie02)

	fmt.Fprintf(w, "Cookieの設定ができたよ")
}

// cookie を取得して テンプレートに埋め込む
func showCookie(w http.ResponseWriter, r *http.Request) {
	cookie01, err := r.Cookie("hoge")
	if err != nil {
		log.Fatal("Cookie: ", err)
	}

	cookie02, err := r.Cookie("hello")
	if err != nil {
		log.Fatal("Cookie: ", err)
	}

	// すべての Cookie を取得
	cookies := r.Cookies()
	fmt.Println(cookies)

	d := struct {
		C01 *http.Cookie
		C02 *http.Cookie
	}{
		C01: cookie01,
		C02: cookie02,
	}

	tmpl := template.Must(template.ParseFiles("web/cookie.html"))
	tmpl.Execute(w, d)
}

func MainCookie() {
	http.HandleFunc("/cookie-set", setCookies)
	http.HandleFunc("/cookie-get", showCookie)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
