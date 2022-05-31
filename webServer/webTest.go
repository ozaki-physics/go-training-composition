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
