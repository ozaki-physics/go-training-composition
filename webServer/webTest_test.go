package webServer_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/ozaki-physics/go-training-composition/webServer"
)

func TestHttpServer(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		log.Println(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(TryTestServer)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf(
			"テストが通らなかったよ: 得られたのは (%v) 欲しいのは (%v)",
			status,
			http.StatusOK,
		)
	}

	expected := "hello world !"
	if rr.Body.String() != expected {
		t.Errorf(
			"テストが通らなかったよ: 得られたのは (%v) 欲しいのは (%v)",
			rr.Body.String(),
			expected,
		)
	}
}
