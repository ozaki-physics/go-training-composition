package webServer

import (
	"fmt"
	"net/http"
)

func TryTestServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world !")
}
