package presentation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Routes 参考記事にはないが echo を使わない代わりに実装する必要が出た
type Routes interface {
	InitRouting(w http.ResponseWriter, r *http.Request)
}

type routes struct {
	taskHandler TaskHandler
}

func NewRoutes(taskHandler TaskHandler) Routes {
	return &routes{taskHandler: taskHandler}
}

// InitRouting routes の初期化
// http.HandleFunc の引数に渡せるのは http.ResponseWriter, *http.Request だけを引数に持ったメソッド
// でも 作った Get() や Post() と紐付けるためには taskHandler struct を渡す必要があった
// だから レシーバで渡そうとしたら 依存性の注入まで気になり いっそ 他のファイルと同じような構造にしてしまおうってなった
func (ro *routes) InitRouting(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		fmt.Fprintln(w, "GET だよ")
		ro.taskHandler.Get(w, r)
	case http.MethodPost:
		body := r.Body
		defer body.Close()

		buf := new(bytes.Buffer)
		io.Copy(buf, body)

		var hello helloJSON
		json.Unmarshal(buf.Bytes(), &hello)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "POST: %+v \n", hello)

		fmt.Fprintln(w, "POST だよ")
		ro.taskHandler.Post(w, r)
	case http.MethodPut:
		fmt.Fprintln(w, "PUT だよ")
		ro.taskHandler.Put(w, r)
	case http.MethodDelete:
		fmt.Fprintln(w, "DELETE だよ")
		ro.taskHandler.Delete(w, r)
	default:
		fmt.Fprintln(w, "許可されたメソッドじゃないよ")
	}
}

type helloJSON struct {
	UserName string `json:"user_name"`
	Content  string `json:"content"`
}
