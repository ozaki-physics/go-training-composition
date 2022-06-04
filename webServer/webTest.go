package webServer

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
)

func TryTestServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world !")
}

// TryResponseRecorder 公式ドキュメントをちょっと改良
// See: https://pkg.go.dev/net/http/httptest#example-ResponseRecorder
// 今はこのメソッドを直接 main で実行して テストする(本当はテストコードに移行する)
func TryResponseRecorder() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "<html><body>Hello World!</body></html>")
	}

	r := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	handler(w, r)

	resp := w.Result()
	greeting, _ := io.ReadAll(resp.Body)

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))

	fmt.Printf("%s\n", greeting)
	// fmt.Println(string(greeting))
}

// TryServer 公式ドキュメントをちょっと改良
// See: https://pkg.go.dev/net/http/httptest#example-Server
// 今はこのメソッドを直接 main で実行して テストする(本当はテストコードに移行する)
func TryServer() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	res, _ := http.Get(ts.URL)
	greeting, _ := io.ReadAll(res.Body)
	res.Body.Close()

	fmt.Printf("%s\n", greeting)
}

// TryServerHttp2 公式ドキュメントをちょっと改良
// See: https://pkg.go.dev/net/http/httptest#example-Server-HTTP2
// 今はこのメソッドを直接 main で実行して テストする(本当はテストコードに移行する)
func TryServerHttp2() {
	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s", r.Proto)
	}))
	ts.EnableHTTP2 = true
	ts.StartTLS()
	defer ts.Close()

	res, _ := ts.Client().Get(ts.URL)
	greeting, _ := io.ReadAll(res.Body)
	res.Body.Close()

	fmt.Printf("%s\n", greeting)
}

// TryTLSServer 公式ドキュメントをちょっと改良
// See: https://pkg.go.dev/net/http/httptest#example-NewTLSServer
// 今はこのメソッドを直接 main で実行して テストする(本当はテストコードに移行する)
func TryTLSServer() {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	res, _ := ts.Client().Get(ts.URL)
	greeting, _ := io.ReadAll(res.Body)
	res.Body.Close()

	fmt.Printf("%s\n", greeting)
}
