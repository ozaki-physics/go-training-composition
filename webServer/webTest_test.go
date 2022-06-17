package webServer_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/ozaki-physics/go-training-composition/webServer"
)

// TestHttpServer
func TestHttpServer(t *testing.T) {
	handler := http.HandlerFunc(TryTestServer)

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if status := w.Code; status != http.StatusOK {
		t.Errorf(
			"テストが通らなかったよ: 得られたのは (%v) 欲しいのは (%v)",
			status,
			http.StatusOK,
		)
	}

	expected := "hello world !"
	if w.Body.String() != expected {
		t.Errorf(
			"テストが通らなかったよ: 得られたのは (%v) 欲しいのは (%v)",
			w.Body.String(),
			expected,
		)
	}
}

// TestHttpServer02 ちょっと書き方を変えた
// *httptest.ResponseRecorder から Result() で取り出してテストする
func TestHttpServer02(t *testing.T) {
	handler := http.HandlerFunc(TryTestServer)

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	resp := w.Result()
	greeting, _ := io.ReadAll(resp.Body)

	if status := w.Code; status != http.StatusOK {
		t.Errorf(
			"テストが通らなかったよ: 得られたのは (%v) 欲しいのは (%v)",
			status,
			http.StatusOK,
		)
	}

	expected := "hello world !"
	if string(greeting) != expected {
		t.Errorf(
			"テストが通らなかったよ: 得られたのは (%v) 欲しいのは (%v)",
			string(greeting),
			expected,
		)
	}
}
