package presentation

import (
	"fmt"
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
		err := ro.taskHandler.Get(w, r)
		if err != nil {
			http.Error(w, "Internal Server Error\n"+err.Error(), 500)
		}
	case http.MethodPost:
		fmt.Fprintln(w, "POST だよ")
		err := ro.taskHandler.Post(w, r)
		if err != nil {
			http.Error(w, "Internal Server Error\n"+err.Error(), 500)
		}
	case http.MethodPut:
		fmt.Fprintln(w, "PUT だよ")
		err := ro.taskHandler.Put(w, r)
		if err != nil {
			http.Error(w, "Internal Server Error\n"+err.Error(), 500)
		}
	case http.MethodDelete:
		fmt.Fprintln(w, "DELETE だよ")
		err := ro.taskHandler.Delete(w, r)
		if err != nil {
			http.Error(w, "Internal Server Error\n"+err.Error(), 500)
		}
	default:
		fmt.Fprintln(w, "許可されたメソッドじゃないよ")
	}

}
