package package02

import (
	"net/http"
)

var Const_big string = "OK_const pkg02_file01"
var const_small string = "NO_const pkg02_file01"

func Sample_server() {
	http.Handle("/", http.FileServer(http.Dir("web")))
	http.ListenAndServe(":8080", nil)
}
