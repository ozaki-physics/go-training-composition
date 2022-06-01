package webServer

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
)

func TryTestServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world !")
}

// TryResponseRecorder
// See: https://pkg.go.dev/net/http/httptest#example-ResponseRecorder
func TryResponseRecorder() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "<html><body>Hello World!</body></html>")
	}

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))
}

// TryServer
// See: https://pkg.go.dev/net/http/httptest#example-Server
func TryServer() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}
	greeting, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", greeting)
}
